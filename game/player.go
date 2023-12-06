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

	bulletSpawnOffset = 50.0
)

type Player struct {
	game *Game

	tank *Tank

	shootCooldown *Timer
}

func NewPlayer(game *Game) *Player {
	return &Player{
		game:          game,
		tank:          NewTank(game),
		shootCooldown: NewTimer(shootCooldown),
	}
}

func (p *Player) Update() {

}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.tank.body.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	position := Vector{
		X: screenWidth/2 - halfW,
		Y: screenHeight/2 - halfH,
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(position.X, position.Y)
	screen.DrawImage(p.tank.body, op)
	op.GeoM.Translate(float64(p.tank.barrel.Bounds().Dx()-1), float64(-p.tank.barrel.Bounds().Dy()/2))
	screen.DrawImage(p.tank.barrel, op)
}
