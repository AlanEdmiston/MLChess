package main

func rewardHeuristic(board Board) float64 {
	if board.isBlackCheckmated() {
		return -10000.0
	}

	if board.isWhiteCheckmated() {
		return 10000.0
	}

	if board.isStalemate() {
		return 0.0
	}

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
		if piece.pieceType.sign == "Kn" {
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
		total += nextPieceValue * colourMult
	}
	return total
}
