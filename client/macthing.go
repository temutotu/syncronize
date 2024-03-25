package main

func (manager *GameManager) matchingStart() error {
	data := []byte{0, 0, 0, 0, 0}
	manager.sendChan <- data

	var isJoined bool = false
	for !isJoined {
		select {
		case packet := <-manager.recvChan:
			if packet[0] == 0 {
				if packet[1] == 1 {

					if packet[2] == 1 {
						manager.isHost = true
					} else {
						manager.isHost = false
					}

					manager.sessionID = packet[3]
					isJoined = true
				} else {
					DebugLog("join session failed")
					return nil
				}
			}
		default:
		}
	}

	// recv process
	if !manager.isHost {
		drawManager.NewDrawLog(0, "you are guest so match start soon")
		return nil
	} else {
		drawManager.NewDrawLog(0, "you are host so watining guest")
	}

	var isJoinedGuest bool = false
	for !isJoinedGuest {
		packet := <-manager.recvChan
		DebugLog(packet)
		if packet[0] == NoticeJoinSession {
			if packet[1] == 1 && packet[2] == 1 {
				drawManager.NewDrawLog(0, "guest joined session so start match")
				isJoinedGuest = true
			}
		}
	}

	return nil
}
