package gmk

import "math"

type Bot struct {
	game   *Board
	Flag   int
	temper float64 //(0,1] smart, dumbass otherwise
}

func BotInit(game *Board, flag int, temper float64) Bot {
	b := Bot{game, flag, temper}
	return b
}
func (this *Bot) Play() {
	p := this.MakeDecision()
	this.game.Apply(p, this.Flag)
}
func (this *Bot) MakeDecision() int {
	MaxUtil := math.Inf(-1)
	var localPolicy int
	for i := 0; i < this.game.NumOfMeta(); i++ {
		_, m := this.game.NextMeta()
		m.Flip(this.Flag)
		p, u := BestLocalPolicyAndUtil(m, SELF, this.temper)
		if MaxUtil < u {
			MaxUtil = u
			localPolicy = p
		}
	}
	return this.game.Globalize(localPolicy)
}

/*****************MinMax******************/
func BestLocalPolicyAndUtil(glance MetaBoard, flag int, decayRate float64) (int, float64) {

	maxUtil := math.Inf(-1)
	var bestPolicy int
	var util float64
	hasEmpty, winner := glance.State()

	if hasEmpty && winner == 0 {
		policies := glance.EmptyPositions()
		var p int
		var u float64
		for _, policy := range policies {
			glance.Apply(policy, flag)
			p, u = WorstLocalPolicyAndUtil(glance, -flag, decayRate)
			glance.Cancel(policy)
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
func WorstLocalPolicyAndUtil(glance MetaBoard, flag int, decayRate float64) (int, float64) {

	minUtil := math.Inf(1)
	var worstPolicy int
	var util float64
	hasEmpty, winner := glance.State()

	if hasEmpty && winner == 0 {
		policies := glance.EmptyPositions()
		var p int
		var u float64
		for _, policy := range policies {
			glance.Apply(policy, flag)
			p, u = BestLocalPolicyAndUtil(glance, -flag, decayRate)
			glance.Cancel(policy)
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
