package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"boardgame/gogame"
)

func main() {
	size := flag.Int("size", 9, "board size (commonly 9, 13, or 19)")
	flag.Parse()

	game, err := gogame.NewGame(*size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot start game: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Go game on %dx%d board. Coordinates like D4, row numbers from bottom.\n", *size, *size)
	fmt.Println("Commands: coordinate to play, 'pass' to pass, 'quit' to exit.")
	printBoard(game)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\nMove %d - %s to play: ", game.MoveNumber()+1, game.ToPlay)
		raw, _ := reader.ReadString('\n')
		input := strings.TrimSpace(strings.ToLower(raw))
		switch input {
		case "q", "quit", "exit":
			fmt.Println("Exiting.")
			return
		case "pass":
			game.Pass()
			fmt.Println("Player passed.")
			if game.ConsecutivePasses >= 2 {
				fmt.Println("Both players passed. Game over.")
				printBoard(game)
				return
			}
			printBoard(game)
			continue
		}

		pos, err := gogame.ParseCoord(input, game.Size)
		if err != nil {
			fmt.Printf("Invalid move: %v\n", err)
			continue
		}
		result, err := game.PlayMove(pos)
		if err != nil {
			fmt.Printf("Invalid move: %v\n", err)
			continue
		}
		if result.Captured > 0 {
			fmt.Printf("Captured %d stones.\n", result.Captured)
		}
		printBoard(game)
	}
}

func printBoard(game *gogame.Game) {
	fmt.Println(gogame.RenderBoardASCII(game))
	fmt.Printf("Captures - Black: %d, White: %d. To play: %s.\n", game.Captures[gogame.Black], game.Captures[gogame.White], game.ToPlay)
}
