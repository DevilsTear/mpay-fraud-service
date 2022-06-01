package config

import (
	"fraud-service/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	// RedisConfig holds global Redis config settings
	RedisConfig model.RedisConfig
	// RedisClient holds global scope Redis Client
	RedisClient redis.Client
)

var (
	// MySQLConfig holds global MySQL config settings
	MySQLConfig model.MySQLConfig
	// MySQLDb holds global scope Gorm Instance
	MySQLDb gorm.DB
)

//var ChannelRuleSetPayload = make(chan []model.RuleSet)

var (
	// PendingCountThreshold is a fraud control param expresses to check pending request counts
	PendingCountThreshold int64 = 10

	// PendingAllowanceByTimeInterval is a fraud control param expresses to check pending request allowance in a predefined time slot
	PendingAllowanceByTimeInterval int64 = 10

	// ApprovedAllowanceByTimeInterval is a fraud control param expresses to check pending request approval in a predefined time slot
	ApprovedAllowanceByTimeInterval int64 = 30

	// MaxDailyAllowancePerUser is a fraud control param expresses to check max daily request allowance limit per user
	MaxDailyAllowancePerUser int64 = 5

	// MinTransactionAmount is a fraud control param expresses to check min daily total transaction allowance limit per user
	MinTransactionAmount float64 = 50.0

	// MaxTransactionAmount is a fraud control param expresses to check max daily total transaction allowance limit per user
	MaxTransactionAmount float64 = 1000.0
)

// Redis pub/sub channel names
const (
	SubRuleSetChanged  = "fraud:rule_sets_changed"
	PubBlackList       = "fraud:blacklist"
	PubClearCounter    = "fraud:clear_counter"
	PubIncreaseCounter = "fraud:increase_counter"
)

// HTTP endpoints
const (
	FraudEndpoint = "fraud"
	RulesEndpoint = "rules"
)
