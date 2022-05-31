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

// LoadInitials sets application specific global settings
func LoadInitials(ctx context.Context) {
	LoadRedisSettings(ctx)
	LoadMySQLSettings(ctx)
}

// LoadRedisSettings sets Redis settings on a global scope
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

// CheckRedisServer checks if Redis server is ready and responding
func CheckRedisServer(ctx context.Context, client *redis.Client) error {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println(pong, err)
	// Output: PONG <nil>

	return nil
}

// LoadMySQLSettings sets MySql settings on a global scope
func LoadMySQLSettings(ctx context.Context) {
	mysqlConfigFile, err := os.Open("config/mysqlSettings.json")
	utils.CheckError(err)
	byteValue, err := ioutil.ReadAll(mysqlConfigFile)
	utils.CheckError(err)
	json.Unmarshal(byteValue, &MySQLConfig)
	defer mysqlConfigFile.Close()

	db, err := CheckMySQLServer()
	if err != nil {
		panic(err)
	}

	MySQLDb = *db
}

// CheckMySQLServer checks if MySql server is ready and responding
func CheckMySQLServer() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       MySQLConfig.DSN,                       // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         MySQLConfig.DefaultStringSize,         // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  MySQLConfig.DisableDatetimePrecision,  // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &MySQLConfig.DefaultDatetimePrecision, // default datetime precision
		DontSupportRenameIndex:    MySQLConfig.DontSupportRenameIndex,    // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   MySQLConfig.DontSupportRenameColumn,   // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: MySQLConfig.SkipInitializeWithVersion, // smart configure based on used version
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
