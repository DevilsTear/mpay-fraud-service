package utils

import (
	"fmt"
	"fraud-service/model"
	"sort"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

type rulesetList []model.RuleSet

func (ruleset rulesetList) Len() int {
	return len(ruleset)
}

func (ruleset rulesetList) Less(i, j int) bool {
	return ruleset[i].Priority > ruleset[j].Priority
}

func (ruleset rulesetList) Swap(i, j int) {
	ruleset[i], ruleset[j] = ruleset[j], ruleset[i]
}

func SortRuleSetsByPriority(payload *model.RuleSetPayload) {
	sort.Sort(rulesetList(payload.Data))
}
