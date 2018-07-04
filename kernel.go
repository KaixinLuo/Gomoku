package gmk

import (
	"math"
)

const (
	FIRSTHAND = 1
	LASTHAND  = -1
	TIE       = 0
	SELF      = 1
	ENEMY     = -1
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
func IsPolicyLigit(game Board, policy int) bool {
	result := policy < len(game.board)
	result = result && game.IsPolicyAvailable(policy)
	return result
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
func (this *Board) Apply(policy int, flag int) {
	this.board[policy] = flag
}
func (this *Board) Cancel(policy int) {
	this.board[policy] = 0
}
func (this *Board) IsPlaybale() (bool, int) {
	canContinue := true
	winner := 0
	for i := 0; i < this.NumOfMeta() && canContinue; i++ {
		_, m := this.NextMeta()
		canContinue, winner = m.State()
	}
	return canContinue, winner
}
func (this *Board) IsEmpty() bool {
	result := true
	for i := 0; i < len(this.board); i++ {
		result = (this.board[i] == 0)
	}
	return result
}
func (this *Board) Globalize(localPolicy int) int {
	vertical, horizontal := IndexDimIncrease(localPolicy, this.metaSize)
	metaVertical, metaHorizontal := IndexDimIncrease(this.metaPosition, this.size)

	return (vertical+metaVertical)*this.size + (horizontal + metaHorizontal)
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
func (this *Board) IsPolicyAvailable(policy int) bool {
	return this.board[policy] == 0
}

/********************************meta board function**********************************/
func (b *MetaBoard) EmptyPositions() []int {
	result := make([]int, 0, b.size*b.size)
	nextPosition := 0
	for i := 0; i < len(b.board); i++ {
		if b.board[i] == 0 {
			result = append(result, i)
			nextPosition++
		}
	}
	return result
}
func (b *MetaBoard) State() (bool, int) {
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
func (this *MetaBoard) Flip(flag int) {
	if flag == -1 {
		for i := 0; i < len(this.board); i++ {
			this.board[i] /= flag
		}
	}

}
func (this *MetaBoard) Apply(policy int, flag int) {
	this.board[policy] = flag
}
func (this *MetaBoard) Cancel(policy int) {
	this.board[policy] = 0
}
func (this *MetaBoard) Localize(globalPolicy int) int {
	vertical, horizontal := IndexDimIncrease(globalPolicy, this.super.size)
	metaVertical, metaHorizontal := IndexDimIncrease(this.super.metaPosition, this.super.size)

	return (vertical-metaVertical)*this.size + (horizontal - metaHorizontal)

}

func (this *MetaBoard) IsPolicyAvailable(policy int) bool {
	return this.board[policy] == 0
}

//has to flip before call this method
func (this *MetaBoard) BestLocalPolicyAndUtil(flag int, decayRate float64) (int, float64) {

	maxUtil := math.Inf(-1)
	var bestPolicy int
	var util float64
	hasEmpty, winner := this.State()

	if hasEmpty && winner == 0 {
		policies := this.EmptyPositions()
		var p int
		var u float64
		for _, policy := range policies {
			this.Apply(policy, flag)
			p, u = this.WorstLocalPolicyAndUtil(-flag, decayRate)
			this.Cancel(policy)
		}
		if maxUtil < u {
			bestPolicy = p
			maxUtil = u
			util = decayRate * u
		}
	} else {
		util = float64(winner * 10)
	}

	return bestPolicy, util

}
func (this *MetaBoard) WorstLocalPolicyAndUtil(flag int, decayRate float64) (int, float64) {

	minUtil := math.Inf(1)
	var worstPolicy int
	var util float64
	hasEmpty, winner := this.State()

	if hasEmpty && winner == 0 {
		policies := this.EmptyPositions()
		var p int
		var u float64
		for _, policy := range policies {
			this.Apply(policy, flag)
			p, u = this.BestLocalPolicyAndUtil(-flag, decayRate)
			this.Cancel(policy)
		}
		if minUtil > u {
			worstPolicy = p
			minUtil = u
			util = decayRate * u
		}
	} else {
		util = float64(winner * 10)
	}
	return worstPolicy, util

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
