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
	game         *Game
	tank         *Tank
	bulletSprite *ebiten.Image
	bullets      []*Bullet

	shootCooldown *Timer
}

func NewPlayer(game *Game) *Player {
	return &Player{
		game:          game,
		tank:          NewRandomTank(game),
		shootCooldown: NewTimer(shootCooldown),
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
		p.tank.rotation -= rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.tank.rotation += rotationSpeed
	}

	// rotate barrel
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.tank.barrel.rotation -= rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.tank.barrel.rotation += rotationSpeed
	}

	// move
	movementX := 0.0
	movementY := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		movementX += math.Sin(p.tank.rotation) * movementSpeed
		movementY -= math.Cos(p.tank.rotation) * movementSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		movementX -= math.Sin(p.tank.rotation) * movementSpeed
		movementY += math.Cos(p.tank.rotation) * movementSpeed
	}
	p.tank.position.X += movementX
	p.tank.position.Y += movementY

	// shoot
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.tank.barrel.slope += slopeSpeed
		if p.tank.barrel.slope > maxSlope {
			p.tank.barrel.slope = maxSlope
		}
		fmt.Println(p.tank.barrel.slope)
	}

	if p.tank.barrel.slope > 0.0 && inpututil.IsKeyJustReleased(ebiten.KeySpace) || p.tank.barrel.slope >= maxSlope {

		p.bullets = append(p.bullets, p.tank.Fire())
		// p.bullets = append(p.bullets, p.tank.Fire())
		p.tank.barrel.slope = 0.0
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
	p.tank.Draw(screen, p.tank.rotation)
	for _, bullet := range p.bullets {
		bullet.Draw(screen)
	}
	// fmt.Println("Bullets", len(p.bullets))
}
