package main

func packetStartSynchronize(sessionID byte) []byte {
	packet := make([]byte, 1024)
	packet[0] = byte(200) // INIT_SYNCHRONIZE
	packet[1] = sessionID
	return packet
}
