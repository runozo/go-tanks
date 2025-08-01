package game

import "github.com/hajimehoshi/ebiten/v2"

type Explosion struct {
	age      float64
	frameSet []*ebiten.Image
}

func NewExplosion(frames []*ebiten.Image) *Explosion {
	return &Explosion{
		age:      0,
		frameSet: frames,
	}
}

func (e *Explosion) Update(tps float64) {
	e.age += 1 / tps * 10
}

func (e *Explosion) Draw(screen *ebiten.Image, x, y float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)
	screen.DrawImage(e.frameSet[int(e.age)%len(e.frameSet)], op)
}
