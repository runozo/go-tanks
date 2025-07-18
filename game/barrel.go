package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Barrel struct {
	sprite           *ebiten.Image
	bulletSprite     *ebiten.Image
	position         Vector
	relativeRotation float64
	absoluteRotation float64
	slope            float64
	tank             *Tank
}

func NewBarrel(sprite, bulletSprite *ebiten.Image, tank *Tank) *Barrel {
	return &Barrel{
		sprite:           sprite,
		bulletSprite:     bulletSprite,
		position:         tank.position,
		relativeRotation: 0.0,
		absoluteRotation: tank.rotation,
		slope:            0.0,
		tank:             tank,
	}
}

func (b *Barrel) Fire() *Bullet {
	// barrelBounds := b.sprite.Bounds()
	// bulletBounds := b.bulletSprite.Bounds()
	/*
		halfWBullet := bulletBounds.Dx() / 2
		halfHBullet := bulletBounds.Dy() / 2
		halfW := float64(barrelBounds.Dx()) / 2
		halfH := float64(barrelBounds.Dy()) / 2
		bulletRotation := b.rotation + tank.rotation
	*/

	/*spawnPos := Vector{
		b.position.X + halfW - math.Cos(bulletRotation)*float64(halfWBullet) + math.Sin(bulletRotation)*bulletSpawnOffset,
		b.position.Y + halfH - math.Sin(bulletRotation)*float64(halfHBullet) + math.Cos(bulletRotation)*-bulletSpawnOffset,
	}*/

	return NewBullet(b.bulletSprite, b)
}

func (b *Barrel) Update() {
	b.absoluteRotation = b.tank.rotation + b.relativeRotation
}

func (b *Barrel) Draw(screen *ebiten.Image) {
	// barrel
	tankBounds := b.tank.bodySprite.Bounds()
	barrellBounds := b.sprite.Bounds()
	barrellHalfW := float64(barrellBounds.Dx() / 2)
	barrellHeight := float64(barrellBounds.Dy())
	op_barrel := &ebiten.DrawImageOptions{}
	op_barrel.GeoM.Translate(-barrellHalfW, -barrellHeight)
	op_barrel.GeoM.Rotate(b.absoluteRotation)
	op_barrel.GeoM.Translate(barrellHalfW, barrellHeight)
	op_barrel.GeoM.Translate(
		b.tank.position.X+float64(tankBounds.Dx())/2-barrellHalfW,
		b.tank.position.Y+float64(tankBounds.Dy())/2-barrellHeight,
	)
	screen.DrawImage(b.sprite, op_barrel)
}
