package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

var numNodes int
var numChannels int
var numTx int
var alpha float64

var numHub int
var participantsPerHub int

func main() {
	numNodes, _ = strconv.Atoi(os.Args[1])
	channelRatio := 4
	numChannels = numNodes * channelRatio
	numTx = numNodes * 100
	alpha, _ = strconv.ParseFloat(os.Args[2], 64)
	numJoin := int(alpha * float64(numNodes))
	participantsPerHub, _ = strconv.Atoi(os.Args[3])
	numHub = int(float64(numJoin+participantsPerHub-1)/float64(participantsPerHub))
	numLm := numNodes / 100

	rand.Seed(time.Now().UTC().UnixNano())	

	// pnPN := NewPaymentNetwork2(numNodes)
	pnPH := NewPaymentNetwork2(numNodes)
	pnCH := NewPaymentNetwork2(numNodes)
	pnVC := NewPaymentNetwork2(numNodes)

	nodes := make([]*Node, numNodes)
	for i := 0; i < numNodes; i++ {
		nodes[i] = NewNode()
	}

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
		// pnPN.addEdge(from, to, 1)
		pnPH.addEdge(from, to, 1)
		pnCH.addEdge(from, to, 1)
		pnVC.addEdge(from, to, 1)
	}

	unionset := makeUnionSet(channels)
	roots := makeRootSet(unionset)

	for i := 1; i < len(roots); i++ {
		from := roots[i-1]
		to := roots[i]
		c := NewChannel(from, to, 0)
		channels = append(channels, c)
		nodes[from].registerChannel(c)
        nodes[to].registerChannel(c)
		// pnPN.addEdge(from, to, 1)
		pnPH.addEdge(from, to, 1)
		pnCH.addEdge(from, to, 1)
		pnVC.addEdge(from, to, 1)
	}

	var wg  sync.WaitGroup
	wg.Add(3)

	// // payment network
	// pnPN.calculateShortestPath()
	// lmsPN := make([]*Landmark, numLm)
	// go func() {
	// 	defer wg.Done()
	// 	for l := range lmsPN {
	// 		lmsPN[l] = NewLandmark(pnPN)
	// 		lmsPN[l].setRoutes(pnPN)
	// 	}
	// }()

	// virtual channel
	for i := 0; i < numJoin; i++ {
		from := rand.Intn(numNodes)
		to := rand.Intn(numNodes)
		if from == to {
			i--
			continue
		}
		pnVC.addEdge(from, to, 1)
	}
	pnVC.calculateShortestPath()
	lmsVC := make([]*Landmark, numLm)
	go func() {
		defer wg.Done()
		for l := range lmsVC {
			lmsVC[l] = NewLandmark(pnVC)
			lmsVC[l].setRoutes(pnVC)
		}
	}()

	// payment hub
	paymentHubs := make([]*PaymentHub, numHub)
	for i := range paymentHubs {
		paymentHubs[i] = NewPaymentHub(pnPH)
	}
	paymentHubNodes := make(map[int]bool)
	joinedPH := 0
	for i := range paymentHubs {
		for j := 0; j < participantsPerHub && joinedPH <= numJoin; {
			nodeId := rand.Intn(numNodes)
			if _, ok := paymentHubNodes[nodeId]; !ok {
				paymentHubNodes[nodeId] = true
				paymentHubs[i].join(nodes[nodeId])
				j++
				joinedPH++
			}
		}
	}
	pnPH.calculateShortestPath()
	lmsPH := make([]*Landmark, numLm)
	go func() {
		defer wg.Done()
		for l := range lmsPH {
			lmsPH[l] = NewLandmark(pnPH)
			lmsPH[l].setRoutes(pnPH)
		}
	}()

	// channel hub
	channelHubs := make([]*ChannelHub, numHub)
	for i := range channelHubs {
		channelHubs[i] = NewChannelHub(pnCH)
	}
	channelHubChannels := make(map[int]bool)
	joinedCH := 0
	for i := range channelHubs {
		for j := 0; j < participantsPerHub && joinedCH <= numJoin; {
			nodeId := rand.Intn(numNodes)
			channel := nodes[nodeId].getRandomChannel()
			if _, ok := channelHubChannels[channel.id]; !ok {
				channelHubChannels[channel.id] = true
				channelHubs[i].join(channel)
				j++
				joinedCH++
			}
		}
	}
	pnCH.calculateShortestPath()
	lmsCH := make([]*Landmark, numLm)
	go func() {
		defer wg.Done()
		for l := range lmsCH {
			lmsCH[l] = NewLandmark(pnCH)
			lmsCH[l].setRoutes(pnCH)
		}
	}()

	wg.Wait()

	totalLenPNFW, totalLenPHFW, totalLenCHFW, totalLenVCFW := int64(0), int64(0), int64(0), int64(0)
	totalLenPNSM, totalLenPHSM, totalLenCHSM, totalLenVCSM := int64(0), int64(0), int64(0), int64(0)
	_ = totalLenPNFW
	_ = totalLenPNSM
	txChan := make(chan Tx, numTx)
	for i := 0; i < numTx; i++ {
		from := rand.Intn(numNodes)
		to := rand.Intn(numNodes)
		if from == to {
			i--
			continue
		}
		txChan <- Tx{from: from, to: to}
		// // totalLenPNFW += pnPN.getDistance(from, to)
		// totalLenVCFW += pnVC.getDistance(from, to)
		// totalLenPHFW += pnPH.getDistance(from, to)
		// totalLenCHFW += pnCH.getDistance(from, to)
		// for l := 0; l < numLm; l++ {
		// 	// totalLenPNSM += lmsPN[l].getDistance(from, to, pnPN, false)
		// 	totalLenVCSM += lmsVC[l].getDistance(from, to, pnVC, false)
		// 	totalLenPHSM += lmsPH[l].getDistance(from, to, pnPH, false)
		// 	totalLenCHSM += lmsCH[l].getDistance(from, to, pnCH, false)
		// }
	}
	wg.Add(12)
	for g := 0; g < 12; g++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case tx := <- txChan:			
					atomic.AddInt64(&totalLenVCFW, int64(pnVC.getDistance(tx.from, tx.to)))
					atomic.AddInt64(&totalLenPHFW, int64(pnPH.getDistance(tx.from, tx.to)))
					atomic.AddInt64(&totalLenCHFW, int64(pnCH.getDistance(tx.from, tx.to)))
					for l := 0; l < numLm; l++ {
						atomic.AddInt64(&totalLenVCSM, int64(lmsVC[l].getDistance(tx.from, tx.to, pnVC, false)))
						atomic.AddInt64(&totalLenPHSM, int64(lmsPH[l].getDistance(tx.from, tx.to, pnPH, false)))
						atomic.AddInt64(&totalLenCHSM, int64(lmsCH[l].getDistance(tx.from, tx.to, pnCH, false)))
					}
				default:
					return
				}
			}
		}()
	}
	wg.Wait()
	fmt.Printf("%d %.2f %d :  %.2f  %.2f  %.2f  %.1f%%  %.1f%%  %.2f  %.2f  %.2f  %.1f%%  %.1f%%\n", 
				numNodes, alpha, participantsPerHub,
				float64(atomic.LoadInt64(&totalLenVCFW)) / float64(numTx),
				float64(atomic.LoadInt64(&totalLenPHFW)) / float64(numTx),
				float64(atomic.LoadInt64(&totalLenCHFW)) / float64(numTx),
				100.0 * float64(atomic.LoadInt64(&totalLenPHFW) - atomic.LoadInt64(&totalLenCHFW)) / float64(atomic.LoadInt64(&totalLenPHFW)),
				100.0 * float64(atomic.LoadInt64(&totalLenVCFW) - atomic.LoadInt64(&totalLenCHFW)) / float64(atomic.LoadInt64(&totalLenVCFW)),
				float64(atomic.LoadInt64(&totalLenVCSM)) / float64(numTx) / float64(numLm),
				float64(atomic.LoadInt64(&totalLenPHSM)) / float64(numTx) / float64(numLm),
				float64(atomic.LoadInt64(&totalLenCHSM)) / float64(numTx) / float64(numLm),
				100.0 * float64(atomic.LoadInt64(&totalLenPHSM) - atomic.LoadInt64(&totalLenCHSM)) / float64(atomic.LoadInt64(&totalLenPHSM)),
				100.0 * float64(atomic.LoadInt64(&totalLenVCSM) - atomic.LoadInt64(&totalLenCHSM)) / float64(atomic.LoadInt64(&totalLenVCSM)))
}