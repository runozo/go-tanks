package game

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	gravity         = 9.8
	bulletSpeed     = 12.0
	bulletMinScale  = 1.0
	scaleCoeff      = 1.8
	initialAltitude = 0.2
)

type Bullet struct {
	position                   resolv.Vector
	sprite                     *ebiten.Image
	verticalSpeed, altitude    float64
	currentSlope, initialSlope float64
	rotation, scale            float64
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

	return &Bullet{
		position:                 position,
		rotation:                 -barrel.solid.Rotation(),
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
}

func (b *Bullet) Update(tps float64) {

	sinRot, cosRot := math.Sincos(b.rotation)
	dt := 1.0 / tps
	b.elapsedTime += dt

	if b.altitude >= initialAltitude {
		b.position.X += sinRot * bulletSpeed
		b.position.Y -= cosRot * bulletSpeed

		gravityEffect := 0.5 * gravity * dt * dt
		b.altitude += b.verticalSpeed*dt - gravityEffect
		b.verticalSpeed -= gravity * dt

		actualSpeed := bulletSpeed * math.Cos(b.initialSlope)
		b.currentSlope = math.Atan2(b.verticalSpeed, actualSpeed)
		b.scale = b.altitude*scaleCoeff + bulletMinScale
		// fmt.Println(b.currentSlope)
	} else {
		b.exploding = true
	}

	if b.exploding {
		// check collision with players
		for _, p := range b.game.players {
			if doesIntersect(b.position, b.sprite.Bounds(), p.tank.solid.Position(), p.tank.bodySprite.Bounds()) {
				b.hasHitTarget = true
				break
			}
		}
		// check collision with enemies
		for _, e := range b.game.enemies {
			if doesIntersect(b.position, b.sprite.Bounds(), e.tank.solid.Position(), e.tank.bodySprite.Bounds()) {
				b.hasHitTarget = true
				break
			}
		}

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
		op.GeoM.Rotate(b.rotation)
		op.GeoM.Translate(bulletHalfW, bulletAndBarrellHeight)

		// true position of the bullet to draw
		op.GeoM.Translate(b.position.X, b.position.Y)
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
		op.GeoM.Translate(b.position.X, b.position.Y)
		screen.DrawImage(frame, op)
	} else {
		b.exploded = true
	}

}

func (b *Bullet) drawExplosion(screen *ebiten.Image) {
	if int(b.explosionElapsedTime) < len(b.explosionFrames) {
		frame := b.explosionFrames[int(b.explosionElapsedTime)]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.position.X, b.position.Y)
		screen.DrawImage(frame, op)
	} else {
		b.exploded = true
	}

}
