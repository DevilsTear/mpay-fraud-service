package main

import (
	"context"
	"flag"
	"fraud-service/config"
	"fraud-service/endpoint"
	"fraud-service/pubsub"
	"log"
	"net/http"
)

var ctx = context.Background()

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	config.LoadInitials(ctx)
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
