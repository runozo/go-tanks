package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	b2 "github.com/oliverbestmann/box2d-go"
)

type ProgressBar struct {
	Width       int
	Height      int
	Position    b2.Vec2
	Percent     float64
	BorderWidth int
	Font        *text.GoTextFace
	Text        string
	txtPos      b2.Vec2
}

const (
	barBorderWidth = 4
)

func NewProgressBar(width, height int, str string, position b2.Vec2, font *text.GoTextFace) *ProgressBar {
	txtwidth, txtheight := text.Measure(str+" 100%", font, 0)

	x := position.X + float32(width)/2 - float32(txtwidth)/2
	y := position.Y + float32(height)/2 - float32(txtheight)/2

	return &ProgressBar{
		Width:       width,
		Height:      height,
		Position:    position,
		Percent:     1.0,
		BorderWidth: barBorderWidth,
		Font:        font,
		Text:        str,
		txtPos: b2.Vec2{
			X: x,
			Y: y,
		},
	}
}

func (p *ProgressBar) Update(percent float64) {
	p.Percent = percent
}

func (p *ProgressBar) Draw(screen *ebiten.Image) {
	backgroundImage := ebiten.NewImage(p.Width, p.Height) // default image is black
	backgroundImage.Fill(color.RGBA{R: 0, A: 128})        // dark gray
	opBackground := &ebiten.DrawImageOptions{}
	opBackground.GeoM.Translate(float64(p.Position.X), float64(p.Position.Y))
	screen.DrawImage(backgroundImage, opBackground)

	barWidth := int(float64(p.Width) * p.Percent / 100)
	if barWidth < p.BorderWidth*2+1 {
		barWidth = p.BorderWidth*2 + 1
	}

	barImage := ebiten.NewImage(barWidth-p.BorderWidth*2, p.Height-p.BorderWidth*2)
	barImage.Fill(color.RGBA{B: 255, A: 255}) // red
	opBar := &ebiten.DrawImageOptions{}
	opBar.GeoM.Translate(float64(p.Position.X)+float64(p.BorderWidth), float64(p.Position.Y)+float64(p.BorderWidth))
	screen.DrawImage(barImage, opBar)

	if p.Text != "" {
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(p.txtPos.X), float64(p.txtPos.Y))
		op.ColorScale.ScaleWithColor(color.RGBA{R: 255, A: 255}) // red
		text.Draw(screen, fmt.Sprintf("%s %d%%", p.Text, int(p.Percent)), p.Font, op)
	}
}
