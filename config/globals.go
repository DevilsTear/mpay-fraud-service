package config

import (
	"fraud-service/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	RedisConfig model.RedisConfig
	RedisClient redis.Client
)

var (
	MySqlConfig model.MySqlConfig
	MySqlDB     gorm.DB
)

var ChannelRuleSetPayload = make(chan []model.RuleSet)

var (
	// PendingCountThreshold is a fraud control param expresses to check pending request counts
	PendingCountThreshold int = 10

	// PendingAllowanceByTimeInterval is a fraud control param expresses to check pending request allowance in a predefined time slot
	PendingAllowanceByTimeInterval int = 10

	// ApprovedAllowanceByTimeInterval is a fraud control param expresses to check pending request approval in a predefined time slot
	ApprovedAllowanceByTimeInterval int = 30

	// MaxDailyAllowancePerUser is a fraud control param expresses to check max daily request allowance limit per user
	MaxDailyAllowancePerUser int = 5

	// MinTransactionAmount is a fraud control param expresses to check min daily total transaction allowance limit per user
	MinTransactionAmount float64 = 50.0

	// MaxTransactionAmount is a fraud control param expresses to check max daily total transaction allowance limit per user
	MaxTransactionAmount float64 = 1000.0
)

// Redis pub/sub channel names
const (
	SUB_RULE_SET_CHANGED = "fraud:rule_sets_changed"
	PUB_BLACK_LIST       = "fraud:blacklist"
	PUB_CLEAR_COUNTER    = "fraud:clear_counter"
	PUB_INCREASE_COUNTER = "fraud:increase_counter"
)

const (
	FRAUD_ENDPOINT = "fraud"
	RULES_ENDPOINT = "rules"
)
