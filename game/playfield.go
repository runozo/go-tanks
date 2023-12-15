package game

import (
	"fmt"
	_ "image/png"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tileWidth  = 64
	tileHeight = 64
)

var allTiles = []string{
	"tileGrass_roadCornerLL.png",
	"tileGrass_roadCornerLR.png",
	"tileGrass_roadCornerUL.png",
	"tileGrass_roadCornerUR.png",
	"tileGrass_roadCrossing.png",
	"tileGrass_roadCrossingRound.png",
	"tileGrass_roadEast.png",
	"tileGrass_roadNorth.png",
	"tileGrass_roadSplitE.png",
	"tileGrass_roadSplitN.png",
	"tileGrass_roadSplitS.png",
	"tileGrass_roadSplitW.png",
	"tileGrass_roadTransitionE.png",
	"tileGrass_roadTransitionE_dirt.png",
	"tileGrass_roadTransitionN.png",
	"tileGrass_roadTransitionN_dirt.png",
	"tileGrass_roadTransitionS.png",
	"tileGrass_roadTransitionS_dirt.png",
	"tileGrass_roadTransitionW.png",
	"tileGrass_roadTransitionW_dirt.png",
	"tileGrass_transitionE.png",
	"tileGrass_transitionN.png",
	"tileGrass_transitionS.png",
	"tileGrass_transitionW.png",
	"tileSand1.png",
	"tileSand2.png",
	"tileSand_roadCornerLL.png",
	"tileSand_roadCornerLR.png",
	"tileSand_roadCornerUL.png",
	"tileSand_roadCornerUR.png",
	"tileSand_roadCrossing.png",
	"tileSand_roadCrossingRound.png",
	"tileSand_roadEast.png",
	"tileSand_roadNorth.png",
	"tileSand_roadSplitE.png",
	"tileSand_roadSplitN.png",
	"tileSand_roadSplitS.png",
	"tileSand_roadSplitW.png",
}

type Playfield struct {
	game   *Game
	tiles  []*ebiten.Image
	width  int
	height int
}

func NewPlayfield(game *Game, width, height int) *Playfield {
	var tiles []*ebiten.Image
	var i int
	for y := 0; y < height; y += tileHeight {
		for x := 0; x < width; x += tileWidth {
			if x%tileWidth == 0 && y%tileHeight == 0 {
				tileName := allTiles[rand.Intn(len(allTiles))]
				tiles = append(tiles, game.assets.GetSprite(tileName))
				fmt.Println(i, tileName)
				i++
			}
		}
	}
	// var PlayfieldTile = mustLoadImage("png/tileGrass2.png")
	return &Playfield{
		game:   game,
		tiles:  tiles,
		width:  width,
		height: height,
	}
}

func (p *Playfield) Update() {

}

func (p *Playfield) Draw(screen *ebiten.Image) {
	var i int
	for y := 0; y < p.height; y += tileHeight {
		for x := 0; x < p.width; x += tileWidth {
			if x%tileWidth == 0 && y%tileHeight == 0 {
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(p.tiles[i], ops)
				i++
			}
		}
	}
}
