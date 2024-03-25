package main

import (
	"errors"
	"fmt"
)

type Session struct {
	ID         byte    // マッチID
	Host       *Player // ホストユーザ
	Guest      *Player // ゲストユーザー
	Syncronize *GmaeSyncronize
	Status     int16
}

const (
	Init = iota + 1
	MatchMaking
	GuestJoining
	Playing
	Finished
)

func NewSession(host *Player) (game *Session) {
	game = &Session{0, host, nil, nil, Init}
	return
}

func (session *Session) JoinGuestSession(guest *Player) error {
	if session.Guest != nil {
		return errors.New("joining game is full player")
	}

	session.Guest = guest
	session.Status = Playing
	// hostにguestのjoin通知をする
	if session.Host == nil {
		return errors.New("host is nil")
	}

	if session.Host.sendChan == nil {
		return errors.New("Host chan is nil")
	}

	session.Host.sendChan <- packetNoticeGuestJoined()
	return nil
}

func (session *Session) StartSyncronize() {

	if err := InitSynchronize(session); err != nil {
		fmt.Println(err)
		return
	}

	if err := session.Syncronize.syncronize.Process(session.Syncronize); err != nil {
		fmt.Println(err)
		return
	}
}
