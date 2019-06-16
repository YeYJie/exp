package main

//import "fmt"

type PHIOU struct {
	from 			int
	to 				int
	tx 				*Transaction
}

type PHReceipt struct {
	from 			int
	to				int
	tx 				*Transaction
}

type PHConfirm struct {
	from 			int
	to 				int
	tx 				*Transaction
}

type PaymentHub struct {
	paymentNetwork 	*PaymentNetwork2
	nodeMap 		map[int]*Node

	iouChan 		chan PHIOU
	receiptChan 	chan PHReceipt
}

func NewPaymentHub(pn *PaymentNetwork2) *PaymentHub {
	return &PaymentHub{paymentNetwork: pn,
					nodeMap: make(map[int]*Node),
					iouChan: make(chan PHIOU),
					receiptChan: make(chan PHReceipt)}
}

func (ph *PaymentHub) join(node *Node) {
	for _, n := range ph.nodeMap {
		ph.paymentNetwork.addEdge(node.id, n.id, 1)
	}
	ph.nodeMap[node.id] = node
}

// func (ph *PaymentHub) contain(nodeId int) bool {
// 	_, ok := ph.nodeMap[nodeId]
// 	return ok
// }

// func (ph *PaymentHub) iou(iou PHIOU) {
// 	communicationDelay()
// 	ph.iouChan <- iou
// }

// func (ph *PaymentHub) receipt(receipt PHReceipt) {
// 	communicationDelay()
// 	ph.receiptChan <- receipt
// }

// func (ph *PaymentHub) PaymentHubCron() {
// 	for {
// 		select {
// 		case iou := <-ph.iouChan:
// 			//fmt.Println("PaymentHub.PaymentHubCron", "PHiou:", iou)
// 			to := ph.paymentNetwork.getNode(iou.to)
// 			to.PHiou(iou)
// 		case receipt := <-ph.receiptChan:
// 			//fmt.Println("PaymentHub.PaymentHubCron", "receipt:", receipt)
// 			conf := PHConfirm{from: receipt.from, to: receipt.to, tx: receipt.tx}
// 			from := ph.paymentNetwork.getNode(receipt.from)
// 			from.PHconf(conf)
// 			to := ph.paymentNetwork.getNode(receipt.to)
// 			to.PHconf(conf)
// 		}
// 	}
// }