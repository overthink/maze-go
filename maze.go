package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// Direction flags are used to indicate which grid walls have openings.  e.g.
// if grid[r][c] == 4 then the South wall in cell (r,c) has been removed.
type Direction int

const (
	N = 1 << iota
	E
	S
	W
)

var opposite = map[Direction]Direction{N: S, E: W, S: N, W: E}

// Offsets deescribe what to add to the row/col to move the given direction in
// a grid.
var rowOffset = map[Direction]int{N: -1, E: 0, S: 1, W: 0}
var colOffset = map[Direction]int{N: 0, E: 1, S: 0, W: -1}

// Prefer NewGrid to create instances of this struct.
type Grid struct {
	RowCount int
	ColCount int
	data     [][]int
}

func NewGrid(rowCount, colCount int) Grid {
	data := make([][]int, rowCount)
	for i := range data {
		data[i] = make([]int, colCount)
	}
	return Grid{rowCount, colCount, data}
}

func (g *Grid) CellId(row, col int) int {
	return row*g.ColCount + col
}

// MazifyRec turns the grid into a maze using recursive backtracking.
func (g *Grid) MazifyRec(row, col int) {
	dirs := []Direction{N, E, S, W}
	rand.Shuffle(len(dirs), func(i, j int) { dirs[i], dirs[j] = dirs[j], dirs[i] })
	for _, d := range dirs {
		nextRow := row + rowOffset[d]
		nextCol := col + colOffset[d]
		// Carve through the wall in direction d if it's available and we
		// haven't already been there.
		if nextRow >= 0 && nextRow < g.RowCount &&
			nextCol >= 0 && nextCol < g.ColCount &&
			g.data[nextRow][nextCol] == 0 {
			g.data[row][col] |= int(d)
			g.data[nextRow][nextCol] |= int(opposite[d])
			g.MazifyRec(nextRow, nextCol)
		}
	}
}

// For Kruskal impl
type edge struct {
	row int
	col int
	d   Direction // other end of edge is in this direction
}

// MazifyKruskal turns grid into a maze using Kruskal's algorithm.
func (g *Grid) MazifyKruskal() {
	// 1. Generate all the possible edges in the grid graph.
	//   - our representation of an edge will be (row, col, direction)
	//     e.g. (3, 4, N) means an edge between cell (3, 4) and (2, 4), since
	//     (2, 4) is North of (3, 4)
	// 2. Shuffle the set of edges.
	// 3. Execute Kruskal's algorithm on the set of shuffled edges.
	//    - use a disjoint set union data structure
	//    - each edge starts in a disjoint subset all by itself
	//    - for each edge (u, v), if u and v are not in the same disjoint
	//      subset
	//      - update the grid allowing a path between u and v
	//      - union the representative sets for u and v

	dirs := []Direction{N, E, S, W}
	var edges []edge
	for row := 0; row < g.RowCount; row++ {
		for col := 0; col < g.ColCount; col++ {
			for _, d := range dirs {
				// If (row, col, d) is a valid edge, add it to our list.
				otherRow := row + rowOffset[d]
				otherCol := col + colOffset[d]
				if otherRow >= 0 && otherRow < g.RowCount &&
					otherCol >= 0 && otherCol < g.ColCount {
					edges = append(edges, edge{row, col, d})
				}
			}
		}
	}

	rand.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})

	// Parent pointers for DSU; initially each elements points to itself
	parent := make([]int, g.RowCount*g.ColCount)
	for i := range parent {
		parent[i] = i
	}

	var find func(int) int
	find = func(id int) int {
		if parent[id] == id {
			return id
		}
		// path compression
		parent[id] = find(parent[id])
		return parent[id]
	}

	union := func(idA, idB int) {
		setA := find(idA)
		setB := find(idB)
		if setA != setB {
			parent[setB] = setA
		}
	}

	for _, edge := range edges {
		otherRow := edge.row + rowOffset[edge.d]
		otherCol := edge.col + colOffset[edge.d]
		setA := find(g.CellId(edge.row, edge.col))
		setB := find(g.CellId(otherRow, otherCol))
		if setA != setB {
			g.data[edge.row][edge.col] |= int(edge.d)
			g.data[otherRow][otherCol] |= int(opposite[edge.d])
			union(setA, setB)
		}
	}
}

func (g *Grid) Print() {
	// print top border
	fmt.Printf(" ")
	fmt.Println(strings.Repeat("_", g.ColCount*2-1))
	for row := 0; row < g.RowCount; row++ {
		// print far left border
		fmt.Printf("|")
		for col := 0; col < g.ColCount; col++ {
			// print south wall if not open
			if g.data[row][col]&S != 0 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("_")
			}
			// handle east wall
			if g.data[row][col]&E != 0 {
				// Checking the east neighbour's southern opening is just done
				// to make the output prettier -- it's not for correctness.
				if (g.data[row][col]|g.data[row][col+1])&S != 0 {
					fmt.Printf(" ")
				} else {
					fmt.Printf("_")
				}
			} else {
				fmt.Printf("|")
			}
		}
		fmt.Println()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var rows int = 10
	var cols int = 10
	var err error
	if len(os.Args) > 1 {
		rows, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}
	if len(os.Args) > 2 {
		cols, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
	}

	grid := NewGrid(rows, cols)
	// grid.MazifyRec(0, 0)
	grid.MazifyKruskal()
	grid.Print()
}
