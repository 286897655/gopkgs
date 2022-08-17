package av

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Extension RTP Header extension
type RtpExtension struct {
	id      uint8
	payload []byte
}

// RtpHeader represents an RTP packet header
type RtpHeader struct {
	Version          uint8
	Padding          bool
	Extension        bool
	Marker           bool
	PayloadType      uint8
	SequenceNumber   uint16
	Timestamp        uint32
	SSRC             uint32
	CSRC             []uint32
	ExtensionProfile uint16
	Extensions       []RtpExtension
}

// Packet represents an RTP Packet
type RtpPacket struct {
	RtpHeader
	Payload     []byte
	PaddingSize byte
}

const (
	headerLength            = 4
	versionShift            = 6
	versionMask             = 0x03
	paddingShift            = 5
	paddingMask             = 0x01
	extensionShift          = 4
	extensionMask           = 0x01
	extensionProfileOneByte = 0xBEDE
	extensionProfileTwoByte = 0x1000
	extensionIDReserved     = 0x0F
	ccMask                  = 0x0F
	markerShift             = 7
	markerMask              = 0x01
	ptMask                  = 0x7F
	seqNumOffset            = 2
	seqNumLength            = 2
	timestampOffset         = 4
	timestampLength         = 4
	ssrcOffset              = 8
	ssrcLength              = 4
	csrcOffset              = 12
	csrcLength              = 4
)

// Unmarshal parses the passed byte slice and stores the result in the Header.
// It returns the number of bytes read n and any error.
func (h *RtpHeader) Unmarshal(buf []byte) (n int, err error) { //nolint:gocognit
	if len(buf) < headerLength {
		return 0, fmt.Errorf("Rtp header size insufficient:%d < %d", len(buf), headerLength)
	}

	/*
	 *  0                   1                   2                   3
	 *  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |V=2|P|X|  CC   |M|     PT      |       sequence number         |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |                           timestamp                           |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |           synchronization source (SSRC) identifier            |
	 * +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
	 * |            contributing source (CSRC) identifiers             |
	 * |                             ....                              |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 */

	h.Version = (buf[0] >> versionShift) & versionMask
	h.Padding = ((buf[0] >> paddingShift) & paddingMask) > 0
	h.Extension = ((buf[0] >> extensionShift) & extensionMask) > 0
	nCSRC := int(buf[0] & ccMask)
	if cap(h.CSRC) < nCSRC || h.CSRC == nil {
		h.CSRC = make([]uint32, nCSRC)
	} else {
		h.CSRC = h.CSRC[:nCSRC]
	}

	n = csrcOffset + (nCSRC * csrcLength)
	if len(buf) < n {
		return n, fmt.Errorf("Rtp header size insufficient:%d < %d", len(buf), n)
	}

	h.Marker = (buf[1] >> markerShift & markerMask) > 0
	h.PayloadType = buf[1] & ptMask

	h.SequenceNumber = binary.BigEndian.Uint16(buf[seqNumOffset : seqNumOffset+seqNumLength])
	h.Timestamp = binary.BigEndian.Uint32(buf[timestampOffset : timestampOffset+timestampLength])
	h.SSRC = binary.BigEndian.Uint32(buf[ssrcOffset : ssrcOffset+ssrcLength])

	for i := range h.CSRC {
		offset := csrcOffset + (i * csrcLength)
		h.CSRC[i] = binary.BigEndian.Uint32(buf[offset:])
	}

	if h.Extensions != nil {
		h.Extensions = h.Extensions[:0]
	}

	if h.Extension {
		if expected := n + 4; len(buf) < expected {
			return n, fmt.Errorf("RTP header insufficient size for extension:%d < %d", len(buf), expected)
		}

		h.ExtensionProfile = binary.BigEndian.Uint16(buf[n:])
		n += 2
		extensionLength := int(binary.BigEndian.Uint16(buf[n:])) * 4
		n += 2

		if expected := n + extensionLength; len(buf) < expected {
			return n, fmt.Errorf("RTP header insufficient size for extension:%d < %d", len(buf), expected)
		}

		switch h.ExtensionProfile {
		// RFC 8285 RTP One Byte Header Extension
		case extensionProfileOneByte:
			end := n + extensionLength
			for n < end {
				if buf[n] == 0x00 { // padding
					n++
					continue
				}

				extid := buf[n] >> 4
				len := int(buf[n]&^0xF0 + 1)
				n++

				if extid == extensionIDReserved {
					break
				}

				extension := RtpExtension{
					id:      extid,
					payload: buf[n : n+len]}
				h.Extensions = append(h.Extensions, extension)
				n += len
			}

		// RFC 8285 RTP Two Byte Header Extension
		case extensionProfileTwoByte:
			end := n + extensionLength
			for n < end {
				if buf[n] == 0x00 { // padding
					n++
					continue
				}

				extid := buf[n]
				n++

				len := int(buf[n])
				n++

				extension := RtpExtension{
					id:      extid,
					payload: buf[n : n+len]}
				h.Extensions = append(h.Extensions, extension)
				n += len
			}

		default: // RFC3550 Extension
			if len(buf) < n+extensionLength {
				return n, fmt.Errorf("RTP header insufficient size for extension:%d < %d", len(buf), n+extensionLength)
			}

			extension := RtpExtension{
				id:      0,
				payload: buf[n : n+extensionLength]}
			h.Extensions = append(h.Extensions, extension)
			n += len(h.Extensions[0].payload)
		}
	}
	return n, nil
}

