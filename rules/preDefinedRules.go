package rules

import (
	"errors"
	"fmt"
	"fraud-service/model"
	rulesets "fraud-service/ruleset"
	"fraud-service/utils"
	"sync"
)

type requestPayload struct {
	Data model.RequestPayload `json:"data"`
	sync.RWMutex
}

var instance requestPayload

func GetRequestPayloadInstance() *requestPayload {
	return &instance
}

func (payload *requestPayload) SetPayload(data model.RequestPayload) {
	payload.Lock()
	defer payload.Unlock()
	payload.Data = data
}

func (payload *requestPayload) GetPayload() model.RequestPayload {
	payload.RLock()
	defer payload.RUnlock()
	return payload.Data
}

func (payload *requestPayload) ProcessRules() (bool, error) {
	ruleSets := rulesets.GetInstance().GetPayload()
	isOK := false
	var err error = nil
	if !anyRuleExists(ruleSets) {
		return false, errors.New("please, define your rule sets first")
	}

	for _, ruleSet := range ruleSets {
		switch ruleSet.Key {
		case "PendingCountThreshold":
			isOK, err = payload.checkPendingCountThreshold()
		case "PendingAllowanceByTimeInterval":
			isOK, err = payload.checkPendingAllowanceByTimeInterval()
		case "ApprovedAllowanceByTimeInterval":
			isOK, err = payload.checkApprovedAllowanceByTimeInterval()
		case "MaxDailyAllowancePerUser":
			isOK, err = payload.checkMaxDailyAllowancePerUser()
		case "MinTransactionAmount":
			isOK, err = payload.checkMinTransactionAmount()
		case "MaxTransactionAmount":
			isOK, err = payload.checkMaxTransactionAmount()
			utils.CheckError(err)
		}

		if !isOK || err != nil {
			return false, errors.New(fmt.Sprintf("%v check is failed!\nError Details: %v", ruleSet.Key, err))
		}
	}
	return true, nil
}

func (payload *requestPayload) checkPendingCountThreshold() (bool, error) {

	return true, nil
}

func (payload *requestPayload) checkPendingAllowanceByTimeInterval() (bool, error) {

	return true, nil
}

func (payload *requestPayload) checkApprovedAllowanceByTimeInterval() (bool, error) {

	return true, nil
}

func (payload *requestPayload) checkMaxDailyAllowancePerUser() (bool, error) {

	return true, nil
}

func (payload *requestPayload) checkMinTransactionAmount() (bool, error) {

	return true, nil
}

func (payload *requestPayload) checkMaxTransactionAmount() (bool, error) {

	return true, nil
}

func anyRuleExists(ruleSets []model.RuleSet) bool {
	preDefinedRuleKeys := []string{
		"PendingCountThreshold",
		"PendingAllowanceByTimeInterval",
		"ApprovedAllowanceByTimeInterval",
		"MaxDailyAllowancePerUser",
		"MinTransactionAmount",
		"MaxTransactionAmount",
	}
	var ruleKeys []string
	for _, ruleSet := range ruleSets {
		ruleKeys = append(ruleKeys, ruleSet.Key)
	}

	return len(utils.Intersection(ruleKeys, preDefinedRuleKeys)) > 0
}
