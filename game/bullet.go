package game

import (
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	b2 "github.com/oliverbestmann/box2d-go"
)

const (
	gravity         = 9.8
	bulletSpeed     = 2.0
	bulletMinScale  = 1.0
	scaleCoeff      = 1.8
	initialAltitude = 0.2
)

type Bullet struct {
	bulletBody                 b2.Body
	sprite                     *ebiten.Image
	verticalSpeed, altitude    float64
	currentSlope, initialSlope float64
	scale                      float64
	elapsedTime                float64
	spriteWidth, spriteHeight  float32
	barrelWidth                float32
	barrelHeight               float32
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
	bulletSpriteWidth := float32(bulletSprite.Bounds().Dx())
	bulletSpriteHeight := float32(bulletSprite.Bounds().Dy())

	position := b2.Vec2{
		X: barrel.position.X + barrel.spriteWidth/2 - bulletSpriteWidth/2,
		Y: barrel.position.Y - bulletSpriteHeight,
	}

	// create box2d body
	bodyDef := b2.DefaultBodyDef()
	bodyDef.Position = position

	s, c := math.Sincos(barrel.absoluteRotation)
	bodyDef.Rotation = b2.Rot{C: float32(c), S: float32(s)}
	bodyDef.Type1 = b2.DynamicBody
	shape := b2.MakeBox(float32(bulletSpriteWidth/2.0), float32(bulletSpriteHeight/2.0))
	body := game.world.CreateBody(bodyDef)
	// attach shape for collisions
	body.CreatePolygonShape(b2.DefaultShapeDef(), shape)
	// fmt.Println("Barrel position", barrel.position.X, barrel.position.Y, "Bullet position", position.X, position.Y)

	return &Bullet{
		bulletBody:               body,
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

	sinRot, cosRot := math.Sincos(float64(b.bulletBody.GetRotation().Angle()))
	dt := 1.0 / tps
	b.elapsedTime += dt

	//if b.altitude >= initialAltitude {
	b.bulletBody.SetLinearVelocity(b2.Vec2{X: float32(sinRot * bulletSpeed * 0.1), Y: float32(cosRot * bulletSpeed * -0.1)})

	gravityEffect := 0.5 * gravity * dt * dt
	b.altitude += b.verticalSpeed*dt - gravityEffect
	b.verticalSpeed -= gravity * dt

	actualSpeed := bulletSpeed * math.Cos(b.initialSlope)
	b.currentSlope = math.Atan2(b.verticalSpeed, actualSpeed)
	b.scale = b.altitude*scaleCoeff + bulletMinScale
	// fmt.Println(b.currentSlope)
	//} else {
	//	b.exploding = true
	//}

	// TODO: check collision with players
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
		op.GeoM.Translate(float64(-bulletHalfW), float64(-bulletHalfH))
		op.GeoM.Scale(b.scale, b.scale-math.Abs(b.currentSlope)*scaleCoeff) // simulate bullet deflection (poorly)
		op.GeoM.Translate(float64(bulletHalfW), float64(bulletHalfH))

		// center the bullet and the barrel than rotate
		op.GeoM.Translate(float64(-bulletHalfW), float64(-bulletAndBarrellHeight))
		op.GeoM.Rotate(float64(b.bulletBody.GetRotation().Angle()))
		op.GeoM.Translate(float64(bulletHalfW), float64(bulletAndBarrellHeight))

		// true position of the bullet to draw
		pos := b.bulletBody.GetPosition()
		op.GeoM.Translate(float64(pos.X), float64(pos.Y))
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
		pos := b.bulletBody.GetPosition()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(pos.X), float64(pos.Y))
		screen.DrawImage(frame, op)
	} else {
		b.exploded = true
	}

}

func (b *Bullet) drawExplosion(screen *ebiten.Image) {
	if int(b.explosionElapsedTime) < len(b.explosionFrames) {
		frame := b.explosionFrames[int(b.explosionElapsedTime)]
		pos := b.bulletBody.GetPosition()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(pos.X), float64(pos.Y))
		screen.DrawImage(frame, op)
	} else {
		b.exploded = true
	}

}
