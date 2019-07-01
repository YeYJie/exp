package main

import (
	"time"
	"math/rand"
	"sort"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func hash(s int) int {
	return s+100
}

func communicationDelay() {
	ms := 50 + rand.Intn(20) - 10
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func makeUnionSet(channels []*Channel) (unionset []int) {
	for i := 0; i < numNodes; i++ {
		unionset = append(unionset, i)
	}
	for _, channel := range channels {
		iA := channel.A
		rootA := unionset[iA]
		for iA != rootA {
			iA = rootA
			rootA = unionset[rootA]
		}
		iB := channel.B
		rootB := unionset[iB]
		for iB != rootB {
			iB = rootB
			rootB = unionset[rootB]
		}
		unionset[rootA] = rootB
	}
	return unionset
}

func makeRootSet(unionset []int) []int {
	roots := make([]int, numNodes)
	for i := 0; i < numNodes; i++ {
		index := i
		root := unionset[index]
		for index != root {
			index = root
			root = unionset[root]
		}
		roots[i] = root
	}
	//fmt.Println(roots)

	set := make(map[int]int)
	for i := 0; i < numNodes; i++ {
		if _, ok := set[roots[i]]; !ok {
			set[roots[i]] = 1
		}
	}
	var res []int
	for k := range set {
		res = append(res, k)
	}
	return res
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

func contain(arr []int, ele int) bool {
	res := false
	for _, p := range arr {
		if ele == p {
			res = true
			break
		}
	}
	return res
}