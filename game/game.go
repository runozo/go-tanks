package game

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/runozo/go-wave-function-collapse/assets"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth   = 1960
	screenHeight  = 1088
	fontDpi       = 72
	fontSizeSmall = 22
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
	fontSmall font.Face
	// bullets   []*Bullet

	// velocityTimer *Timer
}

func NewGame() *Game {
	// ebiten.SetWindowSize(screenWidth, screenHeight)
	// ebiten.SetFullscreen(true)
	ass := assets.NewAssets(
		"assets"+string(os.PathSeparator)+"allSprites_default.png",
		"assets"+string(os.PathSeparator)+"mapped_tiles.json",
	)

	// Load fonts
	fontData, err := os.ReadFile("assets" + string(os.PathSeparator) + "gomarice_no_continue.ttf")
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	gameFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSizeSmall,
		DPI:     fontDpi,
		Hinting: font.HintingVertical,
	})

	if err != nil {
		log.Fatal(err)
	}

	g := &Game{
		assets:    ass,
		width:     screenWidth,
		height:    screenHeight,
		playfield: NewPlayfield(screenWidth, screenHeight, ass),
		fontSmall: gameFont,
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

	text.Draw(screen, "CURSOR KEYS: move tank, A/D: rotate barrel, SPACE: shoot, T: new random tank, P: generate new playfield", g.fontSmall, 10, 20, color.Black)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
