package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Tank struct {
	body   *ebiten.Image
	barrel *ebiten.Image
}

func NewTank(game *Game) *Tank {
	body := game.assets.GetSprite("tankBody_red.png")
	barrel := game.assets.GetSprite("tankSand_barrel1_outline.png")

	return &Tank{
		body:   body,
		barrel: barrel,
	}
}

func (p *Tank) Draw(screen *ebiten.Image, pos Vector, rotation float64) {
	// Draw the tank
	// body
	bodyBounds := p.body.Bounds()
	bodyHalfW := float64(bodyBounds.Dx() / 2)
	bodyHalfH := float64(bodyBounds.Dy() / 2)
	op_body := &ebiten.DrawImageOptions{}
	op_body.GeoM.Translate(-bodyHalfW, -bodyHalfH)
	op_body.GeoM.Rotate(rotation)
	op_body.GeoM.Translate(bodyHalfW, bodyHalfH)
	op_body.GeoM.Translate(pos.X, pos.Y)

	// barrel
	barrellBounds := p.barrel.Bounds()
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

	screen.DrawImage(p.body, op_body)
	screen.DrawImage(p.barrel, op_barrell)
}
