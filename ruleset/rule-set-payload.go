package ruleset

import (
	"errors"
	"fmt"
	"fraud-service/config"
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
	payload.SetGlobalParams()

	return
}

func (payload *ruleSetPayload) SetGlobalParams() {
	for _, ruleSet := range payload.Data {
		switch ruleSet.Key {
		case "PendingCountThreshold":
			config.PendingCountThreshold = ruleSet.Value.(int64)
		case "PendingAllowanceByTimeInterval":
			config.PendingAllowanceByTimeInterval = ruleSet.Value.(int64)
		case "ApprovedAllowanceByTimeInterval":
			config.ApprovedAllowanceByTimeInterval = ruleSet.Value.(int64)
		case "MaxDailyAllowancePerUser":
			config.MaxDailyAllowancePerUser = ruleSet.Value.(int64)
		case "MinTransactionAmount":
			config.MinTransactionAmount = ruleSet.Value.(float64)
		case "MaxTransactionAmount":
			config.MaxTransactionAmount = ruleSet.Value.(float64)
		}
	}
}

// GetPayload gets ruleset slice from the instance
func (payload *ruleSetPayload) GetPayload() (out []model.RuleSet) {
	payload.RLock()
	defer payload.RUnlock()
	return payload.Data
}

// GetPayloadKeyMapping gets ruleset mappings by Key
func (payload *ruleSetPayload) GetPayloadKeyMapping() (out map[string]model.RuleSet) {
	for _, v := range payload.GetPayload() {
		out[v.Key] = v
	}
	return
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
