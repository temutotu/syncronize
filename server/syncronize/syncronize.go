package syncronize

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Syncronizer interface {
	SendResult() error
	ExecSyncronize(hostPacket []byte, guestpacket []byte) error
}

type SyncronizeBase struct {
	syncNum         byte
	isRecvHost      bool
	isRecvGuest     bool
	recvHostPcaket  []byte
	recvGuestPacket []byte
	packetBody      []byte
	readyChan       chan int
	finishChan      chan int
}

func NewSyncronizeBase() *SyncronizeBase {

	return &SyncronizeBase{
		syncNum:         0,
		isRecvHost:      false,
		isRecvGuest:     false,
		recvHostPcaket:  make([]byte, 1024),
		recvGuestPacket: make([]byte, 1024),
		packetBody:      make([]byte, 1024),
		readyChan:       make(chan int),
		finishChan:      make(chan int),
	}
}

func (sync *SyncronizeBase) GetSyncNum() byte {
	return sync.syncNum
}

func (sync *SyncronizeBase) SetRecvPacket(packet []byte) error {
	len := len(packet[4:])
	fmt.Println(packet)
	if int(packet[3]) == 1 {
		copy(sync.recvHostPcaket[:len], packet[4:])
	} else {
		copy(sync.recvGuestPacket[:len], packet[4:])
	}
	return nil
}

func (sync *SyncronizeBase) SetPacketBody(packet []byte) error {
	copy(sync.packetBody, packet)
	return nil
}

func (sync *SyncronizeBase) GetPacketBody() []byte {
	return sync.packetBody
}

func (syncronize *SyncronizeBase) Process(syncronier Syncronizer) error {
	for {
		select {
		case <-syncronize.readyChan:
			if err := syncronier.ExecSyncronize(syncronize.recvHostPcaket, syncronize.recvGuestPacket); err != nil {
				return err
			}
		case <-syncronize.finishChan:
			return nil
		}
	}
}

func (syncronize *SyncronizeBase) RecvSynchronizePacket(pakcet []byte) error {
	if pakcet[0] != SYNCHRONIZE {
		return errors.New("syncronize: prtocool is not mactch")
	}

	if pakcet[2] != syncronize.syncNum {
		return errors.New("syncNum is missMatch")
	}

	if pakcet[3] == 0 {
		syncronize.isRecvHost = true
		//copy(pakcet, syncronize.recvHostPcaket)
	} else if pakcet[3] == 1 {
		syncronize.isRecvGuest = true
		//copy(pakcet, syncronize.recvGuestPacket)
	}
	if syncronize.isRecvGuest && syncronize.isRecvHost {
		syncronize.readyChan <- 1
	}
	return nil
}

func (syncronize *SyncronizeBase) SendSynchronozeResult(syncronier Syncronizer) error {
	if err := syncronier.SendResult(); err != nil {
		return err
	}
	fmt.Println("sync sended")
	syncronize.isRecvHost = false
	syncronize.isRecvGuest = false

	if int(syncronize.syncNum) >= 255 {
		syncronize.syncNum = 0
	} else {
		fmt.Println("sync is increment")
		syncronize.syncNum = syncronize.syncNum + byte(1)
	}
	return nil
}

func ConvertToByte(object any) []byte {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.LittleEndian, object)
	return buf.Next(buf.Len())
}

func ConvertToStruct(packet []byte, object any) {
	reader := bytes.NewReader(packet[:])
	binary.Read(reader, binary.LittleEndian, object)
}
