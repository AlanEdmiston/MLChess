package main

import (
	"math/rand"
	"time"
)

type diminishingMultiplier struct {
	coefficient float64
	exponent    float64
}

// type pieceDupleFormation map[string]map[string]map[Vector]map[Vector]float64

type pieceValueConfig struct {
	baseValues                     map[string]float64
	positionMod                    map[string]map[Vector]float64
	remainingAlliedPiecesMod       map[string]map[int]float64
	remainingOpponentPiecesMod     map[string]map[int]float64
	remainingAlliedPiecesTypeMod   map[string]map[string]float64
	remainingOpponentPiecesTypeMod map[string]map[string]float64
	// possibleMovesMult              map[string]float64
	// possibleMovesAdder             map[string]float64
	// attackMult                     map[string]map[string]float64
	// attackMultExp                  map[string]float64
	// attackAdder                    map[string]map[string]float64
	// attackAdderCoeff               map[string]float64
	// defendMult                     map[string]map[string]float64
	// defendMultExp                  map[string]float64
	// defendAdder                    map[string]map[string]float64
	// defendAdderCoeff               map[string]float64
	// skewerMult                     map[string]map[string]map[string]float64
	// skewerMultExp                  map[string]map[string]float64
	// skewerAdder                    map[string]map[string]map[string]float64
	// skewerAdderCoeff               map[string]map[string]float64
	squareBaseValues map[Vector]float64
	coveredByMod     map[string]float64
	// pieceDupleFormation            pieceDupleFormation
	// pieceFormationChainer          map[*pieceDupleFormation]map[*pieceDupleFormation]float64
}

func generateRandomGrid(sd float64, mean float64) [8][8]float64 {
	grid := [8][8]float64{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			grid[i][j] = rand.NormFloat64()*sd + mean
		}
	}
	return grid
}

func generateRandomVectorMap(sd float64, mean float64) map[Vector]float64 {
	vectorMap := map[Vector]float64{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			vectorMap[Vector{x: i, y: j}] = rand.NormFloat64()*sd + mean
		}
	}
	return vectorMap
}

func generateRandomIntMap(length int, sd float64, mean float64) map[int]float64 {
	intMap := map[int]float64{}
	for i := 0; i < length; i++ {
		intMap[i] = rand.NormFloat64()*sd + mean
	}
	return intMap
}

func generateRandomPieceMap(sd float64, mean float64) map[string]float64 {
	return map[string]float64{
		"P":  rand.NormFloat64()*sd + mean,
		"R":  rand.NormFloat64()*sd + mean,
		"Kn": rand.NormFloat64()*sd + mean,
		"B":  rand.NormFloat64()*sd + mean,
		"Q":  rand.NormFloat64()*sd + mean,
		"K":  rand.NormFloat64()*sd + mean,
	}
}

