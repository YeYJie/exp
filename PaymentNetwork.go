package main

import "strconv"

type HTLC struct {
	from 		int
	to 			int
	tx 			*Transaction
}

type HTLCReply struct {
	from 		int
	to 			int
	tx 			*Transaction
}

type PaymentNetwork struct {
	vert  		[]int
	edges 		map[int]map[int]int
	numEdges 	int
	nodeMap 	map[int]*Node
	channelMap 	map[int]map[int]*Channel
}



func (pn *PaymentNetwork) getNode(nodeId int) *Node {
	return pn.nodeMap[nodeId]
}

func (pn *PaymentNetwork) makeNodeMap(nodes []*Node) {
	pn.nodeMap = make(map[int]*Node)
	for _, node := range nodes {
		pn.nodeMap[node.id] = node
	}
}

func (pn *PaymentNetwork) makeChannelMap(channels []*Channel) {
	pn.channelMap = make(map[int]map[int]*Channel)
	for _, channel := range channels {
		if pn.channelMap[channel.A] == nil {
			pn.channelMap[channel.A] = make(map[int]*Channel)
		}
		pn.channelMap[channel.A][channel.B] = channel
		if pn.channelMap[channel.B] == nil {
			pn.channelMap[channel.B] = make(map[int]*Channel)
		}
		pn.channelMap[channel.B][channel.A] = channel
	}
}

func (pn *PaymentNetwork) containChannel(A, B int) bool {
	_, ok := pn.channelMap[A][B]
	return ok
}

func (pn *PaymentNetwork) getChannel(A, B int) *Channel {
	return pn.channelMap[A][B]
}

func (pn *PaymentNetwork) edge(from, to int, weight int) {
	if _, ok := pn.edges[from]; !ok {
		pn.edges[from] = make(map[int]int)
	}
	pn.edges[from][to] = weight
}

func (pn *PaymentNetwork) addVert(id int) {
	//fmt.Println("addVert", id)
	pn.vert = append(pn.vert, int(id))
}

func (pn *PaymentNetwork) addEdge(from, to, weight int) {
	//fmt.Println("addEdge", from, to, weight)
	pn.edge(from, to, weight)
	pn.edge(to, from, weight)
	pn.numEdges++
}


func (pn PaymentNetwork) Vertices() []int {
	return pn.vert
}

func (pn PaymentNetwork) Neighbors(v int) (neighbors []int) {
	for k := range pn.edges[v] {
		neighbors = append(neighbors, k)
	}
	return neighbors
}

func (pn PaymentNetwork) Weight(u, v int) int {
	return pn.edges[u][v]
}

func (pn PaymentNetwork) path(verties []int) (res string) {
	if len(verties) == 0 {
		return ""
	}
	res = strconv.Itoa(int(verties[0]))
	for _, v := range verties[1:] {
		res += " -> " + strconv.Itoa(int(v))
	}
	return res
}