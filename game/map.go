package game

import (
	"github.com/runozo/go-tanks/internal/maps"
)

// Re-exported types so the rest of the game package keeps using the previous
// names (Map, SpawnPoint, ObstacleSpawn).
type Map = maps.Map
type SpawnPoint = maps.SpawnPoint
type ObstacleSpawn = maps.ObstacleSpawn
