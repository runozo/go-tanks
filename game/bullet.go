package game

import (
	"fmt"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	gravity           = 8.0
	bulletSpawnOffset = 20.0
	bulletSpeed       = 10.0
	bulletMinScale    = 1.0
	bulletMaxScale    = 8.0
)

type Bullet struct {
	position      Vector
	rotation      float64
	sprite        *ebiten.Image
	verticalSpeed float64
	slope         float64
	altitude      float64
	scale         float64
	elapsedTime   float64
	spriteWidth   float64
	spriteHeight  float64
	barrelWidth   float64
	barrelHeight  float64
}

func NewBullet(barrel *Barrel) *Bullet {
	spriteWidth := float64(barrel.bulletSprite.Bounds().Dx())
	spriteHeight := float64(barrel.bulletSprite.Bounds().Dy())

	position := Vector{
		X: barrel.position.X + barrel.spriteWidth/2 - spriteWidth/2,
		Y: barrel.position.Y - spriteHeight,
	}

	fmt.Println("Barrel position", barrel.position.X, barrel.position.Y, "Bullet position", position.X, position.Y)

	return &Bullet{
		position:      position,
		rotation:      barrel.absoluteRotation,
		sprite:        barrel.bulletSprite,
		verticalSpeed: bulletSpeed * math.Sin(barrel.slope),
		slope:         barrel.slope,
		altitude:      0.2,
		scale:         bulletMinScale,
		elapsedTime:   0.0,
		spriteWidth:   spriteWidth,
		spriteHeight:  spriteHeight,
		barrelWidth:   barrel.spriteWidth,
		barrelHeight:  barrel.spriteHeight,
	}
}

func (b *Bullet) Update() {
	if b.altitude > 0.0 {
		b.position.X += math.Sin(b.rotation) * bulletSpeed
		b.position.Y -= math.Cos(b.rotation) * bulletSpeed
		b.altitude += b.verticalSpeed*b.elapsedTime - 0.5*gravity*math.Pow(b.elapsedTime, 2)
		b.scale = b.altitude
		if b.scale < bulletMinScale {
			b.scale = bulletMinScale
		}
		/*
			if b.scale > bulletMaxScale {
				b.scale = bulletMaxScale
			}
		*/
		b.elapsedTime += 1.0 / float64(ebiten.TPS())
		// b.slope = b.verticalSpeed / bulletSpeed
		// fmt.Println(verticalSpeed, b.slope, b.altitude, b.elapsedTime)
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if b.altitude > 0.0 {
		bulletHalfW := b.spriteWidth / 2
		bulletAndBarrellheight := b.barrelHeight + b.spriteHeight

		fmt.Println(b.altitude) // , "Scale", b.scale)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-bulletHalfW, -bulletAndBarrellheight)
		op.GeoM.Rotate(b.rotation)
		op.GeoM.Scale(b.scale, b.scale)
		op.GeoM.Translate(bulletHalfW, bulletAndBarrellheight)
		op.GeoM.Translate(b.position.X, b.position.Y)
		screen.DrawImage(b.sprite, op)
	}
}
