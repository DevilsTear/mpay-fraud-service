package ruleset

import (
	"errors"
	"fmt"
	"fraud-service/config"
	model "fraud-service/model"
	"sort"
	"strconv"
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
			if val, err := strconv.ParseInt(ruleSet.Value.(string), 10, 64); err == nil {
				config.PendingCountThreshold = val
			}
		case "PendingAllowanceByTimeInterval":
			if val, err := strconv.ParseInt(ruleSet.Value.(string), 10, 64); err == nil {
				config.PendingAllowanceByTimeInterval = val
			}
		case "ApprovedAllowanceByTimeInterval":
			if val, err := strconv.ParseInt(ruleSet.Value.(string), 10, 64); err == nil {
				config.ApprovedAllowanceByTimeInterval = val
			}
		case "MaxDailyAllowancePerUser":
			if val, err := strconv.ParseInt(ruleSet.Value.(string), 10, 64); err == nil {
				config.MaxDailyAllowancePerUser = val
			}
		case "MinTransactionAmount":
			if val, err := strconv.ParseFloat(ruleSet.Value.(string), 64); err == nil {
				config.MinTransactionAmount = val
			}
		case "MaxTransactionAmount":
			if val, err := strconv.ParseFloat(ruleSet.Value.(string), 64); err == nil {
				config.MaxTransactionAmount = val
			}
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
