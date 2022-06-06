package endpoint

import (
	"encoding/json"
	"fmt"
	"fraud-service/config"
	"fraud-service/model"
	"fraud-service/rules"
	rulesets "fraud-service/ruleset"
	"log"
	"net/http"
)

var (
	successResponse = model.ResponsePayload{
		Status:  model.SuccessResponse,
		Code:    http.StatusOK,
		Message: "Success",
	}
	failResponse = model.ResponsePayload{
		Status: model.FailResponse,
		Code:   -100,
	}
	errorResponse = model.ResponsePayload{
		Status: model.ErrorResponse,
		Code:   http.StatusBadRequest,
	}
)

// ServeEndpoint includes listeners for HTTP connections
func ServeEndpoint(w http.ResponseWriter, r *http.Request, endpoint string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch endpoint {
	case config.FraudEndpoint:
		successResponse.Data = true
		var payload model.RequestPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Println(err)
			errorResponse.Message = err.Error()
			respJson, _ := json.Marshal(errorResponse)
			http.Error(w, string(respJson), http.StatusBadRequest)

			return
		}

		// Fraud checks
		requestPayloadInstance := rules.GetRequestPayloadInstance()
		requestPayloadInstance.SetPayload(payload)
		isPassed, err := requestPayloadInstance.ProcessRules()
		if err != nil || !isPassed {
			log.Println(err)
			failResponse.Message = err.Error()
			failResponse.Data = isPassed
			respJson, _ := json.Marshal(errorResponse)
			http.Error(w, string(respJson), http.StatusConflict)

			return
		}

		//isPassed, err = rules.EvaluateRules(&payload)
		//if err != nil || !isPassed {
		//	log.Println(err)
		//	resPayload.Status = model.FailResponse
		//	resPayload.Code = -100
		//	resPayload.Message = "Fail"
		//	resPayload.Data = isPassed
		//}
	case config.RulesEndpoint:
		var payload model.RuleSetPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Println(err)
			errorResponse.Message = err.Error()
			respJson, _ := json.Marshal(errorResponse)
			http.Error(w, string(respJson), http.StatusBadRequest)

			return
		}
		// log.Println(payload)
		activeRules := rulesets.GetInstance()
		err = activeRules.SetPayload(payload.Data)
		if err != nil {
			log.Println(err)
			errorResponse.Data = err.Error()
			respJson, _ := json.Marshal(errorResponse)
			http.Error(w, string(respJson), http.StatusBadRequest)

			return
		}
		err = activeRules.SortRuleSetsByPriority()
		if err != nil {
			log.Println(err)
			errorResponse.Data = err.Error()
			respJson, _ := json.Marshal(errorResponse)
			http.Error(w, string(respJson), http.StatusBadRequest)

			return
		}

		// log.Println(activeRules.GetPayload())
		successResponse.Data = activeRules.GetPayload()
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(successResponse)
	if err != nil {
		fmt.Printf("an error occurred while writing json to the http response writer! Error Details: %v\n", err)
	}
}
