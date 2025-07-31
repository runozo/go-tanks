package game

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/runozo/go-wave-function-collapse/assets"
	"github.com/runozo/go-wave-function-collapse/wfc"
)

const (
	tileWidth  = 64
	tileHeight = 64
)

type Playfield struct {
	width       int
	height      int
	numOfTilesX int
	numOfTilesY int
	wfc         *wfc.Wfc
	assets      *assets.Assets
}

// NewPlayfield creates a new Playfield with the specified width, height, and assets.
//
// Parameters:
// - width: the width of the Playfield in pixels.
// - height: the height of the Playfield in pixels.
// - assets: a pointer to the Assets struct containing the necessary assets for the Playfield.
//
// Returns:
// - a pointer to the newly created Playfield.
func NewPlayfield(width, height int, ass *assets.Assets) *Playfield {

	tilesX := width/tileWidth + 1
	tilesY := height/tileHeight + 1

	playfield := &Playfield{
		width:  width,
		height: height,

		wfc:         wfc.NewWfc(screenWidth/tileWidth+1, screenHeight/tileHeight+1, ass.TileEntries),
		assets:      ass,
		numOfTilesX: tilesX,
		numOfTilesY: tilesY,
	}

	go playfield.wfc.StartRender()

	return playfield
}

func (p *Playfield) Update(tps float64) {
}

func (p *Playfield) Draw(screen *ebiten.Image) {
	var i int
	for y := 0; y < screenHeight; y += tileHeight {
		for x := 0; x < screenWidth; x += tileWidth {
			ops := &ebiten.DrawImageOptions{}
			ops.GeoM.Translate(float64(x), float64(y))
			if p.wfc.Tiles[i].Collapsed {
				screen.DrawImage(p.assets.GetSprite(p.wfc.Tiles[i].Name), ops)
			} else {
				screen.DrawImage(ebiten.NewImage(tileWidth, tileHeight), ops)
			}
			i++
		}
	}
}
