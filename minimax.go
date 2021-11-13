package main

import (
	"math"
	"sort"
)

type Strategy string

const (
	//Player1
	Player1 Strategy = "Player1"
	//Player2
	Player2 Strategy = "Player2"
)

//MinimaxTree is a min-max tree of possible future moves
type MinimaxTree struct {
	children      []*MinimaxTree
	value         float64
	height        int
	depth         int
	isMaximising  bool
	state         *Board
	nextMove      *Board
	nextMoveIndex int
	worstSibling  *float64
	bestChild     *float64
	strategy      Strategy
}

func createRoot(height int, colour Colour, state *Board, strategy Strategy) MinimaxTree {

	config1 := randomConfigGenerator()
	config2 := randomConfigGenerator()

	root := MinimaxTree{
		children:     []*MinimaxTree{},
		value:        0.0,
		height:       height,
		depth:        0,
		isMaximising: colour == White,
		state:        state,
		nextMove:     nil,
		worstSibling: nil,
		bestChild:    nil,
		strategy:     strategy,
	}

	moveValue := 0.0
	currentMove := root.state

	for i, nextState := range root.state.children {
		child := newMinimaxTree(*nextState, &root, &config1, &config2)
		root.children = append(root.children, &child)
		if i == 0 || (child.value > moveValue && root.isMaximising) || (child.value < moveValue && !root.isMaximising) {
			moveValue = child.value
			currentMove = child.state
		}
	}
	root.nextMove = currentMove
	return root
}

func newMinimaxTree(boardState Board, parent *MinimaxTree, config1 *pieceValueConfig, config2 *pieceValueConfig) MinimaxTree {
	boardState.children = boardState.getPossibleMoves()

	children := []*MinimaxTree{}
	depth := parent.depth + 1
	isMaximising := !parent.isMaximising
	value := 0.0
	worstSibling := math.Inf(1)
	bestChild := math.Inf(-1)

	if !isMaximising {
		worstSibling = math.Inf(-1)
		bestChild = math.Inf(1)
	}

	if parent.bestChild != nil {
		worstSibling = *parent.bestChild
	}

	if depth == parent.height {
		//switch out heuristic here
		switch parent.strategy {
		case Player1:
			value = generalHeuristic(&boardState, config1)
		case Player2:
			value = generalHeuristic(&boardState, config2)
		default:
			value = verySimpleHeuristic(boardState)
		}

		return MinimaxTree{children, value, parent.height, depth, isMaximising, &boardState, nil, -1, &worstSibling, nil, parent.strategy}
	}

	tree := MinimaxTree{children, value, parent.height, depth, isMaximising, &boardState, nil, -1, &worstSibling, nil, parent.strategy}

	//try to sort moves to calculate by most promising
	rankedStates := map[float64]*Board{}
	colourMult := 1.0
	if boardState.colourToMove == Black {
		colourMult = -1.0
	}
	for _, nextState := range boardState.children {
		rankedStates[verySimpleHeuristic(*nextState)*colourMult] = nextState
	}

	keys := []float64{}
	for key := range rankedStates {
		keys = append(keys, key)
	}
	sort.Float64s(keys)

	for _, key := range keys {
		nextState := rankedStates[key]

		child := newMinimaxTree(*nextState, &tree, config1, config2)
		if isMaximising {
			bestChild = math.Max(child.value, bestChild)
		} else {
			bestChild = math.Min(child.value, bestChild)
		}
		tree.bestChild = &bestChild

		children = append(children, &child)
		if isMaximising {
			if bestChild >= worstSibling {
				break
			}
		} else {
			if bestChild <= worstSibling {
				break
			}
		}

	}

	if len(children) == 0 {
		if boardState.isBlackCheckmated() {
			return MinimaxTree{children, math.Inf(-1), parent.height, depth, isMaximising, &boardState, nil, -1, &worstSibling, nil, parent.strategy}
		}
		if boardState.isWhiteCheckmated() {
			return MinimaxTree{children, math.Inf(1), parent.height, depth, isMaximising, &boardState, nil, -1, &worstSibling, nil, parent.strategy}
		}
		if boardState.isStalemate() {
			return MinimaxTree{children, 0, parent.height, depth, isMaximising, &boardState, nil, -1, &worstSibling, nil, parent.strategy}
		}
		print("minimax error: no child trees\n")
	}

	bestValue := children[0].value
	bestMove := children[0].state
	tree.nextMoveIndex = 0
	if len(children) > 1 {
		if isMaximising {
			for i, child := range children[1:] {
				if child.value > bestValue {
					bestValue = child.value
					bestMove = child.state
					tree.nextMoveIndex = i
				}
			}
		} else {
			for i, child := range children[1:] {
				if child.value < bestValue {
					bestValue = child.value
					bestMove = child.state
					tree.nextMoveIndex = i
				}
			}
		}
	}

	tree.value = bestValue
	tree.nextMove = bestMove
	return tree
}

func verySimpleHeuristic(board Board) float64 {
	total := 0.0
	for _, piece := range board.pieces {
		colourMult := 1.0
		if piece.colour == Black {
			colourMult = -1
		}
		nextPieceValue := 0.0
		if piece.pieceType.sign == "P" {
			nextPieceValue = 1.0
			if piece.colour == White {
				nextPieceValue += float64(piece.position.y) * 0.1
			} else {
				nextPieceValue += (7.0 - float64(piece.position.y)) * 0.1
			}
		}
		if piece.pieceType.sign == "N" {
			nextPieceValue = 3.0
		}
		if piece.pieceType.sign == "B" {
			nextPieceValue = 3.25
		}
		if piece.pieceType.sign == "R" {
			nextPieceValue = 5.0
		}
		if piece.pieceType.sign == "Q" {
			nextPieceValue = 9.0
		}
		if piece.pieceType.sign == "K" {
			nextPieceValue = 1000000.0
		}
		total += nextPieceValue * colourMult
	}
	return total
}