// MarshalSize returns the size of the header once marshaled.
func (h *RtpHeader) MarshalSize() int {
	// NOTE: Be careful to match the MarshalTo() method.
	size := 12 + (len(h.CSRC) * csrcLength)

	if h.Extension {
		extSize := 4

		switch h.ExtensionProfile {
		// RFC 8285 RTP One Byte Header Extension
		case extensionProfileOneByte:
			for _, extension := range h.Extensions {
				extSize += 1 + len(extension.payload)
			}
		// RFC 8285 RTP Two Byte Header Extension
		case extensionProfileTwoByte:
			for _, extension := range h.Extensions {
				extSize += 2 + len(extension.payload)
			}
		default:
			extSize += len(h.Extensions[0].payload)
		}

		// extensions size must have 4 bytes boundaries
		size += ((extSize + 3) / 4) * 4
	}

	return size
}

func (h *RtpHeader) MarshalTo(buf []byte) (n int, err error) {
	/*
	 *  0                   1                   2                   3
	 *  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |V=2|P|X|  CC   |M|     PT      |       sequence number         |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |                           timestamp                           |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |           synchronization source (SSRC) identifier            |
	 * +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
	 * |            contributing source (CSRC) identifiers             |
	 * |                             ....                              |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 */

	size := h.MarshalSize()
	if size > len(buf) {
		return 0, io.ErrShortBuffer
	}

	// The first byte contains the version, padding bit, extension bit,
	// and csrc size.
	buf[0] = (h.Version << versionShift) | uint8(len(h.CSRC))
	if h.Padding {
		buf[0] |= 1 << paddingShift
	}

	if h.Extension {
		buf[0] |= 1 << extensionShift
	}

	// The second byte contains the marker bit and payload type.
	buf[1] = h.PayloadType
	if h.Marker {
		buf[1] |= 1 << markerShift
	}

	binary.BigEndian.PutUint16(buf[2:4], h.SequenceNumber)
	binary.BigEndian.PutUint32(buf[4:8], h.Timestamp)
	binary.BigEndian.PutUint32(buf[8:12], h.SSRC)

	n = 12
	for _, csrc := range h.CSRC {
		binary.BigEndian.PutUint32(buf[n:n+4], csrc)
		n += 4
	}

	if h.Extension {
		extHeaderPos := n
		binary.BigEndian.PutUint16(buf[n+0:n+2], h.ExtensionProfile)
		n += 4
		startExtensionsPos := n

		switch h.ExtensionProfile {
		// RFC 8285 RTP One Byte Header Extension
		case extensionProfileOneByte:
			for _, extension := range h.Extensions {
				buf[n] = extension.id<<4 | (uint8(len(extension.payload)) - 1)
				n++
				n += copy(buf[n:], extension.payload)
			}
		// RFC 8285 RTP Two Byte Header Extension
		case extensionProfileTwoByte:
			for _, extension := range h.Extensions {
				buf[n] = extension.id
				n++
				buf[n] = uint8(len(extension.payload))
				n++
				n += copy(buf[n:], extension.payload)
			}
		default: // RFC3550 Extension
			extlen := len(h.Extensions[0].payload)
			if extlen%4 != 0 {
				// the payload must be in 32-bit words.
				return 0, io.ErrShortBuffer
			}
			n += copy(buf[n:], h.Extensions[0].payload)
		}

		// calculate extensions size and round to 4 bytes boundaries
		extSize := n - startExtensionsPos
		roundedExtSize := ((extSize + 3) / 4) * 4

		binary.BigEndian.PutUint16(buf[extHeaderPos+2:extHeaderPos+4], uint16(roundedExtSize/4))

		// add padding to reach 4 bytes boundaries
		for i := 0; i < roundedExtSize-extSize; i++ {
			buf[n] = 0
			n++
		}
	}

	return n, nil
}

