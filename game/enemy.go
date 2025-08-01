package game

import "github.com/hajimehoshi/ebiten/v2"

type Enemy struct {
	game          *Game
	tank          *Tank
	bullets       []*Bullet
	shootCooldown *Timer
}

func NewEnemy(game *Game) *Enemy {
	bodySprite := game.assets.GetSprite("tankBody_darkLarge")
	barrelSprite := game.assets.GetSprite("specialBarrel1_outline")
	bulletSprite := game.assets.GetSprite("bulletRed1_outline")

	return &Enemy{
		game:          game,
		shootCooldown: NewTimer(100),
		tank:          NewTank(game, bodySprite, barrelSprite, bulletSprite, screenWidth/2-screenWidth/4, screenHeight/2, 0),
		bullets:       make([]*Bullet, 0),
	}
}

func (e *Enemy) Update(tps float64) {
	e.tank.Update(tps)
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	e.tank.Draw(screen)
}
