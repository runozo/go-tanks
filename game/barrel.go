package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	b2 "github.com/oliverbestmann/box2d-go"
)

type Barrel struct {
	sprite                            *ebiten.Image
	spriteWidth                       float32
	spriteHeight                      float32
	bulletSprite                      *ebiten.Image
	shootAnimationFrames              []*ebiten.Image
	explosionAnimationFrames          []*ebiten.Image
	explosionHitTargetAnimationFrames []*ebiten.Image
	position                          b2.Vec2
	offset                            b2.Vec2
	relativeRotation                  float64
	absoluteRotation                  float64
	slope                             float64
	tank                              *Tank
	isFiring                          bool
	shootAge                          float64
	game                              *Game
}

func NewBarrel(game *Game, spriteName, bulletSpriteName string, tank *Tank, offset b2.Vec2) *Barrel {
	sprite := game.assets.GetSprite(spriteName)

	if spriteName == "specialBarrel1_outline" {
		sprite = FlipVertical(sprite)
	}

	bulletSprite := game.assets.GetSprite(bulletSpriteName)

	spriteWidth := float32(sprite.Bounds().Dx())
	spriteHeight := float32(sprite.Bounds().Dy())

	pos := tank.tankBody.GetPosition()
	position := b2.Vec2{
		X: pos.X + offset.X - spriteWidth/2, // tank.bodyWidth/2 - spriteWidth/2,
		Y: pos.Y + offset.Y - spriteHeight,  // tank.bodyHeight/2 - spriteHeight,
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

	explosionHitTargetAnimationSprites := []*ebiten.Image{
		game.assets.GetSprite("explosion1"),
		game.assets.GetSprite("explosion2"),
		game.assets.GetSprite("explosion3"),
		game.assets.GetSprite("explosion4"),
		game.assets.GetSprite("explosion5"),
	}

	return &Barrel{
		sprite:                            sprite,
		spriteWidth:                       spriteWidth,
		spriteHeight:                      spriteHeight,
		bulletSprite:                      bulletSprite,
		shootAnimationFrames:              shootAnimationSprites,
		explosionAnimationFrames:          explosionAnimationSprites,
		explosionHitTargetAnimationFrames: explosionHitTargetAnimationSprites,
		position:                          position,
		offset:                            offset,
		relativeRotation:                  0.0,
		absoluteRotation:                  tank.Rotation,
		slope:                             0.0,
		tank:                              tank,
		isFiring:                          false,
		shootAge:                          0.0,
		game:                              game,
	}
}

func (b *Barrel) Fire() *Bullet {
	b.isFiring = true
	b.shootAge = 0
	return NewBullet(b.game, b)
}

func (b *Barrel) Update(tps float64) {
	b.absoluteRotation = b.tank.Rotation + b.relativeRotation

	pos := b.tank.tankBody.GetPosition()

	position := b2.Vec2{
		X: pos.X + b.offset.X - b.spriteWidth/2,
		Y: pos.Y + b.offset.Y - b.spriteHeight,
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
	op_barrel.GeoM.Translate(float64(-b.spriteWidth/2), float64(-b.spriteHeight))
	op_barrel.GeoM.Rotate(b.absoluteRotation)
	op_barrel.GeoM.Translate(float64(b.spriteWidth/2), float64(b.spriteHeight))
	op_barrel.GeoM.Translate(float64(b.position.X), float64(b.position.Y))
	screen.DrawImage(b.sprite, op_barrel)

	// barrel shoot animation
	if b.isFiring {
		shootFrame := b.shootAnimationFrames[int(b.shootAge)%len(b.shootAnimationFrames)]
		shootHalfW := float64(shootFrame.Bounds().Dx()) / 2
		shootFrameHeight := float64(shootFrame.Bounds().Dy())
		shootHalfH := shootFrameHeight / 2
		shootAndBarrellheight := float64(b.spriteHeight) + shootFrameHeight
		op_shoot := &ebiten.DrawImageOptions{}

		// first reverse the frame alone
		op_shoot.GeoM.Translate(-shootHalfW, -shootHalfH)
		op_shoot.GeoM.Rotate(math.Pi)
		op_shoot.GeoM.Translate(shootHalfW, shootHalfH)

		// then translate the frame to the barrel
		op_shoot.GeoM.Translate(-shootHalfW, -shootAndBarrellheight)
		op_shoot.GeoM.Rotate(b.absoluteRotation)
		op_shoot.GeoM.Translate(shootHalfW, shootAndBarrellheight)
		op_shoot.GeoM.Translate(float64(b.position.X+b.spriteWidth/2-float32(shootHalfW)), float64(b.position.Y)-shootFrameHeight)
		screen.DrawImage(shootFrame, op_shoot)
	}

}
