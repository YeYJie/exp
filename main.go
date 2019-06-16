package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
	"strconv"
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
	numHub = int((alpha * float64(numNodes)+float64(participantsPerHub-1))/float64(participantsPerHub))

	rand.Seed(time.Now().UTC().UnixNano())	

	pnPN := NewPaymentNetwork2(numNodes)
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
		pnPN.addEdge(from, to, 1)
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
		pnPN.addEdge(from, to, 1)
		pnPH.addEdge(from, to, 1)
		pnCH.addEdge(from, to, 1)
		pnVC.addEdge(from, to, 1)
	}

	// payment network
	pnPN.calculateShortestPath()

	// virtual channel
	for i := 0; i < int(alpha * float64(numNodes)); i++ {
		from := rand.Intn(numNodes)
		to := rand.Intn(numNodes)
		if from == to {
			i--
			continue
		}
		pnVC.addEdge(from, to, 1)
	}
	pnVC.calculateShortestPath()

	// payment hub
	paymentHubs := make([]*PaymentHub, numHub)
	for i := range paymentHubs {
		paymentHubs[i] = NewPaymentHub(pnPH)
	}
	paymentHubNodes := make(map[int]bool)
	for i := range paymentHubs {
		for j := 0; j < participantsPerHub; {
			nodeId := rand.Intn(numNodes)
			if _, ok := paymentHubNodes[nodeId]; !ok {
				paymentHubNodes[nodeId] = true
				paymentHubs[i].join(nodes[nodeId])
				j++
			}
		}
	}
	pnPH.calculateShortestPath()

	// channel hub
	channelHubs := make([]*ChannelHub, numHub)
	for i := range channelHubs {
		channelHubs[i] = NewChannelHub(pnCH)
	}
	channelHubChannels := make(map[int]bool)
	for i := range channelHubs {
		for j := 0; j < participantsPerHub; {
			nodeId := rand.Intn(numNodes)
			channel := nodes[nodeId].getRandomChannel()
			if _, ok := channelHubChannels[channel.id]; !ok {
				channelHubChannels[channel.id] = true
				channelHubs[i].join(channel)
				j++
			}
		}
	}
	pnCH.calculateShortestPath()

	totalLenPN, totalLenPH, totalLenCH, totalLenVC := 0, 0, 0, 0
	for i := 0; i < numTx; i++ {
		from := rand.Intn(numNodes)
		to := rand.Intn(numNodes)
		if from == to {
			i--
			continue
		}
		totalLenPN += pnPN.getDistance(from, to)
		totalLenPH += pnPH.getDistance(from, to)
		totalLenCH += pnCH.getDistance(from, to)
		totalLenVC += pnVC.getDistance(from, to)
	}
	fmt.Println(float64(totalLenPN) / float64(numTx), "\t",
				float64(totalLenVC) / float64(numTx), "\t",
				float64(totalLenPH) / float64(numTx), "\t",
				float64(totalLenCH) / float64(numTx), "\t")
}