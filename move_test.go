package main

import "testing"

func createBoardLayout(pieces []*Piece, enPassantRank int, castleBK bool, castleBQ bool, castleWK bool, castleWQ bool, colourToMove Colour) Board {
	squares := make([][]*Piece, 8)
	for i := 0; i < 8; i++ {
		squares[i] = make([]*Piece, 8)
	}

	for _, piece := range pieces {
		squares[piece.position.X][piece.position.Y] = piece
	}

	return BoardInitialise(pieces, enPassantRank, colourToMove, castleBK, castleBQ, castleWK, castleWQ, "")
}

//https://lichess.org/editor/8/ppk5/1pb5/3n4/8/6P1/PP2p3/1KR5_w_-_-_0_1
func setTestBoardPosition() Board {
	return createBoardLayout([]*Piece{
		{pawn, White, Vector{X: 0, Y: 1}},
		{pawn, White, Vector{X: 1, Y: 1}},
		{pawn, White, Vector{X: 6, Y: 2}},
		{king, White, Vector{X: 1, Y: 0}},
		{rook, White, Vector{X: 2, Y: 0}},
		{pawn, Black, Vector{X: 0, Y: 6}},
		{pawn, Black, Vector{X: 1, Y: 6}},
		{pawn, Black, Vector{X: 1, Y: 5}},
		{king, Black, Vector{X: 2, Y: 6}},
		{bishop, Black, Vector{X: 2, Y: 5}},
		{knight, Black, Vector{X: 3, Y: 4}},
		{pawn, Black, Vector{X: 4, Y: 1}},
	}, -1, false, false, false, false, White)
}

//https://lichess.org/editor/rn2k2r/p4pp1/1p5p/1pp5/8/5bP1/PPP4P/R3KBNR_b_-_-_0_1
func setTestBoardPosition2() Board {
	return createBoardLayout([]*Piece{
		{pawn, White, Vector{X: 0, Y: 1}},
		{pawn, White, Vector{X: 1, Y: 1}},
		{pawn, White, Vector{X: 2, Y: 1}},
		{pawn, White, Vector{X: 6, Y: 2}},
		{pawn, White, Vector{X: 7, Y: 1}},
		{rook, White, Vector{X: 0, Y: 0}},
		{king, White, Vector{X: 4, Y: 0}},
		{bishop, White, Vector{X: 5, Y: 0}},
		{knight, White, Vector{X: 6, Y: 0}},
		{rook, White, Vector{X: 7, Y: 0}},

		{pawn, Black, Vector{X: 0, Y: 6}},
		{pawn, Black, Vector{X: 1, Y: 5}},
		{pawn, Black, Vector{X: 1, Y: 4}},
		{pawn, Black, Vector{X: 2, Y: 4}},
		{pawn, Black, Vector{X: 5, Y: 6}},
		{pawn, Black, Vector{X: 6, Y: 6}},
		{pawn, Black, Vector{X: 7, Y: 5}},
		{rook, Black, Vector{X: 0, Y: 7}},
		{knight, Black, Vector{X: 1, Y: 7}},
		{bishop, Black, Vector{X: 5, Y: 2}},
		{king, Black, Vector{X: 4, Y: 7}},
		{rook, Black, Vector{X: 7, Y: 7}},
	}, -1, true, true, true, true, White)
}

func TestPiecePositionChanges(t *testing.T) {
	startPosition := setTestBoardPosition()
	selectedPiece := startPosition.pieces[4]
	nextPosition := startPosition.MakeMove(selectedPiece, Vector{X: 2, Y: 1}, nil)

	if (nextPosition.pieces[4].position != Vector{X: 2, Y: 1}) {
		t.Errorf("piece has not been moved to move coordinates")
	}
	if foundPiece := nextPosition.getSquare(2, 1); foundPiece != nil && foundPiece.pieceType.sign != "R" {
		t.Errorf("no piece found at move coordinates")
	}
}

func TestEnpassantStateSets(t *testing.T) {
	startPosition := setTestBoardPosition()
	selectedPiece := startPosition.pieces[0]
	nextPosition := startPosition.MakeMove(selectedPiece, Vector{X: 0, Y: 3}, nil)

	if nextPosition.enPassantRank != 0 {
		t.Errorf("enPassantRank not set")
	}

	selectedPiece = startPosition.pieces[7]

	finalPosition := nextPosition.MakeMove(selectedPiece, Vector{X: 1, Y: 4}, nil)

	if finalPosition.enPassantRank != -1 {
		t.Errorf("enPassantRank not reset")
	}
}

