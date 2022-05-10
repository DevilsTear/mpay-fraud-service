package main

import (
	"context"
	"encoding/json"
	"flag"
	"fraud-service/config"
	"fraud-service/pubsub"
	"log"
	"net/http"
)

var ctx = context.Background()

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	config.LoadInitials(ctx)
	pubsub.SubscribeEvent(ctx, "fraud:ruleset_changed")

	http.HandleFunc("/fraud", func(w http.ResponseWriter, r *http.Request) {
		serverEndpoint(w, r, "fraud")
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func serverEndpoint(w http.ResponseWriter, r *http.Request, endpoint string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload interface{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
