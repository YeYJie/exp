package main

//import "fmt"

type PCC struct {
	from 			int
	to 				int
	tx 				*Transaction
}

type GCC struct {
	from 			int
	to 				int
	tx 				*Transaction
}

type CHIOU struct {
	from 			int
	to 				int
	tx 				*Transaction
	gcc 			GCC
}

type CHReceipt struct {
	from 			int
	to				int
	tx 				*Transaction
}

type CHConf1 struct {
	from 			int
	to				int
	tx 				*Transaction
}

type CHRes struct {
	from 			int
	to				int
	tx 				*Transaction
	reply 			int
}

type CHConf2 struct {

}

type ChannelHub struct {
	paymentNetwork 	*PaymentNetwork2
	//channels 		[]*Channel

	iouChan 		chan CHIOU
	receiptChan 	chan CHReceipt

	nodeMap 		map[int]*Node
	channelMap 		map[int]*Channel
}

func NewChannelHub(pn *PaymentNetwork2) *ChannelHub {
	return &ChannelHub{paymentNetwork: pn,
					channelMap: make(map[int]*Channel),
					nodeMap: make(map[int]*Node),
					iouChan: make(chan CHIOU),
					receiptChan: make(chan CHReceipt)}
}

func (ch *ChannelHub) join(channel *Channel) {
	//fmt.Println("ChannelHub.join", channel)
	for _, c := range ch.channelMap {
		ch.paymentNetwork.addEdge(c.A, channel.A, 1)
		ch.paymentNetwork.addEdge(c.A, channel.B, 1)
		ch.paymentNetwork.addEdge(c.B, channel.A, 1)
		ch.paymentNetwork.addEdge(c.B, channel.B, 1)
	}
	ch.channelMap[channel.id] = channel
	// ch.nodeMap[channel.A] = ch.paymentNetwork.getNode(channel.A)
	// ch.nodeMap[channel.B] = ch.paymentNetwork.getNode(channel.B)
}

// func (ch *ChannelHub) containChannel(channelId int) bool {
// 	_, ok := ch.channelMap[channelId]
// 	return ok
// }

// func (ch *ChannelHub) containNode(nodeId int) bool {
// 	_, ok := ch.nodeMap[nodeId]
// 	return ok
// }

// func (ch *ChannelHub) iou(iou CHIOU) {
// 	communicationDelay()
// 	ch.iouChan <- iou
// }

// func (ch *ChannelHub) receipt(receipt CHReceipt) {
// 	communicationDelay()
// 	ch.receiptChan <- receipt
// }

// func (ch *ChannelHub) ChannelHubCron() {
// 	for {
// 		select {
// 		case iou := <-ch.iouChan:
// 			//fmt.Println("ChannelHub.ChannelHubCron", "iou:", iou)
// 			to := ch.paymentNetwork.getNode(iou.to)
// 			to.CHiou(iou)
// 		case receipt := <-ch.receiptChan:
// 			//fmt.Println("ChannelHub.ChannelHubCron", "receipt:", receipt)
// 			conf := CHConf1{from: receipt.from, to: receipt.to, tx: receipt.tx}
// 			from := ch.paymentNetwork.getNode(receipt.from)
// 			from.CHConf1(conf)
// 			to := ch.paymentNetwork.getNode(receipt.to)
// 			to.CHConf1(conf)
// 		}
// 	}
// }