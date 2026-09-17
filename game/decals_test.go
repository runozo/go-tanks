package game

import (
	"testing"

	"github.com/runozo/go-tanks/internal/maps"
	"github.com/solarlune/resolv"
)

// TestIsSoftGround verifies the terrain rule: tracks are left only on soft
// ground (grass, sand, transitions), never on roads, and never off-map.
func TestIsSoftGround(t *testing.T) {
	m := &maps.Map{
		Width:  4,
		Height: 2,
		Tiles: [][]string{
			{"tileGrass1", "tileSand1", "tileGrass_roadEast", "tileGrass_roadCrossing"},
			{"tileGrass_transitionN", "tileSand_roadNorth", "tileGrass2", "tileSand2"},
		},
	}
	g := &Game{playfield: &Playfield{mapData: m}}

	tests := []struct {
		px, py int
		want   bool
	}{
		{0, 0, true},    // grass
		{64, 0, true},   // sand
		{128, 0, false}, // road
		{192, 0, false}, // crossing
		{0, 64, true},   // transition (soft)
		{64, 64, false}, // sand road
		{128, 64, true}, // grass2
		{192, 64, true}, // sand2

		{-10, 0, false},    // off-map (left)
		{0, -10, false},    // off-map (top)
		{4 * 64, 0, false}, // off-map (right)
		{0, 2 * 64, false}, // off-map (bottom)
	}
	for _, tc := range tests {
		got := g.isSoftGround(resolv.Vector{X: float64(tc.px), Y: float64(tc.py)})
		if got != tc.want {
			t.Errorf("isSoftGround(%d,%d) = %v, want %v", tc.px, tc.py, got, tc.want)
		}
	}
}

// TestTileAt verifies the cell lookup and clamping.
func TestTileAt(t *testing.T) {
	m := &maps.Map{
		Width:  2,
		Height: 2,
		Tiles:  [][]string{{"tileGrass1", "tileSand1"}, {"tileGrass2", "tileSand2"}},
	}
	pf := &Playfield{mapData: m}

	if got := pf.TileAt(10, 10); got != "tileGrass1" {
		t.Errorf("TileAt(10,10) = %q, want tileGrass1", got)
	}
	if got := pf.TileAt(70, 10); got != "tileSand1" {
		t.Errorf("TileAt(70,10) = %q, want tileSand1", got)
	}
	if got := pf.TileAt(10, 70); got != "tileGrass2" {
		t.Errorf("TileAt(10,70) = %q, want tileGrass2", got)
	}
	if got := pf.TileAt(200, 200); got != "" {
		t.Errorf("TileAt(200,200) = %q, want empty", got)
	}
	if got := pf.TileAt(-5, -5); got != "" {
		t.Errorf("TileAt(-5,-5) = %q, want empty", got)
	}
}
