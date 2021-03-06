package main

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// type pieceDupleFormation map[string]map[string]map[Vector]map[Vector]float64

type PieceValueConfig struct {
	BaseValues                     map[string]float64            `json:"baseValues"`
	PositionMod                    map[string]map[Vector]float64 `json:"positionMod"`
	RemainingAlliedPiecesMod       map[string]map[int]float64    `json:"remainingAlliedPiecesMod"`
	RemainingOpponentPiecesMod     map[string]map[int]float64    `json:"remainingOpponentPiecesMod"`
	RemainingAlliedPiecesTypeMod   map[string]map[string]float64 `json:"remainingAlliedPiecesTypeMod"`
	RemainingOpponentPiecesTypeMod map[string]map[string]float64 `json:"remainingOpponentPiecesTypeMod"`
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
	SquareBaseValues map[Vector]float64 `json:"squareBaseValues"`
	CoveredByMod     map[string]float64 `json:"coveredByMod"`
	// pieceDupleFormation            pieceDupleFormation
	// pieceFormationChainer          map[*pieceDupleFormation]map[*pieceDupleFormation]float64
}

type PieceValueConfigJsonified struct {
	BaseValues                     map[string]float64            `json:"baseValues"`
	PositionMod                    map[string]map[string]float64 `json:"positionMod"`
	RemainingAlliedPiecesMod       map[string]map[int]float64    `json:"remainingAlliedPiecesMod"`
	RemainingOpponentPiecesMod     map[string]map[int]float64    `json:"remainingOpponentPiecesMod"`
	RemainingAlliedPiecesTypeMod   map[string]map[string]float64 `json:"remainingAlliedPiecesTypeMod"`
	RemainingOpponentPiecesTypeMod map[string]map[string]float64 `json:"remainingOpponentPiecesTypeMod"`
	SquareBaseValues               map[string]float64            `json:"squareBaseValues"`
	CoveredByMod                   map[string]float64            `json:"coveredByMod"`
}

// func generateRandomGrid(sd float64, mean float64) [8][8]float64 {
// 	grid := [8][8]float64{}
// 	for i := 0; i < 8; i++ {
// 		for j := 0; j < 8; j++ {
// 			grid[i][j] = rand.NormFloat64()*sd + mean
// 		}
// 	}
// 	return grid
// }

func (pieceValueConfig PieceValueConfig) marshalJson() PieceValueConfigJsonified {
	output := PieceValueConfigJsonified{
		BaseValues:                     pieceValueConfig.BaseValues,
		PositionMod:                    map[string]map[string]float64{},
		RemainingAlliedPiecesMod:       pieceValueConfig.RemainingAlliedPiecesMod,
		RemainingOpponentPiecesMod:     pieceValueConfig.RemainingOpponentPiecesMod,
		RemainingAlliedPiecesTypeMod:   pieceValueConfig.RemainingAlliedPiecesTypeMod,
		RemainingOpponentPiecesTypeMod: pieceValueConfig.RemainingOpponentPiecesTypeMod,
		SquareBaseValues:               map[string]float64{},
		CoveredByMod:                   pieceValueConfig.CoveredByMod,
	}

	for key, vectMap := range pieceValueConfig.PositionMod {
		modifiedMap := map[string]float64{}
		for vectKey, val := range vectMap {
			modifiedMap[vectKey.toString()] = val
		}

		output.PositionMod[key] = modifiedMap
	}

	for key, val := range pieceValueConfig.SquareBaseValues {
		output.SquareBaseValues[key.toString()] = val
	}

	return output
}

func (pieceValueConfig PieceValueConfigJsonified) unmarshalJson() PieceValueConfig {
	output := PieceValueConfig{
		BaseValues:                     pieceValueConfig.BaseValues,
		PositionMod:                    map[string]map[Vector]float64{},
		RemainingAlliedPiecesMod:       pieceValueConfig.RemainingAlliedPiecesMod,
		RemainingOpponentPiecesMod:     pieceValueConfig.RemainingOpponentPiecesMod,
		RemainingAlliedPiecesTypeMod:   pieceValueConfig.RemainingAlliedPiecesTypeMod,
		RemainingOpponentPiecesTypeMod: pieceValueConfig.RemainingOpponentPiecesTypeMod,
		SquareBaseValues:               map[Vector]float64{},
		CoveredByMod:                   pieceValueConfig.CoveredByMod,
	}

	for key, vectMap := range pieceValueConfig.PositionMod {
		modifiedMap := map[Vector]float64{}
		for vectKey, val := range vectMap {
			modifiedMap[fromString(vectKey)] = val
		}

		output.PositionMod[key] = modifiedMap
	}

	for key, val := range pieceValueConfig.SquareBaseValues {
		output.SquareBaseValues[fromString(key)] = val
	}

	return output
}

