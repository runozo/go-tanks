package game

import "github.com/hajimehoshi/ebiten/v2"

type Barrel struct {
	sprite   *ebiten.Image
	position Vector
	rotation float64
	slope    float64
}

func NewBarrel(sprite *ebiten.Image, tankPosition Vector, rotation float64) *Barrel {
	return &Barrel{
		sprite:   sprite,
		position: tankPosition,
		rotation: rotation,
		slope:    0.0,
	}
}

func (b *Barrel) Draw(screen *ebiten.Image, tank *Tank) {
	// barrel
	tankBounds := tank.bodySprite.Bounds()
	barrellBounds := b.sprite.Bounds()
	barrellHalfW := float64(barrellBounds.Dx() / 2)
	barrellHeight := float64(barrellBounds.Dy())
	op_barrel := &ebiten.DrawImageOptions{}
	op_barrel.GeoM.Translate(-barrellHalfW, -barrellHeight)
	op_barrel.GeoM.Rotate(tank.rotation + b.rotation)
	op_barrel.GeoM.Translate(barrellHalfW, barrellHeight)
	op_barrel.GeoM.Translate(
		tank.position.X+float64(tankBounds.Dx())/2-barrellHalfW,
		tank.position.Y+float64(tankBounds.Dy())/2-barrellHeight,
	)
	screen.DrawImage(b.sprite, op_barrel)
}
