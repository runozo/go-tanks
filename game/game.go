package game

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/runozo/go-tanks/assets"
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
	g := &Game{
		assets: assets.NewAssets(),
		// velocityTimer: NewTimer(velocitySpeedUpTime),
		width:  screenWidth,
		height: screenHeight,
	}

	g.players = append(g.players, NewPlayer(g))

	start := time.Now()
	g.playfield = NewPlayfield(g)
	elapsed := time.Since(start)
	slog.Info("Time", "seconds", elapsed)

	return g
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		start := time.Now()
		g.playfield = NewPlayfield(g)
		elapsed := time.Since(start)
		slog.Info("Time", "seconds", elapsed)
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
