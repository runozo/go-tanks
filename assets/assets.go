package assets

import (
	"embed"
	"encoding/xml"
	"image"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

//go:embed spritesheet/*
var assets embed.FS

const (
	spriteSheet  = "spritesheet/allSprites_default.png"
	xmlSpriteMap = "spritesheet/allSprites_default.xml"
)

type Assets struct {
	SpriteSheet       *ebiten.Image
	SpriteSheetXMLMap SpriteMap
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

func NewAssets() *Assets {
	var spriteSheet = mustLoadImage(spriteSheet)
	var spriteMap = mustLoadXMLSpriteMap(xmlSpriteMap)

	return &Assets{
		SpriteSheet:       spriteSheet,
		SpriteSheetXMLMap: spriteMap,
	}
}

func (a *Assets) GetSprite(name string) *ebiten.Image {
	for i := 0; i < len(a.SpriteSheetXMLMap.SubTextures); i++ {
		subTexture := a.SpriteSheetXMLMap.SubTextures[i]
		if subTexture.Name == name {
			return a.SpriteSheet.SubImage(
				image.Rect(
					subTexture.X,
					subTexture.Y,
					subTexture.X+subTexture.Width,
					subTexture.Y+subTexture.Height,
				)).(*ebiten.Image)
		}
	}
	return nil
}

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
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

func mustLoadImages(path string) []*ebiten.Image {
	matches, err := fs.Glob(assets, path)
	if err != nil {
		panic(err)
	}

	images := make([]*ebiten.Image, len(matches))
	for i, match := range matches {
		images[i] = mustLoadImage(match)
	}

	return images
}

func mustLoadFont(name string) font.Face {
	f, err := assets.ReadFile(name)
	if err != nil {
		panic(err)
	}

	tt, err := opentype.Parse(f)
	if err != nil {
		panic(err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		panic(err)
	}

	return face
}

func mustLoadXMLSpriteMap(name string) SpriteMap {
	byteValue, _ := fs.ReadFile(assets, name)
	var s SpriteMap
	if err := xml.Unmarshal(byteValue, &s); err != nil {
		panic(err)
	}
	// fmt.Println(s)
	return s
}
