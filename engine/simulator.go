package engine

import "fmt"

// Rule encapsulates game-specific logic for validating and applying moves.
type Rule interface {
	ValidMoves(g *Game) []Move
	ApplyMove(g *Game, m Move) error
	Status(g *Game) (Outcome, bool) // bool indicates game is finished
}

// Play runs a full game using the provided rule and agents until completion.
func Play(g *Game, rule Rule, agents map[int]Agent) (Outcome, error) {
	if rule == nil {
		return Outcome{}, fmt.Errorf("rule is required")
	}
	if len(agents) == 0 {
		return Outcome{}, fmt.Errorf("at least one agent is required")
	}

	// Limit turns to avoid infinite loops if rules never end the game.
	turnLimit := g.Board.Rows*g.Board.Cols*4 + len(g.Players)

	for turn := 0; turn < turnLimit; turn++ {
		if outcome, done := rule.Status(g); done {
			g.EndGame(outcome)
			return outcome, nil
		}

		validMoves := rule.ValidMoves(g)
		if len(validMoves) == 0 {
			outcome := Outcome{Draw: true}
			g.EndGame(outcome)
			return outcome, nil
		}

		current := g.CurrentPlayer()
		agent := agents[current.ID]
		if agent == nil {
			return Outcome{}, fmt.Errorf("no agent registered for player %d", current.ID)
		}
		move, err := agent.ChooseMove(g, validMoves)
		if err != nil {
			return Outcome{}, err
		}

		if err := rule.ApplyMove(g, move); err != nil {
			return Outcome{}, err
		}
		g.AdvanceTurn()
	}

	outcome := Outcome{Draw: true}
	g.EndGame(outcome)
	return outcome, fmt.Errorf("turn limit exceeded; forcing draw")
}
