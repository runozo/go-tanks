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

var tileOptions = map[string][]int{
	"tileGrass1.png":                  {0, 0, 0, 0}, // 0 grass
	"tileGrass2.png":                  {0, 0, 0, 0},
	"tileGrass_roadCornerLL.png":      {0, 0, 1, 1}, // 1 road with grass
	"tileGrass_roadCornerLR.png":      {0, 1, 1, 0},
	"tileGrass_roadCornerUL.png":      {1, 0, 0, 1},
	"tileGrass_roadCornerUR.png":      {1, 1, 0, 0},
	"tileGrass_roadCrossing.png":      {1, 1, 1, 1},
	"tileGrass_roadCrossingRound.png": {1, 1, 1, 1},
	"tileGrass_roadEast.png":          {0, 1, 0, 1},
	"tileGrass_roadNorth.png":         {1, 0, 1, 0},
	"tileGrass_roadSplitE.png":        {1, 1, 1, 0},
	"tileGrass_roadSplitN.png":        {1, 1, 0, 1},
	"tileGrass_roadSplitS.png":        {0, 1, 1, 1},
	"tileGrass_roadSplitW.png":        {1, 0, 1, 1},
	// "tileGrass_roadTransitionE.png":      {4, 3, 4, 1}, //
	// "tileGrass_roadTransitionE_dirt.png": {4, 3, 4, 1},
	// "tileGrass_roadTransitionN.png":      {3, 6, 1, 6},
	// "tileGrass_roadTransitionN_dirt.png": {3, 6, 1, 6},
	// "tileGrass_roadTransitionS.png":      {1, 6, 3, 6},
	// "tileGrass_roadTransitionS_dirt.png": {1, 6, 3, 6},
	// "tileGrass_roadTransitionW.png":      {4, 1, 4, 3},
	// "tileGrass_roadTransitionW_dirt.png": {4, 1, 4, 3},
	// "tileGrass_transitionE.png":          {4, 2, 4, 0},
	// "tileGrass_transitionN.png":          {2, 6, 0, 6},
	// "tileGrass_transitionS.png":          {0, 4, 2, 4},
	// "tileGrass_transitionW.png":          {4, 0, 4, 2},
	// "tileSand1.png":                  {2, 2, 2, 2},
	// "tileSand2.png":                  {2, 2, 2, 2},
	// "tileSand_roadCornerLL.png":      {2, 2, 3, 3},
	// "tileSand_roadCornerLR.png":      {2, 3, 3, 2},
	// "tileSand_roadCornerUL.png":      {3, 2, 2, 3},
	// "tileSand_roadCornerUR.png":      {3, 3, 2, 2},
	// "tileSand_roadCrossing.png":      {3, 3, 3, 3},
	// "tileSand_roadCrossingRound.png": {3, 3, 3, 3},
	// "tileSand_roadEast.png":          {2, 3, 2, 3},
	// "tileSand_roadNorth.png":         {3, 2, 3, 2},
	// "tileSand_roadSplitE.png":        {3, 3, 3, 2},
	// "tileSand_roadSplitN.png":        {3, 3, 2, 3},
	// "tileSand_roadSplitS.png":        {2, 3, 3, 3},
	// "tileSand_roadSplitW.png":        {3, 2, 3, 3},
}

type Tile struct {
	name  string
	image *ebiten.Image
}
type Playfield struct {
	game  *Game
	tiles []*Tile
}

func filterOptions(orig, options []string) []string {
	var filtered []string
	for _, o := range options {
		for _, o2 := range orig {
			if o2 == o {
				filtered = append(filtered, o2)
			}
		}
	}
	return filtered
}

func NewPlayfield(game *Game) *Playfield {
	var cells [][]string

	// Wave function collapse algorithm
	// https://pvs-studio.com/en/blog/posts/csharp/1027/

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

	// collapse a random cell with random option
	index := rand.Intn(len(cells))
	cells[index] = []string{cells[index][rand.Intn(len(cells[index]))]}

	for u := 0; ; u++ { // update the cells
		for y := 0; y < tilesY; y++ {
			for x := 0; x < tilesX; x++ {
				index := y*tilesX + x
				if len(cells[index]) > 0 {
					// Look UP
					if y > 0 {
						upindex := (y-1)*tilesX + x
						if len(cells[upindex]) > 1 { // if it's not collapsed
							optup := tileOptions[cells[index][0]][ruleUP]
							var options []string
							for k, v := range tileOptions {
								if v[ruleDOWN] == optup {
									options = append(options, k)
								}
							}
							cells[upindex] = filterOptions(cells[upindex], options)
						}

					}
					// Look RIGHT
					if x < tilesX-1 {
						rightindex := y*tilesX + x + 1
						if len(cells[rightindex]) > 1 { // if it's not collapsed
							optright := tileOptions[cells[index][0]][ruleRIGHT]
							var options []string
							for k, v := range tileOptions {
								if v[ruleLEFT] == optright {
									options = append(options, k)
								}
							}
							cells[rightindex] = filterOptions(cells[rightindex], options)
						}
					}
					// Look DOWN
					if y < tilesY-1 {
						downindex := (y+1)*tilesX + x
						if len(cells[downindex]) > 1 { // if it's not collapsed
							optdown := tileOptions[cells[index][0]][ruleDOWN]
							var options []string
							for k, v := range tileOptions {
								if v[ruleUP] == optdown {
									options = append(options, k)
								}
							}
							cells[downindex] = filterOptions(cells[downindex], options)
						}
					}
					// Look LEFT
					if x > 0 {
						leftindex := y*tilesX + x - 1
						if len(cells[leftindex]) > 1 { // if it's not collapsed
							optleft := tileOptions[cells[index][0]][ruleLEFT]
							var options []string
							for k, v := range tileOptions {
								if v[ruleRIGHT] == optleft {
									options = append(options, k)
								}
							}
							cells[leftindex] = filterOptions(cells[leftindex], options)
						}
					}
				}
			}
		}

		// pick least leastEntropy
		leastEntropy := len(initialOptions)
		for i := 0; i < len(cells); i++ {
			if len(cells[i]) > 1 && len(cells[i]) < leastEntropy {
				leastEntropy = len(cells[i])
			}
		}
		fmt.Println("Least entropy:", leastEntropy)

		// pick cells with least entropy
		var leastEntropyIndexes []int
		for i := 0; i < len(cells); i++ {
			if len(cells[i]) == leastEntropy {
				leastEntropyIndexes = append(leastEntropyIndexes, i)
			}

		}

		if len(leastEntropyIndexes) <= 0 {
			fmt.Println("Ended with", u, "iterations")
			break
		} else {
			fmt.Println("leastEntropyIndexes", leastEntropyIndexes)
		}

		// collapse random cell
		index = leastEntropyIndexes[rand.Intn(len(leastEntropyIndexes))]
		fmt.Println("Collapse cell:", index)
		cells[index] = []string{cells[index][rand.Intn(len(cells[index]))]}
		cnt := 0
		for i := 0; i < len(cells); i++ {
			if len(cells[i]) == 1 {
				cnt++
			}
		}
		fmt.Println("Total of collapsed:", cnt, "of", len(cells))

	}

	var tiles []*Tile
	for i := 0; i < len(cells); i++ {
		if len(cells[i]) >= 1 {

			tiles = append(tiles, &Tile{
				name:  cells[i][0],
				image: game.assets.GetSprite(cells[i][0]),
			})
		}
	}

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
