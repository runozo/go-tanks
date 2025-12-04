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
	verticalSpeed, altitude    float64
	currentSlope, initialSlope float64
	scale                      float64
	elapsedTime                float64
	bulletHalfW, bulletHalfH   float64
	explosionElapsedTime       float64
	explosionFrames            []*ebiten.Image
	explosionHitTargetFrames   []*ebiten.Image
	exploding                  bool
	exploded                   bool
	hasHitTarget               bool
	game                       *Game
}

func NewBullet(game *Game, barrel *Barrel) *Bullet {
	bulletSpriteWidth := float64(barrel.bulletSprite.Bounds().Dx())
	bulletSpriteHeight := float64(barrel.bulletSprite.Bounds().Dy())

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
	s, c := math.Sincos(barrel.solid.Rotation())
	solid := resolv.NewRectangle(
		barrel.solid.Center().X-(barrel.spriteHeight+bulletSpriteHeight)*s,
		barrel.solid.Center().Y-(barrel.spriteHeight+bulletSpriteHeight)*c,
		bulletSpriteWidth,
		bulletSpriteHeight,
	)

	solid.SetRotation(barrel.solid.Rotation())
	solid.Tags().Set(TagBullet)
	game.space.Add(solid)

	b := Bullet{
		solid:                    solid,
		sprite:                   barrel.bulletSprite,
		verticalSpeed:            bulletSpeed * math.Sin(barrel.slope),
		currentSlope:             barrel.slope,
		initialSlope:             barrel.slope,
		altitude:                 initialAltitude,
		scale:                    bulletMinScale,
		elapsedTime:              0.0,
		bulletHalfW:              bulletSpriteWidth / 2.0,
		bulletHalfH:              bulletSpriteHeight / 2.0,
		explosionElapsedTime:     0,
		explosionFrames:          explosionAnimationSprites,
		explosionHitTargetFrames: explosionHitTargetAnimationSprites,
		exploding:                false,
		exploded:                 false,
		hasHitTarget:             false,
		game:                     game,
	}

	return &b
}

func (b *Bullet) Update(tps float64) {
	nearbyShapes := b.solid.SelectTouchingCells(4).FilterShapes()
	s, c := math.Sincos(b.solid.Rotation())
	dt := 1.0 / tps
	b.elapsedTime += dt

	/*if b.altitude >= initialAltitude {*/
	b.solid.MoveVec(resolv.Vector{X: -s * bulletSpeed, Y: -c * bulletSpeed})
	b.solid.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: nearbyShapes.ByTags(TagEnemy | TagPlayer),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			b.solid.MoveVec(set.MTV)
			b.exploding = true
			b.hasHitTarget = true
			return true
		},
	})
	/*
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
	*/
	if b.exploding {
		b.explosionElapsedTime += dt * float64(len(b.explosionHitTargetFrames))
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if !b.exploding && !b.exploded {
		op := &ebiten.DrawImageOptions{}

		// center the bullet than scale
		/*
			op.GeoM.Translate(-bulletHalfW, -bulletHalfH)
			op.GeoM.Scale(b.scale, b.scale-math.Abs(b.currentSlope)*scaleCoeff) // simulate bullet deflection (poorly)
			op.GeoM.Translate(bulletHalfW, bulletHalfH)
		*/

		// center the bullet and the barrel than rotate
		op.GeoM.Translate(-b.bulletHalfW, -b.bulletHalfH)
		op.GeoM.Rotate(-b.solid.Rotation())
		op.GeoM.Translate(b.bulletHalfW, b.bulletHalfH)

		// actual position of the bullet to draw
		op.GeoM.Translate(b.solid.Center().X-b.bulletHalfW, b.solid.Center().Y-b.bulletHalfH)
		screen.DrawImage(b.sprite, op)
	} else {
		if b.hasHitTarget {
			b.drawExplosionHitTarget(screen)
		} else {
			b.drawExplosion(screen)
		}
		if b.exploded {
			b.game.space.Remove(b.solid)
		}
	}
}

func (b *Bullet) drawExplosionHitTarget(screen *ebiten.Image) {
	if int(b.explosionElapsedTime) < len(b.explosionHitTargetFrames) {
		frame := b.explosionHitTargetFrames[int(b.explosionElapsedTime)]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.solid.Center().X-float64(frame.Bounds().Dx()/2), b.solid.Center().Y-float64(frame.Bounds().Dy()/2))
		screen.DrawImage(frame, op)
	} else {
		b.exploded = true
	}

}

func (b *Bullet) drawExplosion(screen *ebiten.Image) {
	if int(b.explosionElapsedTime) < len(b.explosionFrames) {
		frame := b.explosionFrames[int(b.explosionElapsedTime)]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(b.solid.Center().X-float64(frame.Bounds().Dx()/2), b.solid.Center().Y-float64(frame.Bounds().Dy()/2))
		screen.DrawImage(frame, op)
	} else {
		b.exploded = true
	}

}
