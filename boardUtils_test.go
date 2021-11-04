package main

import "testing"

//https://lichess.org/editor/3k4/Q7/4K3/p5r1/1p6/1P6/6p1/8_w_-_-_0_1
func setUtilsTestBoardPosition() Board {
	return createBoardLayout([]*Piece{
		{queen, White, Vector{x: 0, y: 6}},
		{king, White, Vector{x: 4, y: 5}},
		{pawn, White, Vector{x: 1, y: 2}},
		{rook, Black, Vector{x: 6, y: 4}},
		{king, Black, Vector{x: 3, y: 7}},
		{pawn, Black, Vector{x: 6, y: 1}},
		{pawn, Black, Vector{x: 0, y: 4}},
		{pawn, Black, Vector{x: 1, y: 3}},
	}, -1, false, false, false, false, White)
}

func TestGetCoveredSquares(t *testing.T) {
	state := setUtilsTestBoardPosition()
	blackSquares := state.getCoveredSquares(Black)

	blackSquaresTotal := 0
	for _, row := range blackSquares {
		for _, square := range row {
			if square {
				blackSquaresTotal++
			}
		}
	}

	if blackSquaresTotal != 23 {
		t.Errorf("expected 23 covered squares but got %d", blackSquaresTotal)
	}

	whiteSquares := state.getCoveredSquares(Black)
	whiteSquaresTotal := 0
	for _, row := range whiteSquares {
		for _, square := range row {
			if square {
				whiteSquaresTotal++
			}
		}
	}

	if whiteSquaresTotal != 23 {
		t.Errorf("expected 23 covered squares but got %d", whiteSquaresTotal)
	}
}

func TestCanDetermineCheck(t *testing.T) {
	state := setUtilsTestBoardPosition()
	nextState := state.MakeMove(state.pieces[3], Vector{6, 5}, nil)

	if state.getIsWhiteChecked() {
		t.Errorf("white shouldnt be in check but is")
	}

	if !nextState.getIsWhiteChecked() {
		t.Errorf("white should be in check but isnt")
	}
}

func TestCanDetermineCheckmate(t *testing.T) {
	state := setUtilsTestBoardPosition()
	nextState := state.MakeMove(state.pieces[0], Vector{3, 6}, nil)

	if state.isBlackCheckmated() {
		t.Errorf("black shouldnt be in checkmate but is")
	}

	if !nextState.isBlackCheckmated() {
		t.Errorf("black should be in checkmate but isnt")
	}
}

func TestCanDetermineStalemate(t *testing.T) {
	//https://lichess.org/editor/6k1/4p3/4K2Q/8/8/8/8/8_w_-_-_0_1
	state := createBoardLayout([]*Piece{
		{queen, White, Vector{x: 7, y: 5}},
		{king, White, Vector{x: 4, y: 5}},
		{king, Black, Vector{x: 6, y: 7}},
		{pawn, Black, Vector{x: 4, y: 6}},
	}, -1, false, false, false, false, Black)

	if !state.isStalemate() {
		t.Errorf("should be stalemate for black to move but isnt")
	}
}
