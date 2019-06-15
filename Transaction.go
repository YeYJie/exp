package main

type TransactionType int
const (
	PNTX = iota
	PHTX
	CHTX
	VCTX
)


type Transaction struct {
	id 			int
	t 			TransactionType
	from 		int
	to			int
	path 		[]int
	secret 		int
	secretHash 	int
	revealChan 	chan bool
}

func (t *Transaction) setPath(p []int) {
	t.path = p
}

func (t *Transaction) getPath() []int {
	return t.path
}

func (t *Transaction) setSecretHash(h int) {
	t.secretHash = h
}

func (t *Transaction) setRevealChan(revealChan chan bool) {
	t.revealChan = revealChan
}

func (t *Transaction) getRevealChan() chan bool {
	return t.revealChan
}