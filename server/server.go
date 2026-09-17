// The multiplayer hub server.
//
// Responsibilities:
//   - relay player transforms and fire events between clients (state relay);
//   - simulate the AI enemies (which are shared entities) and broadcast their
//     state, their fire events and their hits on players;
//   - keep the match score.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/runozo/go-tanks/internal/enemyai"
	"github.com/runozo/go-tanks/internal/maps"
	"github.com/runozo/go-tanks/internal/protocol"
	"github.com/solarlune/resolv"
)

const (
	tps             = 60
	enemyStateHz    = 20
	bulletSpeed     = 8.0
	gravity         = 9.8
	initialAltitude = 0.2
	playerHitRadius = 28
	bulletLife      = 8.0 * tps // frames before a bullet is cleaned up
)

var (
	addr     = flag.String("addr", "localhost:8080", "http service address")
	upgrader = websocket.Upgrader{}
)

// client is one connected game.
type client struct {
	conn *websocket.Conn
	id   string
	x    float64
	y    float64
	r    float64
	send chan []byte
}

// serverEnemy mirrors the tanks the client renders, simulated server-side.
type serverEnemy struct {
	id   string
	x, y float64
	r    float64
	aim  float64
	ai   enemyai.State
	fire float64 // seconds remaining before next shot
}

// serverBullet is an enemy projectile simulated on the server.
type serverBullet struct {
	id                      int64
	x, y, rot               float64
	slope                   float64
	altitude, vSpeed, elaps float64
	life                    int
}

type hub struct {
	mu         sync.Mutex
	clients    map[string]*client
	enemies    []*serverEnemy
	bullets    []*serverBullet
	nextBullet int64

	playerScore map[string]int
	compScore   int

	mapData   *maps.Map
	obstacles []resolv.Vector
}

func newHub(m *maps.Map) *hub {
	h := &hub{
		clients:     make(map[string]*client),
		playerScore: make(map[string]int),
		mapData:     m,
	}

	for i, s := range m.EnemySpawns {
		e := &serverEnemy{
			id:   "enemy_" + strconv.Itoa(i),
			x:    s.X,
			y:    s.Y,
			ai:   enemyai.State{OrbitSign: 1},
			fire: 1.5,
		}
		if rand.Intn(2) == 0 {
			e.ai.OrbitSign = -1
		}
		h.enemies = append(h.enemies, e)
	}

	// Obstacle centers (sprite sizes are roughly 28-32px wide).
	for _, o := range m.Obstacles {
		h.obstacles = append(h.obstacles, resolv.Vector{X: o.X + 16, Y: o.Y + 16})
	}

	return h
}

// ---------------------------------------------------------------------------
// WebSocket handling
// ---------------------------------------------------------------------------

func (h *hub) serveWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	// The first message must be a join.
	_, raw, err := conn.ReadMessage()
	if err != nil {
		conn.Close()
		return
	}
	var join protocol.NetMessage
	if err := json.Unmarshal(raw, &join); err != nil || join.Type != protocol.MsgJoin {
		conn.Close()
		return
	}

	id := join.ClientID
	if id == "" {
		id = fmt.Sprintf("player_%d", time.Now().UnixNano())
	}

	c := &client{conn: conn, id: id, send: make(chan []byte, 128)}

	h.mu.Lock()
	h.clients[id] = c
	h.playerScore[id] = 0
	h.mu.Unlock()

	log.Printf("player joined: %s (total %d)", id, len(h.clients))

	// write loop
	go func() {
		for data := range c.send {
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		}
	}()

	h.sendSnapshot(c)

	// read loop
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			h.unregister(c)
			return
		}
		var msg protocol.NetMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}
		h.handleMessage(c, msg)
	}
}

func (h *hub) handleMessage(c *client, msg protocol.NetMessage) {
	switch msg.Type {
	case protocol.MsgTransform:
		h.mu.Lock()
		c.x, c.y, c.r = msg.X, msg.Y, msg.R
		h.mu.Unlock()
		h.broadcastExcept(c.id, msg)

	case protocol.MsgFire:
		h.broadcastExcept(c.id, msg)

	case protocol.MsgHit:
		h.mu.Lock()
		h.playerScore[c.id]++
		score := protocol.ScoreState{Player: h.playerScore[c.id], Comp: h.compScore}
		h.mu.Unlock()
		h.broadcast(scoreMessage(c.id, score))

	case protocol.MsgLeave:
		h.unregister(c)
	}
}

