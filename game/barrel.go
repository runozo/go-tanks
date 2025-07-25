package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Barrel struct {
	sprite           *ebiten.Image
	spriteWidth      float64
	spriteHeight     float64
	bulletSprite     *ebiten.Image
	shootFrames      []*ebiten.Image
	position         Vector
	relativeRotation float64
	absoluteRotation float64
	slope            float64
	tank             *Tank
	isFiring         bool
	shootFrameNumber int
	count            float64
}

func NewBarrel(sprite, bulletSprite *ebiten.Image, tank *Tank, shootSprites []*ebiten.Image) *Barrel {
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
		shootFrames:      shootSprites,
		position:         position,
		relativeRotation: 0.0,
		absoluteRotation: tank.rotation,
		slope:            0.0,
		tank:             tank,
		isFiring:         false,
		shootFrameNumber: 0,
		count:            0.0,
	}
}

func (b *Barrel) Fire() *Bullet {
	b.isFiring = true
	b.shootFrameNumber = 0
	return NewBullet(b)
}

func (b *Barrel) Update(tps float64) {
	b.absoluteRotation = b.tank.rotation + b.relativeRotation
	position := Vector{
		X: b.tank.position.X + b.tank.bodyWidth/2 - b.spriteWidth/2,
		Y: b.tank.position.Y + b.tank.bodyHeight/2 - b.spriteHeight,
	}
	b.position = position
	if b.isFiring {
		b.shootFrameNumber = int(b.count)
		if b.shootFrameNumber >= len(b.shootFrames) {
			b.isFiring = false
			b.shootFrameNumber = 0
			b.count = 0.0
		}
		b.count += 1 / tps * 10
	}
}

func (b *Barrel) Draw(screen *ebiten.Image) {
	// barrel
	op_barrel := &ebiten.DrawImageOptions{}
	op_barrel.GeoM.Translate(-b.spriteWidth/2, -b.spriteHeight)
	op_barrel.GeoM.Rotate(b.absoluteRotation)
	op_barrel.GeoM.Translate(b.spriteWidth/2, b.spriteHeight)
	op_barrel.GeoM.Translate(b.position.X, b.position.Y)
	screen.DrawImage(b.sprite, op_barrel)

	if b.isFiring {
		shootFrame := b.shootFrames[b.shootFrameNumber]
		shootHalfW := float64(shootFrame.Bounds().Dx()) / 2
		shootAndBarrellheight := b.spriteHeight + float64(shootFrame.Bounds().Dy())
		op_shoot := &ebiten.DrawImageOptions{}
		op_shoot.GeoM.Translate(-shootHalfW, -shootAndBarrellheight)
		op_shoot.GeoM.Rotate(b.absoluteRotation)
		op_shoot.GeoM.Translate(shootHalfW, shootAndBarrellheight)
		op_shoot.GeoM.Translate(b.position.X+b.spriteWidth/2-shootHalfW, b.position.Y-float64(shootFrame.Bounds().Dy()))
		screen.DrawImage(shootFrame, op_shoot)
	}

}
