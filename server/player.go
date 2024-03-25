package main

type Player struct {
	PlayerID   [4]byte
	IsHost     bool
	SideBarPos int
	SessionID  byte
	sendChan   chan []byte
}

func NewPlayer(ID [4]byte, sendChan chan []byte) *Player {
	return &Player{ID, false, 23, 0, sendChan}
}
