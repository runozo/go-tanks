package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/runozo/go-wave-function-collapse/assets"
	"github.com/solarlune/resolv"
)

type Obstacle struct {
	Object *resolv.ConvexPolygon
	Sprite *ebiten.Image
}

func NewObstacle(game *Game) *Obstacle {
	choices := make([]assets.TileEntry, 0)
	for _, t := range game.assets.TileEntries {
		if t.Type == "obstacle" {
			choices = append(choices, t)
		}
	}

	// pick random furniture
	t := choices[rand.Intn(len(choices))]
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
