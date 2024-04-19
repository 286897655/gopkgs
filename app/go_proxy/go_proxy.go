package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"time"
)

var proxy_addr string

func proxy(conn net.Conn) {
	defer conn.Close()

	fmt.Println("make conn for:", proxy_addr)
	proxy_conn, err := net.Dial("tcp", proxy_addr)
	if err != nil {
		fmt.Println(err)
		return
	}

	go io.Copy(proxy_conn, conn)
	go io.Copy(conn, proxy_conn)

	select {}
}

func runTcpServer(port int) {
	fmt.Println("proxy Tcp Server listen at:", port)

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), "Server listen error:", err)
		return
	}

	for {
		conn, err := listener.Accept() // accpet
		if err != nil {
			fmt.Println("Accept() error:", err)
			continue
		}

		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), "Server accept client:", conn.RemoteAddr())

		go proxy(conn)
	}
}

func main() {
	var listen_port int

	flag.IntVar(&listen_port, "l", 8080, "proxy listen server port")
	flag.StringVar(&proxy_addr, "h", "127.0.0.1:80", "proxy addr")

	flag.Parse()

	runTcpServer(listen_port)

	select {}
}
