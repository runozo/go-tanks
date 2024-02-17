package game

import (
	_ "image/png"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tileWidth  = 64
	tileHeight = 64
	ruleUP     = 0
	ruleRIGHT  = 1
	ruleDOWN   = 2
	ruleLEFT   = 3
)

var tile_options = []string{
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

type Tile struct {
	name    string
	options []string
}

var rules = map[string][][]string{
	// up, right, down, left
	"tileGrass1.png": {
		{"tileGrass1.png", "tileGrass2.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_transitionN.png"}, {"tileGrass1.png", "tileGrass2.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"}, {"tileGrass1.png", "tileGrass2.png", "tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"}, {"tileGrass1.png", "tileGrass2.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
	},
	"tileGrass2.png": {
		{"tileGrass1.png", "tileGrass2.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_transitionN.png"}, {"tileGrass1.png", "tileGrass2.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"}, {"tileGrass1.png", "tileGrass2.png", "tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"}, {"tileGrass1.png", "tileGrass2.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
	},
	"tileGrass_roadCornerLL.png": {
		{"tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"},
		{"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
		{"tileGrass_roadCornerUL.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png",
			"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitW.png", "tileGrass_roadTransitionS_dirt.png", "tileGrass_roadTransitionS.png"},
		{"tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCrossing.png",
			"tileGrass_roadCrossingRound.png",
			"tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionW_dirt.png", "tileGrass_roadTransitionW.png"},
	},
	"tileGrass_roadCornerLR.png": {
		{"tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"}, {"tileGrass_roadCornerLL.png", "tileGrass_roadCornerUL.png", "tileGrass_roadCrossing.png",
			"tileGrass_roadCrossingRound.png",
			"tileGrass_roadSplitW.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionE_dirt.png", "tileGrass_roadTransitionE.png"}, {"tileGrass_roadCornerUL.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png",
			"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitW.png", "tileGrass_roadTransitionS_dirt.png", "tileGrass_roadTransitionS.png"}, {"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
	},
	"tileGrass_roadCornerUL.png": {
		{"tileGrass_roadCornerUL.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png",
			"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitW.png", "tileGrass_roadTransitionS_dirt.png", "tileGrass_roadTransitionS.png"}, {"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"}, {"tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"}, {"tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCrossing.png",
			"tileGrass_roadCrossingRound.png",
			"tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionW_dirt.png", "tileGrass_roadTransitionW.png"},
	}, "tileGrass_roadCornerUR.png": {
		{"tileGrass_roadCornerUL.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png",
			"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitW.png", "tileGrass_roadTransitionS_dirt.png", "tileGrass_roadTransitionS.png"}, {"tileGrass_roadCornerLL.png", "tileGrass_roadCornerUL.png", "tileGrass_roadCrossing.png",
			"tileGrass_roadCrossingRound.png",
			"tileGrass_roadSplitW.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionE_dirt.png", "tileGrass_roadTransitionE.png"}, {"tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"}, {"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
	}, "tileGrass_roadCrossing.png": {
		{"tileGrass_roadCornerLL.png", "tileGrass_roadCornerLR.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitW.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionN_dirt.png", "tileGrass_roadTransitionN.png"}, {"tileGrass_roadCornerUL.png", "tileGrass_roadCornerLL.png", "tileGrass_roadEast.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionE_dirt.png", "tileGrass_roadTransitionE.png"}, {"tileGrass_roadCornerUL.png", "tileGrass_roadCornerUR.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitW.png", "tileGrass_roadTransitionS_dirt.png", "tileGrass_roadTransitionS.png"}, {"tileGrass_roadCornerUR.png", "tileGrass_roadCornerLR.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionW_dirt.png", "tileGrass_roadTransitionW.png"},
	},
	"tileGrass_roadCrossingRound.png": {{}, {}, {}, {}},
	"tileGrass_roadEast.png":          {{}, {}, {}, {}}, "tileGrass_roadNorth.png": {{}, {}, {}, {}}, "tileGrass_roadSplitE.png": {{}, {}, {}, {}}, "tileGrass_roadSplitN.png": {{}, {}, {}, {}}, "tileGrass_roadSplitS.png": {{}, {}, {}, {}}, "tileGrass_roadSplitW.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionE.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionE_dirt.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionN.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionN_dirt.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionS.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionS_dirt.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionW.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionW_dirt.png": {{}, {}, {}, {}}, "tileGrass_transitionE.png": {{}, {}, {}, {}}, "tileGrass_transitionN.png": {{}, {}, {}, {}}, "tileGrass_transitionS.png": {{}, {}, {}, {}}, "tileGrass_transitionW.png": {{}, {}, {}, {}}, "tileSand1.png": {{}, {}, {}, {}}, "tileSand2.png": {{}, {}, {}, {}}, "tileSand_roadCornerLL.png": {{}, {}, {}, {}}, "tileSand_roadCornerLR.png": {{}, {}, {}, {}}, "tileSand_roadCornerUL.png": {{}, {}, {}, {}}, "tileSand_roadCornerUR.png": {{}, {}, {}, {}}, "tileSand_roadCrossing.png": {{}, {}, {}, {}}, "tileSand_roadCrossingRound.png": {{}, {}, {}, {}}, "tileSand_roadEast.png": {{}, {}, {}, {}}, "tileSand_roadNorth.png": {{}, {}, {}, {}}, "tileSand_roadSplitE.png": {{}, {}, {}, {}}, "tileSand_roadSplitN.png": {{}, {}, {}, {}}, "tileSand_roadSplitS.png": {{}, {}, {}, {}}, "tileSand_roadSplitW.png": {{}, {}, {}, {}},
}

type Playfield struct {
	game  *Game
	tiles []*Tile
}

func NewPlayfield(game *Game) *Playfield {
	// Wave collapse algorithm
	// Initialize playfield data structure
	var tiles []*ebiten.Image

	var i int
	for y := 0; y < game.height; y += tileHeight {
		for x := 0; x < game.width; x += tileWidth {
			if x%tileWidth == 0 && y%tileHeight == 0 {
				tile := allTiles[rand.Intn(len(allTiles))]
				tile.image = game.assets.GetSprite(tile.name)
				tiles = append(tiles, &tile)
				// fmt.Println(i, tile.name)
				i++
			}
		}
	}
	// var PlayfieldTile = mustLoadImage("png/tileGrass2.png")
	return &Playfield{
		game:  game,
		tiles: tiles,
	}
}

func (p *Playfield) Update() {

}

func (p *Playfield) Draw(screen *ebiten.Image) {
	var i int
	for y := 0; y < p.game.height; y += tileHeight {
		for x := 0; x < p.game.width; x += tileWidth {
			if x%tileWidth == 0 && y%tileHeight == 0 {
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(p.tiles[i].image, ops)
				i++
			}
		}
	}
}