func TestTakenPieceRemoved(t *testing.T) {
	startPosition := setTestBoardPosition()
	selectedPiece := startPosition.pieces[4]
	if len(startPosition.pieces) != 12 {
		t.Errorf("expected starting position with 12 pieces; got %d", len(startPosition.pieces))
	}
	nextPosition := startPosition.MakeMove(selectedPiece, Vector{X: 2, Y: 5}, nil)

	if len(nextPosition.pieces) != 11 {
		t.Errorf("expected 12 pieces after piece taken; got %d", len(nextPosition.pieces))
	}
	if (nextPosition.pieces[4].position != Vector{X: 2, Y: 5}) {
		t.Errorf("piece has not been moved to move coordinates")
	}
	if foundPiece := nextPosition.getSquare(2, 5); foundPiece != nil && foundPiece.pieceType.sign != "R" {
		t.Errorf("correct piece not found at move coordinates")
	}
}

func TestPawnPromotesToQueen(t *testing.T) {
	startPosition := setTestBoardPosition()
	selectedPiece := startPosition.pieces[11]
	nextPosition := startPosition.MakeMove(selectedPiece, Vector{X: 4, Y: 0}, &queen)

	if (nextPosition.pieces[11].position != Vector{X: 4, Y: 0}) {
		t.Errorf("piece has not been moved to move coordinates")
	}
	if nextPosition.pieces[11].pieceType.sign != "Q" {
		t.Errorf("piece has not been promoted. Type %s found", nextPosition.pieces[11].pieceType.sign)
	}
	if foundPiece := nextPosition.getSquare(4, 0); foundPiece != nil && foundPiece.pieceType.sign != "Q" {
		t.Errorf("correct piece not found at move coordinates")
	}
}

func TestCastlingStateUpdatedRookMove(t *testing.T) {
	startPosition := setTestBoardPosition2()
	selectedPiece := startPosition.pieces[5]
	nextPosition := startPosition.MakeMove(selectedPiece, Vector{X: 2, Y: 0}, nil)

	if nextPosition.canWhiteQueenSideCastle {
		t.Errorf("castling not updated after rook moved")
	}
}

func TestCastlingStateUpdatedRookTaken(t *testing.T) {
	startPosition := setTestBoardPosition2()
	selectedPiece := startPosition.pieces[19]
	nextPosition := startPosition.MakeMove(selectedPiece, Vector{X: 7, Y: 0}, nil)

	if nextPosition.canWhiteKingSideCastle {
		t.Errorf("castling not updated after rook taken")
	}
}

func TestCastlingStateUpdatedKingMove(t *testing.T) {
	startPosition := setTestBoardPosition2()
	selectedPiece := startPosition.pieces[6]
	nextPosition := startPosition.MakeMove(selectedPiece, Vector{X: 3, Y: 0}, nil)

	if nextPosition.canWhiteQueenSideCastle || nextPosition.canWhiteKingSideCastle {
		t.Errorf("castling not updated after king moved")
	}
}

func TestCanCastle(t *testing.T) {
	startPosition := setTestBoardPosition2()
	selectedPiece := startPosition.pieces[6]
	nextPosition := startPosition.MakeMove(selectedPiece, Vector{X: 2, Y: 0}, nil)

	if (nextPosition.pieces[6].position != Vector{X: 2, Y: 0}) {
		t.Errorf("king on %s has not moved as expected", nextPosition.pieces[6].position.toString())
	}
	if (nextPosition.pieces[5].position != Vector{X: 3, Y: 0}) {
		t.Errorf("rook on %s has not moved as expected", nextPosition.pieces[5].position.toString())
	}
	if nextPosition.canWhiteQueenSideCastle || nextPosition.canWhiteKingSideCastle {
		t.Errorf("castling state not updated")
	}
}

func TestCanBlockCastling(t *testing.T) {
	startPosition := setTestBoardPosition2()
	selectedPiece := startPosition.pieces[6]
	for _, move := range selectedPiece.getPossibleMoves(startPosition) {
		if (move == Vector{X: 2, Y: 0}) {
			t.Errorf("can castle when blocked")
		}
	}
}

func TestCannotMoveIntoCheck(t *testing.T) {
	startPosition := setTestBoardPosition2()
	for _, move := range startPosition.getPossibleMoves() {
		sq1 := move.squares[3][0]
		sq2 := move.squares[4][1]
		if (sq1 != nil && sq1.pieceType.sign == "K") || (sq2 != nil && sq2.pieceType.sign == "K") {
			t.Errorf("can move into check")
		}
	}
}
