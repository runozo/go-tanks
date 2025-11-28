package game

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	gravity         = 9.8
	bulletSpeed     = 1.0
	bulletMinScale  = 1.0
	scaleCoeff      = 1.8
	initialAltitude = 0.2
)

type Bullet struct {
	solid                      *resolv.ConvexPolygon
	sprite                     *ebiten.Image
	moveVec                    resolv.Vector
	verticalSpeed, altitude    float64
	currentSlope, initialSlope float64
	scale                      float64
	elapsedTime                float64
	spriteWidth, spriteHeight  float64
	barrelWidth                float64
	barrelHeight               float64
	explosionElapsedTime       float64
	explosionFrames            []*ebiten.Image
	explosionHitTargetFrames   []*ebiten.Image
	exploding                  bool
	exploded                   bool
	hasHitTarget               bool
	game                       *Game
}

func NewBullet(game *Game, barrel *Barrel) *Bullet {
	bulletSprite := barrel.bulletSprite
	bulletSpriteWidth := float64(bulletSprite.Bounds().Dx())
	bulletSpriteHeight := float64(bulletSprite.Bounds().Dy())

	position := resolv.Vector{
		X: barrel.solid.Center().X + barrel.spriteWidth/2 - bulletSpriteWidth/2,
		Y: barrel.solid.Center().Y - bulletSpriteHeight,
	}

	// fmt.Println("Barrel position", barrel.position.X, barrel.position.Y, "Bullet position", position.X, position.Y)

	b := Bullet{
		solid:                    resolv.NewRectangle(position.X, position.Y, bulletSpriteWidth, bulletSpriteHeight),
		sprite:                   bulletSprite,
		verticalSpeed:            bulletSpeed * math.Sin(barrel.slope),
		currentSlope:             barrel.slope,
		initialSlope:             barrel.slope,
		altitude:                 initialAltitude,
		scale:                    bulletMinScale,
		elapsedTime:              0.0,
		spriteWidth:              bulletSpriteWidth,
		spriteHeight:             bulletSpriteHeight,
		barrelWidth:              barrel.spriteWidth,
		barrelHeight:             barrel.spriteHeight,
		explosionElapsedTime:     0,
		explosionFrames:          barrel.explosionAnimationFrames,
		explosionHitTargetFrames: barrel.explosionHitTargetAnimationFrames,
		exploding:                false,
		exploded:                 false,
		hasHitTarget:             false,
		game:                     game,
	}

	b.solid.SetRotation(-barrel.solid.Rotation())

	return &b
}

func (b *Bullet) Update(tps float64) {
	sinRot, cosRot := math.Sincos(b.solid.Rotation())
	dt := 1.0 / tps
	b.elapsedTime += dt

	if b.altitude >= initialAltitude {
		b.moveVec.X += sinRot * bulletSpeed
		b.moveVec.Y -= cosRot * bulletSpeed

		gravityEffect := 0.0000000005 * gravity * dt
		b.altitude += b.verticalSpeed*dt - gravityEffect
		b.verticalSpeed -= gravity * dt

		actualSpeed := bulletSpeed * math.Cos(b.initialSlope)
		b.currentSlope = math.Atan2(b.verticalSpeed, actualSpeed)
		b.scale = b.altitude*scaleCoeff + bulletMinScale

		b.solid.MoveVec(b.moveVec)
	} else {
		b.exploding = true
	}

	if b.exploding {
		b.explosionElapsedTime += dt * float64(len(b.explosionHitTargetFrames))
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if !b.exploding && !b.exploded {
		bulletHalfW := b.spriteWidth / 2
		bulletHalfH := b.spriteHeight / 2
		bulletAndBarrellHeight := b.barrelHeight + b.spriteHeight

		// fmt.Println(b.altitude) // , "Scale", b.scale)

		op := &ebiten.DrawImageOptions{}

		// center the bullet than scale
		op.GeoM.Translate(-bulletHalfW, -bulletHalfH)
		op.GeoM.Scale(b.scale, b.scale-math.Abs(b.currentSlope)*scaleCoeff) // simulate bullet deflection (poorly)
		op.GeoM.Translate(bulletHalfW, bulletHalfH)

		// center the bullet and the barrel than rotate
		op.GeoM.Translate(-bulletHalfW, -bulletAndBarrellHeight)
		op.GeoM.Rotate(-b.solid.Rotation())
		op.GeoM.Translate(bulletHalfW, bulletAndBarrellHeight)

		// true position of the bullet to draw
		op.GeoM.Translate(b.solid.Position().X, b.solid.Position().Y)
		screen.DrawImage(b.sprite, op)
	} else {
		if b.hasHitTarget {
			b.drawExplosionHitTarget(screen)
		} else {
			b.drawExplosion(screen)
		}

	}
}

func (b *Bullet) drawExplosionHitTarget(screen *ebiten.Image) {
	if int(b.explosionElapsedTime) < len(b.explosionHitTargetFrames) {
		frame := b.explosionHitTargetFrames[int(b.explosionElapsedTime)]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.solid.Position().X, b.solid.Position().Y)
		screen.DrawImage(frame, op)
	} else {
		b.exploded = true
	}

}

func (b *Bullet) drawExplosion(screen *ebiten.Image) {
	if int(b.explosionElapsedTime) < len(b.explosionFrames) {
		frame := b.explosionFrames[int(b.explosionElapsedTime)]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.solid.Position().X, b.solid.Position().Y)
		screen.DrawImage(frame, op)
	} else {
		b.exploded = true
	}

}
