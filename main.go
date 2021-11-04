package main

import (
	"fmt"
	"time"
)

type player struct {
	colour   string
	strategy string
}

func main() {
	state := NewBoard()
	print(state.ToString())

	for i := 0; i < 11; i++ {
		print(time.Now().Unix())
		print("-------------------------------------------------\n")
		tree := createRoot(4, White, &state, Player1)
		nextMove := tree.nextMove
		print(fmt.Sprint(nextMove.ToString()))
		state = *nextMove

		print(time.Now().Unix())
		print("-------------------------------------------------\n")
		tree = createRoot(4, Black, &state, Player2)
		nextMove = tree.nextMove
		print(fmt.Sprint(nextMove.ToString()))
		state = *nextMove
	}
}
