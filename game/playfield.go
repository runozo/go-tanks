package game

import (
	_ "image/png"
	"log/slog"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/runozo/go-tanks/assets"
)

const (
	tileWidth  = 64
	tileHeight = 64
)

type Tile struct {
	collapsed bool
	image     *ebiten.Image
	options   []string
}
type Playfield struct {
	width       int
	height      int
	tiles       []Tile
	numOfTilesX int
	numOfTilesY int
	isRendered  bool
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
func NewPlayfield(width, height int, assets *assets.Assets) *Playfield {

	tilesX := width/tileWidth + 1
	tilesY := height/tileHeight + 1

	// setup tiles with all the options enabled
	tiles := make([]Tile, tilesX*tilesY)
	resetTilesOptions(&tiles)

	playfield := &Playfield{
		width:       width,
		height:      height,
		tiles:       tiles,
		isRendered:  false,
		numOfTilesX: tilesX,
		numOfTilesY: tilesY,
		assets:      assets,
	}
	go renderPlayfield(playfield)

	return playfield
}

func (p *Playfield) Update() {
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

func renderPlayfield(p *Playfield) {
	startTime := time.Now()
	defer func() {
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		slog.Info("Rendering of playfield took", "duration", duration)
	}()

	for !p.isRendered {
		// pick the minimum entropy indexes
		minEntropyIndexes := getMinEntropyIndexes(&p.tiles)

		if len(minEntropyIndexes) <= 0 {
			slog.Info("Playfiled is rendered. No more collapsable cells.", "tiles", len(p.tiles))
			for i := 0; i < len(p.tiles); i++ {
				if !p.tiles[i].collapsed {
					p.tiles[i].image = p.assets.GetSprite(p.tiles[i].options[0])
					p.tiles[i].collapsed = true
				}
			}
			p.isRendered = true
		} else {
			collapsedIndex := collapseRandomCellWithMinEntropy(&p.tiles, &minEntropyIndexes)
			p.tiles[collapsedIndex].image = p.assets.GetSprite(p.tiles[collapsedIndex].options[0])
			p.tiles[collapsedIndex].collapsed = true

			for y := 0; y < p.numOfTilesY; y++ {
				for x := 0; x < p.numOfTilesX; x++ {
					index := y*p.numOfTilesX + x
					if len(p.tiles[index].options) == 0 {
						// we did not found any options, let's restart
						slog.Info("Restarting!")
						resetTilesOptions(&p.tiles)
					}

					if !p.tiles[index].collapsed {
						// Look UP
						if y > 0 {
							p.tiles[index].options = lookAndFilter(ruleUP, ruleDOWN, p.tiles[index].options, p.tiles[(y-1)*p.numOfTilesX+x].options)
						}
						// Look RIGHT
						if x < p.numOfTilesX-1 {
							p.tiles[index].options = lookAndFilter(ruleRIGHT, ruleLEFT, p.tiles[index].options, p.tiles[y*p.numOfTilesX+x+1].options)
						}
						// Look DOWN
						if y < p.numOfTilesY-1 {
							p.tiles[index].options = lookAndFilter(ruleDOWN, ruleUP, p.tiles[index].options, p.tiles[(y+1)*p.numOfTilesX+x].options)
						}
						// Look LEFT
						if x > 0 {
							p.tiles[index].options = lookAndFilter(ruleLEFT, ruleRIGHT, p.tiles[index].options, p.tiles[y*p.numOfTilesX+x-1].options)
						}
					}
				}
			}
		}
	}
}
