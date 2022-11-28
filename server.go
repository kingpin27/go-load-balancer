package main

import (
	"bufio"
	"fmt"
	"github.com/tatsushid/go-fastping"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Addr  string
	Alive bool
}

func pingServers(servers *[]Server, serverCnt int) {
	var wg sync.WaitGroup
	for i := 0; i < serverCnt; i++ {
		wg.Add(1)
		go func(server *[]Server, i int) {
			defer wg.Done()
			host := strings.Split((*servers)[i].Addr, ":")[0]
			p := fastping.NewPinger()
			ra, err := net.ResolveIPAddr("ip4:icmp", host)
			if err != nil {
				log.Println(err)
				log.Println((*servers)[i].Addr, "not responding.")
				(*servers)[i].Alive = false
				os.Exit(1)
			}
			p.AddIPAddr(ra)
			p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
				//fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
			}
			p.OnIdle = func() {
				(*servers)[i].Alive = true
				return
			}
			err = p.Run()
			if err != nil {
				log.Println((*servers)[i].Addr, "not responding.")
				(*servers)[i].Alive = false
				log.Println(err)
			}
		}(servers, i)
	}
	wg.Wait()
}

func readServerList() ([]Server, int) {
	servers := make([]Server, 100)
	serverCnt := 0
	f, err := os.Open("list.txt")
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(f)
	if err != nil {
		log.Panicln(err)
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		serverAddr := scanner.Text()
		servers[serverCnt] = Server{Addr: serverAddr, Alive: false}
		log.Println("server", serverCnt, ":", serverAddr)
		serverCnt++
	}
	return servers, serverCnt
}

func pingRoutine(servers *[]Server, serverCnt int, interval time.Duration) {
	for {
		pingServers(servers, serverCnt)
		time.Sleep(interval)
	}
}

func selectServer(servers *[]Server, serverCnt int) (int, bool) {

	sid := rand.Intn(serverCnt)
	t := 0
	for (*servers)[sid].Alive == false {
		sid++
		if sid == serverCnt {
			sid = 0
		}
		t++
		if t >= serverCnt {
			return -1, true
		}
	}
	return sid, false
}

func Handler(conn net.Conn, servers *[]Server, serverCnt int) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
	}(conn)
	sid, error := selectServer(servers, serverCnt)
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
	servers, serverCnt := readServerList()
	go pingRoutine(&servers, serverCnt, time.Second*10)
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
