package main

import "github.com/gdamore/tcell/v2"

type GameObjects struct {
	ball     *Ball
	hostBar  *SideBar
	guestBar *SideBar
}

type GameManager struct {
	sessionID     byte
	isHost        bool
	recvChan      chan []byte
	sendChan      chan []byte
	commandChan   chan tcell.Key
	commandReady  chan int
	status        int16
	gameObjects   *GameObjects
	ganeObjectPos *SyncronizeObject
	syncManager   *GmaeSyncronize
}

const (
	INIT_PHASE = iota
	WAIT_RESULT_INIT_PHASE
	DRAW_PHASE
	GET_INPUTKEY_PAHSE
	SEND_SYNCHRONIZATION_PACKET_PHASE
	RECV_SYNCHRONIZATION_PACKET_PHASE
	MOVE_OBJECT_PHASE
	END_PHASE
	STAY_PHASE
)

var gameManager *GameManager = nil

func InitGameObjects() *GameObjects {
	gameobjects := &GameObjects{
		ball:     NewBall(0, 0, 45, 160),
		hostBar:  NewSideBar(1, 0, 7, 45),
		guestBar: NewSideBar(157, 0, 7, 45),
	}
	return gameobjects
}

func InitGameManager(conn *Connection, commandChan chan tcell.Key, commandReady chan int) {
	gameManager = &GameManager{0, false, conn.recvChan, conn.sendChan, commandChan, commandReady, INIT_PHASE, InitGameObjects(), &SyncronizeObject{}, &GmaeSyncronize{}}
}

func (manager *GameManager) UpdateGameObjectPos(gameObjectPos *SyncronizeObject) {
	manager.gameObjects.ball.x = int(gameObjectPos.BallXPos)
	manager.gameObjects.ball.y = int(gameObjectPos.BallYPos)
	var xsign int
	if gameObjectPos.BallXVelISNegative == byte(0) {
		xsign = 1
	} else {
		xsign = -1
	}
	manager.gameObjects.ball.xVelocity = int(gameObjectPos.BallXVelocity) * xsign
	var ysign int
	if gameObjectPos.BallYVelISNegative == byte(0) {
		ysign = 1
	} else {
		ysign = -1
	}
	manager.gameObjects.ball.yVelocity = int(gameObjectPos.BallYVelocity) * ysign
	manager.gameObjects.hostBar.y = int(gameObjectPos.HostSideBarYPos)
	manager.gameObjects.guestBar.y = int(gameObjectPos.GuestSideBarYPos)
}
