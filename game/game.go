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
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/runozo/go-wave-function-collapse/assets"
	"github.com/solarlune/resolv"
)

const (
	screenWidth     = 1960
	screenHeight    = 1088
	fontSizeMedium  = 42
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
	TagObstacle  = resolv.NewTag("Obstacle")
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
	scoreLine     *ScoreLine
	obstacles     []*Obstacle
	fontMedium    *text.GoTextFace
	fontSmall     *text.GoTextFace
	serverAddress string
	netClient     *NetClient
	networkTanks  map[string]*Tank
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
	HELP
	OPTIONS
	INTRO
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
		obstacles: []*Obstacle{},
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
		networkTanks:  make(map[string]*Tank),
		space:         resolv.NewSpace(screenWidth, screenHeight, cellWidth, cellHeight),
		state:         RENDERINGPLAYFIELD,
		debug:         false,
	}

	// create new playfield

	g.playfield = NewPlayfield(g)

	if serverAddress != "" {
		g.netClient = NewNetClient(serverAddress)
	}

	return g
}

func (g *Game) Update() error {
	tps := float64(ebiten.TPS())

	if g.state == RENDERINGPLAYFIELD {
		g.playfield.Update(tps)
	}

	// prepare the game after playfield rendered
	if g.playfield.wfc.IsRendered && g.state == RENDERINGPLAYFIELD {

		// add obstacles
		for range rand.Intn(50) {
			// fmt.Println("obstacle")
			g.obstacles = append(g.obstacles, NewObstacle(g))
		}

		// add players
		g.Tanks = []*Tank{NewRandomTank(g, 0, false)}

		// add enemies
		for i := 0; i < numberOfEnemies; i++ {
			g.Tanks = append(g.Tanks, NewRandomTank(g, 0, true))
		}

		g.scoreLine = NewScoreLine(g)

		g.state = PLAYING
	}

	// recreate playfield
	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.state = RENDERINGPLAYFIELD
		g.space.RemoveAll()
		g.obstacles = []*Obstacle{}
		g.playfield = NewPlayfield(g)
	}

	if g.state == PLAYING {
		// --- LETTURA DATI DI RETE ---
		if g.netClient != nil {
		NetworkLoop:
			for {
				select {
				case payload := <-g.netClient.IncomingTanks:
					// Abbiamo ricevuto un dato!
					remoteTank, exists := g.networkTanks[payload.ID]

					if !exists {
						// È un giocatore nuovo, creiamo un tank per lui!
						// Disabilitiamo l'IA passandogli 'true' o creando una logica apposita per i giocatori remoti
						remoteTank = NewRandomTank(g, payload.Rotation, false)
						remoteTank.IsEnemy = payload.IsEnemy // Sovrascriviamo se necessario
						g.networkTanks[payload.ID] = remoteTank
						g.Tanks = append(g.Tanks, remoteTank) // Aggiungiamolo alla lista per disegnarlo
					}

					// Aggiorniamo posizione e rotazione in base a quanto dice il server
					remoteTank.Object.SetPositionVec(resolv.Vector{X: payload.X, Y: payload.Y})
					remoteTank.Object.Rotate(payload.Rotation)

				default:
					// Niente più messaggi in coda, usciamo dal ciclo
					break NetworkLoop
				}
			}
		}
		// update tanks
		for _, t := range g.Tanks {
			t.Update(tps)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.playfield.Draw(screen)

	if g.state == PLAYING {

		for _, f := range g.obstacles {
			f.Draw(screen)
		}
		for _, e := range g.Tanks {
			e.Draw(screen)
		}
		// Draw bullet above all
		for _, t := range g.Tanks {
			for _, bullet := range t.Bullets {
				bullet.Draw(screen)
			}
		}

		g.scoreLine.Draw(screen)
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

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
