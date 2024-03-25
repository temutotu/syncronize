package main

import (
	"module/syncronize"
)

func packetSimpleResult(protocol byte, result byte) []byte {
	packet := make([]byte, 1024)
	packet[0] = protocol
	packet[1] = result
	return packet
}

// 1          2           3       4~
// protocol   joinResult  isHost  sessionID
func packetJoinResult(result *ResultJoin) []byte {
	protocol := make([]byte, 1024)
	protocol[0] = 0
	protocol[1] = result.result
	protocol[2] = result.ishost
	protocol[3] = result.sessionID
	return protocol
}

func packetNoticeGuestJoined() []byte {
	packet := make([]byte, 1024)
	packet[0] = 1
	packet[1] = 1
	packet[2] = 1
	return packet
}

func packageResultInitSyncronize(result byte, syncObj []byte) []byte {
	packet := make([]byte, 1024)
	packet[0] = syncronize.RESULT_INIT_SYNCHRONIZE
	packet[1] = result
	copy(packet[2:], syncObj[:])
	return packet
}

func packetResultSyncronize(syncNum int, packetBody []byte) []byte {
	packet := make([]byte, 1024)
	packet[0] = syncronize.RESULT_SYNCHRONIZE
	packet[1] = byte(syncNum)
	copy(packet[2:], packetBody[:1022])
	return packet
}
