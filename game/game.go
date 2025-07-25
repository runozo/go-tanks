package game

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/runozo/go-wave-function-collapse/assets"
)

const (
	screenWidth  = 1960
	screenHeight = 1088
)

type Vector struct {
	X float64
	Y float64
}

type Game struct {
	width     int
	height    int
	assets    *assets.Assets
	players   []*Player
	playfield *Playfield
	// bullets   []*Bullet

	// velocityTimer *Timer
}

func NewGame() *Game {
	// ebiten.SetWindowSize(screenWidth, screenHeight)
	// ebiten.SetFullscreen(true)
	ass := assets.NewAssets(
		"data"+string(os.PathSeparator)+"allSprites_default.png",
		"data"+string(os.PathSeparator)+"mapped_tiles.json",
	)
	g := &Game{
		assets: ass,
		// velocityTimer: NewTimer(velocitySpeedUpTime),
		width:     screenWidth,
		height:    screenHeight,
		playfield: NewPlayfield(screenWidth, screenHeight, ass),
	}

	g.players = append(g.players, NewPlayer(g))

	return g
}

func (g *Game) Update() error {
	tps := float64(ebiten.TPS())
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.playfield = NewPlayfield(screenWidth, screenHeight, g.assets)
	}
	g.playfield.Update(tps)
	for _, p := range g.players {
		p.Update(tps)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.playfield.Draw(screen)

	for _, p := range g.players {
		p.Draw(screen)
	}

	// text.Draw(screen, fmt.Sprintf("CURSOR KEYS: move tank. SPACE: shoot. T: new random tank"), nil, 10, 10, color.Black)
	ebitenutil.DebugPrint(screen, "CURSOR KEYS: move tank, A/D: rotate barrel, SPACE: shoot, T: new random tank, P: generate new playfield")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
