// Package maps defines the hand-authored level format shared by the game
// client (rendering) and the server (enemy AI, obstacle positions).
package maps

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed maps/*.json
var mapsFS embed.FS

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
type Map struct {
	Name        string            `json:"name"`
	Legend      map[string]string `json:"legend"`
	Grid        []string          `json:"grid"`
	Obstacles   []ObstacleSpawn   `json:"obstacles"`
	PlayerSpawn *SpawnPoint       `json:"playerSpawn"`
	EnemySpawns []SpawnPoint      `json:"enemySpawns"`

	// Derived fields (not serialized).
	Width  int        `json:"-"`
	Height int        `json:"-"`
	Tiles  [][]string `json:"-"`
}

// Parse validates and expands raw map json into a *Map. It checks that:
//   - the grid is rectangular (all rows have the same length),
//   - every character used is present in the legend,
//   - every referenced tile name passes the optional tileExists check.
//
// tileExists may be nil when the caller does not need sprite validation
// (e.g. the server, which never renders the grid).
func Parse(data []byte, tileExists func(string) bool) (*Map, error) {
	var m Map
	if err := json.Unmarshal(data, &m); err != nil {
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
			if tileExists != nil && !tileExists(name) {
				return nil, fmt.Errorf("map %q: tile %q (char %q) not found in sprites", m.Name, name, ch)
			}
			m.Tiles[y][x] = name
		}
	}

	return &m, nil
}

// All loads and parses every embedded map without sprite validation.
func All() ([]*Map, error) {
	entries, err := mapsFS.ReadDir("maps")
	if err != nil {
		return nil, err
	}

	var all []*Map
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := mapsFS.ReadFile("maps/" + e.Name())
		if err != nil {
			return nil, err
		}
		m, err := Parse(data, nil)
		if err != nil {
			return nil, fmt.Errorf("loading %s: %w", e.Name(), err)
		}
		all = append(all, m)
	}

	if len(all) == 0 {
		return nil, fmt.Errorf("no maps found")
	}

	return all, nil
}

// LoadDir loads and parses every *.json map in the given directory (a runtime
// directory where authored/edited maps are stored). A missing directory is not
// an error; it simply yields no maps.
func LoadDir(dir string) ([]*Map, error) {
	if dir == "" {
		return nil, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var all []*Map
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		m, err := Parse(data, nil)
		if err != nil {
			return nil, fmt.Errorf("loading %s/%s: %w", dir, e.Name(), err)
		}
		all = append(all, m)
	}
	return all, nil
}

// NeedsSpriteValidation reports whether the map should be validated against
// the sprite map before being used. (Kept for the game client.)
func (m *Map) Validate(tileExists func(string) bool) error {
	for y, row := range m.Tiles {
		for x, name := range row {
			if !tileExists(name) {
				return fmt.Errorf("map %q: tile %q at (%d,%d) not found", m.Name, name, x, y)
			}
		}
	}
	return nil
}
