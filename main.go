package main

import (
	"context"
	"encoding/json"
	"flag"
	"fraud-service/config"
	"fraud-service/endpoint"
	"fraud-service/model"
	"fraud-service/pubsub"
	rulesets "fraud-service/ruleset"
	"fraud-service/utils"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var ctx = context.Background()

var addr = flag.String("addr", ":8080", "http service address")

func init() {
	if err := config.LoadInitials(ctx); err != nil {
		panic(err)
	}
	if err := loadPredefinedRules(ctx); err != nil {
		log.Println(err)
	}
}

func main() {
	go pubsub.SubscribeEvent(ctx, config.SubRuleSetChanged)

	http.HandleFunc("/"+config.FraudEndpoint, func(w http.ResponseWriter, r *http.Request) {
		endpoint.ServeEndpoint(w, r, config.FraudEndpoint)
	})

	http.HandleFunc("/"+config.RulesEndpoint, func(w http.ResponseWriter, r *http.Request) {
		endpoint.ServeEndpoint(w, r, config.RulesEndpoint)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// loadPredefinedRules sets the default rule parameters via preconfigured json file
func loadPredefinedRules(ctx context.Context) error {
	var payload model.RuleSetPayload
	rulesConfigFile, err := os.Open("config/ruleSetPayload.json")
	defer rulesConfigFile.Close()
	utils.CheckError(err)
	byteValue, err := ioutil.ReadAll(rulesConfigFile)
	utils.CheckError(err)
	err = json.Unmarshal(byteValue, &payload)
	if err != nil {
		return err
	}

	activeRules := rulesets.GetInstance()
	if err := activeRules.SetPayload(payload.Data); err != nil {
		log.Println(err)
		return err
	}

	if err := activeRules.SortRuleSetsByPriority(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
