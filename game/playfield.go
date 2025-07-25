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

type Tile struct {
	image *ebiten.Image
}
type Playfield struct {
	width       int
	height      int
	tiles       []Tile
	numOfTilesX int
	numOfTilesY int
	isRendered  bool
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
		tiles:  make([]Tile, tilesX*tilesY),

		isRendered:  false,
		numOfTilesX: tilesX,
		numOfTilesY: tilesY,
	}

	myWfc := wfc.NewWfc(screenWidth/tileWidth+1, screenHeight/tileHeight+1, ass.TileEntries)
	myWfc.StartRender()
	for y := 0; y < tilesY; y++ {
		for x := 0; x < tilesX; x++ {
			tile := Tile{
				image: ass.GetSprite(myWfc.Tiles[y*tilesX+x].Name),
			}
			playfield.tiles[y*tilesX+x] = tile
		}
	}

	return playfield
}

func (p *Playfield) Update(tps float64) {
}

func (p *Playfield) Draw(screen *ebiten.Image) {
	var i int
	for y := 0; y < p.height; y += tileHeight {
		for x := 0; x < p.width; x += tileWidth {
			if x%tileWidth == 0 && y%tileHeight == 0 && i < len(p.tiles) {
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(p.tiles[i].image, ops)
				i++
			}
		}
	}
}
