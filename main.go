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
	go pubsub.SubscribeEvent(ctx, config.SUB_RULE_SET_CHANGED)

	http.HandleFunc("/"+config.FRAUD_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		endpoint.ServeEndpoint(w, r, config.FRAUD_ENDPOINT)
	})

	http.HandleFunc("/"+config.RULES_ENDPOINT, func(w http.ResponseWriter, r *http.Request) {
		endpoint.ServeEndpoint(w, r, config.RULES_ENDPOINT)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
