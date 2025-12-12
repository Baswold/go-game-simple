package engine

import "fmt"

// Position represents a cell on the board using zero-based coordinates.
type Position struct {
	Row int
	Col int
}

// Board stores a rectangular grid of player IDs.
type Board struct {
	Rows  int
	Cols  int
	cells []int
}

// NewBoard allocates a board with all cells empty (value 0).
func NewBoard(rows, cols int) *Board {
	if rows <= 0 || cols <= 0 {
		panic("board dimensions must be positive")
	}
	return &Board{
		Rows:  rows,
		Cols:  cols,
		cells: make([]int, rows*cols),
	}
}

// Clone produces a deep copy of the board.
func (b *Board) Clone() *Board {
	copyCells := make([]int, len(b.cells))
	copy(copyCells, b.cells)
	return &Board{
		Rows:  b.Rows,
		Cols:  b.Cols,
		cells: copyCells,
	}
}

// index converts a position into a linear index, returning an error for out of range coordinates.
func (b *Board) index(pos Position) (int, error) {
	if pos.Row < 0 || pos.Row >= b.Rows || pos.Col < 0 || pos.Col >= b.Cols {
		return 0, fmt.Errorf("position out of bounds: %+v", pos)
	}
	return pos.Row*b.Cols + pos.Col, nil
}

// Get returns the player occupying the given cell (0 when empty).
func (b *Board) Get(pos Position) (int, error) {
	idx, err := b.index(pos)
	if err != nil {
		return 0, err
	}
	return b.cells[idx], nil
}

// Set writes the player ID into the given cell when empty.
func (b *Board) Set(pos Position, playerID int) error {
	idx, err := b.index(pos)
	if err != nil {
		return err
	}
	if b.cells[idx] != 0 {
		return fmt.Errorf("cell already occupied at %+v", pos)
	}
	b.cells[idx] = playerID
	return nil
}

// SetAt writes the value regardless of whether the cell is occupied. Primarily used by games that need to clear captures.
func (b *Board) SetAt(pos Position, value int) error {
	idx, err := b.index(pos)
	if err != nil {
		return err
	}
	b.cells[idx] = value
	return nil
}

// IsFull reports whether all cells are occupied.
func (b *Board) IsFull() bool {
	for _, v := range b.cells {
		if v == 0 {
			return false
		}
	}
	return true
}

// ForEach iterates over every position and stored value on the board.
func (b *Board) ForEach(fn func(pos Position, value int)) {
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			idx := r*b.Cols + c
			fn(Position{Row: r, Col: c}, b.cells[idx])
		}
	}
}
