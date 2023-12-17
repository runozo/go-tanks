package game

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	position Vector
	rotation float64
	sprite   *ebiten.Image
	speed    float64
}

func NewBullet(sprite *ebiten.Image, pos Vector, rotation, speed float64) *Bullet {
	return &Bullet{
		position: pos,
		rotation: rotation,
		sprite:   sprite,
		speed:    speed,
	}
}

func (b *Bullet) Update() {
	b.position.X += math.Sin(b.rotation) * b.speed
	b.position.Y += math.Cos(b.rotation) * -b.speed
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.position.X, b.position.Y)
	screen.DrawImage(b.sprite, op)
}
