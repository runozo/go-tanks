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

var rules = map[string][][]string{
	// up, right, down, left
	"tileGrass1.png": {
		{"tileGrass1.png", "tileGrass2.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_transitionN.png"},
		{"tileGrass1.png", "tileGrass2.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
		{"tileGrass1.png", "tileGrass2.png", "tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"},
		{"tileGrass1.png", "tileGrass2.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
	},
	"tileGrass2.png": {
		{"tileGrass1.png", "tileGrass2.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_transitionN.png"},
		{"tileGrass1.png", "tileGrass2.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
		{"tileGrass1.png", "tileGrass2.png", "tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"},
		{"tileGrass1.png", "tileGrass2.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
	},
	"tileGrass_roadCornerLL.png": {
		{"tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"},
		{"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
		{"tileGrass_roadCornerUL.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitW.png", "tileGrass_roadTransitionS_dirt.png", "tileGrass_roadTransitionS.png"},
		{"tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionW_dirt.png", "tileGrass_roadTransitionW.png"},
	},
	"tileGrass_roadCornerLR.png": {
		{"tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"},
		{"tileGrass_roadCornerLL.png", "tileGrass_roadCornerUL.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png", "tileGrass_roadSplitW.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionE_dirt.png", "tileGrass_roadTransitionE.png"},
		{"tileGrass_roadCornerUL.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitW.png", "tileGrass_roadTransitionS_dirt.png", "tileGrass_roadTransitionS.png"},
		{"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
	},
	"tileGrass_roadCornerUL.png": {
		{"tileGrass_roadCornerUL.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitW.png", "tileGrass_roadTransitionS_dirt.png", "tileGrass_roadTransitionS.png"},
		{"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
		{"tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"},
		{"tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionW_dirt.png", "tileGrass_roadTransitionW.png"},
	},
	"tileGrass_roadCornerUR.png": {
		{"tileGrass_roadCornerLL.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitW.png", "tileGrass_roadTransitionS_dirt.png", "tileGrass_roadTransitionS.png"},
		{"tileGrass_roadCornerLL.png", "tileGrass_roadCornerUL.png", "tileGrass_roadCrossing.png", "tileGrass_roadCrossingRound.png", "tileGrass_roadSplitW.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionE_dirt.png", "tileGrass_roadTransitionE.png"},
		{"tileGrass_roadEast.png", "tileGrass_roadSplitN.png", "tileGrass_roadCornerUR.png", "tileGrass_roadCornerUL.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionN.png"},
		{"tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadCornerLR.png", "tileGrass_roadCornerUR.png", "tileGrass1.png", "tileGrass2.png", "tileGrass_transitionE.png"},
	},
	"tileGrass_roadCrossing.png": {
		{"tileGrass_roadCornerLL.png", "tileGrass_roadCornerLR.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitW.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionN_dirt.png", "tileGrass_roadTransitionN.png"},
		{"tileGrass_roadCornerUL.png", "tileGrass_roadCornerLL.png", "tileGrass_roadEast.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionE_dirt.png", "tileGrass_roadTransitionE.png"},
		{"tileGrass_roadCornerUL.png", "tileGrass_roadCornerUR.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitW.png", "tileGrass_roadTransitionS_dirt.png", "tileGrass_roadTransitionS.png"},
		{"tileGrass_roadCornerUR.png", "tileGrass_roadCornerLR.png", "tileGrass_roadNorth.png", "tileGrass_roadSplitE.png", "tileGrass_roadSplitN.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionW_dirt.png", "tileGrass_roadTransitionW.png"},
	},
	//"tileGrass_roadCrossingRound.png": {{}, {}, {}, {}},
	//"tileGrass_roadEast.png":          {{}, {}, {}, {}},
	//"tileGrass_roadNorth.png": {{}, {}, {}, {}},
	//"tileGrass_roadSplitE.png": {{}, {}, {}, {}},
	//"tileGrass_roadSplitN.png": {{}, {}, {}, {}},
	//"tileGrass_roadSplitS.png": {{}, {}, {}, {}},
	//"tileGrass_roadSplitW.png": {{}, {}, {}, {}},
	//"tileGrass_roadTransitionE.png": {{}, {}, {}, {}},
	//"tileGrass_roadTransitionE_dirt.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionN.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionN_dirt.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionS.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionS_dirt.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionW.png": {{}, {}, {}, {}}, "tileGrass_roadTransitionW_dirt.png": {{}, {}, {}, {}}, "tileGrass_transitionE.png": {{}, {}, {}, {}}, "tileGrass_transitionN.png": {{}, {}, {}, {}}, "tileGrass_transitionS.png": {{}, {}, {}, {}}, "tileGrass_transitionW.png": {{}, {}, {}, {}}, "tileSand1.png": {{}, {}, {}, {}}, "tileSand2.png": {{}, {}, {}, {}}, "tileSand_roadCornerLL.png": {{}, {}, {}, {}}, "tileSand_roadCornerLR.png": {{}, {}, {}, {}}, "tileSand_roadCornerUL.png": {{}, {}, {}, {}}, "tileSand_roadCornerUR.png": {{}, {}, {}, {}}, "tileSand_roadCrossing.png": {{}, {}, {}, {}}, "tileSand_roadCrossingRound.png": {{}, {}, {}, {}}, "tileSand_roadEast.png": {{}, {}, {}, {}}, "tileSand_roadNorth.png": {{}, {}, {}, {}}, "tileSand_roadSplitE.png": {{}, {}, {}, {}}, "tileSand_roadSplitN.png": {{}, {}, {}, {}}, "tileSand_roadSplitS.png": {{}, {}, {}, {}}, "tileSand_roadSplitW.png": {{}, {}, {}, {}},
}

type Tile struct {
	name  string
	image *ebiten.Image
}
type Playfield struct {
	game  *Game
	tiles []*Tile
}

func NewPlayfield(game *Game) *Playfield {
	// Wave function collapse algorithm
	// https://pvs-studio.com/en/blog/posts/csharp/1027/

	var tiles []*Tile

	// setup cells with all the options
	var initialOptions []string

	for k, r := range rules {
		if len(r[ruleUP]) > 0 && len(r[ruleRIGHT]) > 0 && len(r[ruleDOWN]) > 0 && len(r[ruleLEFT]) > 0 {
			fmt.Println(r)
			initialOptions = append(initialOptions, k)
		}
	}

	tilesX := game.width/tileWidth + 1
	tilesY := game.height/tileHeight + 1

	var cells [][]string

	for j := 0; j < tilesY; j++ {
		for i := 0; i < tilesX; i++ {
			cells = append(cells, initialOptions)
		}
	}

	// collapse a random cell with random options
	index := rand.Intn(len(cells))
	cells[index] = []string{cells[index][rand.Intn(len(cells[index]))]}

	for {
		// update the cells
		for j := 0; j < tilesY; j++ {
			for i := 0; i < tilesX; i++ {
				index = j*tilesX + i
				if len(cells[index]) == 1 { // cell is collapsed
					// Look UP
					if j > 0 {
						upindex := (j-1)*tilesX + i
						cells[upindex] = rules[cells[index][0]][ruleUP]
					}
					// Look RIGHT
					if i < tilesX {
						rightindex := j*tilesX + i + 1
						cells[rightindex] = rules[cells[index][0]][ruleRIGHT]
					}
					// Look DOWN
					if j < tilesY {
						downindex := (j+1)*tilesX + i
						cells[downindex] = rules[cells[index][0]][ruleDOWN]
					}
					// Look LEFT
					if i > 0 {
						leftindex := j*tilesX + i - 1
						cells[leftindex] = rules[cells[index][0]][ruleLEFT]
					}
				}
			}
		}

		// pick least entropy
		entropy := len(cells[0])
		for j := 0; j < tilesY; j++ {
			for i := 0; i < tilesX; i++ {
				index := j*tilesX + i
				if len(cells[index]) > 1 && len(cells[index]) < entropy {
					entropy = len(cells[index])
				}
			}
		}

		if entropy <= 1 {
			break
		}
		// pick cells with least entropy
		var leastEntropyIndexes []int
		for j := 0; j < tilesY; j++ {
			for i := 0; i < tilesX; i++ {
				index := j*tilesX + i
				if len(cells[index]) == entropy {
					leastEntropyIndexes = append(leastEntropyIndexes, index)
				}
			}
		}
		// collapse random cell
		index = leastEntropyIndexes[rand.Intn(len(leastEntropyIndexes))]
		cells[index] = []string{cells[index][rand.Intn(len(cells[index]))]}
	}

	for i := 0; i < len(cells); i++ {
		tiles = append(tiles, &Tile{
			name:  cells[i][0],
			image: game.assets.GetSprite(cells[i][0]),
		})
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
			if x%tileWidth == 0 && y%tileHeight == 0 && i < len(p.tiles) {
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(p.tiles[i].image, ops)
				i++
			}
		}
	}
}
