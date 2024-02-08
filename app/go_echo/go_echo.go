package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"time"
)

func server_pong(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), "read from client:", conn.RemoteAddr(), " failed, err: ", err)
			break
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), "reads from client:", conn.RemoteAddr(), "data:", string(buf[:n]))
		conn.Write(buf[:n]) // 把收到的数据发送数据
	}
}

func runTcpServer(port int) {
	fmt.Println("Tcp Server listen at:", port)
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

		go server_pong(conn)
	}
}

func runUdpServer(port int) {
	fmt.Println("Udp Server listen at:", port)
	udpserver, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	})

	if err != nil {
		fmt.Println("Server listen err:", err)
		return
	}

	defer udpserver.Close()

	for {
		buf := [128]byte{}
		n, addr, err := udpserver.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Println("udp server read err:", err)
			return
		}

		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), "Read from client:", addr, "data:", string(buf[:n]))
		time.Sleep(time.Duration(1) * time.Second)
		_, err = udpserver.WriteToUDP(buf[:n], addr)
		if err != nil {
			fmt.Println("udp server write err:", err)
			return
		}
	}
}

func runTcpClient(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Println("Tcp Connect remote:", addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), "Connect remote:", addr, " failed,err:", err)
		return
	}
	defer conn.Close()

	// client 先输入"Ping-Pong"
	_, err = conn.Write([]byte("Ping-Pong"))

	for {
		buf := [128]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Println("read failed, err:", err)
			return
		}
		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), "Read from server:", string(buf[:n]))
		time.Sleep(time.Duration(1) * time.Second)
		conn.Write(buf[:n])
	}
}

func runUdpClient(host string, port int) {
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Println("Udp send to remote:", addr)

	udp, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Printf("create udp client to remote:%s err.", addr)
		return
	}

	// client 先输入"Ping-Pong"
	_, err = udp.Write([]byte("Ping-Pong"))
	if err != nil {
		fmt.Println("udp client write err:", err)
		return
	}

	for {
		buf := [128]byte{}
		n, err := udp.Read(buf[:])
		if err != nil {
			fmt.Println("udp client read err:", err)
			return
		}

		fmt.Println(time.Now().Format("2006-01-02 15:04:05.000"), "Read from server:", string(buf[:n]))
		time.Sleep(time.Duration(1) * time.Second)
		_, err = udp.Write(buf[:n])
		if err != nil {
			fmt.Println("udp client write err:", err)
			return
		}
	}
}

func main() {
	var server_mode bool
	var dial_port int
	var dial_ip string
	var udp_mode bool
	flag.BoolVar(&server_mode, "s", false, "default is false,use server mode use -s")
	flag.IntVar(&dial_port, "p", 9021, "server listen port or client dial port,default is 9021")
	flag.StringVar(&dial_ip, "h", "127.0.0.1", "dial ip of server for client")
	flag.BoolVar(&udp_mode, "u", false, "use udp mode")

	flag.Parse()

	if server_mode {
		fmt.Println("echo-ping-pong server mode.")
		if udp_mode {
			runUdpServer(dial_port)
		} else {
			runTcpServer(dial_port)
		}
	} else {
		fmt.Println("echo-ping-pong client mode.")
		if udp_mode {
			for i := 0; i < 10; i++ {
				go runUdpClient(dial_ip, dial_port)
			}
		} else {
			runTcpClient(dial_ip, dial_port)
		}
	}

	select {}
}