func (h *hub) unregister(c *client) {
	h.mu.Lock()
	if _, ok := h.clients[c.id]; !ok {
		h.mu.Unlock()
		return
	}
	delete(h.clients, c.id)
	h.mu.Unlock()

	log.Printf("player left: %s", c.id)
	close(c.send)
	h.broadcast(protocol.NetMessage{Type: protocol.MsgLeave, ClientID: c.id})
}

// ---------------------------------------------------------------------------
// Simulation
// ---------------------------------------------------------------------------

func (h *hub) startSimulation() {
	ticker := time.NewTicker(time.Second / tps)
	stateTicker := time.NewTicker(time.Second / enemyStateHz)
	go func() {
		for range ticker.C {
			h.mu.Lock()
			outgoing := h.stepSimulation(1.0 / tps)
			h.mu.Unlock()
			// broadcast only after releasing the lock
			for _, m := range outgoing {
				h.broadcast(m)
			}
		}
	}()
	go func() {
		for range stateTicker.C {
			h.broadcastEnemyState()
		}
	}()
}

// stepSimulation advances enemies and their bullets by dt seconds and returns
// the messages to broadcast to the clients. Caller must hold h.mu.
func (h *hub) stepSimulation(dt float64) []protocol.NetMessage {
	cfg := enemyai.DefaultCfg()
	cfg.BoundsX0, cfg.BoundsY0 = 0, 0
	cfg.BoundsX1 = float64(h.mapData.Width * 64)
	cfg.BoundsY1 = float64(h.mapData.Height * 64)

	var outgoing []protocol.NetMessage

	for _, e := range h.enemies {
		var target *client
		minDist := math.MaxFloat64
		for _, c := range h.clients {
			d := math.Hypot(c.x-e.x, c.y-e.y)
			if d < minDist {
				minDist = d
				target = c
			}
		}

		if target != nil {
			targetPos := resolv.Vector{X: target.x, Y: target.y}
			pos := resolv.Vector{X: e.x, Y: e.y}

			heading := enemyai.Steer(cfg, &e.ai, pos, targetPos, h.obstacles)

			// rotate toward the heading
			targetRotation := enemyai.RotationForHeading(heading)
			diff := math.Atan2(math.Sin(targetRotation-e.r), math.Cos(targetRotation-e.r))
			if math.Abs(diff) > 0.02 {
				rot := cfg.RotationSpeed / tps
				if diff < 0 {
					rot = -rot
				}
				e.r += rot
			}

			// move forward
			fwd := enemyai.Forward(e.r)
			e.x += fwd.X * cfg.BaseSpeed / tps * cfg.SpeedFactor
			e.y += fwd.Y * cfg.BaseSpeed / tps * cfg.SpeedFactor

			// aim and fire
			e.aim = enemyai.BarrelWorldRotation(pos, targetPos)
			e.fire -= dt
			if e.fire <= 0 {
				if slope := enemyai.BallisticSlope(minDist, tps); slope >= 0 {
					outgoing = append(outgoing, h.spawnEnemyBullet(e, slope))
					e.fire = 2.5 + rand.Float64()*1.0
				}
			}
		}
	}

	outgoing = append(outgoing, h.stepBullets(dt)...)
	return outgoing
}

// spawnEnemyBullet creates a server-side bullet and returns the fire message
// to broadcast. Caller must hold h.mu.
func (h *hub) spawnEnemyBullet(e *serverEnemy, slope float64) protocol.NetMessage {
	h.nextBullet++
	b := &serverBullet{
		id:       h.nextBullet,
		x:        e.x,
		y:        e.y,
		rot:      e.aim,
		slope:    slope,
		vSpeed:   bulletSpeed * math.Sin(slope),
		altitude: initialAltitude,
		life:     bulletLife,
	}
	h.bullets = append(h.bullets, b)

	return protocol.NetMessage{
		Type:     protocol.MsgEnemyFire,
		EnemyID:  e.id,
		BulletID: b.id,
		X:        b.x,
		Y:        b.y,
		R:        b.rot,
		Slope:    slope,
		Flavor:   "bulletSand3_outline",
	}
}

