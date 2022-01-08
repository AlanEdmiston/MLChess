package main

import "fmt"

//Piece is a chess piece on a board
type Piece struct {
	pieceType PieceType
	colour    Colour
	position  Vector
}

//PieceType is a type of chess piece (king, queen, bishop etc.)
type PieceType struct {
	sign           string
	moveDirections []Vector
	otherMoves     []Vector
}

func (piece Piece) toString() string {
	return fmt.Sprintf("(%s, %s, %s)", piece.pieceType.sign, piece.colour[:1], piece.position.toString())
}

func (piece Piece) canPawnMoveTwo() bool {
	if (piece.colour == White && piece.position.Y == 1 && piece.pieceType.sign == "P") || (piece.colour == Black && piece.position.Y == 6 && piece.pieceType.sign == "P") {
		return true
	}
	return false
}

func (piece Piece) getCastlingMoves(boardState Board) []Vector {
	castlingMoves := []Vector{}
	if piece.colour == White {
		if boardState.canWhiteQueenSideCastle && boardState.getSquare(0, 1) == nil && boardState.getSquare(0, 2) == nil && boardState.getSquare(0, 3) == nil {
			castlingMoves = append(castlingMoves, Vector{0, 2})
		}
		if boardState.canWhiteKingSideCastle && boardState.getSquare(0, 5) == nil && boardState.getSquare(0, 6) == nil {
			castlingMoves = append(castlingMoves, Vector{0, 6})
		}
	} else {
		if boardState.canBlackQueenSideCastle && boardState.getSquare(7, 1) == nil && boardState.getSquare(7, 2) == nil && boardState.getSquare(7, 3) == nil {
			castlingMoves = append(castlingMoves, Vector{7, 2})
		}
		if boardState.canBlackKingSideCastle && boardState.getSquare(7, 5) == nil && boardState.getSquare(7, 6) == nil {
			castlingMoves = append(castlingMoves, Vector{7, 6})
		}
	}

	return castlingMoves
}

func (piece Piece) getDirectionMoves(boardState Board) []Vector {
	foundMoves := []Vector{}
	for _, direction := range piece.pieceType.moveDirections {
		for i := 1; i < 8; i++ {
			moveTo := piece.position.add(direction.mult(i))
			if moveTo.isOutOfBounds() {
				break
			}
			if occupyingPiece := boardState.getSquare(moveTo.X, moveTo.Y); occupyingPiece != nil {
				if occupyingPiece.colour != piece.colour {
					foundMoves = append(foundMoves, Vector{moveTo.X, moveTo.Y})
				}
				break
			}
			foundMoves = append(foundMoves, Vector{moveTo.X, moveTo.Y})
		}

		for i := -1; i > -8; i-- {
			moveTo := piece.position.add(direction.mult(i))
			if moveTo.isOutOfBounds() {
				break
			}
			if occupyingPiece := boardState.getSquare(moveTo.X, moveTo.Y); occupyingPiece != nil {
				if occupyingPiece.colour != piece.colour {
					foundMoves = append(foundMoves, Vector{moveTo.X, moveTo.Y})
				}
				break
			}
			foundMoves = append(foundMoves, Vector{moveTo.X, moveTo.Y})
		}
	}

	return foundMoves
}

func (piece Piece) getPawnMoves(boardState Board) []Vector {
	foundMoves := []Vector{}

	colourMult := 1
	if piece.colour == Black {
		colourMult = -1
	}

	if piece.position.Y+colourMult > 0 && piece.position.Y+colourMult < 8 && boardState.getSquare(piece.position.X, piece.position.Y+colourMult) == nil {
		foundMoves = append(foundMoves, Vector{piece.position.X, piece.position.Y + colourMult})
	}
	if piece.position.X != 7 && piece.position.Y+colourMult >= 0 && piece.position.Y+colourMult < 8 {
		if occupyingPiece := boardState.getSquare(piece.position.X+1, piece.position.Y+colourMult); occupyingPiece != nil && occupyingPiece.colour != piece.colour {
			foundMoves = append(foundMoves, Vector{piece.position.X + 1, piece.position.Y + colourMult})
		}

		if occupyingPiece := boardState.getSquare(piece.position.X+1, piece.position.Y); occupyingPiece != nil && occupyingPiece.pieceType.sign == "P" && occupyingPiece.position.X == boardState.enPassantRank {
			foundMoves = append(foundMoves, Vector{piece.position.X + 1, piece.position.Y + colourMult})
		}
	}
	if piece.position.X != 0 && piece.position.Y+colourMult >= 0 && piece.position.Y+colourMult < 8 {
		if occupyingPiece := boardState.getSquare(piece.position.X-1, piece.position.Y+colourMult); occupyingPiece != nil && occupyingPiece.colour != piece.colour {
			foundMoves = append(foundMoves, Vector{piece.position.X - 1, piece.position.Y + colourMult})
		}

		if occupyingPiece := boardState.getSquare(piece.position.X-1, piece.position.Y); occupyingPiece != nil && occupyingPiece.pieceType.sign == "P" && occupyingPiece.position.X == boardState.enPassantRank {
			foundMoves = append(foundMoves, Vector{piece.position.X - 1, piece.position.Y + colourMult})
		}
	}
	if piece.canPawnMoveTwo() && boardState.getSquare(piece.position.X, piece.position.Y+colourMult) == nil && boardState.getSquare(piece.position.X, piece.position.Y+2*colourMult) == nil {
		foundMoves = append(foundMoves, Vector{piece.position.X, piece.position.Y + 2*colourMult})
	}

	return foundMoves
}

