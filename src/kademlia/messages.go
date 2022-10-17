package kademlia

/**
 * ping = PING
 * find contact =FICO
 * find data = FIDA
 * store message = STME
 * refresh message = REFR
 */

type Ping struct {
	startMessage string
}

func newPing() *Ping {
	return &Ping{"PING"}
}

type FindContact struct {
	startMessage string
}

func newFindContact() *FindContact {
	return &FindContact{"FICO"}
}

type FindData struct {
	startMessage string
}

func newFindData() *FindData {
	return &FindData{"FIDA"}
}

type StoreMessage struct {
	startMessage string
}

func newStoreMessage() *StoreMessage {
	return &StoreMessage{"STME"}
}

type RefreshMessage struct {
	startMessage string
}

func newRefreshMessage() *RefreshMessage {
	return &RefreshMessage{"REFR"}
}
