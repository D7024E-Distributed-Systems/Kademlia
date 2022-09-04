package network

/**
 * ping = PING
 * find contact =FICO
 * find data = FIDA
 * store message = STME
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
