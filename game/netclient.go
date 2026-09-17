package game

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/runozo/go-tanks/internal/protocol"
)

// Re-exported aliases so game code keeps using the previous local names.
type PlayerState = protocol.PlayerState
type EnemyState = protocol.EnemyState
type ScoreState = protocol.ScoreState
type NetMessage = protocol.NetMessage

const (
	msgJoin           = protocol.MsgJoin
	msgSnapshot       = protocol.MsgSnapshot
	msgTransform      = protocol.MsgTransform
	msgFire           = protocol.MsgFire
	msgHit            = protocol.MsgHit
	msgEnemyState     = protocol.MsgEnemyState
	msgEnemyFire      = protocol.MsgEnemyFire
	msgEnemyHitPlayer = protocol.MsgEnemyHitPlayer
	msgScore          = protocol.MsgScore
	msgLeave          = protocol.MsgLeave
)

// NetClient connects to the server hub and relays typed messages.
type NetClient struct {
	client    *websocket.Conn
	interrupt chan os.Signal
	Incoming  chan NetMessage
	clientID  string
}

func NewNetClient(serverAddress, clientID string) *NetClient {
	u := url.URL{Scheme: "ws", Host: serverAddress, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	nc := &NetClient{
		client:    c,
		interrupt: interrupt,
		Incoming:  make(chan NetMessage, 1024),
		clientID:  clientID,
	}

	go nc.readLoop()
	nc.SendJoin()

	return nc
}

// ClientID returns the stable id this client joined with.
func (c *NetClient) ClientID() string { return c.clientID }

func (c *NetClient) readLoop() {
	for {
		_, message, err := c.client.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		var msg NetMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("decode:", err)
			continue
		}
		c.Incoming <- msg
	}
}

func (c *NetClient) send(msg NetMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("encode error:", err)
		return
	}
	if err := c.client.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println("write:", err)
	}
}

func (c *NetClient) SendJoin() {
	c.send(NetMessage{Type: msgJoin, ClientID: c.clientID})
}

func (c *NetClient) SendTransform(x, y, r float64) {
	c.send(NetMessage{Type: msgTransform, ClientID: c.clientID, X: x, Y: y, R: r})
}

func (c *NetClient) SendFire(x, y, r, slope float64, flavor string) {
	c.send(NetMessage{Type: msgFire, ClientID: c.clientID, X: x, Y: y, R: r, Slope: slope, Flavor: flavor})
}

func (c *NetClient) SendHit(enemyID string) {
	c.send(NetMessage{Type: msgHit, ClientID: c.clientID, EnemyID: enemyID})
}

func (c *NetClient) SendLeave() {
	c.send(NetMessage{Type: msgLeave, ClientID: c.clientID})
}

func (c *NetClient) Close() {
	c.SendLeave()

	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	err := c.client.WriteMessage(websocket.CloseMessage, closeMsg)
	if err != nil {
		return
	}
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	<-timer.C
}
