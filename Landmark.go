package main

import (
	"math/rand"
	// "math"
	"fmt"
)

type Landmark struct {
	coordinate [][]int
}

func NewLandmark(pn * PaymentNetwork2) *Landmark {
	return &Landmark{
		coordinate: make([][]int, pn.size),
	}
}

func (lm *Landmark) setRoutes(pn *PaymentNetwork2) {
	for i := range pn.channels {
		pn.channels[i] = sortAndUnique(pn.channels[i])
	}
	root := rand.Intn(pn.size)
	queue := make(chan []int, pn.size)
	queue <- pn.channels[root]
	lm.coordinate[root] = []int{}
	cords := make(chan []int, pn.size)
	cords <- lm.coordinate[root]
	queueSize := 1
	dis := 1
	for {
		n := queueSize
		queueSize = 0
		for j := 0; j < n; j++ {
			neighbor := <-queue
			cord := <- cords
			// fmt.Println("neighbor ", neighbor, cord, len(cord))
			for k := range neighbor {
				if lm.coordinate[neighbor[k]] != nil {
					continue
				}
				c := make([]int, len(cord))
				copy(c, cord)
				c = append(c, k)
				lm.coordinate[neighbor[k]] = c
				// fmt.Println(neighbor[k], c)
				queue <- pn.channels[neighbor[k]]
				cords <- c
				queueSize++
			}
		}
		dis++
		if queueSize == 0 {
			break;
		}
	}
}

func (lm *Landmark) dist(from, to int, pn *PaymentNetwork2) int {
	fromCord := lm.coordinate[from]
	toCord := lm.coordinate[to]
	minLen := min(len(fromCord), len(toCord))
	cpl := 0
	for ; cpl < minLen; cpl++ {
		if fromCord[cpl] != toCord[cpl] {
			break
		}
	}
	res := len(fromCord) + len(toCord) - 2 * cpl
	// fmt.Println(from, to, fromCord, toCord, res)
	return res
}

func (lm *Landmark) getDistance(from, to int, pn *PaymentNetwork2, print bool) int {
	next := from
	res := 0
	path := []int{}
	for {
		path = append(path, next)
		neighbors := pn.channels[next]
		minDis := pn.size
		for n := range neighbors {
			if next == neighbors[n] || contain(path, neighbors[n]) {
				continue;
			}
			dis := lm.dist(neighbors[n], to, pn)
			if dis < minDis {
				minDis = dis
				next = neighbors[n]
			}
		}
		res++
		// fmt.Println(from, to, next, res)
		if next == to || res > pn.size {
			break
		}
	}
	// fmt.Println(from, to, path)
	if print {
		fmt.Println(from, to, path, res)
	}
	return res
}