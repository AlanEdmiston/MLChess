package main

import (
	"errors"
	"strconv"
	"strings"
)

//Colour is player colour
type Colour string

//WinState is either BlackWon, WhiteWon, Stalemate, or Undecided
type WinState string

const (
	//Black colour
	Black Colour = "Black"
	//White colour
	White Colour = "White"
)

const (
	//BlackWon state
	BlackWon WinState = "Black"
	//WhiteWon state
	WhiteWon WinState = "White"
	//Stalemate state
	Stalemate WinState = "Stalemate"
	//Undecided state
	Undecided WinState = "Undecided"
)

//Board is a representation of a chess board
type Board struct {
	pieces                  []*Piece
	squares                 [][]*Piece
	enPassantRank           int
	canBlackKingSideCastle  bool
	canBlackQueenSideCastle bool
	canWhiteKingSideCastle  bool
	canWhiteQueenSideCastle bool
	colourToMove            Colour
	coveredSquaresWhite     [][]bool
	coveredSquaresBlack     [][]bool
	isWhiteChecked          bool
	isBlackChecked          bool
	winner                  WinState
	children                []*Board
	lastMoveString          string
}

//BoardInitialise initialises a Board
func BoardInitialise(pieces []*Piece, enPassantRank int, colourToMove Colour, canBlackKingSideCastle, canBlackQueenSideCastle, canWhiteKingSideCastle, canWhiteQueenSideCastle bool, lastMoveString string) Board {
	squares := make([][]*Piece, 8)
	for i := 0; i < 8; i++ {
		squares[i] = make([]*Piece, 8)
	}

	for _, piece := range pieces {
		squares[piece.position.X][piece.position.Y] = piece
	}

	coveredSquares := [][]bool{}
	for i := 0; i < 8; i++ {
		coveredSquares = append(coveredSquares, []bool{false, false, false, false, false, false, false, false})
	}

	returnState := Board{
		pieces,
		squares,
		enPassantRank,
		canBlackKingSideCastle,
		canBlackQueenSideCastle,
		canWhiteKingSideCastle,
		canWhiteQueenSideCastle,
		colourToMove,
		coveredSquares,
		coveredSquares,
		false,
		false,
		Undecided,
		[]*Board{},
		lastMoveString,
	}

	returnState.coveredSquaresWhite = returnState.getCoveredSquares(White)
	returnState.coveredSquaresBlack = returnState.getCoveredSquares(Black)
	returnState.isWhiteChecked = returnState.getIsWhiteChecked()
	returnState.isBlackChecked = returnState.getIsBlackChecked()

	isBlackMated := returnState.isBlackCheckmated()
	isWhiteMated := returnState.isWhiteCheckmated()
	isStale := returnState.isStalemate()

	if isBlackMated {
		returnState.winner = WhiteWon
	}
	if isWhiteMated {
		returnState.winner = BlackWon
	}
	if isStale {
		returnState.winner = Stalemate
	}

	returnState.children = returnState.getPossibleMoves()

	return returnState
}

//NewBoard sets up a board in its inital state
func NewBoard() Board {
	pieces := []*Piece{
		{rook, White, Vector{X: 0, Y: 0}},
		{rook, White, Vector{X: 7, Y: 0}},
		{rook, Black, Vector{X: 0, Y: 7}},
		{rook, Black, Vector{X: 7, Y: 7}},
		{knight, White, Vector{X: 1, Y: 0}},
		{knight, White, Vector{X: 6, Y: 0}},
		{knight, Black, Vector{X: 1, Y: 7}},
		{knight, Black, Vector{X: 6, Y: 7}},
		{bishop, White, Vector{X: 2, Y: 0}},
		{bishop, White, Vector{X: 5, Y: 0}},
		{bishop, Black, Vector{X: 2, Y: 7}},
		{bishop, Black, Vector{X: 5, Y: 7}},
		{queen, White, Vector{X: 3, Y: 0}},
		{queen, Black, Vector{X: 3, Y: 7}},
		{king, White, Vector{X: 4, Y: 0}},
		{king, Black, Vector{X: 4, Y: 7}},
	}

	for i := 0; i < 8; i++ {
		pieces = append(pieces, &Piece{pawn, White, Vector{X: i, Y: 1}})
		pieces = append(pieces, &Piece{pawn, Black, Vector{X: i, Y: 6}})
	}

	return BoardInitialise(pieces, -1, White, true, true, true, true, "")
}

func (boardState Board) getCoveredSquares(colour Colour) [][]bool {
	squareMap := [][]bool{}
	for i := 0; i < 8; i++ {
		squareMap = append(squareMap, []bool{false, false, false, false, false, false, false, false})
	}
	for _, piece := range boardState.pieces {
		if piece.colour == colour {
			for _, space := range piece.getCoveredSquares(boardState) {
				if space.X > 7 || space.Y > 7 {
					print("thing")
				}
				squareMap[space.X][space.Y] = true
			}
		}
	}

	return squareMap
}

