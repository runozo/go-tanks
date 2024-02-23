package game

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/runozo/go-tanks/assets"
)

const (
	screenWidth  = 800
	screenHeight = 600

	velocitySpeedUpTime = 5 * time.Second
)

type Vector struct {
	X float64
	Y float64
}

// main is the entry point of the program

type Game struct {
	width     int
	height    int
	assets    *assets.Assets
	players   []*Player
	playfield *Playfield
	bullets   []*Bullet

	velocityTimer *Timer
}

func NewGame() *Game {
	g := &Game{
		assets:        assets.NewAssets(),
		velocityTimer: NewTimer(velocitySpeedUpTime),
		width:         screenWidth,
		height:        screenHeight,
	}

	g.players = append(g.players, NewPlayer(g))
	g.playfield = NewPlayfield(g)

	return g
}

func (g *Game) Update() error {

	if g.players[0].shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeyP) {
		g.playfield = NewPlayfield(g)
		g.players[0].shootCooldown.Reset()
	}

	for _, p := range g.players {
		p.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.playfield.Draw(screen)

	for _, p := range g.players {
		p.Draw(screen)
	}

	// text.Draw(screen, fmt.Sprintf("CURSOR KEYS: move tank. SPACE: shoot. T: new random tank"), nil, 10, 10, color.Black)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("CURSOR KEYS: move tank, SPACE: shoot, T: new random tank, P: generate new playfield"))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
