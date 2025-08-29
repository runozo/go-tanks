package game

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	enemyOffset = 5
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

func NewEnemy(game *Game, flavor string) *Enemy {
	flavors := map[string][]string{
		// body barrel bullet
		"easy":   []string{"tankBody_darkLarge", "specialBarrel1_outline", "bulletRed1_outline"},
		"medium": []string{"tankBody_darkLarge_outline", "specialBarrel1_outline", "bulletRed1_outline"},
		"hard":   []string{"tankBody_huge_outline", "specialBarrel1_outline", "bulletRed1_outline"},
	}

	position := Vector{
		X: screenWidth/2 - screenWidth/4,
		Y: float64(rand.Intn(screenHeight - tileHeight)),
	}

	// don't overlap position with other enemies
	for _, e := range game.enemies {
		for position.Y > e.tank.Position.Y-enemyOffset && position.Y < e.tank.Position.Y+e.tank.bodyHeight+enemyOffset {
			position.Y = float64(rand.Intn(screenHeight - tileHeight))
			fmt.Println("new position", position.Y, len(game.enemies))
		}
	}

	return &Enemy{
		game:          game,
		shootCooldown: NewTimer(time.Millisecond*2500 + time.Millisecond*time.Duration(rand.Intn(1000))),
		tank:          NewTank(game, flavors[flavor][0], flavors[flavor][1], flavors[flavor][2], position, 0),
		bullets:       make([]*Bullet, 0),
	}
}

func (e *Enemy) Update(tps float64) {
	playerPosition := e.game.players[0].Tank.Position

	for i := 0; i < len(e.tank.barrels); i++ {
		e.tank.barrels[i].relativeRotation = math.Atan2(playerPosition.Y-e.tank.Position.Y, playerPosition.X-e.tank.Position.X) + math.Pi/2
	}

	// fires randomly
	if e.shootCooldown.IsReady() {
		e.shootCooldown.Reset()
		randomSlope := rand.Float64() * barrelMaxSlope
		for i := 0; i < len(e.tank.barrels); i++ {
			e.tank.barrels[i].slope = randomSlope
		}
		e.bullets = append(e.bullets, e.tank.Fire()...)
	}

	e.shootCooldown.Update()

	// update tank
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
