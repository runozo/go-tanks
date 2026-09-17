package game

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/runozo/go-tanks/internal/assets"
)

// buildSpriteMap loads the sprite map definition (mapped_tiles.json) without
// creating an ebiten image, so validation can run headlessly.
func buildSpriteMap(t *testing.T) *assets.Assets {
	t.Helper()
	data, err := assetsFS.ReadFile("assets/mapped_tiles.json")
	if err != nil {
		t.Fatalf("read mapped_tiles.json: %v", err)
	}
	var entries []assets.TileEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		t.Fatalf("parse mapped_tiles.json: %v", err)
	}
	sprites := &assets.Assets{TileEntries: map[string]assets.TileEntry{}}
	for _, e := range entries {
		sprites.TileEntries[e.Name] = e
	}
	return sprites
}

// TestMapsValid loads every embedded map and verifies it is well-formed and
// references only real sprites, with spawn/obstacle positions inside bounds.
func TestMapsValid(t *testing.T) {
	sprites := buildSpriteMap(t)

	maps, err := loadMaps(sprites, "")
	if err != nil {
		t.Fatalf("loadMaps: %v", err)
	}
	if len(maps) == 0 {
		t.Fatal("no maps loaded")
	}

	for _, m := range maps {
		m := m
		t.Run(m.Name, func(t *testing.T) {
			boundsW := m.Width * tileWidth
			boundsH := m.Height * tileHeight

			if len(m.Tiles) != m.Height {
				t.Fatalf("tiles height %d != %d", len(m.Tiles), m.Height)
			}
			for y, row := range m.Tiles {
				if len(row) != m.Width {
					t.Fatalf("row %d width %d != %d", y, len(row), m.Width)
				}
				for x, name := range row {
					if name == "" {
						t.Errorf("empty tile at (%d,%d)", x, y)
						continue
					}
					if _, ok := sprites.TileEntries[name]; !ok {
						t.Errorf("tile %q at (%d,%d) not found in sprites", name, x, y)
					}
				}
			}

			for i, o := range m.Obstacles {
				if _, ok := sprites.TileEntries[o.Type]; !ok {
					t.Errorf("obstacle %d: type %q not found in sprites", i, o.Type)
				}
				if o.X < 0 || o.Y < 0 || o.X >= float64(boundsW) || o.Y >= float64(boundsH) {
					t.Errorf("obstacle %d: position (%.0f,%.0f) out of bounds", i, o.X, o.Y)
				}
			}

			checkSpawn := func(label string, sp *SpawnPoint) {
				if sp == nil {
					return
				}
				if sp.X < 0 || sp.Y < 0 || sp.X >= float64(boundsW) || sp.Y >= float64(boundsH) {
					t.Errorf("%s: position (%.0f,%.0f) out of bounds", label, sp.X, sp.Y)
				}
			}
			checkSpawn("playerSpawn", m.PlayerSpawn)
			for i, sp := range m.EnemySpawns {
				checkSpawn("enemySpawn["+strconv.Itoa(i)+"]", &sp)
			}
		})
	}
}
