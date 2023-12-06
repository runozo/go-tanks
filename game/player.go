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

	position Vector
	rotation float64
	sprite   *ebiten.Image

	shootCooldown *Timer
}

func NewPlayer(game *Game) *Player {
	sprite := game.assets.GetSprite("tankBody_red.png")
	// var TracksSmallSprite = mustLoadImage("png/tracksSmall.png")

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := Vector{
		X: screenWidth/2 - halfW,
		Y: screenHeight/2 - halfH,
	}

	return &Player{
		game:          game,
		position:      pos,
		rotation:      0,
		sprite:        sprite,
		shootCooldown: NewTimer(shootCooldown),
	}
}

func (p *Player) Update() {

}

func (p *Player) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.X, p.position.Y)
	screen.DrawImage(p.sprite, op)
}
