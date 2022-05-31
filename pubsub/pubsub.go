package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"fraud-service/config"
	"fraud-service/model"
	rulesets "fraud-service/ruleset"
)

// PublishEvent method publishes a new event to the appropriate Redis channel
func PublishEvent(ctx context.Context, eventType string, payload interface{}) {
	switch eventType {
	case config.PubBlackList:
	case config.PubClearCounter:
	case config.PubIncreaseCounter:
		if pubPayload, err := json.Marshal(payload); err == nil {
			config.RedisClient.Publish(ctx, eventType, pubPayload)
		}
	}
}

// SubscribeEvent method subscribes to the appropriate Redis channel and receives data from the related listerners
// redis-cli sample:
// PUBLISH "fraud:rule_sets_changed" '{"data": [{"name": "Rule 1","key": "Key ","priority": 1,"value": "10","status": true,"operator": "gt"},{"name": "Rule 2","key": "Key ","priority": 4,"value": "10","status": true,"operator": "gt"},{"name": "Rule 3","key": "Key ","priority": 3,"value": "10","status": false,"operator": "gt"},{"name": "Rule 4","key": "Key ","priority":2,"value": "10","status": true,"operator": "gt"}]}'
// TODO: make refactoring about robustness
func SubscribeEvent(ctx context.Context, redisSubChan string) {
	// sub_loop:
	for {
		payload := receiveMessage(ctx, redisSubChan)
		switch redisSubChan {
		case config.SubRuleSetChanged:
			rulesetPayload := model.RuleSetPayload{}
			err := json.Unmarshal(payload, &rulesetPayload)
			if err != nil {
				fmt.Printf("an error occured! Error Details: %v\n", err)
				break
			}
			//log.Println(rulesetPayload)
			activeRules := rulesets.GetInstance()
			err = activeRules.SetPayload(rulesetPayload.Data)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = activeRules.SortRuleSetsByPriority()
			if err != nil {
				fmt.Printf("an error occured! Error Details: %v\n", err)
				break
			}
			//log.Println(rulesetPayload)
			//log.Println(activeRules.GetPayload())
			//config.ChannelRuleSetPayload <- activeRules.Data
			// default:
			// 	break sub_loop
		}
	}
}

func receiveMessage(ctx context.Context, redisSubChan string) []byte {
	subscriber := config.RedisClient.Subscribe(ctx, redisSubChan)
	msg, err := subscriber.ReceiveMessage(ctx)
	if err != nil {
		panic(err)
	}
	payload := []byte(msg.Payload)

	return payload
}
