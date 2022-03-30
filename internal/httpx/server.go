package httpx

import (
	"context"
	"encoding/json"
	"github.com/enverbisevac/xdefi/internal/queue"
	"github.com/enverbisevac/xdefi/pkg/api"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	http.Server
	router   *mux.Router
	queue    queue.Queue
	receiver <-chan string
	close    chan struct{}
}

func NewServer(ctx context.Context, router *mux.Router, q queue.Queue) *Server {
	server := &Server{
		Server: http.Server{
			Addr:    ":8087",
			Handler: router,
		},
		router:   router,
		queue:    q,
		receiver: q.Register(ctx, queue.SignedEvent),
		close:    make(chan struct{}),
	}
	server.SetupRoutes()
	return server
}

func (s *Server) SetupRoutes() {
	s.router.HandleFunc("/health", s.healthHandler)
	s.router.HandleFunc("/sign", s.signHandler).Methods("POST")
	s.router.HandleFunc("/ws", s.websocket).Methods("GET")
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	response(w, http.StatusOK, map[string]bool{
		"healthy": true,
	})
}

func (s *Server) signHandler(w http.ResponseWriter, r *http.Request) {
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

	err = s.queue.PushToQueue(queue.SignEvent, data)
	if err != nil {
		response(w, http.StatusInternalServerError, api.Error{Detail: err.Error()})
		return
	}

	response(w, http.StatusAccepted, nil)
}

func (s *Server) websocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case data := <-s.receiver:
				err := conn.WriteMessage(websocket.TextMessage, []byte(data))
				if err != nil {
					log.Println(err)
				}
			case <-done:
				return
			}
		}
	}()

	// check when websocket get closed
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			done <- struct{}{}
			close(done)
			break
		}
	}
}