func randomConfigGenerator() pieceValueConfig {
	rand.Seed(time.Now().UnixNano())

	baseValues := generateRandomPieceMap(5, 5)

	positionMod := map[string]map[Vector]float64{
		"P":  generateRandomVectorMap(0.1, 1),
		"R":  generateRandomVectorMap(0.1, 1),
		"Kn": generateRandomVectorMap(0.1, 1),
		"B":  generateRandomVectorMap(0.1, 1),
		"Q":  generateRandomVectorMap(0.1, 1),
		"K":  generateRandomVectorMap(0.1, 1),
	}

	remainingAlliedPiecesMod := map[string]map[int]float64{
		"P":  generateRandomIntMap(16, 0.1, 1),
		"R":  generateRandomIntMap(16, 0.1, 1),
		"Kn": generateRandomIntMap(16, 0.1, 1),
		"B":  generateRandomIntMap(16, 0.1, 1),
		"Q":  generateRandomIntMap(16, 0.1, 1),
		"K":  generateRandomIntMap(16, 0.1, 1),
	}

	remainingOpponentPiecesMod := map[string]map[int]float64{
		"P":  generateRandomIntMap(16, 0.1, 1),
		"R":  generateRandomIntMap(16, 0.1, 1),
		"Kn": generateRandomIntMap(16, 0.1, 1),
		"B":  generateRandomIntMap(16, 0.1, 1),
		"Q":  generateRandomIntMap(16, 0.1, 1),
		"K":  generateRandomIntMap(16, 0.1, 1),
	}

	remainingAlliedPiecesTypeMod := map[string]map[string]float64{
		"P":  generateRandomPieceMap(0.2, 1),
		"R":  generateRandomPieceMap(0.2, 1),
		"Kn": generateRandomPieceMap(0.2, 1),
		"B":  generateRandomPieceMap(0.2, 1),
		"Q":  generateRandomPieceMap(0.2, 1),
		"K":  generateRandomPieceMap(0.2, 1),
	}

	remainingOpponentPiecesTypeMod := map[string]map[string]float64{
		"P":  generateRandomPieceMap(0.2, 1),
		"R":  generateRandomPieceMap(0.2, 1),
		"Kn": generateRandomPieceMap(0.2, 1),
		"B":  generateRandomPieceMap(0.2, 1),
		"Q":  generateRandomPieceMap(0.2, 1),
		"K":  generateRandomPieceMap(0.2, 1),
	}

	// possibleMovesMult              map[string]float64
	// possibleMovesAdder             map[string]float64
	// attackMult                     map[string]map[string]float64
	// attackMultExp                  map[string]float64
	// attackAdder                    map[string]map[string]float64
	// attackAdderCoeff               map[string]float64
	// defendMult                     map[string]map[string]float64
	// defendMultExp                  map[string]float64
	// defendAdder                    map[string]map[string]float64
	// defendAdderCoeff               map[string]float64
	// skewerMult                     map[string]map[string]map[string]float64
	// skewerMultExp                  map[string]map[string]float64
	// skewerAdder                    map[string]map[string]map[string]float64
	// skewerAdderCoeff               map[string]map[string]float64

	squareBaseValues := generateRandomVectorMap(0.02, 0.2)

	coveredByMod := generateRandomPieceMap(0.1, 1)

	return pieceValueConfig{
		baseValues:                     baseValues,
		positionMod:                    positionMod,
		remainingAlliedPiecesMod:       remainingAlliedPiecesMod,
		remainingOpponentPiecesMod:     remainingOpponentPiecesMod,
		remainingAlliedPiecesTypeMod:   remainingAlliedPiecesTypeMod,
		remainingOpponentPiecesTypeMod: remainingOpponentPiecesTypeMod,
		// possibleMovesMult: possibleMovesMult
		// possibleMovesAdder: possibleMovesAdder
		// attackMult: attackMult
		// attackMultExp: attackMultExp
		// attackAdder: attackAdder
		// attackAdderCoeff: attackAdderCoeff
		// defendMult: defendMult
		// defendMultExp: defendMultExp
		// defendAdder: defendAdder
		// defendAdderCoeff: defendAdderCoeff
		// skewerMult: skewerMult
		// skewerMultExp: skewerMultExp
		// skewerAdder: skewerAdder
		// skewerAdderCoeff: skewerAdderCoeff
		squareBaseValues: squareBaseValues,
		coveredByMod:     coveredByMod,
	}
}

func generalHeuristic(boardState *Board, config *pieceValueConfig) float64 {
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
				alliedPieceTypeModifiers = append(alliedPieceTypeModifiers, config.remainingAlliedPiecesTypeMod[piece.pieceType.sign][otherPiece.pieceType.sign])
			} else {
				oppPieceTypeModifiers = append(oppPieceTypeModifiers, config.remainingOpponentPiecesTypeMod[piece.pieceType.sign][otherPiece.pieceType.sign])
			}
		}

		pieceValue := config.baseValues[piece.pieceType.sign] * config.positionMod[piece.pieceType.sign][piece.position] *
			config.remainingAlliedPiecesMod[piece.pieceType.sign][noAlliedPieces] * config.remainingOpponentPiecesMod[piece.pieceType.sign][noOppPieces] *
			PI(alliedPieceTypeModifiers) * PI(oppPieceTypeModifiers)

		for _, square := range piece.getCoveredSquares(*boardState) {
			pieceValue += config.squareBaseValues[Vector{x: square.x, y: square.y}] * config.coveredByMod[piece.pieceType.sign]
		}

		total += pieceValue * colourMult
	}

	return total
}
