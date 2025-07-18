package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Barrel struct {
	sprite       *ebiten.Image
	bulletSprite *ebiten.Image
	position     Vector
	rotation     float64
	slope        float64
}

func NewBarrel(sprite, bulletSprite *ebiten.Image, tankPosition Vector, rotation float64) *Barrel {
	return &Barrel{
		sprite:       sprite,
		bulletSprite: bulletSprite,
		position:     tankPosition,
		rotation:     rotation,
		slope:        0.0,
	}
}

func (b *Barrel) Fire(tank *Tank) *Bullet {
	barrelBounds := b.sprite.Bounds()
	bulletBounds := b.bulletSprite.Bounds()
	halfWBullet := bulletBounds.Dx() / 2
	halfHBullet := bulletBounds.Dy() / 2
	halfW := float64(barrelBounds.Dx()) / 2
	halfH := float64(barrelBounds.Dy()) / 2
	bulletRotation := b.rotation + tank.rotation

	spawnPos := Vector{
		tank.position.X + halfW - math.Cos(bulletRotation)*float64(halfWBullet) + math.Sin(bulletRotation)*bulletSpawnOffset,
		tank.position.Y + halfH - math.Sin(bulletRotation)*float64(halfHBullet) + math.Cos(bulletRotation)*-bulletSpawnOffset,
	}

	return NewBullet(b.bulletSprite, spawnPos, bulletRotation, bulletSpeed, b.slope)
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
