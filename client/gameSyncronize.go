package main

import "module/syncronize"

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
	sync    *syncronize.ClientSyncronizeManager
	syncObj SyncronizeObject
}

func UpdateSyncObj(gameObj *GameObjects, syncObj *SyncronizeObject) {
	syncObj.BallXPos = uint8(gameObj.ball.x)
	syncObj.BallYPos = uint8(gameObj.ball.y)
	if gameObj.ball.xVelocity < 0 {
		syncObj.BallXVelISNegative = uint8(1)
		syncObj.BallXVelocity = uint8(gameObj.ball.xVelocity * (-1))
	} else {
		syncObj.BallXVelISNegative = uint8(0)
		syncObj.BallXVelocity = uint8(gameObj.ball.xVelocity)
	}

	if gameObj.ball.yVelocity < 0 {
		syncObj.BallYVelISNegative = uint8(1)
		syncObj.BallYVelocity = uint8(gameObj.ball.yVelocity * (-1))
	} else {
		syncObj.BallYVelISNegative = uint8(0)
		syncObj.BallYVelocity = uint8(gameObj.ball.yVelocity)
	}
	syncObj.HostSideBarYPos = uint8(gameObj.hostBar.y)
	syncObj.GuestSideBarYPos = uint8(gameObj.guestBar.y)
}
