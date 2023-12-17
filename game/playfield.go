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

type Tile struct {
	name         string
	image        *ebiten.Image
	upAllowed    []string
	leftAllowed  []string
	downAllowed  []string
	rightAllowed []string
}

var allTiles = []Tile{
	Tile{
		name:        "tileGrass_roadCornerLL.png",
		leftAllowed: []string{"tileGrass_roadSplitN.png", "tileGrass_roadCrossing.png", "11 tileGrass_roadCornerLR.png", "tileGrass_roadSplitS.png", "tileGrass_roadTransitionW_dirt.png"},
	},
	Tile{
		name:         "tileGrass_roadCornerLR.png",
		rightAllowed: []string{"tileGrass_roadCornerUL.png"},
	},
	Tile{
		name: "tileGrass_roadCornerUL.png",
	},
	Tile{
		name: "tileGrass_roadCornerUR.png",
	},
	Tile{
		name: "tileGrass_roadCrossing.png",
	},
	Tile{
		name: "tileGrass_roadCrossingRound.png",
	},
	Tile{
		name: "tileGrass_roadEast.png",
	},
	Tile{
		name:        "tileGrass_roadNorth.png",
		upAllowed:   []string{"tileGrass_roadCrossing.png", "tileGrass_roadCrossing.png"},
		downAllowed: []string{"tileGrass_roadCrossing.png", "tileGrass_roadCrossing.png"},
	},
	Tile{
		name: "tileGrass_roadSplitE.png",
	},
	Tile{
		name:         "tileGrass_roadSplitN.png",
		leftAllowed:  []string{"tileGrass_roadSplitE.png", "tileGrass_roadTransitionE.png"},
		rightAllowed: []string{},
	},
	Tile{
		name: "tileGrass_roadSplitS.png",
	},
	Tile{
		name: "tileGrass_roadSplitW.png",
	},
	Tile{
		name: "tileGrass_roadTransitionE.png",
	},
	Tile{
		name: "tileGrass_roadTransitionE_dirt.png",
	},
	Tile{
		name: "tileGrass_roadTransitionN.png",
	},
	Tile{
		name: "tileGrass_roadTransitionN_dirt.png",
	},
	Tile{
		name: "tileGrass_roadTransitionS.png",
	},
	Tile{
		name: "tileGrass_roadTransitionS_dirt.png",
	},
	Tile{
		name: "tileGrass_roadTransitionW.png",
	},
	Tile{
		name: "tileGrass_roadTransitionW_dirt.png",
	},
	Tile{
		name: "tileGrass_transitionE.png",
	},
	Tile{
		name: "tileGrass_transitionN.png",
	},
	Tile{
		name: "tileGrass_transitionS.png",
	},
	Tile{
		name: "tileGrass_transitionW.png",
	},
	Tile{
		name:         "tileSand1.png",
		leftAllowed:  []string{"tileGrass_transitionE.png", "tileSand1.png", "tileSand2.png"},
		upAllowed:    []string{"tileGrass_transitionN.png", "tileSand1.png", "tileSand2.png"},
		downAllowed:  []string{"tileGrass_transitionS.png", "tileSand1.png", "tileSand2.png"},
		rightAllowed: []string{"tileGrass_transitionW.png", "tileSand1.png", "tileSand2.png"},
	},
	Tile{
		name:         "tileSand2.png",
		leftAllowed:  []string{"tileGrass_transitionE.png", "tileSand1.png", "tileSand2.png"},
		upAllowed:    []string{"tileGrass_transitionN.png", "tileSand1.png", "tileSand2.png"},
		downAllowed:  []string{"tileGrass_transitionS.png", "tileSand1.png", "tileSand2.png"},
		rightAllowed: []string{"tileGrass_transitionW.png", "tileSand1.png", "tileSand2.png"},
	},
	Tile{
		name: "tileSand_roadCornerLL.png",
	},
	Tile{
		name: "tileSand_roadCornerLR.png",
	},
	Tile{
		name: "tileSand_roadCornerUL.png",
	},
	Tile{
		name: "tileSand_roadCornerUR.png",
	},
	Tile{
		name: "tileSand_roadCrossing.png",
	},
	Tile{
		name: "tileSand_roadCrossingRound.png",
	},
	Tile{
		name: "tileSand_roadEast.png",
	},
	Tile{
		name:      "tileSand_roadNorth.png",
		upAllowed: []string{"tileSand_roadSplitS.png"},
	},
	Tile{
		name:         "tileSand_roadSplitE.png",
		rightAllowed: []string{"tileGrass_roadSplitS.png"},
	},
	Tile{
		name:        "tileSand_roadSplitN.png",
		upAllowed:   []string{"tileSand_roadSplitS.png"},
		leftAllowed: []string{"tileGrass_roadTransitionW_dirt.png"},
	},

	Tile{
		name:        "tileSand_roadSplitS.png",
		downAllowed: []string{"tileSand_roadSplitN.png"},
	},
	Tile{
		name: "tileSand_roadSplitW.png",
	},
}

type Playfield struct {
	game  *Game
	tiles []*Tile
}

func NewPlayfield(game *Game) *Playfield {
	var tiles []*Tile
	var i int
	for y := 0; y < game.height; y += tileHeight {
		for x := 0; x < game.width; x += tileWidth {
			if x%tileWidth == 0 && y%tileHeight == 0 {
				tile := allTiles[rand.Intn(len(allTiles))]
				tile.image = game.assets.GetSprite(tile.name)
				tiles = append(tiles, &tile)
				fmt.Println(i, tile.name)
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
