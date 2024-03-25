package main

import "errors"

const MAX_SESSION int = 10

var sessionManager *SessionManager = nil

type SessionManager struct {
	sessionList []*Session
	gcChan      chan int
	syncChan    chan int
}

// sessionListからマッチング可能なsessionを取得する
func (manager *SessionManager) GetMatchMakingSession() *Session {
	for _, session := range manager.sessionList {
		if session == nil {
			continue
		}

		if session.Status != MatchMaking {
			continue
		}
		session.Status = GuestJoining
		return session
	}
	return nil
}

func (manager *SessionManager) GetSession(sessionID int) *Session {
	if sessionID >= MAX_SESSION {
		return nil
	}

	return manager.sessionList[sessionID]
}

// sessionListの空きに新たなsessionを保存する
func (manager *SessionManager) SetSession(newSession *Session) (int, error) {
	for i, session := range manager.sessionList {
		if session == nil {
			newSession.ID = byte(i)
			newSession.Status = MatchMaking
			manager.sessionList[i] = newSession
			return i, nil
		}
	}
	return -1, errors.New("sessionList is full")
}

// 指定したIDのsessionを解散する
func (manager *SessionManager) DismissSession(sessionID byte) (*Player, *Player, error) {
	if sessionID >= byte(MAX_SESSION) {
		return nil, nil, errors.New("sessionID is invalid")
	}

	session := sessionManager.sessionList[sessionID]
	if session == nil {
		return nil, nil, errors.New("session is not exist")
	}

	host := *(session.Host)
	guest := *(session.Guest)
	sessionManager.sessionList[sessionID] = nil
	return &host, &guest, nil
}
