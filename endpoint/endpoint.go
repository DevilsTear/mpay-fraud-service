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

// ServeEndpoint includes listeners for HTTP connections
func ServeEndpoint(w http.ResponseWriter, r *http.Request, endpoint string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	resPayload := model.ResponsePayload{
		Status:  model.SuccessResponse,
		Code:    http.StatusOK,
		Message: "Success",
	}

	switch endpoint {
	case config.FraudEndpoint:
		resPayload := model.ResponsePayload{
			Status:  model.SuccessResponse,
			Code:    http.StatusOK,
			Message: "Success",
			Data:    true,
		}
		var payload model.RequestPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Fraud checks
		requestPayloadInstance := rules.GetRequestPayloadInstance()
		requestPayloadInstance.SetPayload(payload)
		isPassed, err := requestPayloadInstance.ProcessRules()
		if err != nil || !isPassed {
			log.Println(err)
			resPayload.Status = model.FailResponse
			resPayload.Code = -100
			resPayload.Message = "Fail"
			resPayload.Data = isPassed
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
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// log.Println(payload)
		activeRules := rulesets.GetInstance()
		err = activeRules.SetPayload(payload.Data)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = activeRules.SortRuleSetsByPriority()
		if err != nil {
			fmt.Println(err)
			return
		}
		// log.Println(activeRules.GetPayload())
		resPayload = model.ResponsePayload{
			Status:  model.SuccessResponse,
			Code:    http.StatusOK,
			Message: "Success",
			Data:    activeRules.GetPayload(),
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(resPayload)
	if err != nil {
		fmt.Printf("an error occurred while writing json to the http response writer! Error Details: %v\n", err)
	}
}
