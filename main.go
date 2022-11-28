package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Server struct {
	Addr  string
	Alive bool
}

func Handler(conn net.Conn, servers *[]Server, serverCnt int) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	sid, error := SelectServer(servers, serverCnt)
	if error == true {
		log.Println("no server available!")
		return
	}
	serverConn, err := net.Dial("tcp", (*servers)[sid].Addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(serverConn net.Conn) {
		err := serverConn.Close()
		if err != nil {
			log.Println(err)
		}
	}(serverConn)
	go io.Copy(serverConn, conn)
	io.Copy(conn, serverConn)
}

func main() {
	servers, serverCnt := ReadServerList()
	go PingRoutine(&servers, serverCnt, time.Second*10)
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide port number")
		return
	}
	PORT := ":" + arguments[1]
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go Handler(c, &servers, serverCnt)
	}
}
