package game

import (
	"math"
	"testing"

	"github.com/solarlune/resolv"
)

// TestRotationForHeading verifies that, for any heading, the rotation returned
// by rotationForHeading makes the tank's forward direction (-sin R, -cos R)
// point exactly along that heading. A bug here made enemies drive sideways and
// wander off-screen.
func TestRotationForHeading(t *testing.T) {
	headings := []resolv.Vector{
		{X: 1, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: -1},
		{X: 1, Y: 1},
		{X: -1, Y: -1},
		{X: 0.6, Y: -0.8},
	}
	for _, h := range headings {
		l := math.Hypot(h.X, h.Y)
		u := resolv.Vector{X: h.X / l, Y: h.Y / l}

		r := rotationForHeading(u)
		s, c := math.Sincos(r)
		forward := resolv.Vector{X: -s, Y: -c}

		if math.Abs(forward.X-u.X) > 1e-9 || math.Abs(forward.Y-u.Y) > 1e-9 {
			t.Errorf("heading (%v): rotation %.4f gives forward (%v), want (%v)",
				u, r, forward, u)
		}
	}
}

// TestBallisticSlope verifies the slope computed by ballisticSlope actually
// reaches the requested distance, using the bullet physics of bullet.go:
//
//	horizontalDistance = (bulletSpeed^2 * tps / gravity) * sin(2*slope)
func TestBallisticSlope(t *testing.T) {
	tank := &Tank{}
	tps := 60.0
	maxRange := bulletSpeed * bulletSpeed * tps / gravity

	for _, dist := range []float64{50, 100, 150, 200, 300, maxRange - 1} {
		slope := tank.ballisticSlope(dist, tps)
		if slope < 0 {
			t.Fatalf("slope for dist %.0f should be >= 0", dist)
		}
		reached := maxRange * math.Sin(2*slope)
		if math.Abs(reached-dist) > 0.5 {
			t.Errorf("dist %.0f: slope %.4f reaches %.1f (want %.1f)", dist, slope, reached, dist)
		}
	}

	// out of range must return -1
	if s := tank.ballisticSlope(maxRange+100, tps); s >= 0 {
		t.Errorf("out-of-range dist should return -1, got %.4f", s)
	}
	// dist 0 must return -1
	if s := tank.ballisticSlope(0, tps); s >= 0 {
		t.Errorf("dist 0 should return -1, got %.4f", s)
	}
}
