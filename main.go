package main

import (
	"fmt"
)

type player struct {
	colour   string
	strategy string
}

func main() {
	state := NewBoard()

	for i := 0; i < 11; i++ {
		tree := createRoot(4, White, &state, Player1)
		nextMove := tree.nextMove
		print(i + 1)
		print(". ")
		print(fmt.Sprint(nextMove.lastMoveString))
		print(" ")
		state = *nextMove

		tree = createRoot(4, Black, &state, Player2)
		nextMove = tree.nextMove
		print(fmt.Sprint(nextMove.lastMoveString))
		print(" ")
		state = *nextMove
	}
}
