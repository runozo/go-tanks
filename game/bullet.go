package game

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	gravity         = 9.8
	bulletSpeed     = 8.0
	bulletMinScale  = 1.0
	scaleCoeff      = 1.8
	initialAltitude = 0.2
	bulletType1     = "bulletSand3_outline"
	bulletType2     = "bulletGreen3_outline"
	bulletType3     = "bulletBlue3_outline"
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
	barrel                     *Barrel
}

func NewBullet(barrel *Barrel, flavor string) *Bullet {
	bulletSprite := barrel.tank.game.assets.GetSprite(flavor)
	bulletSpriteWidth := float64(bulletSprite.Bounds().Dx())
	bulletSpriteHeight := float64(bulletSprite.Bounds().Dy())

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

	barrel.tank.game.space.Add(solid)

	b := Bullet{
		solid:         solid,
		sprite:        bulletSprite,
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
		barrel:        barrel,
		explosion:     nil,
	}

	return &b
}

func (b *Bullet) Update(tps float64) {
	if b.explosion != nil {
		b.explosion.Update(tps)
		return
	}

	dt := 1.0 / tps
	b.elapsedTime += dt

	// Physics Calculation (Gravity and Altitude)
	b.verticalSpeed -= gravity * dt
	b.altitude += b.verticalSpeed * dt

	actualSpeed := bulletSpeed * math.Cos(b.initialSlope)
	b.currentSlope = math.Atan2(b.verticalSpeed, actualSpeed)

	// We calculate the visual scale to simulate the height (the higher it is, the bigger it is)
	if b.altitude > initialAltitude {
		b.scale = (b.altitude * scaleCoeff) + bulletMinScale
	} else {
		b.scale = bulletMinScale
	}

	s, c := math.Sincos(b.solid.Rotation())
	// We use actualSpeed ​​instead of bulletSpeed ​​to slow down the horizontal motion if fired very high
	b.solid.MoveVec(resolv.Vector{X: -s * actualSpeed, Y: -c * actualSpeed})

	// Let's prevent it from exploding instantly (elapsedTime > 0.1) and check that it has "fallen"
	if b.altitude <= initialAltitude && b.elapsedTime > 0.1 {

		nearbyShapes := b.solid.SelectTouchingCells(2).FilterShapes().ByTags(TagEnemy | TagPlayer | TagObstacle)

		b.solid.IntersectionTest(resolv.IntersectionTestSettings{
			TestAgainst: nearbyShapes,
			OnIntersect: func(set resolv.IntersectionSet) bool {
				b.solid.MoveVec(resolv.Vector{X: 0, Y: 0})

				if b.barrel.tank.IsEnemy && set.OtherShape.Tags().Has(TagPlayer) {
					b.hasHitTarget = true
					b.barrel.tank.game.scoreLine.IncrementEnemy()
				}
				if !b.barrel.tank.IsEnemy && set.OtherShape.Tags().Has(TagEnemy) {
					b.hasHitTarget = true
					b.barrel.tank.game.scoreLine.IncrementPlayer()
				}
				return true // Stops the test at the first entity hit
			},
		})

		// The projectile fell: it explodes regardless of whether it hit someone or the empty ground
		b.explosion = NewExplosion(b)
		b.barrel.tank.game.space.Remove(b.solid)
	}
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	if b.explosion != nil {
		b.explosion.Draw(screen)
		return
	}

	if b.altitude > 0.1 {
		// Bullet shadow
		opShadow := &ebiten.DrawImageOptions{}

		opShadow.GeoM.Translate(-b.bulletHalfW, -b.bulletHalfH)

		// Scale shadow
		shadowScale := 1.0 / b.scale
		opShadow.GeoM.Scale(shadowScale, shadowScale)

		opShadow.GeoM.Rotate(-b.solid.Rotation())

		opShadow.GeoM.Translate(b.solid.Center().X, b.solid.Center().Y)

		shadowAlpha := float32(1.0 * shadowScale)
		opShadow.ColorScale.Scale(0, 0, 0, shadowAlpha)

		screen.DrawImage(b.sprite, opShadow)
	}

	// Bullet
	scaleY := b.scale - math.Abs(b.currentSlope)*scaleCoeff
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-b.bulletHalfW, -b.bulletHalfH)

	if scaleY < 0.1 {
		scaleY = 0.1
	}
	op.GeoM.Scale(b.scale, scaleY)

	op.GeoM.Rotate(-b.solid.Rotation())

	visualOffsetY := -b.altitude * 12.0
	op.GeoM.Translate(b.solid.Center().X, b.solid.Center().Y+visualOffsetY)

	screen.DrawImage(b.sprite, op)
}
