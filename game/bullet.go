package game

import (
	"fmt"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	gravity        = 9.8
	bulletSpeed    = 12.0
	bulletMinScale = 1.0
	bulletMaxScale = 8.0
	scaleCoeff     = 1.8
)

type Bullet struct {
	position      Vector
	rotation      float64
	sprite        *ebiten.Image
	verticalSpeed float64
	currentSlope  float64
	initialSlope  float64
	altitude      float64
	scale         float64
	elapsedTime   float64
	spriteWidth   float64
	spriteHeight  float64
	barrelWidth   float64
	barrelHeight  float64
}

func NewBullet(barrel *Barrel) *Bullet {
	bulletSprite := barrel.bulletSprite
	bulletSpriteWidth := float64(bulletSprite.Bounds().Dx())
	bulletSpriteHeight := float64(bulletSprite.Bounds().Dy())

	position := Vector{
		X: barrel.position.X + barrel.spriteWidth/2 - bulletSpriteWidth/2,
		Y: barrel.position.Y - bulletSpriteHeight,
	}

	// fmt.Println("Barrel position", barrel.position.X, barrel.position.Y, "Bullet position", position.X, position.Y)

	return &Bullet{
		position:      position,
		rotation:      barrel.absoluteRotation,
		sprite:        bulletSprite,
		verticalSpeed: bulletSpeed * math.Sin(barrel.slope),
		currentSlope:  barrel.slope,
		initialSlope:  barrel.slope,
		altitude:      0.2,
		scale:         bulletMinScale,
		elapsedTime:   0.0,
		spriteWidth:   bulletSpriteWidth,
		spriteHeight:  bulletSpriteHeight,
		barrelWidth:   barrel.spriteWidth,
		barrelHeight:  barrel.spriteHeight,
	}
}

func (b *Bullet) Update() {
	if b.altitude > 0.0 {
		dt := 1.0 / float64(ebiten.TPS())
		sinRot, cosRot := math.Sincos(b.rotation)
		b.position.X += sinRot * bulletSpeed
		b.position.Y -= cosRot * bulletSpeed
		b.elapsedTime += dt

		gravityEffect := 0.5 * gravity * dt * dt
		b.altitude += b.verticalSpeed*dt - gravityEffect
		b.verticalSpeed -= gravity * dt

		actualSpeed := bulletSpeed * math.Cos(b.initialSlope)
		b.currentSlope = math.Atan2(b.verticalSpeed, actualSpeed)
		b.scale = b.altitude*scaleCoeff + bulletMinScale
		fmt.Println(b.currentSlope)
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if b.altitude > 0.0 {
		bulletHalfW := b.spriteWidth / 2
		bulletHalfH := b.spriteHeight / 2
		bulletAndBarrellHeight := b.barrelHeight + b.spriteHeight

		// fmt.Println(b.altitude) // , "Scale", b.scale)

		op := &ebiten.DrawImageOptions{}

		// center the bullet than scale
		op.GeoM.Translate(-bulletHalfW, -bulletHalfH)
		op.GeoM.Scale(b.scale, b.scale-math.Abs(b.currentSlope)*2.0)
		op.GeoM.Translate(bulletHalfW, bulletHalfH)

		// center the bullet and the barrel than rotate
		op.GeoM.Translate(-bulletHalfW, -bulletAndBarrellHeight)
		op.GeoM.Rotate(b.rotation)
		op.GeoM.Translate(bulletHalfW, bulletAndBarrellHeight)

		// true position of the bullet
		op.GeoM.Translate(b.position.X, b.position.Y)
		screen.DrawImage(b.sprite, op)
	}
}
