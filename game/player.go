package game

import (
	_ "image/png"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	shootCooldown     = time.Millisecond * 500
	rotationPerSecond = math.Pi
	pixelPerSecond    = 100.0

	bulletSpawnOffset = 50.0
)

type Player struct {
	game *Game

	tank     *Tank
	rotation float64
	position Vector
	bullets  []*Bullet

	shootCooldown *Timer
}

func NewPlayer(game *Game) *Player {
	return &Player{
		game:          game,
		rotation:      0,
		tank:          NewTank(game),
		position:      Vector{X: screenWidth / 2, Y: screenHeight / 2},
		shootCooldown: NewTimer(shootCooldown),
	}
}

func (p *Player) Update() {
	rotationSpeed := rotationPerSecond / float64(ebiten.TPS())
	movementSpeed := pixelPerSecond / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= rotationSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += rotationSpeed
	}

	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.shootCooldown.Reset()

		bounds := p.tank.body.Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 2

		spawnPos := Vector{
			p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
			p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
		}

		bullet := NewBullet(p.game, spawnPos, p.rotation)
		p.bullets = append(p.bullets, bullet)
	}

	// move towards facing
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		//first get the direction the entity is pointed
		dx := math.Sin(p.rotation)
		dy := math.Cos(p.rotation)
		// if direction.length() > 0 {
		// 	direction = direction.normalise()
		// }
		p.position.X += dx * movementSpeed
		p.position.Y -= dy * movementSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		//first get the direction the entity is pointed
		dx := math.Sin(p.rotation)
		dy := math.Cos(p.rotation)
		// if direction.length() > 0 {
		// 	direction = direction.normalise()
		// }
		p.position.X -= dx * movementSpeed
		p.position.Y += dy * movementSpeed
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.tank.Draw(screen, p.position, p.rotation)
}
