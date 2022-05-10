package main

import (
	"context"
	"encoding/json"
	"flag"
	"fraud-service/config"
	"fraud-service/model"
	"fraud-service/pubsub"
	"log"
	"net/http"
)

var ctx = context.Background()

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	config.LoadInitials(ctx)
	go pubsub.SubscribeEvent(ctx, "fraud:ruleset_changed")

	http.HandleFunc("/fraud", func(w http.ResponseWriter, r *http.Request) {
		serveEndpoint(w, r, "fraud")
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

func serveEndpoint(w http.ResponseWriter, r *http.Request, endpoint string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	resPayload := model.ResponsePayload{
		Status:  model.SuccessResponse,
		Code:    100,
		Message: "Success",
	}

	switch endpoint {
	case "fraud":
		var payload model.RequestPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		resData, err := json.Marshal(payload)
		if err != nil {
			log.Println(err)
		}
		resPayload = model.ResponsePayload{
			Status:  model.SuccessResponse,
			Code:    100,
			Message: "Success",
			Data:    resData,
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resPayload)
}
