package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Barrel struct {
	sprite           *ebiten.Image
	spriteWidth      float64
	spriteHeight     float64
	bulletSprite     *ebiten.Image
	position         Vector
	relativeRotation float64
	absoluteRotation float64
	slope            float64
	tank             *Tank
}

func NewBarrel(sprite, bulletSprite *ebiten.Image, tank *Tank) *Barrel {
	spriteWidth := float64(sprite.Bounds().Dx())
	spriteHeight := float64(sprite.Bounds().Dy())
	position := Vector{
		X: tank.position.X + tank.bodyWidth/2 - spriteWidth/2,
		Y: tank.position.Y + tank.bodyHeight/2 - spriteHeight,
	}
	return &Barrel{
		sprite:           sprite,
		spriteWidth:      spriteWidth,
		spriteHeight:     spriteHeight,
		bulletSprite:     bulletSprite,
		position:         position,
		relativeRotation: 0.0,
		absoluteRotation: tank.rotation,
		slope:            0.0,
		tank:             tank,
	}
}

func (b *Barrel) Fire() *Bullet {
	return NewBullet(b)
}

func (b *Barrel) Update() {
	b.absoluteRotation = b.tank.rotation + b.relativeRotation
	position := Vector{
		X: b.tank.position.X + b.tank.bodyWidth/2 - b.spriteWidth/2,
		Y: b.tank.position.Y + b.tank.bodyHeight/2 - b.spriteHeight,
	}
	b.position = position
}

func (b *Barrel) Draw(screen *ebiten.Image) {
	// barrel
	op_barrel := &ebiten.DrawImageOptions{}
	op_barrel.GeoM.Translate(-b.spriteWidth/2, -b.spriteHeight)
	op_barrel.GeoM.Rotate(b.absoluteRotation)
	op_barrel.GeoM.Translate(b.spriteWidth/2, b.spriteHeight)
	op_barrel.GeoM.Translate(b.position.X, b.position.Y)
	screen.DrawImage(b.sprite, op_barrel)

	// bullet debug
	/*
		bulletHalfW := float64(b.bulletSprite.Bounds().Dx()) / 2
		bulletAndBarrellheight := b.spriteHeight + float64(b.bulletSprite.Bounds().Dy())
		op_bullet := &ebiten.DrawImageOptions{}
		op_bullet.GeoM.Translate(-bulletHalfW, -bulletAndBarrellheight)
		op_bullet.GeoM.Rotate(b.absoluteRotation)
		op_bullet.GeoM.Translate(bulletHalfW, bulletAndBarrellheight)
		op_bullet.GeoM.Translate(b.position.X+b.spriteWidth/2-bulletHalfW, b.position.Y-float64(b.bulletSprite.Bounds().Dy()))
		screen.DrawImage(b.bulletSprite, op_bullet)
	*/
}