// stepBullets advances enemy bullets, detects hits on players and returns the
// messages to broadcast. Caller must hold h.mu.
func (h *hub) stepBullets(dt float64) []protocol.NetMessage {
	var outgoing []protocol.NetMessage
	alive := h.bullets[:0]
	for _, b := range h.bullets {
		b.life--
		b.elaps += dt
		b.vSpeed -= gravity * dt
		b.altitude += b.vSpeed * dt

		actualSpeed := bulletSpeed * math.Cos(b.slope)
		fwd := enemyai.Forward(b.rot)
		b.x += fwd.X * actualSpeed
		b.y += fwd.Y * actualSpeed

		if b.life <= 0 {
			continue
		}
		// Bullet "landed": check if it hit a player near the impact point.
		if b.altitude <= initialAltitude && b.elaps > 0.1 {
			for _, c := range h.clients {
				if math.Hypot(c.x-b.x, c.y-b.y) < playerHitRadius {
					h.compScore++
					outgoing = append(outgoing,
						protocol.NetMessage{
							Type:     protocol.MsgEnemyHitPlayer,
							ClientID: c.id,
							BulletID: b.id,
							X:        c.x,
							Y:        c.y,
						},
						scoreMessage(c.id, protocol.ScoreState{Player: h.playerScore[c.id], Comp: h.compScore}),
					)
					break
				}
			}
			continue
		}
		alive = append(alive, b)
	}
	h.bullets = alive
	return outgoing
}

// ---------------------------------------------------------------------------
// Broadcast helpers
// ---------------------------------------------------------------------------

func scoreMessage(forClient string, s protocol.ScoreState) protocol.NetMessage {
	return protocol.NetMessage{Type: protocol.MsgScore, ClientID: forClient, Score: s}
}

func (h *hub) sendSnapshot(c *client) {
	h.mu.Lock()
	players := make([]protocol.PlayerState, 0, len(h.clients))
	for id, other := range h.clients {
		players = append(players, protocol.PlayerState{ClientID: id, X: other.x, Y: other.y, R: other.r})
	}
	enemies := make([]protocol.EnemyState, 0, len(h.enemies))
	for _, e := range h.enemies {
		enemies = append(enemies, protocol.EnemyState{ID: e.id, X: e.x, Y: e.y, R: e.r, Aim: e.aim})
	}
	score := protocol.ScoreState{Player: h.playerScore[c.id], Comp: h.compScore}
	msg := protocol.NetMessage{
		Type:    protocol.MsgSnapshot,
		MapName: h.mapData.Name,
		Players: players,
		Enemies: enemies,
		Score:   score,
	}
	h.mu.Unlock()

	h.sendTo(c, msg)
}

func (h *hub) broadcastEnemyState() {
	h.mu.Lock()
	enemies := make([]protocol.EnemyState, 0, len(h.enemies))
	for _, e := range h.enemies {
		enemies = append(enemies, protocol.EnemyState{ID: e.id, X: e.x, Y: e.y, R: e.r, Aim: e.aim})
	}
	msg := protocol.NetMessage{Type: protocol.MsgEnemyState, Enemies: enemies}
	clients := make([]*client, 0, len(h.clients))
	for _, c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.Unlock()

	for _, c := range clients {
		h.sendTo(c, msg)
	}
}

func (h *hub) broadcast(msg protocol.NetMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, c := range h.clients {
		select {
		case c.send <- data:
		default:
			log.Printf("dropping message for slow client %s", c.id)
		}
	}
}

// broadcastExcept relays a message to every client except `exceptID`.
func (h *hub) broadcastExcept(exceptID string, msg protocol.NetMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	for id, c := range h.clients {
		if id == exceptID {
			continue
		}
		select {
		case c.send <- data:
		default:
			log.Printf("dropping message for slow client %s", c.id)
		}
	}
}

func (h *hub) sendTo(c *client, msg protocol.NetMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	select {
	case c.send <- data:
	default:
		log.Printf("dropping message for slow client %s", c.id)
	}
}

func main() {
	flag.Parse()

	all, err := maps.All()
	if err != nil {
		log.Fatal("load maps:", err)
	}
	m := all[rand.Intn(len(all))]
	log.Printf("session map: %q (%dx%d cells)", m.Name, m.Width, m.Height)

	h := newHub(m)

	http.HandleFunc("/", h.serveWS)
	log.Printf("listening on %s", *addr)

	h.startSimulation()
	log.Fatal(http.ListenAndServe(*addr, nil))
}
