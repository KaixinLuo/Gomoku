package gmk

import (
	"fmt"
)

type LigitPlayer interface {
	MakeDecision() int
	Play()
}

type Human struct {
	game *Board
	Flag int
}

func (this *Human) MakeDecision() int {
	var policyCol, policyRow int
	fmt.Print("Please enter your decision:")
	fmt.Scanf("%d %d", &policyCol, &policyRow)
	policy := policyCol*this.game.size + policyRow
	for !IsPolicyLigit(*this.game, policy) {
		fmt.Print("This position is not available, try again:")
		fmt.Scanf("%d %d", &policyCol, &policyRow)
		policy := policyCol*this.game.size + policyRow
	}
	return policy
}

func
