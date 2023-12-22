package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Tank struct {
	BodySprite   *ebiten.Image
	BarrelSprite *ebiten.Image
}

func NewTank(game *Game, bodySprite, barrelSprite *ebiten.Image) *Tank {
	return &Tank{
		BodySprite:   bodySprite,
		BarrelSprite: barrelSprite,
	}
}

func (p *Tank) Draw(screen *ebiten.Image, pos Vector, rotation float64) {
	// Draw the tank
	// body
	bodyBounds := p.BodySprite.Bounds()
	bodyHalfW := float64(bodyBounds.Dx() / 2)
	bodyHalfH := float64(bodyBounds.Dy() / 2)
	op_body := &ebiten.DrawImageOptions{}
	op_body.GeoM.Translate(-bodyHalfW, -bodyHalfH)
	op_body.GeoM.Rotate(rotation)
	op_body.GeoM.Translate(bodyHalfW, bodyHalfH)
	op_body.GeoM.Translate(pos.X, pos.Y)

	// barrel
	barrellBounds := p.BarrelSprite.Bounds()
	barrellHalfW := float64(barrellBounds.Dx() / 2)
	barrellHeight := float64(barrellBounds.Dy())
	op_barrell := &ebiten.DrawImageOptions{}
	op_barrell.GeoM.Translate(-barrellHalfW, -barrellHeight)
	op_barrell.GeoM.Rotate(rotation)
	op_barrell.GeoM.Translate(barrellHalfW, barrellHeight)
	op_barrell.GeoM.Translate(
		pos.X+float64(bodyBounds.Dx())/2-barrellHalfW,
		pos.Y+float64(bodyBounds.Dy())/2-barrellHeight,
	)

	screen.DrawImage(p.BodySprite, op_body)
	screen.DrawImage(p.BarrelSprite, op_barrell)
}
