package main

import (
	"fmt"
	"net"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	tcplistener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	// ゼロ値はnil
	sessionManager = &SessionManager{make([]*Session, MAX_SESSION), make(chan int), make(chan int)}

	// main loop
	go func() {
		for {
			select {
			case sessionID := <-sessionManager.syncChan:
				session := sessionManager.GetSession(sessionID)
				if session != nil {
					go session.StartSyncronize()
				}
			}
		}
	}()

	// make connection loop
	func() {
		for {
			// tcpconnを複数所持する
			tcpconn, err := tcplistener.AcceptTCP()
			if err != nil {
				fmt.Println(err)
				return
			}

			conn := NewConnection(tcpconn)
			if conn == nil {
				fmt.Println("connection is nil")
				return
			}

			go conn.process()
		}
	}()
}
