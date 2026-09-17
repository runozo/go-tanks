package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/runozo/go-tanks/internal/assets"
	"github.com/solarlune/resolv"
)

type Obstacle struct {
	Object *resolv.ConvexPolygon
	Sprite *ebiten.Image
}

// NewObstacleAt places an obstacle from an authored map entry (ObstacleSpawn).
// If the requested sprite type does not exist it falls back to a random one.
func NewObstacleAt(game *Game, spawn ObstacleSpawn) *Obstacle {
	t, ok := game.assets.TileEntries[spawn.Type]
	if !ok {
		t = randomObstacleEntry(game)
	}
	o := &Obstacle{
		Object: resolv.NewRectangle(spawn.X, spawn.Y, float64(t.Width), float64(t.Height)),
		Sprite: game.assets.GetSprite(t.Name),
	}
	o.Object.Tags().Set(TagObstacle)
	game.space.Add(o.Object)

	return o
}

func randomObstacleEntry(game *Game) assets.TileEntry {
	choices := make([]assets.TileEntry, 0)
	for _, t := range game.assets.TileEntries {
		if t.Type == "obstacle" {
			choices = append(choices, t)
		}
	}
	return choices[rand.Intn(len(choices))]
}

func NewObstacle(game *Game) *Obstacle {
	t := randomObstacleEntry(game)
	o := &Obstacle{
		Object: resolv.NewRectangle(float64(rand.Intn(screenWidth-tileWidth)), float64(rand.Intn(screenHeight-tileHeight)), float64(t.Width), float64(t.Height)),
		Sprite: game.assets.GetSprite(t.Name),
	}
	o.Object.Tags().Set(TagObstacle)
	game.space.Add(o.Object)

	return o
}

func (f *Obstacle) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(f.Object.Center().X-float64(f.Sprite.Bounds().Dx()/2), f.Object.Center().Y-float64(f.Sprite.Bounds().Dy()/2))
	screen.DrawImage(f.Sprite, op)
}
