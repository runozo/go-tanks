package game

import (
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type Enemy struct {
	game          *Game
	tank          *Tank
	bullets       []*Bullet
	shootCooldown *Timer
}

func FlipVertical(source *ebiten.Image) *ebiten.Image {
	flipped := ebiten.NewImage(source.Bounds().Dx(), source.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(1, -1)
	op.GeoM.Translate(0, float64(source.Bounds().Dy()))
	flipped.DrawImage(source, op)
	return flipped
}

func NewEnemy(game *Game, flavor string) *Enemy {
	// flavors of enemies
	flavors := map[string][]string{
		// body name, barrel name, bullet name
		"easy":   []string{"tankBody_darkLarge", "specialBarrel1_outline", "bulletRed1_outline"},
		"medium": []string{"tankBody_darkLarge_outline", "specialBarrel1_outline", "bulletRed1_outline"},
		"hard":   []string{"tankBody_huge_outline", "specialBarrel1_outline", "bulletRed1_outline"},
	}

	// create a new enemy tank
	newTank := NewTank(game, flavors[flavor][0], flavors[flavor][1], resolv.Vector{
		X: float64(rand.Intn(screenWidth)),
		Y: float64(rand.Intn(screenHeight - tileHeight)),
	}, 0)

	// don't overlap position with other entities
	for newTank.solid.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: game.space.Shapes(),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			newTank.solid.SetPositionVec(resolv.Vector{
				X: float64(rand.Intn(screenWidth)),
				Y: float64(rand.Intn(screenHeight - tileHeight)),
			})
			return false
		},
	}) {
	}

	return &Enemy{
		game:          game,
		shootCooldown: NewTimer(time.Millisecond*2500 + time.Millisecond*time.Duration(rand.Intn(1000))),
		tank:          newTank,
		bullets:       make([]*Bullet, 0),
	}
}

func (e *Enemy) Update(tps float64) {
	playerPosition := e.game.players[0].tank.solid.Position()

	for i := 0; i < len(e.tank.barrels); i++ {
		e.tank.barrels[i].relativeRotation = -math.Atan2(playerPosition.Y-e.tank.solid.Position().Y, playerPosition.X-e.tank.solid.Position().X) - math.Pi/2
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
