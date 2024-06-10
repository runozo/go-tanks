package game

import (
	_ "image/png"
	"log/slog"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/runozo/go-tanks/assets"
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

func resetTilesOptions(tiles *[]Tile) {
	// create a slice of all the options available
	initialOptions := make([]string, len(tileOptions))
	i := 0
	for k := range tileOptions {
		initialOptions[i] = k
		i++
	}

	// setup tiles with all the options enabled
	for i := 0; i < len(*tiles); i++ {
		(*tiles)[i].options = initialOptions
		(*tiles)[i].image = ebiten.NewImage(tileWidth, tileHeight)
		(*tiles)[i].collapsed = false
	}
}

func NewPlayfield(width, height int, assets *assets.Assets) *Playfield {

	tilesX := width/tileWidth + 1
	tilesY := height/tileHeight + 1

	// setup tiles with all the options enabled
	tiles := make([]Tile, tilesX*tilesY)
	resetTilesOptions(&tiles)

	return &Playfield{
		width:       width,
		height:      height,
		tiles:       tiles,
		isRendered:  false,
		numOfTilesX: tilesX,
		numOfTilesY: tilesY,
		assets:      assets,
	}
}

func (p *Playfield) Update() {

	if !p.isRendered {
		// pick the minimum entropy indexes
		minEntropyIndexes := getMinEntropyIndexes(&p.tiles)

		if len(minEntropyIndexes) <= 0 {
			slog.Info("Playfiled is rendered. No more collapsable cells.")
			for i := 0; i < len(p.tiles); i++ {
				if !p.tiles[i].collapsed {
					p.tiles[i].image = p.assets.GetSprite(p.tiles[i].options[0])
					p.tiles[i].collapsed = true
				}
			}
			p.isRendered = true

		} else {

			collapsedIndex := collapseRandomCellWithMinEntropy(&p.tiles, &minEntropyIndexes)
			// slog.Info("Collapsed", "index", collapsedIndex, "name", p.tiles[collapsedIndex].options[0])
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
							upindex := (y-1)*p.numOfTilesX + x
							rules := []int{}
							for i := 0; i < len(p.tiles[upindex].options); i++ {
								rule := tileOptions[p.tiles[upindex].options[i]][ruleDOWN]
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

							p.tiles[index].options = filterOptions(p.tiles[index].options, optsup)
						}

						// Look RIGHT
						if x < p.numOfTilesY-1 {
							rightindex := y*p.numOfTilesX + x + 1
							rules := []int{}
							for i := 0; i < len(p.tiles[rightindex].options); i++ {
								rule := tileOptions[p.tiles[rightindex].options[i]][ruleLEFT]
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

							p.tiles[index].options = filterOptions(p.tiles[index].options, optright)
						}

						// Look DOWN
						if y < p.numOfTilesY-1 {
							downindex := (y+1)*p.numOfTilesX + x
							rules := []int{}
							for i := 0; i < len(p.tiles[downindex].options); i++ {
								rule := tileOptions[p.tiles[downindex].options[i]][ruleUP]
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

							p.tiles[index].options = filterOptions(p.tiles[index].options, optdown)
						}
						// Look LEFT
						if x > 0 {
							leftindex := y*p.numOfTilesX + x - 1
							rules := []int{}
							for i := 0; i < len(p.tiles[leftindex].options); i++ {
								rule := tileOptions[p.tiles[leftindex].options[i]][ruleRIGHT]
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

							p.tiles[index].options = filterOptions(p.tiles[index].options, optleft)
						}
					}
				}

			}
		}

	}
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
