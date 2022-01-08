package main

import (
	"fmt"
	"io/ioutil"
)

func tournament(dir string) {
	files, _ := ioutil.ReadDir(dir)

	scores := map[string]int{}
	for _, file := range files {
		scores[file.Name()] = 0
	}

	for _, file := range files {
		for _, otherFile := range files {
			if file.Name() != otherFile.Name() {
				result := playMatch(file.Name(), otherFile.Name(), dir)
				print(result)
				if result == WhiteWon {
					scores[file.Name()] = scores[file.Name()] + 2
				}
				if result == BlackWon {
					scores[otherFile.Name()] = scores[otherFile.Name()] + 2
				}
				if result == Stalemate {
					scores[file.Name()] = scores[file.Name()] + 1
					scores[otherFile.Name()] = scores[otherFile.Name()] + 1
				}
			}
		}
	}

	print(scores)
}

func playMatch(file1, file2, dir string) WinState {
	config1 := readConfigJson(file1, dir)
	config2 := readConfigJson(file2, dir)
	state := NewBoard()

	pastStates := []Board{}
	print("----------\n")
	for i := 0; i < 200; i++ {
		pastStates = append(pastStates, state)

		tree := createRoot(4, White, &state, Player1, config1, config2)
		nextMove := tree.nextMove
		print(i + 1)
		print(". ")
		print(fmt.Sprint(nextMove.lastMoveString))
		print(" ")
		state = *nextMove

		tree = createRoot(4, Black, &state, Player2, config1, config2)
		nextMove = tree.nextMove
		print(fmt.Sprint(nextMove.lastMoveString))
		print(" ")
		state = *nextMove

		if state.winner != Undecided {
			return state.winner
		}

		threemoveRuleCount := 0
		for _, pastState := range pastStates {
			if pastState.ToString() == state.ToString() {
				threemoveRuleCount++
			}
		}
		if threemoveRuleCount >= 2 {
			print(state.ToString())
			return Stalemate
		}
	}
	return Stalemate
}
