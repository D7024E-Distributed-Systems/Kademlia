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

/*
Returns a Ping reference instance
*/
func newPing() *Ping {
	return &Ping{"PING"}
}

type FindContact struct {
	startMessage string
}

/*
Returns a find contact reference instance
*/
func newFindContact() *FindContact {
	return &FindContact{"FICO"}
}

type FindData struct {
	startMessage string
}

/*
Returns a find data reference instance
*/
func newFindData() *FindData {
	return &FindData{"FIDA"}
}

type StoreMessage struct {
	startMessage string
}

/*
Returns a store message reference instance
*/
func newStoreMessage() *StoreMessage {
	return &StoreMessage{"STME"}
}

type RefreshMessage struct {
	startMessage string
}

/*
Returns a refresh message reference instance
*/
func newRefreshMessage() *RefreshMessage {
	return &RefreshMessage{"REFR"}
}
