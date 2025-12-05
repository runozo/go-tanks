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
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	TagBarrel    = resolv.NewTag("Barrel")
	TagEnemy     = resolv.NewTag("Enemy")
	TagBullet    = resolv.NewTag("Bullet")
	TagExplosion = resolv.NewTag("Explosion")
)

func FlipVertical(source *ebiten.Image) *ebiten.Image {
	flipped := ebiten.NewImage(source.Bounds().Dx(), source.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(source.Bounds().Dy()))
	flipped.DrawImage(source, op)
	return flipped
}

type Game struct {
	width         int
	height        int
	assets        *assets.Assets
	Tanks         []*Tank
	playfield     *Playfield
	fontMedium    *text.GoTextFace
	fontSmall     *text.GoTextFace
	serverAddress string
	space         *resolv.Space
	state         int
	debug         bool
}

// game state
const (
	PLAYING = iota
	PAUSED
	RENDERINGPLAYFIELD
	GAMEOVER
)

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
		state:         RENDERINGPLAYFIELD,
		debug:         false,
	}

	// create new playfield

	g.playfield = NewPlayfield(g)

	return g
}

func (g *Game) Update() error {
	tps := float64(ebiten.TPS())

	if g.state == RENDERINGPLAYFIELD {
		g.playfield.Update(tps)
	}

	// prepare the game after playfield rendered
	if g.playfield.wfc.IsRendered && g.state == RENDERINGPLAYFIELD {
		// add players
		g.Tanks = []*Tank{NewRandomTank(g, 0, false)}
		// add enemies

		for i := 0; i < numberOfEnemies; i++ {
			g.Tanks = append(g.Tanks, NewRandomTank(g, 0, true))
		}

		g.state = PLAYING
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.state = RENDERINGPLAYFIELD
		g.playfield = NewPlayfield(g)
	}

	if g.state == PLAYING {
		for _, t := range g.Tanks {
			t.Update(tps)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.playfield.Draw(screen)
	if !g.playfield.wfc.IsRunning {
		for _, e := range g.Tanks {
			e.Draw(screen)
		}
		// Draw bullet above all
		for _, t := range g.Tanks {
			for _, bullet := range t.Bullets {
				bullet.Draw(screen)
			}
		}
	}

	// debug shapes
	if g.debug {
		g.space.ForEachShape(func(shape resolv.IShape, index, maxCount int) bool {

			var drawColor color.Color = color.White

			tags := shape.Tags()

			if tags.Has(TagEnemy) && !tags.Has(TagBarrel) {
				drawColor = color.RGBA{255, 128, 35, 255}
			}
			if tags.Has(TagPlayer) {
				drawColor = color.RGBA{32, 255, 128, 255}
			}
			switch o := shape.(type) {
			case *resolv.Circle:
				vector.StrokeCircle(screen, float32(o.Position().X), float32(o.Position().Y), float32(o.Radius()), 2, drawColor, false)
			case *resolv.ConvexPolygon:

				for _, l := range o.Lines() {
					vector.StrokeLine(screen, float32(l.Start.X), float32(l.Start.Y), float32(l.End.X), float32(l.End.Y), 2, drawColor, false)
				}
			}

			return true

		})
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
