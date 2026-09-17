// Package protocol defines the JSON wire messages shared by the game client
// and the server hub.
package protocol

// Message types (client <-> server).
const (
	MsgJoin           = "join"             // c2s: client_id, name
	MsgSnapshot       = "snapshot"         // s2c: map_name, players, enemies, score
	MsgTransform      = "transform"        // c2s + relay: client_id, x, y, r
	MsgFire           = "fire"             // c2s + relay: client_id, barrel world rotation, slope, flavor
	MsgHit            = "hit"              // c2s: client_id, enemy_id (player bullet hit an enemy)
	MsgEnemyState     = "enemy_state"      // s2c: enemies (server simulates the AI)
	MsgEnemyFire      = "enemy_fire"       // s2c: enemy_id, bullet_id, x, y, r, slope, flavor
	MsgEnemyHitPlayer = "enemy_hit_player" // s2c: client_id hit by an enemy bullet
	MsgScore          = "score"            // s2c: score
	MsgLeave          = "leave"            // both: client_id
)

// PlayerState is one player transform sent around.
type PlayerState struct {
	ClientID string  `json:"client_id"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	R        float64 `json:"r"`
}

// EnemyState is a server-simulated enemy transform.
type EnemyState struct {
	ID  string  `json:"id"`
	X   float64 `json:"x"`
	Y   float64 `json:"y"`
	R   float64 `json:"r"`
	Aim float64 `json:"aim"` // world barrel rotation
}

// ScoreState is the match score (player kills / AI kills).
type ScoreState struct {
	Player int `json:"player"`
	Comp   int `json:"comp"`
}

// NetMessage is the single JSON envelope used for every message type.
type NetMessage struct {
	Type string `json:"type"`

	ClientID string  `json:"client_id,omitempty"`
	Name     string  `json:"name,omitempty"`
	X        float64 `json:"x,omitempty"`
	Y        float64 `json:"y,omitempty"`
	R        float64 `json:"r,omitempty"`
	Slope    float64 `json:"slope,omitempty"`
	Flavor   string  `json:"flavor,omitempty"`
	EnemyID  string  `json:"enemy_id,omitempty"`
	BulletID int64   `json:"bullet_id,omitempty"`

	MapName string        `json:"map_name,omitempty"`
	Players []PlayerState `json:"players,omitempty"`
	Enemies []EnemyState  `json:"enemies,omitempty"`
	Score   ScoreState    `json:"score,omitempty"`
}
