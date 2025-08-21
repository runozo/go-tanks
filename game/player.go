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
	barrelMaxSlope    = math.Pi / 4
)

type Player struct {
	game          *Game
	tank          *Tank
	bullets       []*Bullet
	shootCooldown *Timer
}

func NewPlayer(game *Game) *Player {

	position := Vector{
		X: screenWidth/2 + screenWidth/4,
		Y: screenHeight / 2,
	}

	return &Player{
		game:          game,
		tank:          NewRandomTank(game, position, 0),
		shootCooldown: NewTimer(shootCooldown),
	}
}

func (p *Player) Update(tps float64) {
	rotationSpeed := rotationPerSecond / tps
	movementSpeed := tankSpeed / tps
	slopeSpeed := barrelMaxSlope / tps
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
		for i := 0; i < len(p.tank.barrels); i++ {
			p.tank.barrels[i].relativeRotation -= rotationSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		for i := 0; i < len(p.tank.barrels); i++ {
			p.tank.barrels[i].relativeRotation += rotationSpeed
		}
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
		for i := 0; i < len(p.tank.barrels); i++ {
			p.tank.barrels[i].slope += slopeSpeed
			if p.tank.barrels[i].slope > barrelMaxSlope {
				p.tank.barrels[i].slope = barrelMaxSlope
			}

		}
		// fmt.Println(p.tank.barrel.slope)
	}

	// fire
	for i := 0; i < len(p.tank.barrels); i++ {
		if (p.tank.barrels[i].slope > 0.0 && inpututil.IsKeyJustReleased(ebiten.KeySpace) || p.tank.barrels[i].slope >= barrelMaxSlope) && p.shootCooldown.IsReady() {
			p.shootCooldown.Reset()
			p.bullets = append(p.bullets, p.tank.Fire()...)
			p.tank.barrels[i].slope = 0.0
			// fmt.Println(len(p.bullets))
		}
	}

	// new tank
	if ebiten.IsKeyPressed(ebiten.KeyT) && inpututil.IsKeyJustPressed(ebiten.KeyT) {
		p.tank = NewRandomTank(p.game, p.tank.position, p.tank.rotation)
	}

	// update tank(s)
	p.tank.Update(tps)

	// update bullets
	var activeBullets []*Bullet
	for _, bullet := range p.bullets {
		bullet.Update(tps)
		if !bullet.exploded {
			activeBullets = append(activeBullets, bullet)
		}
	}
	p.bullets = activeBullets
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.tank.Draw(screen)
	for _, bullet := range p.bullets {
		bullet.Draw(screen)
	}
	// fmt.Println("Bullets", len(p.bullets))
}
