package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/runozo/go-tanks/internal/assets"
)

const (
	tileWidth  = 64
	tileHeight = 64
)

// Playfield is a hand-authored level (a Map) rendered on top of the screen.
// Unlike the previous WFC-based generation, tiles are fixed and authored, so
// the layout is always deterministic.
type Playfield struct {
	assets  *assets.Assets
	mapData *Map
	ready   bool
}

// NewPlayfield builds a Playfield from the map chosen on the Game (by name, or
// a random one among those embedded).
func NewPlayfield(game *Game) *Playfield {
	pf := &Playfield{
		assets:  game.assets,
		mapData: game.pickMap(),
		ready:   true,
	}
	log.Printf("Playfield: using map %q (%dx%d)", pf.mapData.Name, pf.mapData.Width, pf.mapData.Height)
	return pf
}

func (p *Playfield) Update(tps float64) {}

func (p *Playfield) Draw(screen *ebiten.Image) {
	for y := 0; y < p.mapData.Height; y++ {
		for x := 0; x < p.mapData.Width; x++ {
			name := p.mapData.Tiles[y][x]
			ops := &ebiten.DrawImageOptions{}
			ops.GeoM.Translate(float64(x*tileWidth), float64(y*tileHeight))
			screen.DrawImage(p.assets.GetSprite(name), ops)
		}
	}
}
