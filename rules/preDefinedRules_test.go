package rules

import (
	"context"
	"encoding/json"
	"fraud-service/config"
	"fraud-service/model"
	rulesets "fraud-service/ruleset"
	"log"
	"os"
	"path"
	"runtime"
	"testing"
)

var jsonRuleSetPayload = `{
    "data": [
        {
            "name": "Rule 1",
            "key": "PendingCountThreshold",
            "priority": 1,
            "value": "10",
            "status": true,
            "operator": "gt"
        },{
            "name": "Rule 2",
            "key": "Key ",
            "priority": 4,
            "value": "10",
            "status": true,
            "operator": "gt"
        },{
            "name": "Rule 3",
            "key": "Key ",
            "priority": 3,
            "value": "10",
            "status": false,
            "operator": "gt"
        },{
            "name": "Rule 4",
            "key": "Key ",
            "priority": 2,
            "value": "10",
            "status": true,
            "operator": "gt"
        }
    ]
}`

var jsonRequestPayload = `{
    "client": {
        "id": "1",
        "cc_user_perm_check": true,
        "fullname_cc_match": true
    },
    "data": {
        "amount": "250.00",
        "trx": "oaisufklafasfl1111112d1233",
        "card_number": "4943141412612271",
        "expiration_month": "10",
        "expiration_year": "2022",
        "cardholders_name": "Serdar ürgün",
        "cvv": "402",
        "return_url": "https://envoysoft3.net/deposit/mpayReturn"
    },
    "user": {
        "username": "srdr16",
        "userID": "17206739184",
        "yearofbirth": "1983",
        "fullname": "Serdar ürgün",
        "email": "mawiay16@hotmail.com",
        "tckn": "17206739184",
        "ip_address": "95.10.24.238"
    }
}`

var activeRulesInstance = rulesets.GetInstance()

func init() {
	_, filename, _, _ := runtime.Caller(0)
	// The ".." may change depending on you folder structure
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func init() {
	if err := config.LoadInitials(context.TODO()); err != nil {
		panic(err)
	}
}

func init() {
	var payload model.RuleSetPayload
	if err := json.Unmarshal([]byte(jsonRuleSetPayload), &payload); err != nil {
		log.Println(err)
		return
	}

	if err := activeRulesInstance.SetPayload(payload.Data); err != nil {
		log.Println(err)
		return
	}
	if err := activeRulesInstance.SortRuleSetsByPriority(); err != nil {
		log.Println(err)
		return
	}

	//activeRules = activeRulesInstance.GetPayloadKeyMapping()
}

func init() {
	var payload model.RequestPayload
	if err := json.Unmarshal([]byte(jsonRequestPayload), &payload); err != nil {
		log.Println(err)
		return
	}
	requestPayloadInstance := GetRequestPayloadInstance()
	requestPayloadInstance.SetPayload(payload)
}

func TestCheckCardBIN(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkCardBIN(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkCardBIN", err)
	}
}

func TestCheckThreeUniqueCardsAllowed(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkThreeUniqueCardsAllowed(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkThreeUniqueCardsAllowed", err)
	}
}

func TestCheckFifteenCountClearance(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkFifteenCountClearance(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkFifteenCountClearance", err)
	}
}

func TestCheckOneApprovedAllowedByThirtyMinuteInterval(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkOneApprovedAllowedByThirtyMinuteInterval(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkOneApprovedAllowedByThirtyMinuteInterval", err)
	}
}

func TestCheckOneCardPerBank(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkOneCardPerBank(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkOneCardPerBank", err)
	}
}

func TestCheckOneTcknPerUser(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkOneTcknPerUser(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkOneTcknPerUser", err)
	}
}

func TestCheckLastTenTransactions(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkLastTenTransactions(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkLastTenTransactions", err)
	}
}

func TestCheckPendingCountThreshold(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkPendingCountThreshold(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkPendingCountThreshold", err)
	}
}

func TestCheckPendingAllowanceByTimeInterval(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkPendingAllowanceByTimeInterval(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkPendingAllowanceByTimeInterval", err)
	}
}

func TestCheckApprovedAllowanceByTimeInterval(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkApprovedAllowanceByTimeInterval(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkApprovedAllowanceByTimeInterval", err)
	}
}

func TestCheckMaxDailyAllowancePerUser(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkMaxDailyAllowancePerUser(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkMaxDailyAllowancePerUser", err)
	}
}

func TestCheckMinTransactionAmount(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkMinTransactionAmount(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkMinTransactionAmount", err)
	}
}

func TestCheckMaxTransactionAmount(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkMaxTransactionAmount(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkMaxTransactionAmount", err)
	}
}

func TestCheckCardholdersNameMatch(t *testing.T) {
	if isOK, err := requestPayloadInstance.checkCardholdersNameMatch(); !isOK || err != nil {
		t.Errorf("%v check is failed!\nError Details: %v", "checkCardholdersNameMatch", err)
	}
}
