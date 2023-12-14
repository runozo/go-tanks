package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type GameServer struct {
	Name string
	Host string
	Port int
}

const wsServerEndpoint = "ws://localhost:40000"

func NewGameServer(name string, host string, port int) *GameServer {
	return &GameServer{
		Name: name,
		Host: host,
		Port: port,
	}
}

func (s *GameServer) startHTTP() {
	log.Println("Starting HTTP server on port %d", g.Port)
	go func() {
		http.HandleFunc("/ws", s.handleWS)
		http.ListenAndServe(fmt.Sprintf("%s:%d", g.Host, g.Port), nil)
	}()
}

func (s *GameServer) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}
	log.Println("Client connected")
}

func main() {
	gs := NewGameServer("Game Server", "localhost", 40000)
}
