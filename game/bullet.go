package game

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	game *Game

	position Vector
	rotation float64
	sprite   *ebiten.Image

	shootCooldown *Timer
}

func NewBullet(game *Game, position Vector, rotation float64) *Bullet {
	sprite := game.assets.GetSprite("bulletRed1.png")

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := Vector{
		X: screenWidth/2 - halfW,
		Y: screenHeight/2 - halfH,
	}

	return &Bullet{
		game:          game,
		position:      pos,
		rotation:      0,
		sprite:        sprite,
		shootCooldown: NewTimer(shootCooldown),
	}
}

func (p *Bullet) Update() {

}

func (p *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.position.X, p.position.Y)
	screen.DrawImage(p.sprite, op)
}
