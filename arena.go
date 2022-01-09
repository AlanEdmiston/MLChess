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
				playMatchWithResult(file.Name(), otherFile.Name(), dir, scores)
			}
		}
	}

	print(scores)
}

func playMatchWithResult(file1, file2, dir string, scores map[string]int) {
	result := playMatch(file1, file2, dir)
	print(result)
	if result == WhiteWon {
		scores[file1] = scores[file1] + 2
	}
	if result == BlackWon {
		scores[file2] = scores[file2] + 2
	}
	if result == Stalemate {
		scores[file1] = scores[file1] + 1
		scores[file2] = scores[file2] + 1
	}
}

func playMatch(file1, file2, dir string) WinState {
	config1 := readConfigJson(file1, dir)
	config2 := readConfigJson(file2, dir)
	state := NewBoard()

	gameString := "\n----------\n"

	pastStates := []Board{}
	for i := 0; i < 200; i++ {
		pastStates = append(pastStates, state)

		tree := createRoot(4, White, &state, Player1, config1, config2)
		nextMove := tree.nextMove
		gameString += fmt.Sprintf("%d. %s ", i+1, nextMove.lastMoveString)
		state = *nextMove

		tree = createRoot(4, Black, &state, Player2, config1, config2)
		nextMove = tree.nextMove
		gameString += fmt.Sprint(nextMove.lastMoveString)
		state = *nextMove

		if state.winner != Undecided {
			print(gameString)
			return state.winner
		}

		threemoveRuleCount := 0
		for _, pastState := range pastStates {
			if pastState.ToString() == state.ToString() {
				threemoveRuleCount++
			}
		}
		if threemoveRuleCount >= 2 {
			print(gameString)
			return Stalemate
		}
	}
	print(gameString)
	return Stalemate
}
