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

func (b *Barrel) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(b.rotation)
	op.GeoM.Translate(b.position.X, b.position.Y)
	screen.DrawImage(b.sprite, op)
}
