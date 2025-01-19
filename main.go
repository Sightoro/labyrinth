// Lonskiy Y.V. 2025
// code on github:

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	Name      []int
	Neighbors []*Node
	Way       []string
	FullCost  int
	Cost      int
}

type Item struct {
	node *Node
	cost int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].cost < pq[j].cost }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func writeError(out *bufio.Writer, errStr string) {
	fmt.Fprintln(os.Stderr, errStr) // Запись ошибки в stderr
	out.WriteString("-1\n.\n")
	out.Flush()
	os.Exit(0)
}

func parseLine(line string, count int, out *bufio.Writer) []string {
	fields := strings.Fields(line)
	if len(fields) != count {
		writeError(out, fmt.Sprintf("Expected %d fields, got %d: %s", count, len(fields), line))
	}
	return fields
}

func parseInt(value string, out *bufio.Writer) int {
	num, err := strconv.Atoi(value)
	if err != nil {
		writeError(out, fmt.Sprintf("Invalid integer: %s", value))
	}
	return num
}

func dijkstra(n, m int, start, end *Node, out *bufio.Writer) {
	start.FullCost = start.Cost
	start.Way = append(start.Way, fmt.Sprintf("%d %d", start.Name[0], start.Name[1]))
	pq := &PriorityQueue{}
	heap.Push(pq, Item{start, start.Cost})

	for pq.Len() > 0 {
		curr := heap.Pop(pq).(Item)
		node, cost := curr.node, curr.cost

		if cost > node.FullCost {
			continue
		}

		for _, edge := range node.Neighbors {
			newCost := node.FullCost + edge.Cost
			if newCost < edge.FullCost {
				edge.FullCost = newCost
				heap.Push(pq, Item{edge, newCost})
				edge.Way = append([]string(nil), node.Way...)
				edge.Way = append(edge.Way, fmt.Sprintf("%d %d", edge.Name[0], edge.Name[1]))
			}
		}
	}
	if end.FullCost == math.MaxInt32 {
		writeError(out, "No path found")
	} else {
		out.WriteString(strings.Join(end.Way, "\n") + "\n.")
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	line, err := in.ReadString('\n')
	if err != nil {
		writeError(out, "Failed to read dimensions")
	}
	fields := parseLine(strings.TrimSpace(line), 2, out)
	n := parseInt(fields[0], out)
	m := parseInt(fields[1], out)

	labyrinth := make([][]*Node, n+2)
	for i := 0; i < n+2; i++ {
		labyrinth[i] = make([]*Node, m+2)
	}

	for i := 1; i <= n; i++ {
		line, err := in.ReadString('\n')
		if err != nil {
			writeError(out, fmt.Sprintf("Failed to read row %d", i))
		}
		fields := parseLine(strings.TrimSpace(line), m, out)
		for j := 1; j <= m; j++ {
			cost := parseInt(fields[j-1], out)
			if cost != 0 {
				labyrinth[i][j] = &Node{
					Name:      []int{i - 1, j - 1},
					Neighbors: make([]*Node, 0),
					Cost:      cost,
					FullCost:  math.MaxInt32,
				}
			}
		}
	}

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			node := labyrinth[i][j]
			if node == nil {
				continue
			}
			if neighbor := labyrinth[i-1][j]; neighbor != nil {
				node.Neighbors = append(node.Neighbors, neighbor)
			}
			if neighbor := labyrinth[i][j-1]; neighbor != nil {
				node.Neighbors = append(node.Neighbors, neighbor)
			}
			if neighbor := labyrinth[i+1][j]; neighbor != nil {
				node.Neighbors = append(node.Neighbors, neighbor)
			}
			if neighbor := labyrinth[i][j+1]; neighbor != nil {
				node.Neighbors = append(node.Neighbors, neighbor)
			}
		}
	}

	line, err = in.ReadString('\n')
	if err != nil {
		writeError(out, "Failed to read start and end points")
	}
	fields = parseLine(strings.TrimSpace(line), 4, out)
	start_r := parseInt(fields[0], out)
	start_c := parseInt(fields[1], out)
	end_r := parseInt(fields[2], out)
	end_c := parseInt(fields[3], out)

	start_node := labyrinth[start_r+1][start_c+1]
	end_node := labyrinth[end_r+1][end_c+1]

	if start_node == nil || end_node == nil {
		writeError(out, "Invalid start or end point")
	}

	dijkstra(n, m, start_node, end_node, out)
}
