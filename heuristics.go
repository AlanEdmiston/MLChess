package main

func generalHeuristic(boardState *Board, config *PieceValueConfig) float64 {
	total := 0.0

	for _, piece := range boardState.pieces {
		colourMult := 1.0
		if piece.colour == Black {
			colourMult = -1
		}

		noAlliedPieces := 0
		noOppPieces := 0
		if piece.colour == White {
			noAlliedPieces = boardState.noWhitePieces()
			noOppPieces = boardState.noBlackPieces()
		} else {
			noAlliedPieces = boardState.noBlackPieces()
			noOppPieces = boardState.noWhitePieces()
		}

		alliedPieceTypeModifiers := []float64{}
		oppPieceTypeModifiers := []float64{}

		for _, otherPiece := range boardState.pieces {
			if otherPiece.colour == piece.colour {
				alliedPieceTypeModifiers = append(alliedPieceTypeModifiers, config.RemainingAlliedPiecesTypeMod[piece.pieceType.sign][otherPiece.pieceType.sign])
			} else {
				oppPieceTypeModifiers = append(oppPieceTypeModifiers, config.RemainingOpponentPiecesTypeMod[piece.pieceType.sign][otherPiece.pieceType.sign])
			}
		}

		pieceValue := config.BaseValues[piece.pieceType.sign] * config.PositionMod[piece.pieceType.sign][piece.position] *
			config.RemainingAlliedPiecesMod[piece.pieceType.sign][noAlliedPieces] * config.RemainingOpponentPiecesMod[piece.pieceType.sign][noOppPieces] *
			PI(alliedPieceTypeModifiers) * PI(oppPieceTypeModifiers)

		for _, square := range piece.getCoveredSquares(*boardState) {
			pieceValue += config.SquareBaseValues[Vector{X: square.X, Y: square.Y}] * config.CoveredByMod[piece.pieceType.sign]
		}

		total += pieceValue * colourMult
	}

	return total
}
