package main

import (
	"github.com/tatsushid/go-fastping"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func PingServers(servers *[]Server, serverCnt int) {
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

func PingRoutine(servers *[]Server, serverCnt int, interval time.Duration) {
	for {
		PingServers(servers, serverCnt)
		time.Sleep(interval)
	}
}
