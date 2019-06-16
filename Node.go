package main

import (
	"math/rand"
	//"fmt"
)

type Node struct {
	id 					int

	channels 			[]*Channel

	pathPN 				map[int][]int
	pathPH 				map[int][]int
	pathCH 				map[int][]int
	pathVC 				map[int][]int

	txChan 				chan Transaction
	doneChan 			chan bool

	lockedChannel		map[int]*Channel
	lockedChChannel		map[int]*Channel

	pn 					*PaymentNetwork

	htlcChan 			chan HTLC
	htlcReplyChan 		chan HTLCReply

	ph 					*PaymentHub
	phiouChan  			chan PHIOU
	phconfChan			chan PHConfirm

	ch 					*ChannelHub
	chchannels 			[]*Channel

	pccChan 			chan PCC
	gccChan 			chan GCC
	chiouChan 			chan CHIOU
	chconf1Chan 		chan CHConf1
	chresChan 			chan CHRes
	chconf2Chan 		chan CHConf2

	secretChan 			chan Transaction
}

var nodeIdGenerator int = 0

func NewNode() *Node {
	node := Node{id: nodeIdGenerator,
				channels: []*Channel{},
				pathPN: make(map[int][]int),
				pathPH: make(map[int][]int),
				pathCH: make(map[int][]int),
				pathVC: make(map[int][]int),
				txChan: make(chan Transaction),
				lockedChannel: make(map[int]*Channel),
				lockedChChannel: make(map[int]*Channel),
				htlcChan: make(chan HTLC),
				htlcReplyChan: make(chan HTLCReply),
				phiouChan: make(chan PHIOU),
				phconfChan: make(chan PHConfirm),
				secretChan: make(chan Transaction),
				chchannels: []*Channel{},
				pccChan: make(chan PCC),
				gccChan: make(chan GCC),
				chiouChan: make(chan CHIOU),
				chconf1Chan: make(chan CHConf1),
				chresChan: make(chan CHRes),
				chconf2Chan: make(chan CHConf2)}
	nodeIdGenerator++
	return &node
}

// func (n *Node) setPNPath(destination int, path []int) {
// 	n.pathPN[destination] = path
// }

// func (n *Node) setPHPath(destination int, path []int) {
// 	n.pathPH[destination] = path
// }

// func (n *Node) setCHPath(destination int, path []int) {
// 	n.pathCH[destination] = path
// }

// func (n *Node) setVCPath(destination int, path []int) {
// 	n.pathVC[destination] = path
// }

// func (n *Node) getPathPN(destination int) []int {
// 	return n.pathPN[destination]
// }

// func (n *Node) getPathPH(destination int) []int {
// 	return n.pathPH[destination]
// }

// func (n *Node) getPathCH(destination int) []int {
// 	return n.pathCH[destination]
// }

// func (n *Node) getPathVC(destination int) []int {
// 	return n.pathVC[destination]
// }

// func (n *Node) setDoneChan(done chan bool) {
// 	n.doneChan = done
// }

// func (n *Node) setPaymentNetwork(pn *PaymentNetwork) {
// 	n.pn = pn
// }

// func (n *Node) setPaymentHub(ph *PaymentHub) {
// 	n.ph = ph
// }

// func (n *Node) setChannelHub(ch *ChannelHub) {
// 	n.ch = ch
// }

// func (n *Node) execTx(tx Transaction) {
// 	//fmt.Println("execTxPN", n.id, tx)
// 	n.txChan <- tx
// }

// func (n *Node) lockChannel(tx *Transaction, channel *Channel) {
// 	//fmt.Println("lock", n.id, tx.id, channel)
// 	channel.lock()
// 	channel.lockTx = tx
// 	n.lockedChannel[tx.id] = channel
// }

// func (n *Node) lockCHChannel(tx *Transaction, channel *Channel) {
// 	//fmt.Println("lockch", n.id, tx.id, channel)
// 	channel.lock()
// 	channel.lockTx = tx
// 	n.lockedChChannel[tx.id] = channel
// }

// func (n *Node) unlockChannel(tx *Transaction) {
// 	if channel, ok := n.lockedChannel[tx.id]; ok {
// 		//fmt.Println("unlock", n.id, tx.id, channel)
// 		channel.lockTx = nil
// 		channel.unLock()
// 		delete(n.lockedChannel, tx.id)
// 	}
// }

// func (n *Node) unlockCHChannel(tx *Transaction) {
// 	if channel, ok := n.lockedChChannel[tx.id]; ok {
// 		//fmt.Println("unlockch", n.id, tx.id, channel)
// 		channel.lockTx = nil
// 		channel.unLock()
// 		delete(n.lockedChChannel, tx.id)
// 	}
// }

// func (n *Node) revealSecret(tx Transaction) {
// 	communicationDelay()
// 	n.secretChan <- tx
// }

