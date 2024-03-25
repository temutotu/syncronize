package main

import (
	"net"
	"time"
)

type Connection struct {
	conn     net.Conn
	recvChan chan []byte
	sendChan chan []byte
	quitChan chan struct{}
}

func NewConnection() (*Connection, error) {
	tcpconn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return nil, err
	}
	if err := tcpconn.SetReadDeadline(time.Now().Add(300 * time.Second)); err != nil {
		return nil, err
	}

	recvchan := make(chan []byte, 10)
	sendchan := make(chan []byte, 10)
	quitchan := make(chan struct{})

	return &Connection{tcpconn, recvchan, sendchan, quitchan}, nil
}

func (conn *Connection) Process() {
	// send process
	go func() {
		for {
			select {
			case data := <-conn.sendChan:
				//DebugLog(data)
				_, err := conn.conn.Write(data)
				if err != nil {
					close(conn.quitChan)
				}
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// recv process
	for {
		buf := make([]byte, 1024)
		packetSize, err := conn.conn.Read(buf)
		if err != nil {
			close(conn.quitChan)
		}

		if packetSize >= 1 {
			conn.recvChan <- buf
		}
		time.Sleep(100 * time.Millisecond)
	}

	//<-conn.quitChan
}
