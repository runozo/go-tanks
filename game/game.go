package game

import (
	"bytes"
	"embed"
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/runozo/go-wave-function-collapse/assets"
	"github.com/solarlune/resolv"
)

const (
	screenWidth     = 1960
	screenHeight    = 1088
	fontSizeMedium  = 32
	fontSizeSmall   = 22
	numberOfEnemies = 5 // *3
	cellWidth       = 64
	cellHeight      = 64
)

//go:embed assets/*
var assetsFS embed.FS

var (
	TagPlayer    = resolv.NewTag("Player")
	TagSolidWall = resolv.NewTag("SolidWall")
	TagEnemy     = resolv.NewTag("Enemy")
)

type Vector struct {
	X float64
	Y float64
}

type Game struct {
	width         int
	height        int
	assets        *assets.Assets
	players       []*Player
	enemies       []*Enemy
	playfield     *Playfield
	fontMedium    *text.GoTextFace
	fontSmall     *text.GoTextFace
	serverAddress string
	space         *resolv.Space
}

func NewGame(serverAddress string) *Game {
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
	assets := assets.NewAssets(
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
		assets:    assets,
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
		serverAddress: serverAddress,
		space:         resolv.NewSpace(screenWidth, screenHeight, cellWidth, cellHeight),
	}

	// create new playfield

	g.playfield = NewPlayfield(g)

	// add players
	playerSolids := resolv.ShapeCollection{}
	p := NewPlayer(g)
	g.players = append(g.players, p)
	playerSolids = append(playerSolids, p.tank.solid)
	playerSolids.SetTags(TagPlayer)
	g.space.Add(playerSolids...)

	// add enemies

	typesOfEnemies := []string{"hard", "medium", "easy"}
	enemySolids := resolv.ShapeCollection{}
	// add enemies one at a time so they don't overlap
	for i := 0; i < numberOfEnemies; i++ {
		index := rand.Intn(len(typesOfEnemies))
		e := NewEnemy(g, typesOfEnemies[index]) // random enemy
		g.enemies = append(g.enemies, e)
		enemySolids = append(enemySolids, e.tank.solid)
	}

	enemySolids.SetTags(TagEnemy)

	g.space.Add(enemySolids...)

	return g
}

func (g *Game) Update() error {
	tps := float64(ebiten.TPS())

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.playfield = NewPlayfield(g)
	}

	g.playfield.Update(tps)

	if !g.playfield.wfc.IsRunning {
		for _, e := range g.enemies {
			e.Update(tps)
		}

		for _, p := range g.players {
			p.Update(tps)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.playfield.Draw(screen)
	if !g.playfield.wfc.IsRunning {
		for _, e := range g.enemies {
			e.Draw(screen)
		}
		for _, p := range g.players {
			p.Draw(screen)
		}
		// Draw bullet above all
		for _, e := range g.enemies {
			for _, bullet := range e.bullets {
				bullet.Draw(screen)
			}
		}
		for _, p := range g.players {
			for _, bullet := range p.Bullets {
				bullet.Draw(screen)
			}
		}
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
