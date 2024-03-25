package syncronize

type ClientSyncronizeManager struct {
	syncNum byte
}

func NewClientSyncronizeManager() *ClientSyncronizeManager {
	return &ClientSyncronizeManager{0}
}

func (sync *ClientSyncronizeManager) RecvServerPacket() error {
	sync.syncNum = sync.syncNum + byte(1)
	return nil
}

func (sync *ClientSyncronizeManager) GetSyncNum() byte {
	return sync.syncNum
}