func (boardState Board) getIsWhiteChecked() bool {
	for _, piece := range boardState.pieces {
		if piece.pieceType.sign == "K" && piece.colour == White {
			return boardState.coveredSquaresBlack[piece.position.X][piece.position.Y]
		}
	}
	return false
}

func (boardState Board) getIsBlackChecked() bool {
	for _, piece := range boardState.pieces {
		if piece.pieceType.sign == "K" && piece.colour == Black {
			return boardState.coveredSquaresWhite[piece.position.X][piece.position.Y]
		}
	}
	return false
}

func (boardState Board) isWhiteCheckmated() bool {
	return boardState.isWhiteChecked && len(boardState.getPossibleMoves()) == 0
}

func (boardState Board) isBlackCheckmated() bool {
	return boardState.isBlackChecked && len(boardState.getPossibleMoves()) == 0
}

func (boardState Board) isStalemate() bool {
	return !boardState.isBlackChecked && !boardState.isWhiteChecked && len(boardState.getPossibleMoves()) == 0
}

//ToString converts board to string
func (boardState Board) ToString() string {
	out := ""
	blankSpaceCounter := 0
	for y, column := range boardState.squares {
		for x := range column {

			if boardState.squares[x][y] != nil {
				if blankSpaceCounter > 0 {
					out += strconv.Itoa(blankSpaceCounter)
					blankSpaceCounter = 0
				}
				if boardState.squares[x][y].colour == White {
					out += strings.ToUpper(boardState.squares[x][y].pieceType.sign)
				} else {
					out += strings.ToLower(boardState.squares[x][y].pieceType.sign)
				}
			} else {
				blankSpaceCounter++
			}
		}
		if blankSpaceCounter > 0 {
			out += strconv.Itoa(blankSpaceCounter)
			blankSpaceCounter = 0
		}
		out += "/"
	}

	if boardState.colourToMove == White {
		out += " w "
	} else {
		out += " b "
	}
	if boardState.canWhiteKingSideCastle {
		out += "K"
	}
	if boardState.canWhiteQueenSideCastle {
		out += "Q"
	}
	if boardState.canBlackKingSideCastle {
		out += "k"
	}
	if boardState.canBlackQueenSideCastle {
		out += "q"
	}

	if !boardState.canWhiteKingSideCastle && !boardState.canWhiteQueenSideCastle && !boardState.canBlackKingSideCastle && !boardState.canBlackQueenSideCastle {
		out += "-"
	}

	return out
}

func (boardState Board) clone() Board {
	pieceClones := []*Piece{}
	squareClones := make([][]*Piece, 8)
	for i := 0; i < 8; i++ {
		squareClones[i] = make([]*Piece, 8)
	}

	for _, piece := range boardState.pieces {
		nextPiece := piece.clone()
		pieceClones = append(pieceClones, &nextPiece)
		squareClones[piece.position.X][piece.position.Y] = &nextPiece
	}

	return Board{
		pieceClones,
		squareClones,
		boardState.enPassantRank,
		boardState.canBlackKingSideCastle,
		boardState.canBlackQueenSideCastle,
		boardState.canWhiteKingSideCastle,
		boardState.canWhiteQueenSideCastle,
		boardState.colourToMove,
		boardState.coveredSquaresWhite,
		boardState.coveredSquaresBlack,
		boardState.isWhiteChecked,
		boardState.isBlackChecked,
		boardState.winner,
		boardState.children,
		boardState.lastMoveString,
	}
}

