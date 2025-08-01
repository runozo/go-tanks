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
	frame := e.frameSet[int(e.age)%len(e.frameSet)]
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frame.Bounds().Dx())/2, -float64(frame.Bounds().Dy())/2)
	op.GeoM.Translate(x, y)
	screen.DrawImage(frame, op)
}
