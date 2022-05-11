package pubsub

import (
	"context"
	"encoding/json"
	"fraud-service/config"
	"fraud-service/model"
	rulesets "fraud-service/ruleset"
	"fraud-service/utils"
	"log"
)

// PublishEvent method publishes a new event to the appropriate Redis channel
func PublishEvent(ctx context.Context, eventType string, payload interface{}) {
	switch eventType {
	case config.PUB_BLACK_LIST:
	case config.PUB_CLEAR_COUNTER:
	case config.PUB_INCREASE_COUNTER:
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
		case config.SUB_RULE_SET_CHANGED:
			rulesetPayload := model.RuleSetPayload{}
			json.Unmarshal(payload, &rulesetPayload)
			log.Println(rulesetPayload)
			utils.SortRuleSetsByPriority(&rulesetPayload)
			log.Println(rulesetPayload)
			activeRules := rulesets.GetInstance()
			activeRules.SetPayload(rulesetPayload.Data)
			log.Println(activeRules.GetPayload())
			// config.ChannelRuleSetPayload <- rulesetPayload
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
