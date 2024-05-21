package game

import (
	_ "image/png"
	"log/slog"
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

type Tile struct {
	name  string
	image *ebiten.Image
}
type Playfield struct {
	game  *Game
	tiles []*Tile
}

func NewPlayfield(game *Game) *Playfield {
	var cells [][]string

	// Wave function collapse algorithm
	// https://pvs-studio.com/en/blog/posts/csharp/1027/

	tilesX := game.width/tileWidth + 1
	tilesY := game.height/tileHeight + 1

	// setup cells with all the options
	var initialOptions []string
	for k, _ := range tileOptions {
		initialOptions = append(initialOptions, k)
	}

	for j := 0; j < tilesY; j++ {
		for i := 0; i < tilesX; i++ {
			cells = append(cells, initialOptions)
		}
	}

	for iterations := 0; ; iterations++ { // update the cells

		// pick the minimum entropy indexes
		minEntropyIndexes := getMinEntropyIndexes(&cells)

		if len(minEntropyIndexes) <= 0 {
			slog.Info("Ended with", "iterations", iterations)
			break
		}

		// collapse random cell with least entropy
		index := minEntropyIndexes[rand.Intn(len(minEntropyIndexes))]
		// fmt.Println("Collapse cell:", index)
		cells[index] = []string{cells[index][rand.Intn(len(cells[index]))]}
		// cnt := 0
		// for i := 0; i < len(cells); i++ {
		// 	if len(cells[i]) == 1 {
		// 		cnt++
		// 	}
		// }
		// slog.Info("Cells", "collapsed", cnt, "of", len(cells))

		for y := 0; y < tilesY; y++ {
			for x := 0; x < tilesX; x++ {
				index := y*tilesX + x
				if len(cells[index]) == 0 {
					// we did not found any options, let's restart
					return NewPlayfield(game)
				}

				if len(cells[index]) > 1 { // if the cell it's not collapsed
					// Look UP
					if y > 0 {
						upindex := (y-1)*tilesX + x
						rules := []int{}
						for i := 0; i < len(cells[upindex]); i++ {
							rule := tileOptions[cells[upindex][i]][ruleDOWN]
							if !intInSlice(rule, rules) {
								rules = append(rules, rule)
							}
						}

						optsup := []string{}
						for k, v := range tileOptions {
							if intInSlice(v[ruleUP], rules) {
								optsup = append(optsup, k)
							}
						}

						cells[index] = filterOptions(cells[index], optsup)
					}

					// Look RIGHT
					if x < tilesX-1 {
						rightindex := y*tilesX + x + 1
						rules := []int{}
						for i := 0; i < len(cells[rightindex]); i++ {
							rule := tileOptions[cells[rightindex][i]][ruleLEFT]
							if !intInSlice(rule, rules) {
								rules = append(rules, rule)
							}
						}

						optright := []string{}
						for k, v := range tileOptions {
							if intInSlice(v[ruleRIGHT], rules) {
								optright = append(optright, k)
							}
						}

						cells[index] = filterOptions(cells[index], optright)
					}

					// Look DOWN
					if y < tilesY-1 {
						downindex := (y+1)*tilesX + x
						rules := []int{}
						for i := 0; i < len(cells[downindex]); i++ {
							rule := tileOptions[cells[downindex][i]][ruleUP]
							if !intInSlice(rule, rules) {
								rules = append(rules, rule)
							}
						}

						optdown := []string{}
						for k, v := range tileOptions {
							if intInSlice(v[ruleDOWN], rules) {
								optdown = append(optdown, k)
							}
						}

						cells[index] = filterOptions(cells[index], optdown)
					}
					// Look LEFT
					if x > 0 {
						leftindex := y*tilesX + x - 1
						rules := []int{}
						for i := 0; i < len(cells[leftindex]); i++ {
							rule := tileOptions[cells[leftindex][i]][ruleRIGHT]
							if !intInSlice(rule, rules) {
								rules = append(rules, rule)
							}
						}

						optleft := []string{}
						for k, v := range tileOptions {
							if intInSlice(v[ruleLEFT], rules) {
								optleft = append(optleft, k)
							}
						}

						cells[index] = filterOptions(cells[index], optleft)
					}
				}
			}
		}
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
