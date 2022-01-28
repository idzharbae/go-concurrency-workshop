package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

var visited = new(sync.Map)

func bfs(adjList map[int][]int, root int) (traversedNodes int) {
	q := list.New()

	if _, ok := visited.Load(root); ok {
		return
	}

	visited.Store(root, true)
	q.PushBack(root)
	traversedNodes++

	for q.Len() > 0 {
		e := q.Front()
		front := e.Value.(int)
		q.Remove(e)

		// fmt.Println(root, front)

		for _, child := range adjList[front] {
			if _, ok := visited.Load(child); ok {
				continue
			}

			// if root != 0 {
			// 	fmt.Println("asdf")
			// }

			visited.Store(child, true)
			q.PushBack(child)
			traversedNodes++
		}
	}

	return
}

func bfsConcurrent(adjList map[int][]int, n int) {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		traversedNodes := bfs(adjList, 0)
		fmt.Println("Traversed nodes 1:", traversedNodes)
	}()

	go func() {
		defer wg.Done()
		traversedNodes := bfs(adjList, (n*n)-1)
		fmt.Println("Traversed nodes 2:", traversedNodes)
	}()

	wg.Wait()
}

func generateMazeAdjList(n int) map[int][]int {
	result := make(map[int][]int)

	// node number = i + n*j
	// total nodes = n*n
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			nodeNumber := i + (n * j)

			if (nodeNumber+1)%n != 0 {
				result[nodeNumber] = append(result[nodeNumber], nodeNumber+1)
			}

			if nodeNumber%n != 0 {
				result[nodeNumber] = append(result[nodeNumber], nodeNumber-1)
			}

			if nodeNumber+n < n*n {
				result[nodeNumber] = append(result[nodeNumber], nodeNumber+n)
			}

			if nodeNumber-n > 0 {
				result[nodeNumber] = append(result[nodeNumber], nodeNumber-n)
			}
		}
	}

	return result
}

func main() {
	maze := generateMazeAdjList(800)

	now := time.Now()
	// bfs(maze, 0)
	bfsConcurrent(maze, 800)
	fmt.Println(time.Since(now))
}
