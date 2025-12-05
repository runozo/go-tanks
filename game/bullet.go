package game

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	gravity         = 9.8
	bulletSpeed     = 4.0
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
	explosion                  *Explosion
	exploding                  bool
	exploded                   bool
	hasHitTarget               bool
	game                       *Game
}

func NewBullet(game *Game, barrel *Barrel) *Bullet {
	bulletSpriteWidth := float64(barrel.bulletSprite.Bounds().Dx())
	bulletSpriteHeight := float64(barrel.bulletSprite.Bounds().Dy())

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
		solid:         solid,
		sprite:        barrel.bulletSprite,
		verticalSpeed: bulletSpeed * math.Sin(barrel.slope),
		currentSlope:  barrel.slope,
		initialSlope:  barrel.slope,
		altitude:      initialAltitude,
		scale:         bulletMinScale,
		elapsedTime:   0.0,
		bulletHalfW:   bulletSpriteWidth / 2.0,
		bulletHalfH:   bulletSpriteHeight / 2.0,
		exploding:     false,
		exploded:      false,
		hasHitTarget:  false,
		game:          game,
		explosion:     nil,
	}

	return &b
}

func (b *Bullet) Update(tps float64) {
	s, c := math.Sincos(b.solid.Rotation())
	dt := 1.0 / tps
	b.elapsedTime += dt

	/*if b.altitude >= initialAltitude {*/
	if !b.hasHitTarget {
		b.solid.MoveVec(resolv.Vector{X: -s * bulletSpeed, Y: -c * bulletSpeed})
		nearbyShapes := b.solid.SelectTouchingCells(2).FilterShapes().ByTags(TagEnemy | TagPlayer)
		b.solid.IntersectionTest(resolv.IntersectionTestSettings{
			TestAgainst: nearbyShapes,
			OnIntersect: func(set resolv.IntersectionSet) bool {
				b.solid.MoveVec(resolv.Vector{X: 0, Y: 0})
				b.hasHitTarget = true
				b.explosion = NewExplosion(b)
				b.game.space.Remove(b.solid)
				return true
			},
		})
	} else {
		b.explosion.Update(tps)
	}

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
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if b.explosion == nil {
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
		b.explosion.Draw(screen)
	}
}