func (piece Piece) getPossibleMoves(boardState Board) []Vector {
	foundMoves := piece.getDirectionMoves(boardState)

	for _, move := range piece.pieceType.otherMoves {
		moveTo := piece.position.add(move)
		if !moveTo.isOutOfBounds() {
			if occupyingPiece := boardState.getSquare(moveTo.X, moveTo.Y); occupyingPiece == nil || occupyingPiece.colour != piece.colour {
				foundMoves = append(foundMoves, Vector{moveTo.X, moveTo.Y})
			}
		}
	}

	if piece.pieceType.sign == "K" {
		foundMoves = append(foundMoves, piece.getCastlingMoves(boardState)...)
	}

	if piece.pieceType.sign == "P" {
		foundMoves = append(foundMoves, piece.getPawnMoves(boardState)...)
	}
	return foundMoves
}

func (piece Piece) getCoveredSquares(boardState Board) []Vector {
	coveredSquares := []Vector{}

	for _, direction := range piece.pieceType.moveDirections {
		for i := 1; i < 8; i++ {
			moveTo := piece.position.add(direction.mult(i))
			if moveTo.isOutOfBounds() {
				break
			}
			if boardState.getSquare(moveTo.X, moveTo.Y) != nil {
				coveredSquares = append(coveredSquares, Vector{moveTo.X, moveTo.Y})
				break
			}
			coveredSquares = append(coveredSquares, Vector{moveTo.X, moveTo.Y})
		}

		for i := -1; i > -8; i-- {
			moveTo := piece.position.add(direction.mult(i))
			if moveTo.isOutOfBounds() {
				break
			}
			if boardState.getSquare(moveTo.X, moveTo.Y) != nil {
				coveredSquares = append(coveredSquares, Vector{moveTo.X, moveTo.Y})
				break
			}
			coveredSquares = append(coveredSquares, Vector{moveTo.X, moveTo.Y})
		}
	}

	if piece.pieceType.sign == "P" {
		if piece.position.Y == 7 || piece.position.Y == 0 {
			return []Vector{}
		}
		if piece.colour == White {
			if piece.position.X != 0 {
				coveredSquares = append(coveredSquares, Vector{piece.position.X - 1, piece.position.Y + 1})
			}
			if piece.position.X != 7 {
				coveredSquares = append(coveredSquares, Vector{piece.position.X + 1, piece.position.Y + 1})
			}
		} else {
			if piece.position.X != 0 {
				coveredSquares = append(coveredSquares, Vector{piece.position.X - 1, piece.position.Y - 1})
			}
			if piece.position.X != 7 {
				coveredSquares = append(coveredSquares, Vector{piece.position.X + 1, piece.position.Y - 1})
			}
		}
	} else {
		for _, move := range piece.pieceType.otherMoves {
			moveTo := piece.position.add(move)
			if !moveTo.isOutOfBounds() {
				coveredSquares = append(coveredSquares, Vector{moveTo.X, moveTo.Y})
			}
		}
	}

	return coveredSquares
}

func (piece Piece) isProtecting(otherPiece *Piece, boardState Board) bool {
	if piece.colour != otherPiece.colour {
		return false
	}
	for _, move := range piece.getPossibleMoves(boardState) {
		if move.X == otherPiece.position.X && move.Y == otherPiece.position.Y {
			return true
		}
	}
	return false
}

func (piece Piece) isAttacking(otherPiece *Piece, boardState Board) bool {
	if piece.colour == otherPiece.colour {
		return false
	}
	for _, move := range piece.getPossibleMoves(boardState) {
		if move.X == otherPiece.position.X && move.Y == otherPiece.position.Y {
			return true
		}
	}
	return false
}

