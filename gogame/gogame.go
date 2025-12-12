package gogame

import (
	"fmt"
	"strings"

	"boardgame/engine"
)

// Color represents stone colors.
type Color int

const (
	None  Color = 0
	Black Color = 1
	White Color = 2
)

func (c Color) String() string {
	switch c {
	case Black:
		return "Black"
	case White:
		return "White"
	default:
		return "None"
	}
}

// MoveResult describes the outcome of applying a move.
type MoveResult struct {
	Captured int
}

// Game holds Go state for one board.
type Game struct {
	Board             *engine.Board
	Size              int
	ToPlay            Color
	Captures          map[Color]int
	ConsecutivePasses int
	moveNumber        int
	history           map[string]struct{}
	lastHash          string
}

// NewGame initializes an empty board with Black to play.
func NewGame(size int) (*Game, error) {
	if size < 5 {
		return nil, fmt.Errorf("board size must be at least 5")
	}
	board := engine.NewBoard(size, size)
	g := &Game{
		Board:    board,
		Size:     size,
		ToPlay:   Black,
		Captures: map[Color]int{Black: 0, White: 0},
		history:  map[string]struct{}{},
	}
	hash := serialize(board, g.ToPlay)
	g.lastHash = hash
	g.history[hash] = struct{}{}
	return g, nil
}

// PlayMove places a stone for the current player, enforcing capture, suicide, and simple superko.
func (g *Game) PlayMove(pos engine.Position) (MoveResult, error) {
	if g.ToPlay == None {
		return MoveResult{}, fmt.Errorf("game is finished")
	}
	mover := g.ToPlay
	if pos.Row < 0 || pos.Row >= g.Size || pos.Col < 0 || pos.Col >= g.Size {
		return MoveResult{}, fmt.Errorf("position out of bounds")
	}

	working := g.Board.Clone()
	if err := working.Set(pos, int(g.ToPlay)); err != nil {
		return MoveResult{}, err
	}

	opponent := other(g.ToPlay)
	totalCaptured := 0

	for _, n := range neighbors(g.Size, pos) {
		val, _ := working.Get(n)
		if Color(val) != opponent {
			continue
		}
		group, libs, err := collectGroup(working, n)
		if err != nil {
			return MoveResult{}, err
		}
		if len(libs) == 0 {
			for _, p := range group {
				if err := working.SetAt(p, 0); err != nil {
					return MoveResult{}, err
				}
			}
			totalCaptured += len(group)
		}
	}

	// Check liberties of the newly placed stone (after any captures).
	_, libs, err := collectGroup(working, pos)
	if err != nil {
		return MoveResult{}, err
	}
	if len(libs) == 0 {
		return MoveResult{}, fmt.Errorf("suicide is not allowed")
	}

	nextToPlay := opponent
	newHash := serialize(working, nextToPlay)
	if _, exists := g.history[newHash]; exists {
		return MoveResult{}, fmt.Errorf("move violates superko (repeats a previous position)")
	}

	g.Board = working
	g.ToPlay = nextToPlay
	g.lastHash = newHash
	g.history[newHash] = struct{}{}
	g.moveNumber++
	g.ConsecutivePasses = 0
	if totalCaptured > 0 {
		g.Captures[mover] += totalCaptured
	}

	return MoveResult{Captured: totalCaptured}, nil
}

// Pass ends the current turn without placing a stone.
func (g *Game) Pass() {
	if g.ToPlay == None {
		return
	}
	g.moveNumber++
	g.ConsecutivePasses++
	g.ToPlay = other(g.ToPlay)
	g.lastHash = serialize(g.Board, g.ToPlay)
	g.history[g.lastHash] = struct{}{}
}

// MoveNumber returns the number of moves played.
func (g *Game) MoveNumber() int {
	return g.moveNumber
}

func other(c Color) Color {
	if c == Black {
		return White
	}
	if c == White {
		return Black
	}
	return None
}

func neighbors(size int, pos engine.Position) []engine.Position {
	dirs := []engine.Position{
		{Row: -1, Col: 0},
		{Row: 1, Col: 0},
		{Row: 0, Col: -1},
		{Row: 0, Col: 1},
	}
	out := make([]engine.Position, 0, 4)
	for _, d := range dirs {
		r := pos.Row + d.Row
		c := pos.Col + d.Col
		if r >= 0 && r < size && c >= 0 && c < size {
			out = append(out, engine.Position{Row: r, Col: c})
		}
	}
	return out
}

func collectGroup(b *engine.Board, start engine.Position) ([]engine.Position, map[engine.Position]struct{}, error) {
	color, err := b.Get(start)
	if err != nil {
		return nil, nil, err
	}
	if color == 0 {
		return nil, nil, fmt.Errorf("no stone at %+v", start)
	}

	group := []engine.Position{}
	liberties := map[engine.Position]struct{}{}
	seen := map[engine.Position]struct{}{}
	stack := []engine.Position{start}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if _, ok := seen[current]; ok {
			continue
		}
		seen[current] = struct{}{}
		group = append(group, current)

		for _, n := range neighbors(b.Rows, current) {
			val, _ := b.Get(n)
			switch val {
			case 0:
				liberties[n] = struct{}{}
			case color:
				if _, ok := seen[n]; !ok {
					stack = append(stack, n)
				}
			}
		}
	}
	return group, liberties, nil
}

func serialize(b *engine.Board, toPlay Color) string {
	var sb strings.Builder
	sb.Grow(b.Rows*b.Cols + 1)
	sb.WriteByte(byte(toPlay) + '0')
	b.ForEach(func(_ engine.Position, v int) {
		sb.WriteByte(byte(v) + '0')
	})
	return sb.String()
}
