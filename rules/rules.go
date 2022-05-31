package rules

import (
	"fmt"
	"fraud-service/model"
	"fraud-service/utils"

	"github.com/nikunjy/rules/parser"
)

func evaluate(rule string, items map[string]interface{}) (retVal bool, retErr error) {
	defer func() {
		info := recover()
		if info != nil {
			retErr = fmt.Errorf("%q", info)
		}
	}()
	ev, err := parser.NewEvaluator(rule)
	if err != nil {
		return bool(false), err
	}
	return ev.Process(items)
}

// EvaluateRules triggers the generic rule engine
func EvaluateRules(payload *model.RequestPayload) (bool, error) {
	if payload == nil {
		return false, fmt.Errorf("%q", recover())
	}
	mappedPayload, err := utils.Struct2Map(&payload)
	if err != nil || mappedPayload == nil {
		return false, fmt.Errorf("%q", recover())
	}

	// for _, rule := range activeRules.GetPayload() {

	// 	isPassed, err := Evaluate()
	// }

	return bool(true), nil
}
