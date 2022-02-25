package main

import (
	"container/list"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func bfs(adjList map[int][]int, root int, visited *sync.Map) (traversedNodes int) {
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

		for _, child := range adjList[front] {
			if _, ok := visited.Load(child); ok {
				continue
			}

			visited.Store(child, true)
			q.PushBack(child)
			traversedNodes++
		}
	}

	return
}

func bfsConcurrent(adjList map[int][]int, n int, visited *sync.Map) {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		bfs(adjList, 0, visited)
		// fmt.Println("Traversed nodes 1:", traversedNodes)
	}()

	go func() {
		defer wg.Done()
		bfs(adjList, (n*n)-1, visited)
		// fmt.Println("Traversed nodes 2:", traversedNodes)
	}()

	wg.Wait()
}

func generateMazeAdjList(n int) map[int][]int {
	result := make(map[int][]int)

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

	var wg sync.WaitGroup
	wg.Add(12)

	now := time.Now()

	for i := 0; i < 12; i++ {
		go func(maze map[int][]int) {
			defer wg.Done()
			bfsConcurrent(maze, 800, new(sync.Map))
		}(maze)
	}

	wg.Wait()

	fmt.Println(time.Since(now))

	fmt.Println(runtime.NumCPU())
}
