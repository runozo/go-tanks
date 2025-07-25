package game

import (
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

	maxSlope = math.Pi / 4
)

type Player struct {
	game          *Game
	tank          *Tank
	bullets       []*Bullet
	shootCooldown *Timer
}

func NewPlayer(game *Game) *Player {
	return &Player{
		game:          game,
		tank:          NewRandomTank(game),
		shootCooldown: NewTimer(shootCooldown),
	}
}

func (p *Player) Update(tps float64) {
	rotationSpeed := rotationPerSecond / tps
	movementSpeed := tankSpeed / tps
	slopeSpeed := maxSlope / tps
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
		p.tank.barrel.relativeRotation -= rotationSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.tank.barrel.relativeRotation += rotationSpeed
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

	// charge shoot
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.tank.barrel.slope += slopeSpeed
		if p.tank.barrel.slope > maxSlope {
			p.tank.barrel.slope = maxSlope
		}
		// fmt.Println(p.tank.barrel.slope)
	}

	// fire
	if (p.tank.barrel.slope > 0.0 && inpututil.IsKeyJustReleased(ebiten.KeySpace) || p.tank.barrel.slope >= maxSlope) && p.shootCooldown.IsReady() {

		p.bullets = append(p.bullets, p.tank.Fire())

		p.tank.barrel.slope = 0.0
		p.shootCooldown.Reset()
		// fmt.Println(len(p.bullets))
	}

	// new tank
	if ebiten.IsKeyPressed(ebiten.KeyT) && inpututil.IsKeyJustPressed(ebiten.KeyT) {
		p.tank = NewRandomTank(p.game)
	}

	// update tank(s)
	p.tank.Update(tps)

	// update bullets
	var activeBullets []*Bullet
	for _, bullet := range p.bullets {
		bullet.Update(tps)
		if bullet.altitude > 0.0 {
			activeBullets = append(activeBullets, bullet)
		}
	}
	p.bullets = activeBullets
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.tank.Draw(screen, p.tank.rotation)
	for _, bullet := range p.bullets {
		bullet.Draw(screen)
	}
	// fmt.Println("Bullets", len(p.bullets))
}
