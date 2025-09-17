package main

import (
	"flag"
	"fmt"
	"hash/crc32"
	"log"
	"net"
)

func main() {
	crcdata := "123456789"

	log.Println("crc32:", crc32.Checksum([]byte(crcdata), crc32.IEEETable))
	log.Println("crc32:", crc32.ChecksumIEEE([]byte(crcdata)))
	return
	var port int
	var ip string
	flag.IntVar(&port, "p", 123456, "server listen port or client dial port,default is 12345")
	flag.StringVar(&ip, "i", "0.0.0.0", "server listen ip,default use 0.0.0.0")

	flag.Parse()

	// 只listen 不accept
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Println("tcp listen fail")
	}
	log.Println(listener.Addr().String())

	select {}

}
