package syncronize

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type syncObj interface {
}

type Syncronizer interface {
	SendResult(packet []byte) error
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

	if pakcet[1] != syncronize.syncNum {
		return errors.New("syncNum is missMatch")
	}

	if pakcet[2] == 0 {
		syncronize.isRecvHost = true
		copy(pakcet, syncronize.recvHostPcaket)
	} else if pakcet[2] == 1 {
		syncronize.isRecvGuest = true
		copy(pakcet, syncronize.recvGuestPacket)
	}

	if syncronize.isRecvGuest && syncronize.isRecvHost {
		syncronize.readyChan <- 1
	}
	return nil
}

func (syncronize *SyncronizeBase) SendSynchronozeResult(syncronier Syncronizer) error {
	if err := syncronier.SendResult(syncronize.packetBody); err != nil {
		return err
	}

	syncronize.isRecvHost = false
	syncronize.isRecvGuest = false

	if syncronize.syncNum >= 255 {
		syncronize.syncNum = 0
	} else {
		syncronize.syncNum++
	}
	return nil
}

func ConvertToByte(object any) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, object)
	return buf.Next(1024)
}

func ConverToStruct(buf []byte) {
	// reader := bytes.NewReader(buf[:6])
	// binary.Read(reader, binary.LittleEndian, obj)
	// return obj
}
