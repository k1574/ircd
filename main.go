package main

import(
	"log"
	"net"
	"fmt"
)

var(
	AddrStr = "localhost:6667"
)

func
HandleConn(conn net.Conn) {
}

func
main() {
	ln, err := net.Listen("tcp", AddrStr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(conn.RemoteAddr())
		go HandleConn(conn);
	}
}
