package rules

import (
	"errors"
	"fmt"
	"fraud-service/config"
	"fraud-service/model"
	rulesets "fraud-service/ruleset"
	"fraud-service/utils"
	"strings"
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

	if isOK, err = payload.checkCardBIN(); err != nil {
		return false, errors.New(fmt.Sprintf("%v check is failed!\nError Details: %v", "checkCardBIN", err))
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

func (payload *requestPayload) checkCardBIN() (bool, error) {
	binExists := false
	cardNumber := &payload.Data.Transaction.CardNumber
	*cardNumber = strings.Trim(payload.Data.Transaction.CardNumber, " ")
	if *cardNumber == "" {
		return binExists, errors.New("card number is empty")
	}

	cardBin := (*cardNumber)[:6]
	err := config.MySqlDB.Raw("SELECT 1 as binExists FROM cc_binlist WHERE card_bin = ? LIMIT 1", cardBin).
		Scan(&binExists)

	if err.Error != nil || !binExists {
		errString := "card issuer is not listed in the bin list!"
		if err.Error != nil {
			errString += fmt.Sprintf("\nError Details: %v", err)
		}

		return binExists, errors.New(errString)
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
