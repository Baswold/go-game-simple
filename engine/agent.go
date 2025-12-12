package engine

import (
	"errors"
	"fmt"
	"math/rand"
)

// Agent chooses the next move given the current game and list of valid options.
type Agent interface {
	ChooseMove(g *Game, moves []Move) (Move, error)
}

// RandomAgent picks any available move uniformly at random.
type RandomAgent struct {
	Rand *rand.Rand
}

// ChooseMove implements Agent.
func (a *RandomAgent) ChooseMove(_ *Game, moves []Move) (Move, error) {
	if len(moves) == 0 {
		return Move{}, errors.New("no moves to choose from")
	}
	r := a.Rand
	if r == nil {
		r = rand.New(rand.NewSource(rand.Int63())) //nolint:gosec // best-effort randomness for simulations
	}
	return moves[r.Intn(len(moves))], nil
}

// ScriptedAgent replays a preset sequence of positions. Useful for tests or demos.
type ScriptedAgent struct {
	Positions []Position
	next      int
	Fallback  Agent
}

// ChooseMove selects the next scripted move when available, otherwise defers to Fallback.
func (a *ScriptedAgent) ChooseMove(_ *Game, moves []Move) (Move, error) {
	// try to match a scripted move that is valid
	for a.next < len(a.Positions) {
		pos := a.Positions[a.next]
		a.next++
		for _, m := range moves {
			if m.Pos == pos {
				return m, nil
			}
		}
	}
	if a.Fallback != nil {
		return a.Fallback.ChooseMove(nil, moves)
	}
	return Move{}, fmt.Errorf("no scripted moves left and no fallback provided")
}