func generateRandomVectorMap(sd float64, mean float64) map[Vector]float64 {
	vectorMap := map[Vector]float64{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			vectorMap[Vector{X: i, Y: j}] = rand.NormFloat64()*sd + mean
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
		"P": rand.NormFloat64()*sd + mean,
		"R": rand.NormFloat64()*sd + mean,
		"N": rand.NormFloat64()*sd + mean,
		"B": rand.NormFloat64()*sd + mean,
		"Q": rand.NormFloat64()*sd + mean,
		"K": rand.NormFloat64()*sd + mean,
	}
}

func randomConfigGenerator(name string, dir string) {
	rand.Seed(time.Now().UnixNano())

	baseValues := generateRandomPieceMap(5, 5)

	positionMod := map[string]map[Vector]float64{
		"P": generateRandomVectorMap(0.1, 1),
		"R": generateRandomVectorMap(0.1, 1),
		"N": generateRandomVectorMap(0.1, 1),
		"B": generateRandomVectorMap(0.1, 1),
		"Q": generateRandomVectorMap(0.1, 1),
		"K": generateRandomVectorMap(0.1, 1),
	}

	remainingAlliedPiecesMod := map[string]map[int]float64{
		"P": generateRandomIntMap(16, 0.1, 1),
		"R": generateRandomIntMap(16, 0.1, 1),
		"N": generateRandomIntMap(16, 0.1, 1),
		"B": generateRandomIntMap(16, 0.1, 1),
		"Q": generateRandomIntMap(16, 0.1, 1),
		"K": generateRandomIntMap(16, 0.1, 1),
	}

	remainingOpponentPiecesMod := map[string]map[int]float64{
		"P": generateRandomIntMap(16, 0.1, 1),
		"R": generateRandomIntMap(16, 0.1, 1),
		"N": generateRandomIntMap(16, 0.1, 1),
		"B": generateRandomIntMap(16, 0.1, 1),
		"Q": generateRandomIntMap(16, 0.1, 1),
		"K": generateRandomIntMap(16, 0.1, 1),
	}

	remainingAlliedPiecesTypeMod := map[string]map[string]float64{
		"P": generateRandomPieceMap(0.2, 1),
		"R": generateRandomPieceMap(0.2, 1),
		"N": generateRandomPieceMap(0.2, 1),
		"B": generateRandomPieceMap(0.2, 1),
		"Q": generateRandomPieceMap(0.2, 1),
		"K": generateRandomPieceMap(0.2, 1),
	}

	remainingOpponentPiecesTypeMod := map[string]map[string]float64{
		"P": generateRandomPieceMap(0.2, 1),
		"R": generateRandomPieceMap(0.2, 1),
		"N": generateRandomPieceMap(0.2, 1),
		"B": generateRandomPieceMap(0.2, 1),
		"Q": generateRandomPieceMap(0.2, 1),
		"K": generateRandomPieceMap(0.2, 1),
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

	config := PieceValueConfig{
		BaseValues:                     baseValues,
		PositionMod:                    positionMod,
		RemainingAlliedPiecesMod:       remainingAlliedPiecesMod,
		RemainingOpponentPiecesMod:     remainingOpponentPiecesMod,
		RemainingAlliedPiecesTypeMod:   remainingAlliedPiecesTypeMod,
		RemainingOpponentPiecesTypeMod: remainingOpponentPiecesTypeMod,
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
		SquareBaseValues: squareBaseValues,
		CoveredByMod:     coveredByMod,
	}

	writeConfig(Policy{HeauristicConfig: config.marshalJson()}, name, dir)
}

type Policy struct {
	HeauristicConfig PieceValueConfigJsonified `json:"heauristicConfig"`
}

func writeConfig(data Policy, fileName string, dir string) {

	fileJson, err := json.Marshal(data)
	if err != nil {
		print("error")
	}
	_ = ioutil.WriteFile(dir+fileName+".json", fileJson, 0644)
}

func readConfigJson(fileName, dir string) Policy {
	file, err := ioutil.ReadFile(dir + fileName)

	if err != nil {
		print(err)
	}

	data := Policy{}

	_ = json.Unmarshal([]byte(file), &data)

	return data
}

func writeRandomConfigs(dir string, number int) {
	for i := 0; i < number; i++ {
		randomConfigGenerator(uuid.NewString(), dir)
	}
}
