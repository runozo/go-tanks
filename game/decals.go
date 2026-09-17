package game

import (
	"math/rand"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	trackSpacing  = 26.0 // px of travel between successive track decals
	trackLifetime = 6.0  // seconds before a track decal fades away
	trackCap      = 600  // maximum number of decals kept at once
)

var trackNames = []string{"tracksLarge", "tracksSmall", "tracksDouble"}

// TrackDecal is a single faded tank-track imprint on the terrain.
type TrackDecal struct {
	x, y     float64
	rotation float64
	sprite   *ebiten.Image
	age      float64
}

// leaveTracks may drop a track imprint behind the tank, but only on soft ground
// (anything that is not a road). Each track sprite already shows both treads,
// so a single decal centered on the tank (at its rotation) is enough.
func (g *Game) leaveTracks(t *Tank) {
	if g.playfield == nil {
		return
	}

	center := t.Object.Center()
	if !g.isSoftGround(center) {
		return
	}

	spr := g.assets.GetSprite(trackNames[rand.Intn(len(trackNames))])
	if spr == nil {
		return
	}

	g.addDecal(center.X, center.Y, t.Object.Rotation(), spr)
}

func (g *Game) isSoftGround(p resolv.Vector) bool {
	name := g.playfield.TileAt(int(p.X), int(p.Y))
	if name == "" {
		return false
	}
	// grass, sand and transitions are soft; roads (any "_road*" tile) are not.
	return !strings.Contains(name, "_road")
}

func (g *Game) addDecal(x, y, rot float64, spr *ebiten.Image) {
	g.decals = append(g.decals, TrackDecal{x: x, y: y, rotation: rot, sprite: spr})
	if excess := len(g.decals) - trackCap; excess > 0 {
		g.decals = g.decals[excess:]
	}
}

func (g *Game) updateDecals(tps float64) {
	dt := 1.0 / tps
	keep := g.decals[:0]
	for i := range g.decals {
		d := &g.decals[i]
		d.age += dt
		if d.age < trackLifetime {
			keep = append(keep, *d)
		}
	}
	g.decals = keep
}

func (g *Game) drawDecals(screen *ebiten.Image) {
	for i := range g.decals {
		d := &g.decals[i]
		alpha := 1.0 - d.age/trackLifetime
		if alpha < 0 {
			alpha = 0
		}

		w := float64(d.sprite.Bounds().Dx())
		h := float64(d.sprite.Bounds().Dy())

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Rotate(-d.rotation)
		op.GeoM.Translate(w/2, h/2)
		op.GeoM.Translate(d.x-w/2, d.y-h/2)
		op.ColorScale.ScaleAlpha(float32(alpha))

		screen.DrawImage(d.sprite, op)
	}
}
