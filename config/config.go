package config

import (
	"context"
	"encoding/json"
	"fmt"
	"fraud-service/utils"
	"io/ioutil"
	"os"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func LoadInitials(ctx context.Context) {
	LoadRedisSettings(ctx)
	LoadMySqlSettings(ctx)
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

func LoadMySqlSettings(ctx context.Context) {
	mysqlConfigFile, err := os.Open("config/mysqlSettings.json")
	utils.CheckError(err)
	byteValue, err := ioutil.ReadAll(mysqlConfigFile)
	utils.CheckError(err)
	json.Unmarshal(byteValue, &MySqlConfig)
	defer mysqlConfigFile.Close()

	db, err := CheckMySqlServer()
	if err != nil {
		panic(err)
	}

	MySqlDB = *db
}

func CheckMySqlServer() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       MySqlConfig.DSN,                       // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         MySqlConfig.DefaultStringSize,         // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  MySqlConfig.DisableDatetimePrecision,  // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &MySqlConfig.DefaultDatetimePrecision, // default datetime precision
		DontSupportRenameIndex:    MySqlConfig.DontSupportRenameIndex,    // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   MySqlConfig.DontSupportRenameColumn,   // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: MySqlConfig.SkipInitializeWithVersion, // smart configure based on used version
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
