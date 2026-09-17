package game

import (
	"encoding/json"
	"fmt"

	"github.com/runozo/go-tanks/internal/assets"
)

// SpawnPoint describes a position (and optional rotation) where a tank appears.
type SpawnPoint struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Rotation float64 `json:"rotation"`
}

// ObstacleSpawn describes an authored obstacle placed on the map.
type ObstacleSpawn struct {
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Type string  `json:"type"`
}

// Map is a hand-authored level. The Grid holds one legend character per cell,
// which is expanded (via legend) into the tile names used for rendering.
// Obstacles and spawn points are authored so that levels are deterministic.
type Map struct {
	Name        string            `json:"name"`
	Legend      map[string]string `json:"legend"`
	Grid        []string          `json:"grid"`
	Obstacles   []ObstacleSpawn   `json:"obstacles"`
	PlayerSpawn *SpawnPoint       `json:"playerSpawn"`
	EnemySpawns []SpawnPoint      `json:"enemySpawns"`

	// Derived fields (not serialized).
	Width  int
	Height int
	Tiles  [][]string
}

// parseMap validates and expands a raw map json into a *Map. It checks that:
//   - the grid is rectangular (all rows have the same length),
//   - every character used is present in the legend,
//   - every referenced tile name exists in the asset map.
func parseMap(jsonData []byte, sprites *assets.Assets) (*Map, error) {
	var m Map
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return nil, fmt.Errorf("invalid map json: %w", err)
	}

	if len(m.Grid) == 0 {
		return nil, fmt.Errorf("map %q has an empty grid", m.Name)
	}

	m.Height = len(m.Grid)
	m.Width = len(m.Grid[0])

	m.Tiles = make([][]string, m.Height)
	for y, row := range m.Grid {
		if len(row) != m.Width {
			return nil, fmt.Errorf("map %q row %d has length %d, expected %d", m.Name, y, len(row), m.Width)
		}
		m.Tiles[y] = make([]string, m.Width)
		for x := 0; x < m.Width; x++ {
			ch := string(row[x])
			name, ok := m.Legend[ch]
			if !ok {
				return nil, fmt.Errorf("map %q row %d col %d: character %q not in legend", m.Name, y, x, ch)
			}
			if sprites == nil {
				return nil, fmt.Errorf("map %q: nil sprite map", m.Name)
			}
			if _, exists := sprites.TileEntries[name]; !exists {
				return nil, fmt.Errorf("map %q: tile %q (char %q) not found in sprites", m.Name, name, ch)
			}
			m.Tiles[y][x] = name
		}
	}

	return &m, nil
}
