package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type judge struct {
	*baseclient.BaseJudge
	c *client
}

func (c *client) GetClientJudgePointer() roles.Judge {
	c.Logf("Team 5 became Judge.")
	return &judge{c: c, BaseJudge: &baseclient.BaseJudge{GameState: c.ServerReadHandle.GetGameState()}}
}

// Pardon ourselves and homies
func (j *judge) GetPardonedIslands(currentSanctions map[int][]shared.Sanction) map[int][]bool {
	pardons := make(map[int][]bool)
	for key, sanctionList := range currentSanctions {
		lst := make([]bool, len(sanctionList))
		pardons[key] = lst
		for index, sanction := range sanctionList {
			if j.c.opinions[sanction.ClientID].getScore() > 0.7 {
				pardons[key][index] = true
			} else {
				pardons[key][index] = false
			}
			if sanction.ClientID == shared.Team5 {
				pardons[key][index] = true
			}
		}
	}
	j.c.Logf("TEAM5 pardoned our homies")
	return pardons
}

// Pay president based on the status of our own wealth
// If we are not doing verywell, pay President less so we have more in the CP to take from
func (j *judge) PayPresident() (shared.Resources, bool) {
	PresidentSalaryRule, ok := j.GameState.RulesInfo.CurrentRulesInPlay["salary_cycle_president"]
	var salary shared.Resources = 0
	if ok {
		salary = shared.Resources(PresidentSalaryRule.ApplicableMatrix.At(0, 1))
	}
	if j.c.wealth() == jeffBezos {
		return salary, true
	} else if j.c.wealth() == middleClass {
		salary = salary * 0.8
	} else {
		salary = salary * 0.5
	}
	return salary, true
}

// if the real winner is on our bad side, then we choose our best friend
func (j *judge) DecideNextPresident(winner shared.ClientID) shared.ClientID {
	aliveTeams := j.c.getAliveTeams(false) //not including us
	if j.c.opinions[winner].getScore() < 0 {
		ballot := j.c.VoteForElection(shared.President, aliveTeams)
		winner = ballot[0] //choose the first one in Borda Vote
	}
	return winner
}
