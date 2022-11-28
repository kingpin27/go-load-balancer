package main

import (
	"bufio"
	"log"
	"os"
)

func ReadServerList() ([]Server, int) {
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
