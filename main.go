package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
	"strconv"
	"sync"
)

var numNodes int
var numChannels int
var numTx int
var alpha float64

var numHub int
var participantsPerHub int

func main() {
	numNodes, _ = strconv.Atoi(os.Args[1])
	channelRatio, _ := strconv.Atoi(os.Args[2])
	numChannels = numNodes * channelRatio
	numTx = numNodes * 100
	alpha, _ = strconv.ParseFloat(os.Args[3], 64)
	participantsPerHub, _ = strconv.Atoi(os.Args[4])
	// numHub, _ = strconv.Atoi(os.Args[5])
	numHub = int((alpha * float64(numNodes)+float64(participantsPerHub-1))/float64(participantsPerHub))

	rand.Seed(time.Now().UTC().UnixNano())

	// create payment network
	pnPN := PaymentNetwork{vert: []int{}, edges: make(map[int]map[int]int), numEdges: 0}
	pnPH := PaymentNetwork{vert: []int{}, edges: make(map[int]map[int]int), numEdges: 0}
	pnCH := PaymentNetwork{vert: []int{}, edges: make(map[int]map[int]int), numEdges: 0}
	pnVC := PaymentNetwork{vert: []int{}, edges: make(map[int]map[int]int), numEdges: 0}

	// create nodes
	nodes := make([]*Node, numNodes)
	for i := 0; i < numNodes; i++ {
		nodes[i] = NewNode()
		pnPN.addVert(i)
		pnPH.addVert(i)
		pnCH.addVert(i)
		pnVC.addVert(i)
	}


	// create channels
	var channels []*Channel
	for i := 0; i < numChannels; i++ {
		from := rand.Intn(len(nodes))
		to := rand.Intn(len(nodes))
		if from == to {
			i--
			continue
		}
		c := NewChannel(from, to, 0)
		channels = append(channels, c)
		nodes[from].registerChannel(c)
		nodes[to].registerChannel(c)
		pnPN.addEdge(from, to, 1)
		pnPH.addEdge(from, to, 1)
		pnCH.addEdge(from, to, 1)
		pnVC.addEdge(from, to, 1)
	}


	// add channels so that any two nodes in payment network is connected
	unionset := makeUnionSet(channels)
	//fmt.Println(unionset)

	roots := makeRootSet(unionset)
	//fmt.Println(len(roots), roots)
	//return

	for i := 1; i < len(roots); i++ {
		from := roots[i-1]
		to := roots[i]
		c := NewChannel(from, to, 0)
		channels = append(channels, c)
		nodes[from].registerChannel(c)
		nodes[to].registerChannel(c)
		pnPN.addEdge(from, to, 1)
		pnPH.addEdge(from, to, 1)
		pnCH.addEdge(from, to, 1)
		pnVC.addEdge(from, to, 1)
	}

	//for i := 0; i < len(roots); i++ {
	//	for j := i+1; j < len(roots); j++ {
	//		from := roots[i]
	//		to := roots[j]
	//		channels = append(channels, &Channel{id: channelId,
	//											A: from,
	//											B: to,
	//											capacity: 0,
	//											lockChan: make(chan struct{}, 1),
	//											lockTx: nil})
	//		channelId++
	//		pnPN.addEdge(from, to, 1)
	//		pnPH.addEdge(from, to, 1)
	//		pnCH.addEdge(from, to, 1)
	//	}
	//}

	pnPN.makeNodeMap(nodes[:])
	pnPH.makeNodeMap(nodes[:])
	pnCH.makeNodeMap(nodes[:])

	pnPN.makeChannelMap(channels)
	pnPH.makeChannelMap(channels)
	pnCH.makeChannelMap(channels)


	var waitGroup sync.WaitGroup
	waitGroup.Add(4)

	// payment network
	go func() {
		defer waitGroup.Done()
		dist, next := FloydWarshall(pnPN)
		//fmt.Println("pair\tdist\tpath")
		for u, m := range dist {
			for v := range m {
				if u == v {
					continue
				}
				//fmt.Printf("%d -> %d\t%3d\t\t%s\n", u, v, d, pnPN.path(Path(u, v, next)))
				nodes[u].setPNPath(v, Path(u, v, next))
			}
		}
	}()

	// virtual channel
	go func() {
		defer waitGroup.Done()
		// create virtual channel
		for i := 0; i < int(alpha * float64(len(nodes))); i++ {
			from := rand.Intn(numNodes)
			to := rand.Intn(numNodes)
			if from == to {
				i--
				continue
			}
			c := NewChannel(from, to, 0)
			channels = append(channels, c)
			pnVC.addEdge(from, to, 1)
		}
		pnVC.makeNodeMap(nodes[:])
		pnVC.makeChannelMap(channels)

		dist, next := FloydWarshall(pnVC)
		for u, m := range dist {
			for v := range m {
				if u == v {
					continue
				}
				nodes[u].setVCPath(v, Path(u, v, next))
			}
		}
	}()


	// // payment hub
	// go func() {
	// 	defer waitGroup.Done()
	// 	ph := NewPaymentHub(&pnPH)
	// 	go ph.PaymentHubCron()
	// 	for i := 0; i < int(alpha * float64(len(nodes))); {
	// 		nodeId := rand.Intn(len(nodes))
	// 		if !ph.contain(nodeId) {
	// 			ph.join(nodes[nodeId])
	// 			i++
	// 		}
	// 	}

	// 	dist, next := FloydWarshall(pnPH)
	// 	for u, m := range dist {
	// 		for v := range m {
	// 			if u == v {
	// 				continue
	// 			}
	// 			nodes[u].setPHPath(v, Path(u, v, next))
	// 		}
	// 	}
	// } ()

	// payment hub 2
	go func() {
		defer waitGroup.Done()
		paymentHubs := make([]*PaymentHub, numHub)
		for i := range paymentHubs {
			paymentHubs[i] = NewPaymentHub(&pnPH)
		}
		paymentHubNodes := make(map[int]bool)
		for i := range paymentHubs {
			for j := 0; j < participantsPerHub; {
				nodeId := rand.Intn(len(nodes))
				if _, ok := paymentHubNodes[nodeId]; !ok {
					paymentHubNodes[nodeId] = true
					paymentHubs[i].join(nodes[nodeId])
					j++
				}
			}
		}

		dist, next := FloydWarshall(pnPH)
		for u, m := range dist {
			for v := range m {
				if u == v {
					continue
				}
				nodes[u].setPHPath(v, Path(u, v, next))
			}
		}
	} ()

	// // channel hub
	// go func() {
	// 	defer waitGroup.Done()
	// 	ch := NewChannelHub(&pnCH)
	// 	go ch.ChannelHubCron()
	// 	// 随机选择channel
	// 	//for i := 0; i < int(alpha * float32(len(nodes))); {
	// 	//	channel := channels[rand.Intn(len(channels))]
	// 	//	if !(ch.containNode(channel.A) || ch.containNode(channel.B)) {
	// 	//		ch.join(channel)
	// 	//		nodes[channel.A].joinChannelHub(&ch, channel)
	// 	//		nodes[channel.B].joinChannelHub(&ch, channel)
	// 	//		i++
	// 	//	}
	// 	//}

	// 	// 随机选择节点, 让每个节点自己选择一个channel
	// 	channelHubNodes := make(map[int]bool)
	// 	for i := 0; i < int(alpha * float64(len(nodes))); {
	// 		nodeId := rand.Intn(len(nodes))
	// 		if _, ok := channelHubNodes[nodeId]; !ok {
	// 			channel := nodes[nodeId].getRandomChannel()
	// 			if channel != nil && !ch.containChannel(channel.id) {
	// 				channelHubNodes[nodeId] = true
	// 				ch.join(channel)
	// 				nodes[channel.A].joinChannelHub(ch, channel)
	// 				nodes[channel.B].joinChannelHub(ch, channel)
	// 				i++
	// 			}
	// 		}
	// 	}

	// 	dist, next := FloydWarshall(pnCH)
	// 	for u, m := range dist {
	// 		for v := range m {
	// 			if u == v {
	// 				continue
	// 			}
	// 			nodes[u].setCHPath(v, Path(u, v, next))
	// 		}
	// 	}
	// }()

	// channel hub 2
	go func() {
		defer waitGroup.Done()
		channelHubs := make([]*ChannelHub, numHub)
		for i := range channelHubs {
			channelHubs[i] = NewChannelHub(&pnCH)
		}

		channelHubChannels := make(map[int]bool)
		for i := range channelHubs {
			for j := 0; j < participantsPerHub; {
				nodeId := rand.Intn(len(nodes))
				channel := nodes[nodeId].getRandomChannel()
				if _, ok := channelHubChannels[channel.id]; !ok {
					channelHubChannels[channel.id] = true
					channelHubs[i].join(channel)
					nodes[channel.A].joinChannelHub(channelHubs[i], channel)
					nodes[channel.B].joinChannelHub(channelHubs[i], channel)
					j++
				}
			}
		}

		dist, next := FloydWarshall(pnCH)
		for u, m := range dist {
			for v := range m {
				if u == v {
					continue
				}
				nodes[u].setCHPath(v, Path(u, v, next))
			}
		}
	}()


	waitGroup.Wait()

	// new tx
	txPN := make([]Transaction, numTx)
	txPH := make([]Transaction, numTx)
	txCH := make([]Transaction, numTx)
	txVC := make([]Transaction, numTx)
	totalLenPN, totalLenPH, totalLenCH, totalLenVC := 0, 0, 0, 0
	for i := 0; i < numTx; i++ {
		from := rand.Intn(numNodes)
		to := rand.Intn(numNodes)
		if from == to {
			i--
			continue
		}
		txPN[i] = Transaction{id: i, t: PNTX, from: from, to: to, secret: -1}
		txPH[i] = Transaction{id: i, t: PHTX, from: from, to: to, secret: -1}
		txCH[i] = Transaction{id: i, t: CHTX, from: from, to: to, secret: -1}
		txVC[i] = Transaction{id: i, t: VCTX, from: from, to: to, secret: -1}

		totalLenPN += len(nodes[from].getPathPN(to)) - 1
		totalLenPH += len(nodes[from].getPathPH(to)) - 1
		totalLenCH += len(nodes[from].getPathCH(to)) - 1
		totalLenVC += len(nodes[from].getPathVC(to)) - 1
	}
	fmt.Println(float64(totalLenPN) / float64(numTx), "\t",
				float64(totalLenVC) / float64(numTx), "\t",
				float64(totalLenPH) / float64(numTx), "\t",
				float64(totalLenCH) / float64(numTx), "\t")

	//done := make(chan bool, numTx)
	//
	//// nodeCron
	//for _, node := range nodes {
	//	node.setDoneChan(done)
	//	node.setPaymentNetwork(&pnCH)
	//	node.setPaymentHub(ph)
	//	node.setChannelHub(ch)
	//	go node.NodeCron()
	//}
	//
	//
	////exec tx in payment network
	//startPN := time.Now()
	//for _, tx := range txPN {
	//	nodes[tx.from].execTx(tx)
	//	<-done
	//}
	//elapsePN := time.Since(startPN)
	////fmt.Println("payment network:", elapsePN)
	//
	//
	//// exec tx in payment hub
	//startPH := time.Now()
	//for _, tx := range txPH {
	//	nodes[tx.from].execTx(tx)
	//	<-done
	//}
	//elapsePH := time.Since(startPH)
	////fmt.Println("payment hub:", elapsePH)
	//
	//
	//// exec tx in channel hub
	//startCH := time.Now()
	//for _, tx := range txCH {
	//	nodes[tx.from].execTx(tx)
	//	<-done
	//}
	//elapseCH := time.Since(startCH)
	////fmt.Println("channel hub:", elapseCH)
	//
	//
	//// exec tx in virtual channel
	//for _, node := range nodes {
	//	node.setPaymentNetwork(&pnVC)
	//}
	//startVC := time.Now()
	//for _, tx := range txVC {
	//	nodes[tx.from].execTx(tx)
	//	<-done
	//}
	//elapseVC := time.Since(startVC)
	////fmt.Println("payment network:", elapseVC)
	//
	//
	//fmt.Println(float64(totalLenPN) / numTx, "\t",
	//			float64(totalLenVC) / numTx, "\t",
	//			float64(totalLenPH) / numTx, "\t",
	//			float64(totalLenCH) / numTx, "\t",
	//			elapsePN / numTx, "\t",
	//			elapseVC / numTx, "\t",
	//			elapsePH / numTx, "\t",
	//			elapseCH / numTx)
}
