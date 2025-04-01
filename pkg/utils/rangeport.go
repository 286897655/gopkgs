package utils

import (
	"errors"
	"net"
	"sync"
)

type RangePort struct {
	tcp_ports sync.Map
	udp_ports sync.Map
}

func New(port_min int, port_max int) *RangePort {
	if port_min >= port_max {
		panic("port min is larger than port max")
	}
	range_port := &RangePort{}
	for i := port_min; i <= port_max; i++ {
		range_port.tcp_ports.Store(i, 0)
		range_port.udp_ports.Store(i, 0)
	}
	return range_port
}

func (range_port *RangePort) SelectTcpPort() (int, error) {
	select_port := 0
	range_port.tcp_ports.Range(func(key, value any) bool {
		port := key.(int)
		used := value.(int)
		if used > 0 { //if used this port continue
			return true
		}
		listener, err := net.ListenTCP("tcp", &net.TCPAddr{
			Port: port,
		})
		if err != nil {
			// listen fail,port was used select netx
			return false
		}
		defer listener.Close()
		select_port = port
		return true
	})
	if select_port == 0 {
		return 0, errors.New("can't select unused port")
	}
	range_port.tcp_ports.Store(select_port, 1)
	return select_port, nil
}

func (range_port *RangePort) SelectUdpPort() {

}
