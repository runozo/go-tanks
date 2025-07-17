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
	screenWidth        int
	screenHeight       int
	spriteSheet        *ebiten.Image
	mouseX             int
	mouseY             int
	tileEntries        map[string]assets.TileEntry
	assets             *assets.Assets
	spriteSheetOffsetX int
	spriteSheetOffsetY int
}

func (g *Game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	spriteSheetOptions := &ebiten.DrawImageOptions{}
	spriteSheetOptions.GeoM.Translate(float64(g.spriteSheetOffsetX), float64(g.spriteSheetOffsetY))
	screen.DrawImage(g.spriteSheet, spriteSheetOptions)

	tileName := g.isOverTile()
	if tileName != "" {
		tileEntry := g.tileEntries[tileName]
		text.Draw(screen, fmt.Sprintf(" %s X: %d, Y: %d, W: %d, H: %d", tileName, tileEntry.X, tileEntry.Y, tileEntry.Width, tileEntry.Height), mplusNormalFont, 10, g.spriteSheet.Bounds().Dy()+g.spriteSheetOffsetY+20, color.White)

		tileX := tileEntry.X
		tileY := tileEntry.Y
		tileWidth := tileEntry.Width
		tileHeight := tileEntry.Height
		tileImage := g.assets.GetSprite(tileName)

		imageOptions := &ebiten.DrawImageOptions{}
		scaleFactor := 2
		imageOptions.GeoM.Scale(float64(scaleFactor), float64(scaleFactor))
		imageOptions.GeoM.Translate(float64(tileX-tileWidth/2+g.spriteSheetOffsetX), float64(tileY-tileHeight/2+g.spriteSheetOffsetY))
		screen.DrawImage(tileImage, imageOptions)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

func (g *Game) isOverTile() string {
	g.spriteSheet.Bounds()
	for _, tile := range g.tileEntries {
		if g.mouseX > tile.X+g.spriteSheetOffsetX &&
			g.mouseX < tile.X+tile.Width+g.spriteSheetOffsetX &&
			g.mouseY > tile.Y+g.spriteSheetOffsetY && g.mouseY < tile.Y+tile.Height+g.spriteSheetOffsetY {
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

	maxTileWidth := 0
	maxTileHeight := 0
	for _, tile := range ass.TileEntries {
		if tile.Width > maxTileWidth {
			maxTileWidth = tile.Width
		}

		if tile.Height > maxTileHeight {
			maxTileHeight = tile.Height
		}
	}
	// spriteSheet.SetPivot(float64(maxTileWidth/2), float64(ass.TileEntries["grass"].Height))

	g := &Game{
		screenWidth:        spriteSheet.Bounds().Dx() + maxTileWidth,
		screenHeight:       spriteSheet.Bounds().Dy() + maxTileHeight + 50,
		spriteSheet:        spriteSheet,
		tileEntries:        ass.TileEntries,
		assets:             ass,
		spriteSheetOffsetX: maxTileWidth / 2,
		spriteSheetOffsetY: maxTileHeight / 2,
	}

	ebiten.SetWindowSize(g.screenWidth, g.screenHeight)
	ebiten.SetWindowTitle("Font (Ebitengine Demo)")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