// func (piece Piece) isAttacking(otherPiece *Piece, boardState Board) bool {
// 	if piece.colour == otherPiece.colour {
// 		return false
// 	}
// 	for _, direction := range piece.pieceType.moveDirections {
// 		for i := 1; i < 8; i++ {
// 			moveTo := piece.position.add(direction.mult(i))
// 			if moveTo.isOutOfBounds() {
// 				break
// 			}
// 			if occupyingPiece := boardState.getSquare(moveTo.x, moveTo.y); occupyingPiece != nil {
// 				if occupyingPiece.position.x == otherPiece.position.x && occupyingPiece.position.y == otherPiece.position.y {
// 					return true
// 				}
// 				break
// 			}
// 		}

// 		for i := -1; i > -8; i-- {
// 			moveTo := piece.position.add(direction.mult(i))
// 			if moveTo.isOutOfBounds() {
// 				break
// 			}
// 			if occupyingPiece := boardState.getSquare(moveTo.x, moveTo.y); occupyingPiece != nil {
// 				if occupyingPiece.position.x == otherPiece.position.x && occupyingPiece.position.y == otherPiece.position.y {
// 					return true
// 				}
// 				break
// 			}
// 		}
// 	}

// 	for _, move := range piece.pieceType.otherMoves {
// 		moveTo := piece.position.add(move)
// 		if !moveTo.isOutOfBounds() {
// 			if occupyingPiece := boardState.getSquare(moveTo.x, moveTo.y); occupyingPiece != nil &&
// 				occupyingPiece.position.x == otherPiece.position.x && occupyingPiece.position.y == otherPiece.position.y {
// 				return true
// 			}
// 		}
// 	}

// 	if piece.pieceType.sign == "P" {
// 		colourMult := 1
// 		if piece.colour == Black {
// 			colourMult = -1
// 		}

// 		if piece.position.x != 7 && piece.position.y+colourMult >= 0 && piece.position.y+colourMult < 8 {
// 			if occupyingPiece := boardState.getSquare(piece.position.x+1, piece.position.y+colourMult); occupyingPiece.position.x == otherPiece.position.x && occupyingPiece.position.y == otherPiece.position.y {
// 				return true
// 			}

// 			if occupyingPiece := boardState.getSquare(piece.position.x+1, piece.position.y); occupyingPiece != nil &&
// 				occupyingPiece.pieceType.sign == "P" && occupyingPiece.position.x == boardState.enPassantRank &&
// 				occupyingPiece.position.x == otherPiece.position.x && occupyingPiece.position.y == otherPiece.position.y {
// 				return true
// 			}
// 		}
// 		if piece.position.x != 0 && piece.position.y+colourMult >= 0 && piece.position.y+colourMult < 8 {
// 			if occupyingPiece := boardState.getSquare(piece.position.x-1, piece.position.y+colourMult); occupyingPiece != nil &&
// 				occupyingPiece.position.x == otherPiece.position.x && occupyingPiece.position.y == otherPiece.position.y {
// 				return true
// 			}

// 			if occupyingPiece := boardState.getSquare(piece.position.x-1, piece.position.y); occupyingPiece != nil &&
// 				occupyingPiece.pieceType.sign == "P" && occupyingPiece.position.x == boardState.enPassantRank &&
// 				occupyingPiece.position.x == otherPiece.position.x && occupyingPiece.position.y == otherPiece.position.y {
// 				return true
// 			}
// 		}
// 	}
// 	return false
// }

func (piece Piece) clone() Piece {
	return Piece{piece.pieceType, piece.colour, piece.position}
}

func (piece Piece) canTake(otherPiece Piece, boardState *Board) bool {
	moves := piece.getPossibleMoves(*boardState)
	for _, move := range moves {
		if move.X == otherPiece.position.X && move.Y == otherPiece.position.Y {
			return true
		}
	}
	return false
}

var pawn = PieceType{"P", []Vector{}, []Vector{}}

var rook = PieceType{"R", []Vector{{0, 1}, {1, 0}}, []Vector{}}

var knight = PieceType{"N", []Vector{}, []Vector{{2, 1}, {-2, 1}, {2, -1}, {-2, -1}, {1, 2}, {-1, 2}, {1, -2}, {-1, -2}}}

var bishop = PieceType{"B", []Vector{{-1, 1}, {1, 1}}, []Vector{}}

var queen = PieceType{"Q", []Vector{{0, 1}, {1, 0}, {-1, 1}, {1, 1}}, []Vector{}}

var king = PieceType{"K", []Vector{}, []Vector{{1, -1}, {1, 0}, {1, 1}, {0, -1}, {0, 1}, {-1, -1}, {-1, 0}, {-1, 1}}}
