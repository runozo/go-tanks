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

func (c *NetClient) Close() {
	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	err := c.client.WriteMessage(websocket.CloseMessage, closeMsg)
	if err != nil {
		return
	}
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	<-timer.C
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
