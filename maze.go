package main

import (
	"fmt"
	"math/rand"
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
	grid := NewGrid(10, 10)
	grid.MazifyRec(0, 0)
	grid.Print()
}
