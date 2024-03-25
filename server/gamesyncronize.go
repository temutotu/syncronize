package main

import (
	"errors"
	"fmt"
	"module/syncronize"
)

const INIT_BALL_X_POS uint8 = 120
const INIT_BALL_Y_POS uint8 = 20
const INIT_BALL_X_VELOCITY uint8 = 1
const INIT_BALL_Y_VELOCITY uint8 = 1
const INIT_SIDEBAR_Y_POS uint8 = 23

var GAME_INITIAL_STATE SyncronizeObject = SyncronizeObject{INIT_BALL_X_POS, INIT_BALL_Y_POS, 0, INIT_BALL_X_VELOCITY, 0, INIT_BALL_Y_VELOCITY, INIT_SIDEBAR_Y_POS, INIT_SIDEBAR_Y_POS}

type SyncronizeObject struct {
	BallXPos           uint8
	BallYPos           uint8
	BallXVelISNegative uint8
	BallXVelocity      uint8
	BallYVelISNegative uint8
	BallYVelocity      uint8
	HostSideBarYPos    uint8
	GuestSideBarYPos   uint8
}

type GmaeSyncronize struct {
	sessionID  byte
	object     *SyncronizeObject
	syncronize *syncronize.SyncronizeBase
}

func InitSynchronize(session *Session) error {
	session.Syncronize = nil

	session.Syncronize = &GmaeSyncronize{
		sessionID:  session.ID,
		object:     &GAME_INITIAL_STATE,
		syncronize: syncronize.NewSyncronizeBase(),
	}

	// host guetsゲームの初期状態を送信
	packetBody := syncronize.ConvertToByte(session.Syncronize.object)
	resultPacket := packageResultInitSyncronize(SUCCESS, packetBody)
	fmt.Println(resultPacket)
	session.Host.sendChan <- resultPacket
	session.Guest.sendChan <- resultPacket
	return nil
}

func (sync *GmaeSyncronize) SendResult() error {
	packetBody := sync.syncronize.GetPacketBody()
	fmt.Println(int(sync.syncronize.GetSyncNum()))
	packet := packetResultSyncronize(int(sync.syncronize.GetSyncNum()), packetBody)
	session := sessionManager.GetSession(int(sync.sessionID))
	if session == nil {
		return errors.New("session is nil")
	}
	fmt.Println(packet)
	session.Host.sendChan <- packet
	session.Guest.sendChan <- packet
	return nil
}

// 同期処理を実行
func (sync *GmaeSyncronize) ExecSyncronize(hostPacket []byte, guestPacket []byte) error {
	var hostSyncObj SyncronizeObject
	var guestSyncObj SyncronizeObject
	fmt.Println(hostPacket)
	fmt.Println(guestPacket)
	syncronize.ConvertToStruct(hostPacket[:], &hostSyncObj)
	syncronize.ConvertToStruct(guestPacket[:], &guestSyncObj)

	sync.object = &hostSyncObj
	sync.object.GuestSideBarYPos = guestSyncObj.GuestSideBarYPos

	packageBody := syncronize.ConvertToByte(*sync.object)
	fmt.Println(packageBody)
	sync.syncronize.SetPacketBody(packageBody)
	if err := sync.syncronize.SendSynchronozeResult(sync); err != nil {
		return err
	}

	return nil
}
