package game

import (
	"fmt"
	_ "image/png"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	shootCooldown     = time.Millisecond * 500
	rotationPerSecond = math.Pi

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
	speed := rotationPerSecond / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= speed
		fmt.Println(p.rotation, p.position)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += speed
		fmt.Println(p.rotation, p.position)
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
}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.tank.body.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2
	// fmt.Println(bounds)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(p.position.X, p.position.Y)
	screen.DrawImage(p.tank.body, op)
	op.GeoM.Translate(float64(p.tank.barrel.Bounds().Dx()-1), float64(-p.tank.barrel.Bounds().Dy()/2))
	screen.DrawImage(p.tank.barrel, op)
}
