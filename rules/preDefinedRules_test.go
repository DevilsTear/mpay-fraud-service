package rules

import (
	"encoding/json"
	"fraud-service/model"
	rulesets "fraud-service/ruleset"
	"log"
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

	activeRules = activeRulesInstance.GetPayloadKeyMapping()
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
