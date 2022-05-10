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
	case config.PUB_BLACK_LIST:
	case config.PUB_CLEAR_COUNTER:
	case config.PUB_INCREASE_COUNTER:
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
		case config.SUB_RULE_SET_CHANGED:
			rulesetPayload := model.RuleSetPayload{}
			payload := receiveMessage(ctx, redisSubChan)
			json.Unmarshal(payload, &rulesetPayload)
			config.ChannelRuleSetPayload <- rulesetPayload

			//TODO - Set global variables..
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
