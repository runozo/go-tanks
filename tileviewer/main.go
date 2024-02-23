package main

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	spriteSheet  = "../assets/spritesheet/allSprites_default.png"
	xmlSpriteMap = "../assets/spritesheet/allSprites_default.xml"
	screenWidth  = 1024
	screenHeight = 900
	tileWidth    = 64
	tileHeight   = 64
)

var tileOptions = map[string][]int{
	"tileGrass1.png":                     {0, 0, 0, 0}, // 0 grass
	"tileGrass2.png":                     {0, 0, 0, 0},
	"tileGrass_roadCornerLL.png":         {0, 0, 1, 1}, // 1 road with grass
	"tileGrass_roadCornerLR.png":         {0, 1, 1, 0},
	"tileGrass_roadCornerUL.png":         {1, 0, 0, 1},
	"tileGrass_roadCornerUR.png":         {1, 1, 0, 0},
	"tileGrass_roadCrossing.png":         {1, 1, 1, 1},
	"tileGrass_roadCrossingRound.png":    {1, 1, 1, 1},
	"tileGrass_roadEast.png":             {0, 1, 0, 1},
	"tileGrass_roadNorth.png":            {1, 0, 1, 0},
	"tileGrass_roadSplitE.png":           {1, 1, 1, 0},
	"tileGrass_roadSplitN.png":           {1, 1, 0, 1},
	"tileGrass_roadSplitS.png":           {0, 1, 1, 1},
	"tileGrass_roadSplitW.png":           {1, 0, 1, 1},
	"tileGrass_roadTransitionE.png":      {4, 3, 4, 1}, //
	"tileGrass_roadTransitionE_dirt.png": {4, 3, 4, 1},
	"tileGrass_roadTransitionN.png":      {3, 6, 1, 6},
	"tileGrass_roadTransitionN_dirt.png": {3, 6, 1, 6},
	"tileGrass_roadTransitionS.png":      {1, 8, 3, 8},
	"tileGrass_roadTransitionS_dirt.png": {1, 8, 3, 8},
	"tileGrass_roadTransitionW.png":      {5, 1, 5, 3},
	"tileGrass_roadTransitionW_dirt.png": {5, 1, 5, 3},
	"tileGrass_transitionE.png":          {4, 2, 4, 0},
	"tileGrass_transitionN.png":          {2, 6, 0, 6},
	"tileGrass_transitionS.png":          {0, 8, 2, 8},
	"tileGrass_transitionW.png":          {5, 0, 5, 2},
	"tileSand1.png":                      {2, 2, 2, 2},
	"tileSand2.png":                      {2, 2, 2, 2},
	"tileSand_roadCornerLL.png":          {2, 2, 3, 3},
	"tileSand_roadCornerLR.png":          {2, 3, 3, 2},
	"tileSand_roadCornerUL.png":          {3, 2, 2, 3},
	"tileSand_roadCornerUR.png":          {3, 3, 2, 2},
	"tileSand_roadCrossing.png":          {3, 3, 3, 3},
	"tileSand_roadCrossingRound.png":     {3, 3, 3, 3},
	"tileSand_roadEast.png":              {2, 3, 2, 3},
	"tileSand_roadNorth.png":             {3, 2, 3, 2},
	"tileSand_roadSplitE.png":            {3, 3, 3, 2},
	"tileSand_roadSplitN.png":            {3, 3, 2, 3},
	"tileSand_roadSplitS.png":            {2, 3, 3, 3},
	"tileSand_roadSplitW.png":            {3, 2, 3, 3},
}

var mplusNormalFont font.Face
var icons []Icon

type Icon struct {
	Name string
	Img  *ebiten.Image
}

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
		icons = append(icons, Icon{
			Name: subTexture.Name,
			Img: origin.SubImage(
				image.Rect(
					subTexture.X,
					subTexture.Y,
					subTexture.X+subTexture.Width,
					subTexture.Y+subTexture.Height,
				)).(*ebiten.Image),
		})
	}

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

	fmt.Println("Found", len(icons), "icons")
}

type Game struct {
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x := 20
	y := 20
	// Draw info
	// msg := fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS())
	// text.Draw(screen, msg, mplusNormalFont, x, 20, color.White)
	for i, max_h := 0, 0; i < len(icons); i++ {
		if strings.HasPrefix(icons[i].Name, "tile") {

			w, h := icons[i].Img.Bounds().Dx()*2, icons[i].Img.Bounds().Dy()
			if h > max_h {
				max_h = h
			}
			ops := &ebiten.DrawImageOptions{}
			ops.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(icons[i].Img, ops)
			text.Draw(screen, fmt.Sprintf("%.28s", icons[i].Name), mplusNormalFont, x, y+12+icons[i].Img.Bounds().Dy()/2, color.White)
			text.Draw(screen, fmt.Sprintf("%d", tileOptions[icons[i].Name][0]), mplusNormalFont, x+icons[i].Img.Bounds().Dx()/2, y+16, color.Black)
			text.Draw(screen, fmt.Sprintf("%d", tileOptions[icons[i].Name][1]), mplusNormalFont, x+icons[i].Img.Bounds().Dx()-8, y+icons[i].Img.Bounds().Dy()/2, color.Black)
			text.Draw(screen, fmt.Sprintf("%d", tileOptions[icons[i].Name][2]), mplusNormalFont, x+icons[i].Img.Bounds().Dx()/2, y+icons[i].Img.Bounds().Dy(), color.Black)
			text.Draw(screen, fmt.Sprintf("%d", tileOptions[icons[i].Name][3]), mplusNormalFont, x, y+icons[i].Img.Bounds().Dy()/2, color.Black)
			x = x + w*2
			if x > screenWidth {
				x = 20
				y = y + max_h + 10
				max_h = 0
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
