package ruleset

import (
	"errors"
	"fmt"
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

// GetInstance constructs ruleset payload instance
func GetInstance() *ruleSetPayload {
	return &instance
}

// SetPayload sets ruleset slice to the instance
func (payload *ruleSetPayload) SetPayload(data []model.RuleSet) (err error) {
	payload.Lock()
	defer payload.Unlock()
	payload.Data, err = filterActiveOnes(data)

	return
}

// GetPayload gets ruleset slice from the instance
func (payload *ruleSetPayload) GetPayload() (out []model.RuleSet) {
	payload.RLock()
	defer payload.RUnlock()
	return payload.Data
}

func filterActiveOnes(rules []model.RuleSet) (out []model.RuleSet, err error) {
	if len(rules) <= 0 {
		err = fmt.Errorf("rules object empty")
	}
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