//MakeMove decides a move and exports a board with that move having been made
func (boardState Board) MakeMove(piece *Piece, move Vector, promotion *PieceType) Board {
	nextState := boardState.clone()
	pieceDouble := nextState.getSquare(piece.position.X, piece.position.Y)
	nextState.lastMoveString = strings.ToUpper(piece.pieceType.sign) + piece.position.boardPosition() + move.boardPosition() + " "

	//remove taken piece
	if nextState.getSquare(move.X, move.Y) != nil {
		for i, piece2 := range nextState.pieces {
			if piece2.position.X == move.X && piece2.position.Y == move.Y {
				nextState.pieces[i] = nextState.pieces[len(nextState.pieces)-1]
				nextState.pieces = nextState.pieces[:len(nextState.pieces)-1]
				break
			}
		}
	}

	//update board en passant rank
	if pieceDouble.pieceType.sign == "P" && piece.position.Y%5 == 1 && (move.Y == 3 || move.Y == 4) {
		nextState.enPassantRank = piece.position.X
	} else {
		nextState.enPassantRank = -1
	}

	//change piece position
	nextState.squares[move.X][move.Y] = pieceDouble
	nextState.squares[piece.position.X][piece.position.Y] = nil
	pieceDouble.position = move

	//remove en passant taken piece
	if pieceDouble.pieceType.sign == "P" && nextState.enPassantRank == move.X {
		dir := 1
		if piece.colour == Black {
			dir = -1
		}
		for i, piece2 := range nextState.pieces {
			if piece2.position.X == move.X && piece2.position.Y == move.Y-dir {
				nextState.pieces[i] = nextState.pieces[len(nextState.pieces)-1]
				nextState.pieces = nextState.pieces[:len(nextState.pieces)-1]
				break
			}
		}
	}

	//do promotions
	if ((move.Y == 0 && piece.colour == Black) || (move.Y == 7 && piece.colour == White)) && piece.pieceType.sign == "P" {
		if promotion == nil {
			println(errors.New("promotion needs to be defined"))
			return nextState
		}
		//promote pawn
		pieceDouble.pieceType = *promotion
		nextState.lastMoveString += "=" + strings.ToUpper(promotion.sign)
	}

	//castling
	if pieceDouble.pieceType.sign == "K" {
		if move.X == 2 {
			if pieceDouble.colour == White && nextState.canWhiteQueenSideCastle {
				nextState.getSquare(0, 0).position.X = 3
			}
			if pieceDouble.colour == Black && nextState.canBlackQueenSideCastle {
				nextState.getSquare(0, 7).position.X = 3
			}
			nextState.lastMoveString = "O-O-O"
		}
		if move.X == 6 {
			if pieceDouble.colour == White && nextState.canWhiteKingSideCastle {
				nextState.getSquare(7, 0).position.X = 5
			}
			if pieceDouble.colour == Black && nextState.canBlackKingSideCastle {
				nextState.getSquare(7, 7).position.X = 5
			}
			nextState.lastMoveString = "O-O"
		}
	}

	//set new covered squares
	nextState.coveredSquaresBlack = nextState.getCoveredSquares(Black)
	nextState.coveredSquaresWhite = nextState.getCoveredSquares(White)

	//is checked
	nextState.isBlackChecked = nextState.getIsBlackChecked()
	nextState.isWhiteChecked = nextState.getIsWhiteChecked()

	//update castling state
	if pieceDouble.pieceType.sign == "K" {
		if pieceDouble.colour == White {
			nextState.canWhiteKingSideCastle = false
			nextState.canWhiteQueenSideCastle = false
		} else {
			nextState.canBlackKingSideCastle = false
			nextState.canBlackQueenSideCastle = false
		}
	}

	if pieceDouble.pieceType.sign == "R" {
		if pieceDouble.colour == White {
			if piece.position.X == 7 && pieceDouble.position.Y == 0 {
				nextState.canWhiteKingSideCastle = false
			}
			if piece.position.X == 0 && pieceDouble.position.Y == 0 {
				nextState.canWhiteQueenSideCastle = false
			}
		} else {
			if piece.position.X == 7 && pieceDouble.position.Y == 7 {
				nextState.canBlackKingSideCastle = false
			}
			if piece.position.X == 0 && pieceDouble.position.Y == 7 {
				nextState.canBlackQueenSideCastle = false
			}
		}
	}

	if move.X == 0 && move.Y == 0 {
		nextState.canWhiteQueenSideCastle = false
	}
	if move.X == 7 && move.Y == 0 {
		nextState.canWhiteKingSideCastle = false
	}
	if move.X == 0 && move.Y == 7 {
		nextState.canBlackQueenSideCastle = false
	}
	if move.X == 7 && move.Y == 7 {
		nextState.canBlackKingSideCastle = false
	}

	if nextState.colourToMove == White {
		nextState.colourToMove = Black
	} else {
		nextState.colourToMove = White
	}

	return nextState
}

func (boardState Board) getPossibleMoves() []*Board {
	returnStates := []*Board{}
	for _, piece := range boardState.pieces {
		if piece.colour == boardState.colourToMove {
			moves := piece.getPossibleMoves(boardState)
			for _, move := range moves {
				if piece.pieceType.sign == "P" && (move.Y == 0 || move.Y == 7) {
					nextState := boardState.MakeMove(piece, move, &knight)
					if nextState.verifyBoardState() {
						returnStates = append(returnStates, &nextState)
					}
					nextState = boardState.MakeMove(piece, move, &bishop)
					if nextState.verifyBoardState() {
						returnStates = append(returnStates, &nextState)
					}
					nextState = boardState.MakeMove(piece, move, &rook)
					if nextState.verifyBoardState() {
						returnStates = append(returnStates, &nextState)
					}
					nextState = boardState.MakeMove(piece, move, &queen)
					if nextState.verifyBoardState() {
						returnStates = append(returnStates, &nextState)
					}
				} else {
					nextState := boardState.MakeMove(piece, move, nil)
					if nextState.verifyBoardState() {
						returnStates = append(returnStates, &nextState)
					}
				}
			}
		}
	}
	return returnStates
}

func (boardState Board) verifyBoardState() bool {
	return !(boardState.colourToMove == White && boardState.isBlackChecked) && !(boardState.colourToMove == Black && boardState.isWhiteChecked)
}

func (boardState Board) getSquare(x int, y int) *Piece {
	return boardState.squares[x][y]
}

func (boardState Board) noWhitePieces() int {
	count := 0
	for _, piece := range boardState.pieces {
		if piece.colour == White {
			count++
		}
	}
	return count
}

func (boardState Board) noBlackPieces() int {
	count := 0
	for _, piece := range boardState.pieces {
		if piece.colour == Black {
			count++
		}
	}
	return count
}
