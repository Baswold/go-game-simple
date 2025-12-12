package tictactoe

import (
	"fmt"
	"strings"

	"boardgame/engine"
)

// Rules implements Tic-Tac-Toe on a square board (default 3x3).
type Rules struct {
	Size int
}

// NewRules returns standard 3x3 Tic-Tac-Toe rules.
func NewRules() Rules {
	return Rules{Size: 3}
}

// NewGame constructs a game state for Tic-Tac-Toe.
func (r Rules) NewGame(players []engine.Player) (*engine.Game, error) {
	if len(players) != 2 {
		return nil, fmt.Errorf("tic-tac-toe requires 2 players, got %d", len(players))
	}
	return engine.NewGame(engine.NewBoard(r.Size, r.Size), players)
}

// ValidMoves returns empty cells for the current player.
func (r Rules) ValidMoves(g *engine.Game) []engine.Move {
	if g.Outcome.Winner != nil || g.Outcome.Draw {
		return nil
	}
	current := g.CurrentPlayer().ID
	moves := make([]engine.Move, 0, g.Board.Rows*g.Board.Cols)
	g.Board.ForEach(func(pos engine.Position, value int) {
		if value == 0 {
			moves = append(moves, engine.Move{
				PlayerID: current,
				Pos:      pos,
			})
		}
	})
	return moves
}

// ApplyMove writes a player's mark if the move is valid.
func (r Rules) ApplyMove(g *engine.Game, m engine.Move) error {
	if m.PlayerID != g.CurrentPlayer().ID {
		return fmt.Errorf("it is not player %d's turn", m.PlayerID)
	}
	if err := g.RecordMove(m); err != nil {
		return err
	}
	// If the move completes the game, store the outcome now so Status can surface it.
	if winnerID := r.findWinner(g.Board); winnerID != 0 {
		g.EndGame(engine.Outcome{Winner: lookupPlayer(g.Players, winnerID)})
	} else if g.Board.IsFull() {
		g.EndGame(engine.Outcome{Draw: true})
	}
	return nil
}

// Status reports whether the current position is terminal.
func (r Rules) Status(g *engine.Game) (engine.Outcome, bool) {
	if g.Outcome.Winner != nil || g.Outcome.Draw {
		return g.Outcome, true
	}
	if winnerID := r.findWinner(g.Board); winnerID != 0 {
		return engine.Outcome{Winner: lookupPlayer(g.Players, winnerID)}, true
	}
	if g.Board.IsFull() {
		return engine.Outcome{Draw: true}, true
	}
	return engine.Outcome{}, false
}

// findWinner returns the player ID that has aligned a full row, column, or diagonal.
func (r Rules) findWinner(b *engine.Board) int {
	size := r.Size
	lines := make([][]engine.Position, 0, size*2+2)

	// rows
	for row := 0; row < size; row++ {
		line := make([]engine.Position, 0, size)
		for col := 0; col < size; col++ {
			line = append(line, engine.Position{Row: row, Col: col})
		}
		lines = append(lines, line)
	}
	// columns
	for col := 0; col < size; col++ {
		line := make([]engine.Position, 0, size)
		for row := 0; row < size; row++ {
			line = append(line, engine.Position{Row: row, Col: col})
		}
		lines = append(lines, line)
	}
	// diagonals
	diag := make([]engine.Position, 0, size)
	for i := 0; i < size; i++ {
		diag = append(diag, engine.Position{Row: i, Col: i})
	}
	lines = append(lines, diag)

	diag = make([]engine.Position, 0, size)
	for i := 0; i < size; i++ {
		diag = append(diag, engine.Position{Row: i, Col: size - i - 1})
	}
	lines = append(lines, diag)

	for _, line := range lines {
		first, _ := b.Get(line[0])
		if first == 0 {
			continue
		}
		allMatch := true
		for _, pos := range line[1:] {
			value, _ := b.Get(pos)
			if value != first {
				allMatch = false
				break
			}
		}
		if allMatch {
			return first
		}
	}
	return 0
}

// RenderBoard returns an ASCII board using player tokens for readability.
func RenderBoard(g *engine.Game) string {
	token := func(id int) string {
		if id == 0 {
			return " "
		}
		p := lookupPlayer(g.Players, id)
		if p == nil || p.Token == "" {
			return fmt.Sprintf("%d", id)
		}
		return p.Token
	}

	var sb strings.Builder
	for r := 0; r < g.Board.Rows; r++ {
		for c := 0; c < g.Board.Cols; c++ {
			val, _ := g.Board.Get(engine.Position{Row: r, Col: c})
			sb.WriteString(token(val))
			if c < g.Board.Cols-1 {
				sb.WriteString(" | ")
			}
		}
		if r < g.Board.Rows-1 {
			sb.WriteString("\n")
			sb.WriteString(strings.Repeat("-", g.Board.Cols*4-3))
			sb.WriteString("\n")
		}
	}
	return sb.String()
}

func lookupPlayer(players []engine.Player, id int) *engine.Player {
	for i := range players {
		if players[i].ID == id {
			return &players[i]
		}
	}
	return nil
}
