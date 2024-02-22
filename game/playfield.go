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

var tileOptions = map[string][]int{
	"tileGrass1.png":                     {0, 0, 0, 0},
	"tileGrass2.png":                     {0, 0, 0, 0},
	"tileGrass_roadCornerLL.png":         {0, 0, 1, 1},
	"tileGrass_roadCornerLR.png":         {0, 1, 1, 0},
	"tileGrass_roadCornerUL.png":         {1, 0, 0, 1},
	"tileGrass_roadCornerUR.png":         {1, 1, 0, 0},
	"tileGrass_roadCrossing.png":         {1, 1, 1, 1},
	"tileGrass_roadCrossingRound.png":    {1, 1, 1, 1},
	"tileGrass_roadEast.png":             {0, 1, 0, 1},
	"tileGrass_roadNorth.png":            {1, 0, 1, 0},
	"tileGrass_roadSplitE.png":           {1, 1, 1, 0},
	"tileGrass_roadSplitN.png":           {1, 1, 0, 1},
	"tileGrass_roadSplitS.png":           {0, 1, 1, 1},
	"tileGrass_roadSplitW.png":           {1, 0, 1, 1},
	"tileGrass_roadTransitionE.png":      {4, 3, 4, 1},
	"tileGrass_roadTransitionE_dirt.png": {4, 3, 4, 1},
	"tileGrass_roadTransitionN.png":      {3, 6, 1, 6},
	"tileGrass_roadTransitionN_dirt.png": {3, 6, 1, 6},
	"tileGrass_roadTransitionS.png":      {1, 6, 3, 6},
	"tileGrass_roadTransitionS_dirt.png": {1, 6, 3, 6},
	"tileGrass_roadTransitionW.png":      {4, 1, 4, 3},
	"tileGrass_roadTransitionW_dirt.png": {4, 1, 4, 3},
	"tileGrass_transitionE.png":          {4, 2, 4, 0},
	"tileGrass_transitionN.png":          {2, 6, 0, 6},
	"tileGrass_transitionS.png":          {0, 4, 2, 4},
	"tileGrass_transitionW.png":          {4, 0, 4, 2},
	"tileSand1.png":                      {2, 2, 2, 2},
	"tileSand2.png":                      {2, 2, 2, 2},
	"tileSand_roadCornerLL.png":          {2, 2, 3, 3},
	"tileSand_roadCornerLR.png":          {2, 3, 3, 2},
	"tileSand_roadCornerUL.png":          {3, 2, 2, 3},
	"tileSand_roadCornerUR.png":          {3, 3, 2, 2},
	"tileSand_roadCrossing.png":          {3, 3, 3, 3},
	"tileSand_roadCrossingRound.png":     {3, 3, 3, 3},
	"tileSand_roadEast.png":              {2, 3, 2, 3},
	"tileSand_roadNorth.png":             {3, 2, 3, 2},
	"tileSand_roadSplitE.png":            {3, 3, 3, 2},
	"tileSand_roadSplitN.png":            {3, 3, 2, 3},
	"tileSand_roadSplitS.png":            {2, 3, 3, 3},
	"tileSand_roadSplitW.png":            {3, 2, 3, 3},
}

type Tile struct {
	name  string
	image *ebiten.Image
}
type Playfield struct {
	game  *Game
	tiles []*Tile
}

var cells [][]string

func NewPlayfield(game *Game) *Playfield {
	// Wave function collapse algorithm
	// https://pvs-studio.com/en/blog/posts/csharp/1027/

	var tiles []*Tile

	// setup cells with all the options
	var initialOptions []string

	for k, o := range tileOptions {
		if len(o) == 4 {
			// fmt.Println(o)
			initialOptions = append(initialOptions, k)
		}
	}

	tilesX := game.width/tileWidth + 1
	tilesY := game.height/tileHeight + 1

	for j := 0; j < tilesY; j++ {
		for i := 0; i < tilesX; i++ {
			cells = append(cells, initialOptions)
		}
	}

	for i := 0; i < len(cells); i++ {
		tiles = append(tiles, &Tile{
			name:  "tileGrass1.png",
			image: game.assets.GetSprite("tileGrass1.png"),
		})
	}

	// collapse a random cell with random options
	index := rand.Intn(len(cells))
	cells[index] = []string{cells[index][rand.Intn(len(cells[index]))]}

	// var PlayfieldTile = mustLoadImage("png/tileGrass2.png")
	return &Playfield{
		game:  game,
		tiles: tiles,
	}
}

func (p *Playfield) Update() {
	tilesX := p.game.width/tileWidth + 1
	tilesY := p.game.height/tileHeight + 1
	var index int
	// update the cells
	for j := 0; j < tilesY; j++ {
		for i := 0; i < tilesX; i++ {
			index = j*tilesX + i
			if len(cells[index]) == 1 { // cell is collapsed
				// Look UP
				if j > 0 {
					upindex := (j-1)*tilesX + i
					optup := tileOptions[cells[index][0]][ruleUP]
					var options []string
					for k, v := range tileOptions {
						if v[ruleDOWN] == optup {
							options = append(options, k)
						}
					}
					cells[upindex] = options

				}
				// Look RIGHT
				if i < tilesX-1 {
					rightindex := j*tilesX + i + 1
					optright := tileOptions[cells[index][0]][ruleRIGHT]
					var options []string
					for k, v := range tileOptions {
						if v[ruleLEFT] == optright {
							options = append(options, k)
						}
					}
					cells[rightindex] = options
				}
				// Look DOWN
				if j < tilesY-1 {
					downindex := (j+1)*tilesX + i
					optdown := tileOptions[cells[index][0]][ruleDOWN]
					var options []string
					for k, v := range tileOptions {
						if v[ruleUP] == optdown {
							options = append(options, k)
						}
					}
					cells[downindex] = options
				}
				// Look LEFT
				if i > 0 {
					leftindex := j*tilesX + i - 1
					optleft := tileOptions[cells[index][0]][ruleLEFT]
					var options []string
					for k, v := range tileOptions {
						if v[ruleRIGHT] == optleft {
							options = append(options, k)
						}
					}
					cells[leftindex] = options
				}
			}
		}
		var tiles []*Tile
		for i := 0; i < len(cells); i++ {
			tiles = append(tiles, &Tile{
				name:  cells[i][0],
				image: p.game.assets.GetSprite(cells[i][0]),
			})
		}
		p.tiles = tiles
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
		return
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
