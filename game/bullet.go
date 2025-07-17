package game

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	gravity = 3.0
)

type Bullet struct {
	position    Vector
	rotation    float64
	sprite      *ebiten.Image
	speed       float64
	slope       float64
	altitude    float64
	elapsedTime float64
}

func NewBullet(sprite *ebiten.Image, pos Vector, rotation, speed, slope float64) *Bullet {
	return &Bullet{
		position:    pos,
		rotation:    rotation,
		sprite:      sprite,
		speed:       speed,
		slope:       slope,
		altitude:    0.2,
		elapsedTime: 0.0,
	}
}

func (b *Bullet) Update() {
	if b.altitude > 0.0 {
		b.elapsedTime += 1.0 / float64(ebiten.TPS())
		b.position.X += math.Sin(b.rotation) * b.speed
		b.position.Y += math.Cos(b.rotation) * -b.speed
		verticalSpeed := b.speed * math.Sin(b.slope)
		b.slope = verticalSpeed / b.speed
		b.altitude += verticalSpeed*b.elapsedTime - gravity*math.Pow(b.elapsedTime, 2)
		// fmt.Println(verticalSpeed, b.slope, b.altitude, b.elapsedTime)
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if b.altitude > 0.0 {
		scale := b.altitude * 0.1
		if scale < 1.0 {
			scale = 1.0
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Rotate(b.rotation)

		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(b.position.X, b.position.Y)
		screen.DrawImage(b.sprite, op)
	}
}
