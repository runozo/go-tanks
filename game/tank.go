package game

import "github.com/hajimehoshi/ebiten/v2"

type Tank struct {
	body   *ebiten.Image
	barrel *ebiten.Image
}

func NewTank(game *Game) *Tank {
	body := game.assets.GetSprite("tankBody_red.png")
	barrel := game.assets.GetSprite("tankGreen_barrel1.png")

	return &Tank{
		body:   body,
		barrel: barrel,
	}
}

func (p *Tank) Draw(screen *ebiten.Image, pos Vector, rotation float64) {
	// Draw the tank
	// body
	bodyBounds := p.body.Bounds()
	op_body := &ebiten.DrawImageOptions{}
	op_body.GeoM.Translate(-float64(bodyBounds.Dx())/2, -float64(bodyBounds.Dy())/2)
	op_body.GeoM.Rotate(rotation)
	op_body.GeoM.Translate(float64(bodyBounds.Dx())/2, float64(bodyBounds.Dy())/2)
	op_body.GeoM.Translate(pos.X, pos.Y)
	// barrel
	barrellBounds := p.barrel.Bounds()
	op_barrell := &ebiten.DrawImageOptions{}
	op_barrell.GeoM.Translate(-float64(barrellBounds.Dx())/2, -float64(barrellBounds.Dy()+1))
	op_barrell.GeoM.Rotate(rotation)
	op_barrell.GeoM.Translate(float64(barrellBounds.Dx())/2, float64(barrellBounds.Dy()+1))
	op_barrell.GeoM.Translate(
		pos.X+float64(p.barrel.Bounds().Dx()-1),
		pos.Y+float64(-p.barrel.Bounds().Dy()/2))
	screen.DrawImage(p.body, op_body)
	screen.DrawImage(p.barrel, op_barrell)
}
