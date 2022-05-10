package pubsub

import (
	"context"
	"encoding/json"
	"fraud-service/config"
	"fraud-service/model"
)

// PublishEvent method publishes a new event to the appropriate Redis channel
func PublishEvent(ctx context.Context, eventType string, payload interface{}) {
	switch eventType {
	case "fraud:blacklist":
	case "fraud:clear_counter":
	case "fraud:increase_counter":
		if pubPayload, err := json.Marshal(payload); err == nil {
			config.RedisClient.Publish(ctx, eventType, pubPayload)
		}
	}
}

// SubscribeEvent method subscribes to the appropriate Redis channel and receives data from the related listerners
func SubscribeEvent(ctx context.Context, redisSubChan string) {
sub_loop:
	for {
		switch redisSubChan {
		case "fraud:rule_sets_changed":
			rulesetPayload := model.RuleSetPayload{}
			payload := receiveMessage(ctx, redisSubChan)
			json.Unmarshal(payload, &rulesetPayload)
			config.ChannelRuleSetPayload <- rulesetPayload
		default:
			break sub_loop
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
