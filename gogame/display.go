package gogame

import (
	"fmt"
	"strconv"
	"strings"

	"boardgame/engine"
)

// ColumnLabels returns board column labels (skipping I like traditional Go coordinates).
func ColumnLabels(size int) []string {
	labels := make([]string, 0, size)
	for ch := 'A'; len(labels) < size; ch++ {
		if ch == 'I' { // skip I
			continue
		}
		labels = append(labels, string(ch))
	}
	return labels
}

// ParseCoord converts user text (e.g., "D4") into a board position.
// Rows count from the bottom (1) upwards to the board size.
func ParseCoord(input string, size int) (engine.Position, error) {
	s := strings.TrimSpace(strings.ToUpper(input))
	if len(s) < 2 {
		return engine.Position{}, fmt.Errorf("coordinate must include column and row (e.g., D4)")
	}

	labels := ColumnLabels(size)
	colLabel := s[:1]
	col := -1
	for i, l := range labels {
		if l == colLabel {
			col = i
			break
		}
	}
	if col == -1 {
		return engine.Position{}, fmt.Errorf("invalid column %q", colLabel)
	}

	rowNum, err := strconv.Atoi(s[1:])
	if err != nil {
		return engine.Position{}, fmt.Errorf("invalid row number")
	}
	if rowNum < 1 || rowNum > size {
		return engine.Position{}, fmt.Errorf("row must be between 1 and %d", size)
	}

	// Row index 0 is the top, so invert from bottom-based numbering.
	row := size - rowNum
	return engine.Position{Row: row, Col: col}, nil
}

// RenderBoardASCII prints the board with X (Black), O (White), and . (empty).
func RenderBoardASCII(g *Game) string {
	labels := ColumnLabels(g.Size)
	var sb strings.Builder
	sb.Grow(g.Size * g.Size * 2)

	writeLabels := func() {
		sb.WriteString("   ")
		for _, l := range labels {
			sb.WriteString(l)
			sb.WriteByte(' ')
		}
	}

	writeLabels()
	sb.WriteByte('\n')

	for row := 0; row < g.Size; row++ {
		displayRow := g.Size - row
		sb.WriteString(fmt.Sprintf("%2d ", displayRow))
		for col := 0; col < g.Size; col++ {
			val, _ := g.Board.Get(engine.Position{Row: row, Col: col})
			ch := "."
			if val == int(Black) {
				ch = "X"
			} else if val == int(White) {
				ch = "O"
			}
			sb.WriteString(ch)
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%2d\n", displayRow))
	}

	writeLabels()
	return sb.String()
}
