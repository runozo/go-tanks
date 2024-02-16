package main

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	spriteSheet  = "../assets/spritesheet/allSprites_default.png"
	xmlSpriteMap = "../assets/spritesheet/allSprites_default.xml"
	screenWidth  = 800
	screenHeight = 800
	tileWidth    = 64
	tileHeight   = 64
)

var mplusNormalFont font.Face
var icons []*ebiten.Image

type SubTexture struct {
	XMLName xml.Name `xml:"SubTexture"`
	Name    string   `xml:"name,attr"`
	X       int      `xml:"x,attr"`
	Y       int      `xml:"y,attr"`
	Width   int      `xml:"width,attr"`
	Height  int      `xml:"height,attr"`
}

type SpriteMap struct {
	XMLName     xml.Name     `xml:"TextureAtlas"`
	ImagePath   string       `xml:"imagePath,attr"`
	SubTextures []SubTexture `xml:"SubTexture"`
}

func init() {
	origin := mustLoadImage(spriteSheet)
	xml := mustLoadXMLSpriteMap(xmlSpriteMap)

	for i := 0; i < len(xml.SubTextures); i++ {
		subTexture := xml.SubTextures[i]
		icons = append(icons, origin.SubImage(
			image.Rect(
				subTexture.X,
				subTexture.Y,
				subTexture.X+subTexture.Width,
				subTexture.Y+subTexture.Height,
			)).(*ebiten.Image))
	}

	// Load fonts
	const dpi = 72
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found", len(icons), "icons")
}

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	const x = 20
	// Draw info
	msg := fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS())
	text.Draw(screen, msg, mplusNormalFont, x, 20, color.White)
	var i int
	for y := 0; y < screenHeight; y += tileHeight {
		for x := 0; x < screenWidth; x += tileWidth {
			if x%tileWidth == 0 && y%tileHeight == 0 {
				ops := &ebiten.DrawImageOptions{}
				ops.GeoM.Translate(float64(x), float64(y))
				screen.DrawImage(icons[i], ops)
				i++
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Font (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func mustLoadXMLSpriteMap(name string) SpriteMap {
	byteValue, _ := os.ReadFile(name)
	var s SpriteMap
	if err := xml.Unmarshal(byteValue, &s); err != nil {
		panic(err)
	}
	// fmt.Println(s)
	return s
}

func mustLoadImage(name string) *ebiten.Image {
	fmt.Println(name)
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}
