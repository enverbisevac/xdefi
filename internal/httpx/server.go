package httpx

import (
	"encoding/json"
	"github.com/enverbisevac/xdefi/internal/queue"
	"github.com/enverbisevac/xdefi/pkg/api"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Server struct {
	router *mux.Router
	queue  queue.Producer
}

func NewServer(router *mux.Router, queue queue.Producer) Server {
	server := Server{
		router: router,
		queue:  queue,
	}
	server.SetupRoutes()
	return server
}

func (s Server) SetupRoutes() {
	s.router.HandleFunc("/health", s.healthHandler)
	s.router.HandleFunc("/sign", s.signHandler).Methods("POST")
}

func (s Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	response(w, http.StatusOK, map[string]bool{
		"healthy": true,
	})
}

func (s Server) signHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response(w, http.StatusBadRequest, api.Error{Detail: err.Error()})
		return
	}

	var reqBody api.EventBody
	err = json.Unmarshal(data, &reqBody)
	if err != nil {
		response(w, http.StatusBadRequest, api.Error{Detail: err.Error()})
		return
	}

	data, err = json.Marshal(reqBody)
	if err != nil {
		response(w, http.StatusInternalServerError, api.Error{Detail: err.Error()})
		return
	}

	err = s.queue.PushToQueue(data)
	if err != nil {
		response(w, http.StatusInternalServerError, api.Error{Detail: err.Error()})
		return
	}

	response(w, http.StatusAccepted, nil)
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