// Marshal serializes the header into bytes.
func (h *RtpHeader) Marshal() (buf []byte, err error) {

	buf = make([]byte, h.MarshalSize())
	n, err := h.MarshalTo(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

// Unmarshal parses the passed byte slice and stores the result in the Packet.
func (p *RtpPacket) Unmarshal(buf []byte) error {
	n, err := p.RtpHeader.Unmarshal(buf)
	if err != nil {
		return err
	}
	end := len(buf)
	if p.Padding {
		p.PaddingSize = buf[end-1]
		end -= int(p.PaddingSize)
	}
	if end < n {
		return errors.New("buffer too small")
	}
	p.Payload = buf[n:end]
	return nil
}

// Marshal serializes the packet into bytes.
func (p *RtpPacket) Marshal() (buf []byte, err error) {
	buf = make([]byte, p.MarshalSize())

	n, err := p.MarshalTo(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

// MarshalTo serializes the packet and writes to the buffer.
func (p *RtpPacket) MarshalTo(buf []byte) (n int, err error) {
	p.RtpHeader.Padding = p.PaddingSize != 0
	n, err = p.RtpHeader.MarshalTo(buf)
	if err != nil {
		return 0, err
	}

	// Make sure the buffer is large enough to hold the packet.
	if n+len(p.Payload)+int(p.PaddingSize) > len(buf) {
		return 0, io.ErrShortBuffer
	}

	m := copy(buf[n:], p.Payload)
	if p.RtpHeader.Padding {
		buf[n+m+int(p.PaddingSize-1)] = p.PaddingSize
	}

	return n + m + int(p.PaddingSize), nil
}

// MarshalSize returns the size of the packet once marshaled.
func (p *RtpPacket) MarshalSize() int {
	return p.RtpHeader.MarshalSize() + len(p.Payload) + int(p.PaddingSize)
}

// String helps with debugging by printing packet information in a readable way
func (p *RtpPacket) FmtString() string {
	out := "RTP PACKET:\n"

	out += fmt.Sprintf("\tVersion: %v\n", p.Version)
	out += fmt.Sprintf("\tMarker: %v\n", p.Marker)
	out += fmt.Sprintf("\tPayload Type: %d\n", p.PayloadType)
	out += fmt.Sprintf("\tSequence Number: %d\n", p.SequenceNumber)
	out += fmt.Sprintf("\tTimestamp: %d\n", p.Timestamp)
	out += fmt.Sprintf("\tSSRC: %d (%x)\n", p.SSRC, p.SSRC)
	out += fmt.Sprintf("\tPayload Length: %d\n", len(p.Payload))

	return out
}
