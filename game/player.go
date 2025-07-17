package game

import (
	"fmt"
	_ "image/png"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	shootCooldown     = time.Millisecond * 250
	rotationPerSecond = math.Pi
	tankSpeed         = 120.0

	bulletSpawnOffset = 20.0
	bulletSpeed       = 10.0
	maxSlope          = 1.0
)

type Player struct {
	game *Game

	tank           *Tank
	tankRotation   float64
	barrelRotation float64
	barrelSlope    float64
	position       Vector
	bulletSprite   *ebiten.Image
	bullets        []*Bullet

	shootCooldown *Timer
}

func NewPlayer(game *Game) *Player {
	bulletSprite := game.assets.GetSprite("bulletRed2")

	return &Player{
		game:           game,
		tankRotation:   0,
		barrelRotation: 0,
		tank:           NewRandomTank(game),
		position:       Vector{X: screenWidth / 2, Y: screenHeight / 2},
		shootCooldown:  NewTimer(shootCooldown),
		bulletSprite:   bulletSprite,
		barrelSlope:    0.0,
	}
}

func (p *Player) Update() {
	tps := float64(ebiten.TPS())
	rotationSpeed := rotationPerSecond / tps
	movementSpeed := tankSpeed / tps
	slopeSpeed := maxSlope / tps / 2
	p.shootCooldown.Update()

	// rotate tank
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.tankRotation -= rotationSpeed
		p.barrelRotation -= rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.tankRotation += rotationSpeed
		p.barrelRotation += rotationSpeed
	}

	// rotate barrel
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.barrelRotation -= rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.barrelRotation += rotationSpeed
	}

	// move
	movementX := 0.0
	movementY := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		movementX += math.Sin(p.tankRotation) * movementSpeed
		movementY -= math.Cos(p.tankRotation) * movementSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		movementX -= math.Sin(p.tankRotation) * movementSpeed
		movementY += math.Cos(p.tankRotation) * movementSpeed
	}
	p.position.X += movementX
	p.position.Y += movementY

	// shoot
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.barrelSlope += slopeSpeed
		if p.barrelSlope > maxSlope {
			p.barrelSlope = maxSlope
		}
		fmt.Println(p.barrelSlope)
	}

	if p.barrelSlope > 0.0 && inpututil.IsKeyJustReleased(ebiten.KeySpace) || p.barrelSlope >= maxSlope {

		tankBounds := p.tank.BodySprite.Bounds()
		bulletBounds := p.bulletSprite.Bounds()
		halfWBullet := bulletBounds.Dx() / 2
		halfHBullet := bulletBounds.Dy() / 2
		halfW := float64(tankBounds.Dx()) / 2
		halfH := float64(tankBounds.Dy()) / 2

		spawnPos := Vector{
			p.position.X + halfW - math.Cos(p.barrelRotation)*float64(halfWBullet) + math.Sin(p.barrelRotation)*bulletSpawnOffset,
			p.position.Y + halfH - math.Sin(p.barrelRotation)*float64(halfHBullet) + math.Cos(p.barrelRotation)*-bulletSpawnOffset,
		}

		p.bullets = append(p.bullets, NewBullet(p.bulletSprite, spawnPos, p.barrelRotation, bulletSpeed, p.barrelSlope))
		p.barrelSlope = 0.0
		p.shootCooldown.Reset()
		// fmt.Println(len(p.bullets))
	}

	// new tank
	if ebiten.IsKeyPressed(ebiten.KeyT) && inpututil.IsKeyJustPressed(ebiten.KeyT) {
		p.tank = NewRandomTank(p.game)
	}

	var visibleBullets []*Bullet
	for _, bullet := range p.bullets {
		bullet.Update()
		if bullet.altitude > 0.0 {
			visibleBullets = append(visibleBullets, bullet)
		}
	}
	p.bullets = visibleBullets
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.tank.Draw(screen, p.position, p.tankRotation, p.barrelRotation)
	for _, bullet := range p.bullets {
		bullet.Draw(screen)
	}
	// fmt.Println("Bullets", len(p.bullets))
}
