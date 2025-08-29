package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Barrel struct {
	sprite                   *ebiten.Image
	spriteWidth              float64
	spriteHeight             float64
	bulletSprite             *ebiten.Image
	shootAnimationFrames     []*ebiten.Image
	explosionAnimationFrames []*ebiten.Image
	position                 Vector
	offset                   Vector
	relativeRotation         float64
	absoluteRotation         float64
	slope                    float64
	tank                     *Tank
	isFiring                 bool
	shootAge                 float64
}

func NewBarrel(game *Game, spriteName, bulletSpriteName string, tank *Tank, offset Vector) *Barrel {
	sprite := game.assets.GetSprite(spriteName)

	if spriteName == "specialBarrel1_outline" {
		sprite = FlipVertical(sprite)
	}

	bulletSprite := game.assets.GetSprite(bulletSpriteName)

	spriteWidth := float64(sprite.Bounds().Dx())
	spriteHeight := float64(sprite.Bounds().Dy())

	position := Vector{
		X: tank.Position.X + offset.X - spriteWidth/2, // tank.bodyWidth/2 - spriteWidth/2,
		Y: tank.Position.Y + offset.Y - spriteHeight,  // tank.bodyHeight/2 - spriteHeight,
	}

	shootAnimationSprites := []*ebiten.Image{
		game.assets.GetSprite("shotThin"),
		game.assets.GetSprite("shotLarge"),
		game.assets.GetSprite("shotOrange"),
		game.assets.GetSprite("shotRed"),
		game.assets.GetSprite("shotOrange"),
		game.assets.GetSprite("shotLarge"),
	}

	explosionAnimationSprites := []*ebiten.Image{
		game.assets.GetSprite("explosionSmoke1"),
		game.assets.GetSprite("explosionSmoke2"),
		game.assets.GetSprite("explosionSmoke3"),
		game.assets.GetSprite("explosionSmoke4"),
		game.assets.GetSprite("explosionSmoke5"),
	}

	return &Barrel{
		sprite:                   sprite,
		spriteWidth:              spriteWidth,
		spriteHeight:             spriteHeight,
		bulletSprite:             bulletSprite,
		shootAnimationFrames:     shootAnimationSprites,
		explosionAnimationFrames: explosionAnimationSprites,
		position:                 position,
		offset:                   offset,
		relativeRotation:         0.0,
		absoluteRotation:         tank.Rotation,
		slope:                    0.0,
		tank:                     tank,
		isFiring:                 false,
		shootAge:                 0.0,
	}
}

func (b *Barrel) Fire() *Bullet {
	b.isFiring = true
	b.shootAge = 0
	return NewBullet(b)
}

func (b *Barrel) Update(tps float64) {
	b.absoluteRotation = b.tank.Rotation + b.relativeRotation
	position := Vector{
		X: b.tank.Position.X + b.offset.X - b.spriteWidth/2,
		Y: b.tank.Position.Y + b.offset.Y - b.spriteHeight,
	}
	b.position = position
	if b.isFiring {
		b.shootAge += 1 / tps * 20
		if int(b.shootAge) >= len(b.shootAnimationFrames) {
			b.isFiring = false
			b.shootAge = 0
		}
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

	// barrel shoot animation
	if b.isFiring {
		shootFrame := b.shootAnimationFrames[int(b.shootAge)%len(b.shootAnimationFrames)]
		shootHalfW := float64(shootFrame.Bounds().Dx()) / 2
		shootFrameHeight := float64(shootFrame.Bounds().Dy())
		shootHalfH := shootFrameHeight / 2
		shootAndBarrellheight := b.spriteHeight + shootFrameHeight
		op_shoot := &ebiten.DrawImageOptions{}

		// first reverse the frame alone
		op_shoot.GeoM.Translate(-shootHalfW, -shootHalfH)
		op_shoot.GeoM.Rotate(math.Pi)
		op_shoot.GeoM.Translate(shootHalfW, shootHalfH)

		// then translate the frame to the barrel
		op_shoot.GeoM.Translate(-shootHalfW, -shootAndBarrellheight)
		op_shoot.GeoM.Rotate(b.absoluteRotation)
		op_shoot.GeoM.Translate(shootHalfW, shootAndBarrellheight)
		op_shoot.GeoM.Translate(b.position.X+b.spriteWidth/2-shootHalfW, b.position.Y-shootFrameHeight)
		screen.DrawImage(shootFrame, op_shoot)
	}

}
