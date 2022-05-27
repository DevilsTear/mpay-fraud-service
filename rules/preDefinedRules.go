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

	if isOK, err = payload.checkThreeUniqueCardsAllowed(); err != nil {
		return false, errors.New(fmt.Sprintf("%v check is failed!\nError Details: %v", "checkThreeUniqueCardsAllowed", err))
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
	*cardNumber = strings.Trim(*cardNumber, " ")
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

func (payload *requestPayload) checkThreeUniqueCardsAllowed() (bool, error) {
	cardAllowance := false
	tckn := &payload.Data.User.TCKN
	*tckn = strings.Trim(*tckn, " ")
	if *tckn == "" {
		return false, errors.New("tckn is empty")
	}

	err := config.MySqlDB.Raw(`SELECT COUNT(1) < ? AS cardAllowance FROM request_jetpay_registrations AS rjr 
			INNER JOIN request AS r ON rjr.request_id = r.ID 
			WHERE created_at >= CAST(CURDATE() AS DATETIME) 
			AND created_at <= DATE_SUB(CAST(DATE_ADD(CURDATE(), INTERVAL 1 DAY) AS DATETIME), INTERVAL 1 SECOND) 
			AND user_tckn = ? AND r.Status = 1 GROUP BY rjr.crypted_cc`, 3, *tckn).
		Scan(&cardAllowance)

	if err.Error != nil || !cardAllowance {
		errString := "daily card allowance limit reached!"
		if err.Error != nil {
			errString += fmt.Sprintf("\nError Details: %v", err)
		}

		return false, errors.New(errString)
	}

	return true, nil
}

func (payload *requestPayload) checkFifteenCountClearance() (bool, error) {
	fifteenNeedsClearance := false
	userId := &payload.Data.User.UserId
	clientId := &payload.Data.ClientId
	*userId = strings.Trim(*userId, " ")
	*clientId = strings.Trim(*clientId, " ")
	if *userId == "" || *clientId == "" {
		return false, errors.New("userId and/or clientId is empty")
	}

	err := config.MySqlDB.Raw(`SELECT * FROM cc_fraud WHERE user_id = ? AND client_id = ?`, *userId, *clientId).
		Scan(&fifteenNeedsClearance)

	if err.Error != nil || fifteenNeedsClearance {
		errString := "kredi Kartı harcama güvenliğinizin sağlanması kapsamında, sitemizin çağrı merkezi ile iletişime geçerek işlemlerin sizin tarafınızdan yapıldığını doğrulamanız gerekmektedir!"
		if err.Error != nil {
			errString += fmt.Sprintf("\nError Details: %v", err)
		}

		return false, errors.New(errString)
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
