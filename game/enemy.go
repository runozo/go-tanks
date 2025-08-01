package game

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	game          *Game
	tank          *Tank
	bullets       []*Bullet
	shootCooldown *Timer
}

func FlipVertical(source *ebiten.Image) *ebiten.Image {
	result := ebiten.NewImage(source.Bounds().Dx(), source.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(source.Bounds().Dy()))
	result.DrawImage(source, op)
	return result
}

func NewEnemy(game *Game) *Enemy {
	bodySprite := game.assets.GetSprite("tankBody_darkLarge")
	barrelSprite := game.assets.GetSprite("specialBarrel1_outline")
	barrelSpriteReversed := FlipVertical(barrelSprite)
	bulletSprite := game.assets.GetSprite("bulletRed1_outline")

	return &Enemy{
		game:          game,
		shootCooldown: NewTimer(time.Millisecond * 250),
		tank:          NewTank(game, bodySprite, barrelSpriteReversed, bulletSprite, screenWidth/2-screenWidth/4, screenHeight/2, 0),
		bullets:       make([]*Bullet, 0),
	}
}

func (e *Enemy) Update(tps float64) {
	playerPosition := e.game.players[0].tank.position
	e.tank.barrel.relativeRotation = math.Atan2(playerPosition.Y-e.tank.position.Y, playerPosition.X-e.tank.position.X) + math.Pi/2
	// fmt.Println(e.tank.barrel.absoluteRotation)

	if e.shootCooldown.IsReady() {
		e.shootCooldown.Reset()
		e.tank.Fire()
		e.bullets = append(e.bullets, e.tank.Fire())
	}
	e.shootCooldown.Update()
	e.tank.Update(tps)
	// update bullets
	var activeBullets []*Bullet
	for _, bullet := range e.bullets {
		bullet.Update(tps)
		if !bullet.exploded {
			activeBullets = append(activeBullets, bullet)
		}
	}
	e.bullets = activeBullets

}

func (e *Enemy) Draw(screen *ebiten.Image) {
	e.tank.Draw(screen)
	for _, bullet := range e.bullets {
		bullet.Draw(screen)
	}
}
