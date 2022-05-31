package endpoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fraud-service/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRulesEndpoint(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ServeEndpoint(w, r, "rules")
	}

	rulesPayload := model.RuleSetPayload{}
	err := json.Unmarshal([]byte(reqRulesJSON), &rulesPayload)
	if err != nil {
		fmt.Println(err)
	}

	reqBody, _ := json.Marshal(rulesPayload)

	req := httptest.NewRequest("POST", "/rules", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	handler(w, req)

	// Status code test
	if w.Code != 200 {
		t.Error("Http test request failed!")
	}

	rulesResponsePayload := model.RuleSetPayload{}
	err = json.Unmarshal(w.Body.Bytes(), &rulesResponsePayload)
	if err != nil {
		t.Log(err)
		t.Log(w.Body.String())
		t.Error("Unexpected response, test failed!")
	}
}

func TestFraudEndpoint(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ServeEndpoint(w, r, "fraud")
	}

	reqPayload := model.RequestPayload{}
	err := json.Unmarshal([]byte(reqFraudJSON), &reqPayload)
	if err != nil {
		t.Log(err)
		t.Error("Unexpected response, test failed!")
	}

	reqBody, _ := json.Marshal(reqPayload)

	req := httptest.NewRequest("POST", "/fraud", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	handler(w, req)

	// Status code test
	if w.Code != 200 {
		t.Error("Http test request failed!")
	}

	if w.Body == nil {
		t.Error("Unexpected response, test failed!")
	}
}

var reqFraudJSON = `{
	"data": {
		"amount": "250.00",
		"trx": "oaisufklafasfl1111112d1233",
		"card_number": "4943******2271",
		"expiration_month": "10",
		"expiration_year": "2022",
		"cardholders_name": "*** *****",
		"cvv": "402",
		"return_url": "https://envoysoft3.net/deposit/mpayReturn"
	},
	"user": {
		"username": "srdr16",
		"userID": "17206739184",
		"yearofbirth": "1983",
		"fullname": "*** *****",
		"email": "****@hotmail.com",
		"tckn": "17206739184",
		"ip_address": "95.10.24.238"
	}
}`

var resFraudExpectedJSON = ``

var reqRulesJSON = `{
	"data": [
		{
			"name": "Rule 1",
			"key": "Key ",
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

// var resRulesExpectedJSON = `{"Status":"success","Data":[{"name":"Rule 1","key":"Key ","priority":1,"value":"10","status":true,"operator":"gt"},{"name":"Rule 4","key":"Key ","priority":2,"value":"10","status":true,"operator":"gt"},{"name":"Rule 2","key":"Key ","priority":4,"value":"10","status":true,"operator":"gt"}],"Code":200,"Message":"Success"}`
