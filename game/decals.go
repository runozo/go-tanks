package game

import (
	"math"
	"math/rand"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

const (
	trackSpacing    = 26.0 // px of travel between successive track decals
	trackOffsetFrac = 0.55 // how far the two tracks sit from the tank center line
	trackLifetime   = 6.0  // seconds before a track decal fades away
	trackCap        = 600  // maximum number of decals kept at once
)

var trackNames = []string{"tracksLarge", "tracksSmall", "tracksDouble"}

// TrackDecal is a single faded tank-track imprint on the terrain.
type TrackDecal struct {
	x, y     float64
	rotation float64
	sprite   *ebiten.Image
	age      float64
}

// leaveTracks may drop a pair of track imprints (left + right track) behind the
// tank, but only on soft ground (anything that is not a road).
func (g *Game) leaveTracks(t *Tank) {
	if g.playfield == nil {
		return
	}

	center := t.Object.Center()
	rot := t.Object.Rotation()

	// perpendicular to the facing direction (-sin R, -cos R)
	p := resolv.Vector{X: -math.Cos(rot), Y: math.Sin(rot)}
	off := float64(t.Width) / 2 * trackOffsetFrac

	left := resolv.Vector{X: center.X - p.X*off, Y: center.Y - p.Y*off}
	right := resolv.Vector{X: center.X + p.X*off, Y: center.Y + p.Y*off}

	spr := g.assets.GetSprite(trackNames[rand.Intn(len(trackNames))])
	if spr == nil {
		return
	}

	if g.isSoftGround(left) {
		g.addDecal(left.X, left.Y, rot, spr)
	}
	if g.isSoftGround(right) {
		g.addDecal(right.X, right.Y, rot, spr)
	}
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
