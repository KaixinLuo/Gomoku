package main

import (
	"fmt"
)

/********************************board structure**********************************/
// the game Board
type Board struct {
	board        []int
	size         int
	metaSize     int
	metaPosition int
}

//a piece of game board
type MetaBoard struct {
	super    *Board
	board    []int
	size     int
	position int
}

/********************************help function**********************************/
func IndexDimIncrease(position int, size int) (int, int) {
	verticalIndex := position / size
	horizontalIndex := position % size
	return verticalIndex, horizontalIndex
}
func IndexDimReduction(dim1 int, dim2 int, metric int) int {
	return dim1*metric + dim2
}

/********************************board function**********************************/
func BoardInit(size, metaSize int) Board {
	b := Board{nil, size, metaSize, 0}
	b.board = make([]int, size*size)
	for index := range b.board {
		b.board[index] = 0
	}
	return b
}
func (board *Board) PosiReshape() (int, int) {
	verticalIndex := board.metaPosition / board.size
	horizontalIndex := board.metaPosition % board.size

	return verticalIndex, horizontalIndex
}
func (board *Board) NumOfMeta() int {
	return (board.size - board.metaSize + 1) * (board.size - board.metaSize + 1)
}
func (board *Board) NextMeta() (bool, MetaBoard) {
	columnIndex, rowIndex := board.PosiReshape()
	var m MetaBoard
	canProceed := true
	if board.size-columnIndex >= board.metaSize {
		m = MetaBoard{board,
			make([]int, board.metaSize*board.metaSize),
			board.metaSize,
			board.metaPosition}
		globalIndex := board.metaPosition
		for i := 0; i < board.metaSize; i++ {
			oldIndex := globalIndex
			for j := 0; j < board.metaSize; j++ {
				m.board[i*board.metaSize+j] = board.board[globalIndex]
				globalIndex++
			}
			globalIndex = oldIndex + board.size
		}
		if board.size-rowIndex >= board.metaSize {
			board.metaPosition++
		} else {
			board.metaPosition += board.metaSize
		}
	} else {
		canProceed = false
	}
	return canProceed, m
}

func (b *Board) ToString() string {
	result := ""
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			if b.board[i*b.size+j] == 1 {
				result += "O "
			} else if b.board[i*b.size+j] == 1 {
				result += "X "
			} else {
				result += "- "
			}
		}
		result += "\n"
	}
	result += "\n"
	return result
}

/********************************meta board function**********************************/
func (b *MetaBoard) EmptyPositions() []int {
	result := make([]int, 0, b.size*b.size)
	nextPosition := 0
	for i := 0; i < len(b.board); i++ {
		if b.board[i] == 0 {
			result[nextPosition] = i
			nextPosition++
		}
	}
	return result
}
func (b *Board) State() (bool, int) {
	hasEmpty := false
	winner := 0
	//check empty
	for _, elem := range b.board {
		if elem == 0 {
			hasEmpty = true
		}
	}

	for i := 0; i < b.size; i++ {
		winner += b.board[i*b.size+i]
	}
	winner /= b.size
	if winner == 0 {
		for i := 0; i < b.size; i++ {
			winner += b.board[i*b.size+(b.size-i-1)]
		}
		winner /= b.size
	} else {
		for i := 0; i < b.size && (winner == 0); i++ {
			for j := 0; j < b.size; j++ {
				winner += b.board[i*b.size+j]
			}
			winner /= b.size
		}

		if winner == 0 {
			for i := 0; i < b.size && (winner == 0); i++ {
				for j := 0; j < b.size; j++ {
					winner += b.board[j*b.size+i]
				}
				winner /= b.size
			}
		}
	}
	return hasEmpty, winner
}
func (this *MetaBoard) Flip() {

}
func (this *MetaBoard) Apply(policy int) {

}
func (this *MetaBoard) Cancel(policy int) {

}
func (this *MetaBoard) Localize(globalPolicy int) int {

}
func (this *MetaBoard) Globalize(localPolicy int) int {

}
func (this *MetaBoard) IsPolicyAvailable(policy int) {

}
func (b *MetaBoard) ToString() string {
	result := ""
	for i := 0; i < b.size; i++ {
		for j := 0; j < b.size; j++ {
			if b.board[i*b.size+j] == 1 {
				result += "O "
			} else if b.board[i*b.size+j] == 1 {
				result += "X "
			} else {
				result += "- "
			}
		}
		result += "\n"
	}
	//result += "\n"
	return result
}

//policy generalization
//MinMax
func main() {
	b := BoardInit(6, 5)
	fmt.Println(b.ToString())
	for i := 0; i < b.NumOfMeta(); i++ {
		_, m := b.NextMeta()
		fmt.Println(m.ToString())
	}
}