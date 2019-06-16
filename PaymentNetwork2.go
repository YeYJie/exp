package main

import (
	"sync"
	"sort"
)

const INF_DISTANCE = 1000000

type PaymentNetwork2 struct {
	size int
	channels [][]int
	distance [][]int
}

func NewPaymentNetwork2(size int) *PaymentNetwork2 {
	paymentNetwork2 := PaymentNetwork2{
		size: size,
		channels: make([][]int, size),
		distance: make([][]int, size)}
	for i := range paymentNetwork2.distance {
		paymentNetwork2.distance[i] = make([]int, size)
		for j := range paymentNetwork2.distance[i] {
			paymentNetwork2.distance[i][j] = INF_DISTANCE
		}
	}
	return &paymentNetwork2
}

func (pn *PaymentNetwork2) addEdge(from, to, weight int) {
	pn.channels[from] = append(pn.channels[from], to)
	pn.channels[to] = append(pn.channels[to], from)
}

func sortAndUnique(intSlice []int) []int {
    keys := make(map[int]bool)
    list := []int{} 
    for _, entry := range intSlice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }    
	sort.Ints(list)
	return list
}

func (pn *PaymentNetwork2) calculateShortestPath() {
	for i := range pn.channels {
		pn.channels[i] = sortAndUnique(pn.channels[i])
	}
	iChan := make(chan int, pn.size)
	for i := 0; i < pn.size; i++ {
		iChan <- i
	}
	var wg sync.WaitGroup
	numG := 32
	wg.Add(numG)
	for g := 0; g < numG; g++ {
		go func() {
			defer wg.Done()
			queue := make(chan []int, pn.size)
			for {
				select {
				case i := <-iChan:
					queue <- pn.channels[i]
					queueSize := 1
					dis := 1
					for {
						n := queueSize
						queueSize = 0
						for j := 0; j < n; j++ {
							neighbor := <-queue
							for k := range neighbor {
								if i != neighbor[k] && pn.distance[i][neighbor[k]] > dis {
									pn.distance[i][neighbor[k]] = dis
									queue <- pn.channels[neighbor[k]]
									queueSize++
								}
							}
						}
						dis++
						if queueSize == 0 {
							break;
						}
					}
				default:
					return
				}
			}
		}()
	}
	wg.Wait()
	pn.channels = nil
}

func (pn *PaymentNetwork2) getDistance(from, to int) int {
	if pn.distance[from][to] >= INF_DISTANCE {
		return 0
	}
	return pn.distance[from][to]
}