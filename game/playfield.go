package game

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/runozo/go-wave-function-collapse/assets"
	"github.com/runozo/go-wave-function-collapse/wfc"
	"github.com/solarlune/resolv"
)

const (
	tileWidth  = 64
	tileHeight = 64
)

type Playfield struct {
	width, height, numOfTilesX, numOfTilesY int
	wfc                                     *wfc.Wfc
	assets                                  *assets.Assets
	progressBar                             *ProgressBar
}

func NewPlayfield(game *Game) *Playfield {

	playfield := &Playfield{
		width:  game.width,
		height: game.height,

		wfc:         wfc.NewWfc(screenWidth/tileWidth+1, screenHeight/tileHeight+1, game.assets.TileEntries),
		assets:      game.assets,
		progressBar: NewProgressBar(400, 36, "Generating playfield", resolv.Vector{X: (screenWidth - 400) / 2, Y: (screenHeight - 36) / 2}, game.fontSmall),
		numOfTilesX: game.width/tileWidth + 1,
		numOfTilesY: game.height/tileHeight + 1,
	}

	playfield.progressBar.Update(0)

	go playfield.wfc.StartRender()

	return playfield
}

func (p *Playfield) Update(tps float64) {
	if p.wfc.IsRunning {
		p.progressBar.Update(float64(p.wfc.ProcessedTiles) / float64(p.wfc.TotalTiles) * 100)
	}
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

	if p.wfc.IsRunning {
		p.progressBar.Draw(screen)
	}

}