// func (n *Node) doSecret(tx Transaction) {
// 	//fmt.Println("dosecret", n.id, tx)
// 	if n.id == tx.from {
// 		n.unlockChannel(&tx)
// 		n.doneChan <- true
// 		return
// 	}
// 	prevNode := -1
// 	for i, node := range tx.path {
// 		if n.id == node {
// 			prevNode = tx.path[i-1]
// 			break
// 		}
// 	}
// 	if tx.t == PHTX && !n.pn.containChannel(n.id, prevNode) && n.ph.contain(n.id) && n.ph.contain(prevNode) {
// 		n.ph.receipt(PHReceipt{from: prevNode, to: n.id, tx: &tx})
// 		<-n.phconfChan
// 		//conf := <-n.phconfChan
// 		//fmt.Println("phconf", n.id, conf)
// 	}
// 	if tx.t == CHTX && !n.pn.containChannel(n.id, prevNode) && n.ch.containNode(n.id) && n.ch.containNode(prevNode) {
// 		n.ch.receipt(CHReceipt{from: prevNode, to: n.id, tx: &tx})
// 		conf := <-n.chconf1Chan
// 		go n.CHdoConf1(conf)
// 		//n.pn.getNode(prevNode).revealSecret(tx)
// 		//return
// 	}
// 	n.unlockChannel(&tx)
// 	n.pn.getNode(prevNode).revealSecret(tx)
// }



// func (n *Node) setTxPath(tx *Transaction) {
// 	switch tx.t {
// 	case PNTX:
// 		tx.setPath(n.getPathPN(tx.to))
// 	case PHTX:
// 		tx.setPath(n.getPathPH(tx.to))
// 	case CHTX:
// 		tx.setPath(n.getPathCH(tx.to))
// 	case VCTX:
// 		tx.setPath(n.getPathVC(tx.to))
// 	}
// }

// func (n *Node) doTx(tx Transaction) {
// 	if n.id == tx.from {
// 		n.setTxPath(&tx)
// 		//fmt.Println("path", n.id, tx.id, tx.getPath())
// 		secret := int(rand.Int31())
// 		secretHash := hash(secret)
// 		tx.setSecretHash(secretHash)
// 		revealChan := make(chan bool, 1)
// 		tx.setRevealChan(revealChan)
// 		go func() {
// 			<-revealChan
// 			tx.secret = secret
// 			n.pn.getNode(tx.to).revealSecret(tx)
// 		}()
// 	}

// 	prevNode := -1
// 	nextNode := -1
// 	pathIndex := 0
// 	for ; pathIndex < len(tx.path); pathIndex++ {
// 		if n.id == tx.path[pathIndex] {
// 			break
// 		}
// 	}
// 	if n.id != tx.from {
// 		prevNode = tx.path[pathIndex-1]
// 	}
// 	if n.id != tx.to {
// 		nextNode = tx.path[pathIndex+1]
// 	}


// 	usePaymentChannel := true
// 	switch tx.t {
// 	case PHTX:
// 		if prevNode >= 0 && !n.pn.containChannel(n.id, prevNode) && n.ph.contain(n.id) && n.ph.contain(prevNode) {
// 			<-n.phiouChan
// 		}
// 		if !n.pn.containChannel(n.id, nextNode) && n.ph.contain(n.id) && n.ph.contain(nextNode) {
// 			usePaymentChannel = false
// 			n.ph.iou(PHIOU{from: n.id, to: nextNode, tx: &tx})
// 		}
// 	case CHTX:
// 		if prevNode >= 0 && !n.pn.containChannel(n.id, prevNode) && n.ch.containNode(n.id) && n.ch.containNode(prevNode) {
// 			iou := <- n.chiouChan
// 			//fmt.Println("iou", n.id, iou)
// 			channelBD := n.getCHChannel()
// 			n.lockCHChannel(iou.tx, channelBD)
// 			peer := channelBD.getPeerId(n.id)
// 			go func() {
// 				n.pn.getNode(peer).CHpcc(PCC{from: n.id, to: peer, tx: iou.tx})
// 				<-n.gccChan
// 			} ()
// 			//gcc := <-n.gccChan
// 			//fmt.Println("gcc", n.id, gcc)
// 		}
// 		if !n.pn.containChannel(n.id, nextNode) && n.ch.containNode(n.id) && n.ch.containNode(nextNode) {
// 			usePaymentChannel = false
// 			channelAC := n.getCHChannel()
// 			n.lockCHChannel(&tx, channelAC)
// 			peer := channelAC.getPeerId(n.id)
// 			n.pn.getNode(peer).CHpcc(PCC{from: n.id, to: peer, tx: &tx})
// 			gcc := <-n.gccChan
// 			//fmt.Println("gcc", n.id, gcc)
// 			n.ch.iou(CHIOU{from: n.id, to: nextNode, tx: &tx, gcc: gcc})
// 		}
// 	}
// 	if n.id == tx.to {
// 		tx.getRevealChan() <- true
// 		return
// 	}
// 	if usePaymentChannel {
// 		channel := n.pn.getChannel(n.id, nextNode)
// 		n.lockChannel(&tx, channel)
// 		n.pn.getNode(nextNode).htlc(HTLC{from: n.id, to: nextNode, tx: &tx})
// 		<-n.htlcReplyChan
// 	}
// 	n.pn.getNode(nextNode).execTx(tx)
// }

