package main

import (
	"fmt"
)

func joinSession(buff []byte, sendChan chan []byte) (*ResultJoin, error) {
	// パケットの長さチェック
	playerID := buff[1:5]
	newPlayer := NewPlayer([4]byte(playerID), sendChan)
	// sessionManagerから参加可能なsessionを探す
	joinSession := sessionManager.GetMatchMakingSession()
	if joinSession == nil {
		// session作成してホストになる
		newPlayer.IsHost = true
		newSession := NewSession(newPlayer)
		sessionID, err := sessionManager.SetSession(newSession)
		if err != nil {
			return nil, err
		}
		// sessionIDをクライアントに返す
		fmt.Println("make new session and host join sessionID:", sessionID)
		return &ResultJoin{SUCCESS, HOST, byte(sessionID)}, nil
	} else {
		// guestとして参加
		if err := joinSession.JoinGuestSession(newPlayer); err != nil {
			return nil, err
		}
		fmt.Println("guest join sessionID:", joinSession.ID)
		return &ResultJoin{SUCCESS, GUSET, byte(joinSession.ID)}, nil
	}
}

func leaveSession(buff []byte) error {
	// パケットの長さチェック

	sessionID := buff[0]
	reason := buff[1]
	host, guest, err := sessionManager.DismissSession(sessionID)
	if err != nil {
		return err
	}

	fmt.Println(host, guest, reason)
	return nil
}
