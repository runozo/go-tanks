package game

import (
	"bytes"
	"embed"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/runozo/go-wave-function-collapse/assets"
)

const (
	screenWidth    = 1960
	screenHeight   = 1088
	fontSizeMedium = 32
	fontSizeSmall  = 22
)

//go:embed assets/*
var assetsFS embed.FS

type Vector struct {
	X float64
	Y float64
}

type Game struct {
	width      int
	height     int
	assets     *assets.Assets
	players    []*Player
	playfield  *Playfield
	fontMedium *text.GoTextFace
	fontSmall  *text.GoTextFace
}

func NewGame() *Game {
	// Load sprite sheet
	spriteSheetData, err := assetsFS.ReadFile("assets/allSprites_default.png")
	if err != nil {
		log.Fatal(err)
	}

	// Load tile map
	jsonData, err := assetsFS.ReadFile("assets/mapped_tiles.json")
	if err != nil {
		log.Fatal(err)
	}
	ass := assets.NewAssets(
		spriteSheetData,
		jsonData,
	)

	// Load fonts
	fontData, err := assetsFS.ReadFile("assets/gomarice_no_continue.ttf")
	if err != nil {
		log.Fatal(err)
	}

	textFS, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	g := &Game{
		assets:    ass,
		width:     screenWidth,
		height:    screenHeight,
		playfield: nil,
		fontMedium: &text.GoTextFace{
			Source:    textFS,
			Direction: text.DirectionLeftToRight,
			Size:      fontSizeMedium,
		},
		fontSmall: &text.GoTextFace{
			Source:    textFS,
			Direction: text.DirectionLeftToRight,
			Size:      fontSizeSmall,
		},
	}

	g.players = append(g.players, NewPlayer(g))
	g.playfield = NewPlayfield(g)

	return g
}

func (g *Game) Update() error {
	tps := float64(ebiten.TPS())
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.playfield = NewPlayfield(g)
	}
	g.playfield.Update(tps)
	for _, p := range g.players {
		p.Update(tps)
	}
	// fmt.Println(g.playfield.tiles)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.playfield.Draw(screen)

	for _, p := range g.players {
		p.Draw(screen)
	}

	str := "CURSOR KEYS: move tank, A/D: rotate barrel, SPACE: shoot, T: new random tank, P: generate new playfield"
	width, height := text.Measure(str, g.fontMedium, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64((screen.Bounds().Dx()-int(width))/2), float64(height))
	op.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, str, g.fontMedium, op)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
