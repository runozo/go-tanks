package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type Explosion struct {
	bullet       *Bullet
	solid        *resolv.ConvexPolygon
	frameCounter float64
	frames       []*ebiten.Image
	hasHitTarget bool
}

func NewExplosion(bullet *Bullet) *Explosion {
	var frames []*ebiten.Image

	if bullet.hasHitTarget {
		frames = []*ebiten.Image{
			bullet.game.assets.GetSprite("explosion1"),
			bullet.game.assets.GetSprite("explosion2"),
			bullet.game.assets.GetSprite("explosion3"),
			bullet.game.assets.GetSprite("explosion4"),
			bullet.game.assets.GetSprite("explosion5"),
		}
	} else {
		frames = []*ebiten.Image{
			bullet.game.assets.GetSprite("explosionSmoke1"),
			bullet.game.assets.GetSprite("explosionSmoke2"),
			bullet.game.assets.GetSprite("explosionSmoke3"),
			bullet.game.assets.GetSprite("explosionSmoke4"),
			bullet.game.assets.GetSprite("explosionSmoke5"),
		}
	}

	e := &Explosion{
		bullet:       bullet,
		solid:        nil,
		frameCounter: 0.0,
		frames:       frames,
		hasHitTarget: bullet.hasHitTarget,
	}

	e.solid = resolv.NewRectangle(bullet.solid.Center().X, bullet.solid.Center().Y, float64(frames[1].Bounds().Dx()), float64(frames[1].Bounds().Dy()))
	e.solid.Tags().Set(TagExplosion)
	bullet.game.space.Add(e.solid)
	e.solid.SetData(e)

	return e
}

func (e *Explosion) Update(tps float64) {
	e.frameCounter += 1.0 / tps * 5

	nearbyShapes := e.solid.SelectTouchingCells(2).FilterShapes().ByTags(TagEnemy | TagPlayer)
	e.solid.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: nearbyShapes,
		OnIntersect: func(set resolv.IntersectionSet) bool {
			return true
		},
	})

	if int(e.frameCounter) > len(e.frames) {
		e.Destroy()
		e.bullet.exploded = true
	}
}

func (e *Explosion) Destroy() {
	e.bullet.game.space.Remove(e.solid)
}

func (e *Explosion) Draw(screen *ebiten.Image) {
	if int(e.frameCounter) < len(e.frames) {
		frame := e.frames[int(e.frameCounter)]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(e.solid.Center().X-float64(frame.Bounds().Dx()/2), e.solid.Center().Y-float64(frame.Bounds().Dy()/2))
		screen.DrawImage(frame, op)
	}
}
