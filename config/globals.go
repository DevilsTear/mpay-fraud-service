package config

import (
	"fraud-service/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var RedisConfig model.RedisConfig
var RedisClient redis.Client

var MySqlConfig model.MySqlConfig
var MySqlDB gorm.DB

var ChannelRuleSetPayload = make(chan model.RuleSetPayload)

// Fraud control params
var PendingCountThreshold int = 10
var PendingAllowanceByTimeInterval int = 10
var ApprovedAllowanceByTimeInterval int = 30
var MaxDailyAllowancePerUser int = 5
var MinTransactionAmount float64 = 50.0
var MaxTransactionAmount float64 = 1000.0

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
