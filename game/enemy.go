package game

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	b2 "github.com/oliverbestmann/box2d-go"
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

	newPosition := b2.Vec2{
		X: float32(rand.Intn(screenWidth)),
		Y: float32(rand.Intn(screenHeight - tileHeight)),
	}

	// randomRotation := rand.Float64() * 2 * math.Pi
	newTank := NewTank(game, flavors[flavor][0], flavors[flavor][1], flavors[flavor][2], newPosition, 0.0)

	// TODO: don't overlap position with other enemies

	return &Enemy{
		game:          game,
		shootCooldown: NewTimer(time.Millisecond*2500 + time.Millisecond*time.Duration(rand.Intn(1000))),
		tank:          newTank,
		bullets:       make([]*Bullet, 0),
	}
}

func (e *Enemy) Update(tps float64) {
	playerPosition := e.game.players[0].Tank.tankBody.GetPosition()

	for i := 0; i < len(e.tank.barrels); i++ {
		pos := e.tank.barrels[i].tank.tankBody.GetPosition()
		e.tank.barrels[i].relativeRotation = math.Atan2(float64(playerPosition.Y-pos.Y), float64(playerPosition.X-pos.X)) + math.Pi/2.0
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
}
