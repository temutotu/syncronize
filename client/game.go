package main

import (
	"bytes"
	"encoding/binary"
	"module/syncronize"
	"time"

	"github.com/gdamore/tcell/v2"
)

type GameObject struct {
	hostSideBar  *SideBar
	guestSideBar *SideBar
	ball         *Ball
}

func (gamemanager *GameManager) process() {
	// init game process
	if gameManager.isHost {
		gameManager.sendChan <- packetStartSynchronize(gameManager.sessionID)
	}

	gamemanager.status = WAIT_RESULT_INIT_PHASE

	for {
		switch gameManager.status {
		case WAIT_RESULT_INIT_PHASE:
			select {
			case packet := <-gameManager.recvChan:
				if packet[0] == 201 {
					if packet[1] == 1 {
						packetBody := packet[2:]
						reader := bytes.NewReader(packetBody[:])
						err := binary.Read(reader, binary.LittleEndian, &gamemanager.syncManager.syncObj)
						if err != nil {
							DebugLog(err)
						}
						DebugLog("in recv phase")
						DebugLog(gamemanager.syncManager.syncObj)
						gamemanager.syncManager.sync = syncronize.NewClientSyncronizeManager()
						gameManager.status = MOVE_OBJECT_PHASE
					} else {
						DebugLog("Init game failure")
						return
					}
				}
			}
		case GET_INPUTKEY_PAHSE:
			key := GetInputKey(gameManager.commandChan)
			if key == nil {

			} else {
				switch *key {
				case tcell.KeyUp:
					if gamemanager.isHost {
						gameManager.gameObjects.hostBar.SlideUpSideBar(1)
					} else {
						gameManager.gameObjects.guestBar.SlideUpSideBar(1)
					}
				case tcell.KeyDown:
					if gamemanager.isHost {
						gameManager.gameObjects.hostBar.SlideDownSideBar(1)
					} else {
						gameManager.gameObjects.guestBar.SlideDownSideBar(1)
					}
				default:
				}
			}
			EmptyCommandChannel(gameManager.commandChan)
			gamemanager.status = SEND_SYNCHRONIZATION_PACKET_PHASE
		case SEND_SYNCHRONIZATION_PACKET_PHASE:
			UpdateSyncObj(gameManager.gameObjects, &gamemanager.syncManager.syncObj)
			DebugLog(gamemanager.syncManager.syncObj)
			packetBody := syncronize.ConvertToByte(gamemanager.syncManager.syncObj)
			packet := make([]byte, 1024)
			packet[0] = syncronize.SYNCHRONIZE
			packet[1] = gamemanager.sessionID
			packet[2] = gameManager.syncManager.sync.GetSyncNum()
			if gameManager.isHost {
				packet[3] = 1
			} else {
				packet[3] = 0
			}
			len := len(packetBody)
			copy(packet[4:], packetBody[:len])
			gamemanager.sendChan <- packet
			gamemanager.status = RECV_SYNCHRONIZATION_PACKET_PHASE
		case RECV_SYNCHRONIZATION_PACKET_PHASE:
			select {
			case packet := <-gameManager.recvChan:
				if packet[0] == syncronize.RESULT_SYNCHRONIZE {
					if packet[1] == gameManager.syncManager.sync.GetSyncNum() { // check syncNum
						packetBody := packet[2:]
						reader := bytes.NewReader(packetBody[:])
						err := binary.Read(reader, binary.LittleEndian, &gamemanager.syncManager.syncObj)
						if err != nil {
							DebugLog(err)
						}
						DebugLog(gamemanager.syncManager.syncObj)
						gamemanager.syncManager.sync.RecvServerPacket()
						gameManager.status = MOVE_OBJECT_PHASE
					} else {
						DebugLog("sync num is not match")
						return
					}
				}
			}
		case MOVE_OBJECT_PHASE:
			gamemanager.UpdateGameObjectPos(&gamemanager.syncManager.syncObj)
			gamemanager.gameObjects.ball.MoveBall()
			gameManager.status = DRAW_PHASE
		case DRAW_PHASE:
			drawManager.CreanGameScreen()
			drawManager.DrawBall(nil, *gameManager.gameObjects.ball)
			drawManager.DrawSideBar(nil, *gameManager.gameObjects.hostBar)
			drawManager.DrawSideBar(nil, *gamemanager.gameObjects.guestBar)
			gamemanager.status = STAY_PHASE
		case END_PHASE:
			return
		case STAY_PHASE:
			if len(gamemanager.commandReady) <= 0 {
				gamemanager.commandReady <- 1
			}
			time.Sleep(15 * time.Millisecond)
			gamemanager.status = GET_INPUTKEY_PAHSE
		}
	}
}

func GetInputKey(commandChan chan tcell.Key) *tcell.Key {
	select {
	case command := <-commandChan:
		return &command
	default:
		return nil
	}
}

func EmptyCommandChannel(commandChan chan tcell.Key) {
	for {
		select {
		case <-commandChan:

		default:
			return
		}
	}
}
