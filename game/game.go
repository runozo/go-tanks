package game

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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

	for _, b := range g.bullets {
		b.Draw(screen)
	}

	// text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, screenWidth/2-100, 50, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
