package engine

import "fmt"

// Player represents a participant in a game.
type Player struct {
	ID    int
	Name  string
	Token string // short symbol for display, e.g. "X"
}

// Move represents an action on the board.
type Move struct {
	PlayerID int
	Pos      Position
}

// Outcome captures the result of a completed game.
type Outcome struct {
	Winner *Player
	Draw   bool
}

// Game tracks shared state used by rule implementations.
type Game struct {
	Board        *Board
	Players      []Player
	currentIndex int
	Log          []Move
	Outcome      Outcome
}

// NewGame constructs a Game with the provided board and players.
func NewGame(board *Board, players []Player) (*Game, error) {
	if board == nil {
		return nil, fmt.Errorf("board is required")
	}
	if len(players) == 0 {
		return nil, fmt.Errorf("at least one player required")
	}
	playerIDs := make(map[int]bool)
	for _, p := range players {
		if p.ID == 0 {
			return nil, fmt.Errorf("player %q must use non-zero ID", p.Name)
		}
		if playerIDs[p.ID] {
			return nil, fmt.Errorf("duplicate player ID %d", p.ID)
		}
		playerIDs[p.ID] = true
	}
	return &Game{
		Board:   board,
		Players: players,
	}, nil
}

// CurrentPlayer returns the player whose turn it is.
func (g *Game) CurrentPlayer() Player {
	return g.Players[g.currentIndex]
}

// AdvanceTurn increments the turn order when the game is still active.
func (g *Game) AdvanceTurn() {
	if g.Outcome.Winner != nil || g.Outcome.Draw {
		return
	}
	g.currentIndex = (g.currentIndex + 1) % len(g.Players)
}

// RecordMove appends a move to the log and updates the board.
func (g *Game) RecordMove(m Move) error {
	if m.PlayerID != g.CurrentPlayer().ID {
		return fmt.Errorf("not player %d's turn", m.PlayerID)
	}
	if err := g.Board.Set(m.Pos, m.PlayerID); err != nil {
		return err
	}
	g.Log = append(g.Log, m)
	return nil
}

// EndGame stores the final outcome.
func (g *Game) EndGame(outcome Outcome) {
	g.Outcome = outcome
}
