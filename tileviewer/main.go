package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/runozo/go-wave-function-collapse/assets"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var mplusNormalFont font.Face

type Icon struct {
	Name string
	Img  *ebiten.Image
}

type Game struct {
	screenWidth  int
	screenHeight int
	spriteSheet  *ebiten.Image
	mouseX       int
	mouseY       int
	tileEntries  map[string]assets.TileEntry
	assets       *assets.Assets
}

func (g *Game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.spriteSheet, nil)

	tileName := g.isOverTile()
	if tileName != "" {
		tileEntry := g.tileEntries[tileName]
		text.Draw(screen, fmt.Sprintf(" %s X: %d, Y: %d, W: %d, H: %d", tileName, tileEntry.X, tileEntry.Y, tileEntry.Width, tileEntry.Height), mplusNormalFont, 10, g.spriteSheet.Bounds().Dy()+20, color.White)

		tileX := tileEntry.X
		tileY := tileEntry.Y
		// tileWidth := tileEntry.Width
		// tileHeight := tileEntry.Height
		tileImage := g.assets.GetSprite(tileName)

		imageOptions := &ebiten.DrawImageOptions{}
		imageOptions.GeoM.Scale(float64(2), float64(2))
		imageOptions.GeoM.Translate(float64(tileX), float64(tileY))
		screen.DrawImage(tileImage, imageOptions)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

func (g *Game) isOverTile() string {
	for _, tile := range g.tileEntries {
		if g.mouseX > tile.X && g.mouseX < tile.X+tile.Width && g.mouseY > tile.Y && g.mouseY < tile.Y+tile.Height {
			return tile.Name
		}
	}
	return ""
}

func main() {
	spriteSheetPath := ".." + string(os.PathSeparator) + "data" + string(os.PathSeparator) + "allSprites_default.png"
	spriteMapPath := ".." + string(os.PathSeparator) + "data" + string(os.PathSeparator) + "mapped_tiles.json"

	dat, err := os.Open(spriteSheetPath)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(dat)
	if err != nil {
		log.Fatal(err)
	}

	spriteSheet := ebiten.NewImageFromImage(img)

	ass := assets.NewAssets(
		spriteSheetPath,
		spriteMapPath,
	)

	// Load fonts
	const dpi = 72
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	g := &Game{
		screenWidth:  spriteSheet.Bounds().Dx(),
		screenHeight: spriteSheet.Bounds().Dy() + 50,
		spriteSheet:  spriteSheet,
		tileEntries:  ass.TileEntries,
		assets:       ass,
	}

	ebiten.SetWindowSize(g.screenWidth, g.screenHeight)
	ebiten.SetWindowTitle("Font (Ebitengine Demo)")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
