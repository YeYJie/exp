package main

import (
	//"fmt"
	//"strconv"
)

type Graph interface {
	Vertices() []int
	Neighbors(v int) []int
	Weight(u, v int) int
}

const Infinity = int(^uint(0) >> 1)

func FloydWarshall(g Graph) (dist map[int]map[int]int, next map[int]map[int]*int) {
	vert := g.Vertices()
	dist = make(map[int]map[int]int)
	next = make(map[int]map[int]*int)
	for _, u := range vert {
		dist[u] = make(map[int]int)
		next[u] = make(map[int]*int)
		for _, v := range vert {
			dist[u][v] = Infinity
		}
		dist[u][u] = 0
		for _, v := range g.Neighbors(u) {
			v := v
			dist[u][v] = g.Weight(u, v)
			next[u][v] = &v
		}
	}
	for _, k := range vert {
		for _, i := range vert {
			for _, j := range vert {
				if dist[i][k] < Infinity && dist[k][j] < Infinity {
					if dist[i][j] > dist[i][k]+dist[k][j] {
						dist[i][j] = dist[i][k] + dist[k][j]
						next[i][j] = next[i][k]
					}
				}
			}
		}
	}
	return dist, next
}

func Path(u, v int, next map[int]map[int]*int) (path []int) {
	if next[u][v] == nil {
		return
	}
	path = []int{u}
	for u != v {
		u = *next[u][v]
		path = append(path, u)
	}
	return path
}

//func main() {
//	g := ig{[]int{1, 2, 3, 4}, make(map[int]map[int]int)}
//	g.edge(1, 3, 2)
//	g.edge(3, 4, 2)
//	g.edge(4, 2, 1)
//	g.edge(2, 1, 4)
//	g.edge(2, 3, 3)
//
//	dist, next := FloydWarshall(g)
//	fmt.Println("pair\tdist\tpath")
//	for u, m := range dist {
//		for v, d := range m {
//			if u != v {
//				fmt.Printf("%d -> %d\t%3d\t%s\n", u, v, d, g.path(Path(u, v, next)))
//			}
//		}
//	}
//}
