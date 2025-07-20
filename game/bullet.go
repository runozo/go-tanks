package game

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	gravity           = 3.0
	bulletSpawnOffset = 20.0
	bulletSpeed       = 10.0
)

type Bullet struct {
	position    Vector
	rotation    float64
	sprite      *ebiten.Image
	speed       float64
	slope       float64
	altitude    float64
	elapsedTime float64
	width       float64
	height      float64
}

func NewBullet(barrel *Barrel) *Bullet {
	position := Vector{
		X: barrel.position.X + math.Sin(barrel.absoluteRotation),
		Y: barrel.position.Y - math.Cos(barrel.absoluteRotation) + bulletSpawnOffset,
	}

	return &Bullet{
		position:    position,
		rotation:    barrel.absoluteRotation,
		sprite:      barrel.bulletSprite,
		speed:       bulletSpeed,
		slope:       barrel.slope,
		altitude:    0.2,
		elapsedTime: 0.0,
		width:       float64(barrel.bulletSprite.Bounds().Dx()),
		height:      float64(barrel.bulletSprite.Bounds().Dy()),
	}
}

func (b *Bullet) Update() {
	if b.altitude > 0.0 {
		b.position.X += math.Sin(b.rotation) * b.speed
		b.position.Y -= math.Cos(b.rotation) * b.speed
		verticalSpeed := b.speed * math.Sin(b.slope)
		b.altitude += verticalSpeed - gravity*math.Pow(b.elapsedTime, 2)
		b.elapsedTime += 1.0 / float64(ebiten.TPS())
		b.slope = verticalSpeed / b.speed
		// fmt.Println(verticalSpeed, b.slope, b.altitude, b.elapsedTime)
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if b.altitude > 0.0 {
		scale := b.altitude * 0.2
		if scale < 1.0 {
			scale = 1.0
		}

		bulletHalfH := b.height / 2
		bulletHalfW := b.width / 2

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-bulletHalfW, -bulletHalfH)
		op.GeoM.Rotate(b.rotation)
		op.GeoM.Translate(bulletHalfW, bulletHalfH)
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(b.position.X, b.position.Y)
		screen.DrawImage(b.sprite, op)
	}
}
