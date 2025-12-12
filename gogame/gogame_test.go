package gogame

import (
	"testing"
)

func TestCaptureSingleStone(t *testing.T) {
	g, err := NewGame(3)
	if err != nil {
		t.Fatalf("new game: %v", err)
	}
	sequence := []string{
		"A2", // Black
		"B2", // White stone to be captured
		"B1", // Black
		"C3", // White elsewhere
		"C2", // Black
		"C1", // White elsewhere
		"B3", // Black completes surround; captures B2
	}
	for _, coord := range sequence {
		play(t, g, coord)
	}

	pos, _ := ParseCoord("B2", g.Size)
	val, _ := g.Board.Get(pos)
	if val != 0 {
		t.Fatalf("expected captured stone at B2 to be empty, got %d", val)
	}
	if g.Captures[Black] != 1 {
		t.Fatalf("expected Black captures=1, got %d", g.Captures[Black])
	}
}

func TestSuicideRejected(t *testing.T) {
	g, err := NewGame(3)
	if err != nil {
		t.Fatalf("new game: %v", err)
	}
	sequence := []string{
		"A2", // B
		"C1", // W
		"B1", // B
		"C3", // W
		"B3", // B
		"A1", // W
		"C2", // B -> White to play
	}
	for _, coord := range sequence {
		play(t, g, coord)
	}

	// Now White tries to play at B2 with no liberties.
	pos, _ := ParseCoord("B2", g.Size)
	if _, err := g.PlayMove(pos); err == nil {
		t.Fatalf("expected suicide move to be rejected")
	}
}

func play(t *testing.T, g *Game, coord string) {
	pos, err := ParseCoord(coord, g.Size)
	if err != nil {
		t.Fatalf("parse %s: %v", coord, err)
	}
	if _, err := g.PlayMove(pos); err != nil {
		t.Fatalf("play %s: %v", coord, err)
	}
}
