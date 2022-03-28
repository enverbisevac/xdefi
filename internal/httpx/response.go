package httpx

import (
	"encoding/json"
	"log"
	"net/http"
)

func response(w http.ResponseWriter, code int, payload interface{}) {
	if payload == nil {
		w.WriteHeader(code)
		return
	}
	response, err := json.Marshal(payload)
	if err != nil {
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			log.Printf("error writing to writer response")
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Printf("error writing to writer response")
	}
}