func (n *Node) getRandomChannel() *Channel {
	return n.channels[rand.Intn(len(n.channels))]
}

func (n *Node) registerChannel(channel *Channel) {
	n.channels = append(n.channels, channel)
}

// //
// // HTLC
// //
// func (n *Node) htlc(htlc HTLC) {
// 	communicationDelay()
// 	n.htlcChan <- htlc
// }

// func (n *Node) doHTLC(htlc HTLC) {
// 	n.pn.getNode(htlc.from).htlcReply(HTLCReply{from: htlc.from, to: htlc.to, tx: htlc.tx})
// }

// func (n *Node) htlcReply(htlcReply HTLCReply) {
// 	communicationDelay()
// 	n.htlcReplyChan <- htlcReply
// }


// //
// //	Payment Hub
// //
// func (n *Node) PHiou(iou PHIOU) {
// 	communicationDelay()
// 	//fmt.Println("phiou", n.id, iou)
// 	n.phiouChan <- iou
// }

// func (n *Node) PHdoIOU(iou PHIOU) {
// 	//fmt.Println("phdoiou", n.id, iou)
// }

// func (n *Node) PHconf(conf PHConfirm) {
// 	communicationDelay()
// 	n.phconfChan <- conf
// }

// func (n *Node) PHdoConf(conf PHConfirm) {
// 	//fmt.Println("phdoconf", n.id, conf)
// }


// //
// //	Channel Hub
// //
// func (n *Node) joinChannelHub(ch *ChannelHub, channel *Channel) {
// 	n.chchannels = append(n.chchannels, channel)
// }

// func (n *Node) getCHChannel() *Channel {
// 	return n.chchannels[rand.Intn(len(n.chchannels))]
// }

// func (n *Node) CHpcc(pcc PCC) {
// 	communicationDelay()
// 	n.pccChan <- pcc
// }

// func (n *Node) CHgcc(gcc GCC) {
// 	communicationDelay()
// 	n.gccChan <- gcc
// }

// func (n *Node) CHiou(iou CHIOU) {
// 	communicationDelay()
// 	n.chiouChan <- iou
// }

// func (n *Node) CHConf1(conf CHConf1) {
// 	communicationDelay()
// 	n.chconf1Chan <- conf
// }

// func (n *Node) CHRes(res CHRes) {
// 	communicationDelay()
// 	n.chresChan <- res
// }

// func (n *Node) CHConf2(conf CHConf2) {
// 	communicationDelay()
// 	n.chconf2Chan <- conf
// }


// func (n *Node) CHdoPCC(pcc PCC) {
// 	//fmt.Println("dopcc", n.id, pcc)
// 	n.pn.getNode(pcc.from).CHgcc(GCC{from: pcc.from, to: pcc.to, tx: pcc.tx})
// }


// //func (n *Node) CHdoIOU(iou CHIOU) {
// //	//fmt.Println("doiou", n.id, iou)
// //	channelBD := n.getCHChannel()
// //	channelBD.lock()
// //	channelBD.lockTx = iou.tx
// //	peer := channelBD.getPeerId(n.id)
// //	n.pn.getNode(peer).CHpcc(PCC{from: n.id, to: peer, tx: iou.tx})
// //	<-n.gccChan
// //}

// func (n *Node) CHdoRes(res CHRes) {
// 	//fmt.Println("dores", n.id, res)
// 	n.pn.getNode(res.reply).CHConf2(CHConf2{})
// }

// func (n *Node) CHdoConf1(conf CHConf1) {
// 	//fmt.Println("doconf1", n.id, conf)
// 	channel := n.lockedChChannel[conf.tx.id]
// 	peer := channel.getPeerId(n.id)
// 	n.pn.getNode(peer).CHRes(CHRes{from: conf.from, to: conf.to, tx: conf.tx, reply: n.id})
// 	<- n.chconf2Chan
// 	//conf2 := <- n.chconf2Chan
// 	//fmt.Println("conf2", n.id, conf2)
// 	n.unlockCHChannel(conf.tx)
// }


// func (n *Node) NodeCron() {
// 	for {
// 		select {
// 		case tx := <-n.txChan:
// 			n.doTx(tx)

// 		case htlc := <-n.htlcChan:
// 			n.doHTLC(htlc)

// 		case conf := <-n.phconfChan:
// 			n.PHdoConf(conf)

// 		case pcc := <-n.pccChan:
// 			n.CHdoPCC(pcc)
// 		case conf := <-n.chconf1Chan:
// 			go n.CHdoConf1(conf)
// 		case res := <-n.chresChan:
// 			n.CHdoRes(res)

// 		case tx := <-n.secretChan:
// 			n.doSecret(tx)
// 		}
// 	}
// }