package game

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type NetClient struct {
	player    *Player
	client    *websocket.Conn
	interrupt chan os.Signal
}

func NewNetClient(serverAddress string, player *Player) *NetClient {
	u := url.URL{Scheme: "ws", Host: serverAddress, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	nc := &NetClient{
		player:    player,
		client:    c,
		interrupt: interrupt,
	}
	go nc.readLoop()
	return nc
}

func (c *NetClient) close() {
	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := c.client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	select {
	case <-time.After(time.Second):
	}
}

func (c *NetClient) sendData(data []byte) {
	err := c.client.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func (c *NetClient) readLoop() {

	for {
		_, message, err := c.client.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}

}
