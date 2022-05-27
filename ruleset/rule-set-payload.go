package ruleset

import (
	"errors"
	model "fraud-service/model"
	"sort"
	"sync"
)

type ruleSetPayload struct {
	Data []model.RuleSet `json:"data"`
	sync.RWMutex
}

type rulesetList []model.RuleSet

var instance ruleSetPayload

func GetInstance() *ruleSetPayload {
	return &instance
}

func (payload *ruleSetPayload) SetPayload(data []model.RuleSet) {
	payload.Lock()
	defer payload.Unlock()
	payload.Data = filterActiveOnes(data)
}

func (payload *ruleSetPayload) GetPayload() []model.RuleSet {
	payload.RLock()
	defer payload.RUnlock()
	return payload.Data
}

func filterActiveOnes(rules []model.RuleSet) (out []model.RuleSet) {
	for i := range rules {
		if rules[i].Status {
			out = append(out, rules[i])
		}
	}

	return
}

func (payload *ruleSetPayload) SortRuleSetsByPriority() error {
	if payload == nil {
		return errors.New("payload is nil")
	}
	sort.Sort(rulesetList(payload.Data))

	return nil
}

func (ruleset rulesetList) Len() int {
	return len(ruleset)
}

func (ruleset rulesetList) Less(i, j int) bool {
	return ruleset[i].Priority < ruleset[j].Priority
}

func (ruleset rulesetList) Swap(i, j int) {
	ruleset[i], ruleset[j] = ruleset[j], ruleset[i]
}
