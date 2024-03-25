package main

import (
	"errors"
	"fmt"
	"math/rand"
	"module/syncronize"
	"net"
)

var PACKET_BYTE_SIZE int = 1024

type Connection struct {
	ID       int
	conn     *net.TCPConn
	recvChan chan []byte
	sendChan chan []byte
	recvBuf  []byte
	sendBuf  []byte
}

func NewConnection(conn *net.TCPConn) *Connection {
	connection := &Connection{
		ID:       rand.Intn(100),
		conn:     conn,
		recvChan: make(chan []byte, 10),
		sendChan: make(chan []byte, 10),
		recvBuf:  make([]byte, PACKET_BYTE_SIZE),
		sendBuf:  make([]byte, PACKET_BYTE_SIZE),
	}

	return connection
}

func (conn *Connection) process() {
	fmt.Println("process start")
	// send proccess
	go func() {
		for {
			select {
			case packet := <-conn.sendChan:
				if _, err := conn.conn.Write(packet); err != nil {
					fmt.Println("send packet is failed")
				}
			default:
			}
		}
	}()
	// recv process
	for {
		n, err := conn.conn.Read(conn.recvBuf)
		if err != nil {
			fmt.Println(err)
			return
		}

		if n != 0 {
			fmt.Println("recv:", conn.recvBuf[:10])
			if err := conn.handleRecvProtocolProtocol(); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func (conn *Connection) handleRecvProtocolProtocol() error {
	protocol := conn.recvBuf[0]
	fmt.Println("read protocol:", protocol)
	switch protocol {
	case 0:
		result, err := joinSession(conn.recvBuf, conn.sendChan)
		if err != nil {
			fmt.Println(err)
			conn.sendChan <- packetSimpleResult(0, FAILD)
		} else {
			fmt.Println("sendpacket")
			conn.sendChan <- packetJoinResult(result)
			fmt.Println(sessionManager.sessionList[0])
		}
	case syncronize.INIT_SYNCHRONIZE:
		sessionID := conn.recvBuf[1]
		sessionManager.syncChan <- int(sessionID)
	case syncronize.SYNCHRONIZE:
		sessuinID := conn.recvBuf[1]
		session := sessionManager.GetSession(int(sessuinID))
		if session == nil {
			return errors.New("session is nil")
		}
		session.Syncronize.syncronize.SetRecvPacket(conn.recvBuf[:])
		if err := session.Syncronize.syncronize.RecvSynchronizePacket(conn.recvBuf[:]); err != nil {
			return err
		}
	default:
	}
	return nil
}
