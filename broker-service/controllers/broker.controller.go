package controllers

import (
	"broken/helpers"
	"encoding/json"
	"net/http"
)

func Broker(w http.ResponseWriter, r *http.Request) {
	payload := helpers.JsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}
