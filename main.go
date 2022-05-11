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

	http.HandleFunc("/fraud", func(w http.ResponseWriter, r *http.Request) {
		endpoint.ServeEndpoint(w, r, "fraud")
	})

	http.HandleFunc("/rules", func(w http.ResponseWriter, r *http.Request) {
		endpoint.ServeEndpoint(w, r, "rules")
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
