package main


type Channel struct {
	id 			int
	A 			int
	B 			int
	capacity 	int

	lockChan 	chan struct{}
	lockTx 		*Transaction
}

var channelIdGenerator int = 0

func NewChannel(A, B, capacity int) *Channel {
	c := Channel{id: channelIdGenerator,
				A: A,
				B: B,
				capacity: capacity,
				lockChan: make(chan struct{}, 1),
				lockTx: nil}
	channelIdGenerator++
	return &c
}

func (c *Channel) getPeerId(nodeId int) int {
	if c.A == nodeId {
		return c.B
	} else if c.B == nodeId {
		return c.A
	}
	return -1
}

func (c *Channel) lock() {
	c.lockChan <- struct{}{}
}

func (c *Channel) unLock() {
	select {
	case <- c.lockChan:
		return
	default:
		panic("unlock of unlocked channel")
	}
}

func (c *Channel) tryLock() bool {
	select {
	case c.lockChan <- struct{}{}:
		return true
	default:
		return false
	}
}