package game

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type TankPayload struct {
	ID       string  `json:"id"` // <- Fondamentale per riconoscere chi si muove!
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Rotation float64 `json:"r"`
	IsEnemy  bool    `json:"is_enemy"`
}

type NetClient struct {
	client        *websocket.Conn
	interrupt     chan os.Signal
	IncomingTanks chan TankPayload // <- Canale per inviare i dati al loop di gioco
}

func NewNetClient(serverAddress string) *NetClient {
	u := url.URL{Scheme: "ws", Host: serverAddress, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	nc := &NetClient{
		client:        c,
		interrupt:     interrupt,
		IncomingTanks: make(chan TankPayload, 100), // Bufferizziamo fino a 100 aggiornamenti
	}
	go nc.readLoop()
	return nc
}

func (c *NetClient) readLoop() {
	for {
		_, message, err := c.client.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		// Decodifichiamo il JSON in arrivo
		var payload TankPayload
		if err := json.Unmarshal(message, &payload); err == nil {
			// Lo spediamo al thread principale tramite il channel
			c.IncomingTanks <- payload
		}
	}
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

func (c *NetClient) SendTankData(tank *Tank) {
	payload := TankPayload{
		ID:       tank.ID,
		X:        tank.Object.Center().X,
		Y:        tank.Object.Center().Y,
		Rotation: tank.Object.Rotation(),
		IsEnemy:  tank.IsEnemy,
	}

	j, err := json.Marshal(payload)
	if err != nil {
		log.Println("encode error:", err)
		return // Meglio non usare log.Fatal qui, altrimenti il gioco crasha se c'è un glitch di rete
	}

	c.SendData(j)
}

func (c *NetClient) SendData(data []byte) {
	err := c.client.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Println("write:", err)
		return
	}
}
