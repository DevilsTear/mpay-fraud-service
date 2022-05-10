package config

import (
	"fraud-service/model"

	"github.com/go-redis/redis/v8"
)

var RedisConfig model.RedisConfig
var RedisClient redis.Client
var ChannelRuleSetPayload = make(chan model.RuleSetPayload)

// Fraud control params
var PendingCountThreshold int = 10
var PendingAllowanceByTimeInterval int = 10
var ApprovedAllowanceByTimeInterval int = 30
var MaxDailyAllowancePerUser int = 5
var MinTransactionAmount float64 = 50.0
var MaxTransactionAmount float64 = 1000.0
