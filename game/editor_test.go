package game

import (
	"encoding/json"
	"testing"

	"github.com/runozo/go-tanks/internal/maps"
)

// TestEditorToMapRoundTrip verifies that a map produced by the editor survives
// a save-to-JSON / parse-from-JSON round trip (legend generation included).
func TestEditorToMapRoundTrip(t *testing.T) {
	e := &Editor{name: "test", selTile: baseGrass, selObst: "crateWood"}
	e.New(8, 5, "test")

	// paint a few roads (tiles is [y][x])
	e.tiles[1][1] = "tileGrass_roadEast"     // cell (1,1)
	e.tiles[1][2] = "tileGrass_roadEast"     // cell (2,1)
	e.tiles[2][1] = "tileGrass_roadNorth"    // cell (1,2)
	e.tiles[3][4] = "tileGrass_roadCrossing" // cell (4,3)
	e.tiles[4][6] = "tileSand1"              // cell (6,4)

	// add content
	e.obstacles = append(e.obstacles, ObstacleSpawn{X: 0, Y: 0, Type: "crateWood"}, ObstacleSpawn{X: 128, Y: 128, Type: "treeGreen_small"})
	e.playerSpawn = &SpawnPoint{X: 96, Y: 96, Rotation: 0.5}
	e.enemySpawns = append(e.enemySpawns, SpawnPoint{X: 160, Y: 160}, SpawnPoint{X: 224, Y: 224})

	m := e.ToMap("test")

	// simulate a save: marshal to JSON then parse it back (as the loader does)
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// derived fields must not be serialized
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, forbidden := range []string{"Width", "Height", "Tiles"} {
		if _, ok := raw[forbidden]; ok {
			t.Errorf("serialized map should not contain %q", forbidden)
		}
	}

	got, err := maps.Parse(data, nil)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	if got.Width != 8 || got.Height != 5 {
		t.Fatalf("dims = %dx%d, want 8x5", got.Width, got.Height)
	}

	checks := []struct {
		x, y int
		want string
	}{
		{1, 1, "tileGrass_roadEast"},
		{2, 1, "tileGrass_roadEast"},
		{1, 2, "tileGrass_roadNorth"},
		{4, 3, "tileGrass_roadCrossing"},
		{6, 4, "tileSand1"},
		{0, 0, baseGrass},
		{7, 4, baseGrass},
	}
	for _, c := range checks {
		if got.Tiles[c.y][c.x] != c.want {
			t.Errorf("tile(%d,%d) = %q, want %q", c.x, c.y, got.Tiles[c.y][c.x], c.want)
		}
	}

	if len(got.Obstacles) != 2 {
		t.Fatalf("got %d obstacles, want 2", len(got.Obstacles))
	}
	if got.Obstacles[0].Type != "crateWood" {
		t.Errorf("obstacle type = %q", got.Obstacles[0].Type)
	}

	if got.PlayerSpawn == nil || got.PlayerSpawn.X != 96 || got.PlayerSpawn.Rotation != 0.5 {
		t.Errorf("player spawn = %+v", got.PlayerSpawn)
	}
	if len(got.EnemySpawns) != 2 {
		t.Fatalf("got %d enemy spawns, want 2", len(got.EnemySpawns))
	}
}
