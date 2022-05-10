package config

import (
	"context"
	"encoding/json"
	"fmt"
	"fraud-service/utils"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis/v8"
)

func LoadInitials(ctx context.Context) {
	LoadRedisSettings(ctx)
}

func LoadRedisSettings(ctx context.Context) {
	redisConfigFile, err := os.Open("config/redisSettings.json")
	utils.CheckError(err)
	byteValue, err := ioutil.ReadAll(redisConfigFile)
	utils.CheckError(err)
	json.Unmarshal(byteValue, &RedisConfig)
	defer redisConfigFile.Close()

	client := redis.NewClient(&redis.Options{
		Addr:     RedisConfig.Address,
		Password: RedisConfig.Password, // no password set
		DB:       RedisConfig.DB,       // use default DB
	})

	if err := CheckRedisServer(ctx, client); err != nil {
		panic(err)
	}

	RedisClient = *client
}

func CheckRedisServer(ctx context.Context, client *redis.Client) error {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}
