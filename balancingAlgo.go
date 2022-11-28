package main

import "math/rand"

func SelectServer(servers *[]Server, serverCnt int) (int, bool) {
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
