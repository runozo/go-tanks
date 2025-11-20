package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type Barrel struct {
	sprite                            *ebiten.Image
	solid                             *resolv.ConvexPolygon
	spriteWidth                       float64
	spriteHeight                      float64
	bulletSprite                      *ebiten.Image
	shootAnimationFrames              []*ebiten.Image
	explosionAnimationFrames          []*ebiten.Image
	explosionHitTargetAnimationFrames []*ebiten.Image
	offset                            resolv.Vector
	relativeRotation                  float64
	slope                             float64
	tank                              *Tank
	isFiring                          bool
	shootAge                          float64
	game                              *Game
}

func NewBarrel(game *Game, spriteName, bulletSpriteName string, tank *Tank, offset resolv.Vector) *Barrel {
	sprite := game.assets.GetSprite(spriteName)

	if spriteName == "specialBarrel1_outline" {
		sprite = FlipVertical(sprite)
	}

	bulletSprite := game.assets.GetSprite(bulletSpriteName)

	spriteWidth := float64(sprite.Bounds().Dx())
	spriteHeight := float64(sprite.Bounds().Dy())

	position := resolv.Vector{
		X: tank.solid.Position().X + offset.X - spriteWidth/2, // tank.bodyWidth/2 - spriteWidth/2,
		Y: tank.solid.Position().Y + offset.Y - spriteHeight,  // tank.bodyHeight/2 - spriteHeight,
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

	// define resolv shape
	solid := resolv.NewRectangle(position.X, position.Y, spriteWidth, spriteHeight)
	solid.SetRotation(tank.solid.Rotation())
	solid.Tags().Set(TagBarrel)
	game.space.Add(solid)

	return &Barrel{
		sprite:                            sprite,
		solid:                             solid,
		spriteWidth:                       spriteWidth,
		spriteHeight:                      spriteHeight,
		bulletSprite:                      bulletSprite,
		shootAnimationFrames:              shootAnimationSprites,
		explosionAnimationFrames:          explosionAnimationSprites,
		explosionHitTargetAnimationFrames: explosionHitTargetAnimationSprites,
		offset:                            offset,
		relativeRotation:                  0.0,
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
	b.solid.SetRotation(-b.tank.solid.Rotation() - b.relativeRotation)
	b.solid.SetPositionVec(b.tank.solid.Position().Add(b.offset))

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
	spriteHalfW := b.solid.Bounds().Width() / 2
	spriteHalfH := b.solid.Bounds().Height() / 2
	op_barrel := &ebiten.DrawImageOptions{}
	op_barrel.GeoM.Translate(-spriteHalfW, -spriteHalfH)
	op_barrel.GeoM.Rotate(-b.solid.Rotation())
	op_barrel.GeoM.Translate(spriteHalfW, spriteHalfH)
	op_barrel.GeoM.Translate(b.solid.Center().X-spriteHalfW, b.solid.Center().Y-spriteHalfH)
	screen.DrawImage(b.sprite, op_barrel)

	// barrel shoot animation
	if b.isFiring {
		shootFrame := b.shootAnimationFrames[int(b.shootAge)%len(b.shootAnimationFrames)]
		shootHalfW := float64(shootFrame.Bounds().Dx()) / 2
		shootFrameHeight := float64(shootFrame.Bounds().Dy())
		shootHalfH := shootFrameHeight / 2
		shootAndBarrellheight := b.spriteHeight + shootFrameHeight
		op_shoot := &ebiten.DrawImageOptions{}

		// first rotate the frame alone
		op_shoot.GeoM.Translate(-shootHalfW, -shootHalfH)
		op_shoot.GeoM.Rotate(math.Pi)
		op_shoot.GeoM.Translate(shootHalfW, shootHalfH)

		// then translate the frame to the barrel
		op_shoot.GeoM.Translate(-shootHalfW, -shootAndBarrellheight)
		op_shoot.GeoM.Rotate(-b.solid.Rotation())
		op_shoot.GeoM.Translate(shootHalfW, shootAndBarrellheight)
		op_shoot.GeoM.Translate(
			b.solid.Center().X-shootHalfW,
			b.solid.Center().Y-shootFrameHeight,
		)
		screen.DrawImage(shootFrame, op_shoot)
	}

}
