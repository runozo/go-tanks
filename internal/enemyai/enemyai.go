// Package enemyai contains the pure stepping logic used to drive AI tanks.
// It has no dependency on ebiten, so it can be used both by the game client
// and by the server (which simulates the AI enemies in multiplayer).
//
// The formulas mirror the client-side Tank AI in game/tank.go.
package enemyai

import (
	"math"

	"github.com/solarlune/resolv"
)

// Bullet physics constants, mirrored from game/bullet.go.
const (
	bulletSpeed = 8.0
	gravity     = 9.8
)

// Cfg holds the tunable AI behavior parameters.
type Cfg struct {
	SpeedFactor    float64 // fraction of BaseSpeed the tank actually moves at
	BaseSpeed      float64 // px/s at 60 tps
	RotationSpeed  float64 // rad/s
	MinEngageDist  float64 // closer than this: back off
	MaxEngageDist  float64 // farther than this: approach; in between: orbit
	ProbeDist      float64 // how far ahead obstacle probes reach
	HeadingSamples int     // candidate headings around the desired one
	HeadingStep    float64 // angular spread between candidates
	StuckFramesMax int     // frames without real movement before evasive maneuver
	StuckProgress  float64 // px moved per frame below which counts as stuck
	EvasiveTime    int     // frames of sideways evasive movement

	// Playfield bounds (pixels). All zero disables bound steering.
	BoundsX0, BoundsY0, BoundsX1, BoundsY1 float64
	BoundsMargin                           float64
}

// DefaultCfg returns the same values used by the game client AI.
func DefaultCfg() Cfg {
	return Cfg{
		SpeedFactor:    0.85,
		BaseSpeed:      200,
		RotationSpeed:  math.Pi / 2,
		MinEngageDist:  140,
		MaxEngageDist:  360,
		ProbeDist:      72,
		HeadingSamples: 5,
		HeadingStep:    math.Pi / 5,
		StuckFramesMax: 40,
		StuckProgress:  1.5,
		EvasiveTime:    30,
		BoundsMargin:   90,
	}
}

// State is the per-enemy AI memory (persisted between frames).
type State struct {
	StuckFrames   int
	EvasiveFrames int
	LastPos       resolv.Vector
	OrbitSign     float64 // strafe/orbit direction (+1 or -1)
}

// Forward returns the tank's forward direction for a rotation R.
// Tanks move along (-sin R, -cos R).
func Forward(r float64) resolv.Vector {
	return resolv.Vector{X: -math.Sin(r), Y: -math.Cos(r)}
}

// RotationForHeading returns the rotation R whose forward direction matches h.
func RotationForHeading(h resolv.Vector) float64 {
	return math.Atan2(-h.X, -h.Y)
}

// BarrelWorldRotation returns the world rotation a barrel must have so that a
// bullet fired from `from` flies toward `to`.
func BarrelWorldRotation(from, to resolv.Vector) float64 {
	return -math.Atan2(to.Y-from.Y, to.X-from.X) - math.Pi/2
}

// BallisticSlope returns the barrel slope needed to hit a target at horizontal
// distance dist (mirrors game/tank.go), or -1 when out of range.
//
//	horizontalDistance = (bulletSpeed^2 * tps / gravity) * sin(2*slope)
func BallisticSlope(dist, tps float64) float64 {
	if dist <= 0 {
		return -1
	}

	maxRange := bulletSpeed * bulletSpeed * tps / gravity
	if dist > maxRange {
		return -1
	}

	sin2 := dist * gravity / (bulletSpeed * bulletSpeed * tps)
	if sin2 > 1 {
		sin2 = 1
	}
	if sin2 < -1 {
		sin2 = -1
	}

	return 0.5 * math.Asin(sin2)
}

// Steer computes the heading the enemy should steer toward, given its position,
// the target position and the obstacle centers. It updates the AI state
// (stuck detection, evasive maneuvers). The caller applies rotation and
// movement using Forward().
func Steer(cfg Cfg, st *State, pos, target resolv.Vector, obstacles []resolv.Vector) resolv.Vector {
	dx := target.X - pos.X
	dy := target.Y - pos.Y
	dist := math.Hypot(dx, dy)

	desired := resolv.Vector{X: dx, Y: dy}
	if dist > 0 {
		desired.X /= dist
		desired.Y /= dist
	}

	if dist > 0 && dist < cfg.MinEngageDist {
		// too close: back off
		desired.X, desired.Y = -desired.X, -desired.Y
	} else if dist <= cfg.MaxEngageDist {
		// engagement band: orbit the target (mostly sideways, slightly inward)
		if st.OrbitSign == 0 {
			st.OrbitSign = 1
		}
		perp := resolv.Vector{X: -desired.Y * st.OrbitSign, Y: desired.X * st.OrbitSign}
		desired = resolv.Vector{X: perp.X*0.8 + desired.X*0.2, Y: perp.Y*0.8 + desired.Y*0.2}
	}

	// evasive maneuver: briefly move sideways instead
	if st.EvasiveFrames > 0 {
		perp := resolv.Vector{X: -desired.Y * st.OrbitSign, Y: desired.X * st.OrbitSign}
		desired = perp
		st.EvasiveFrames--
	}

	// keep inside the playfield bounds: near an edge, steer back toward center
	if cfg.BoundsX1 > cfg.BoundsX0 && cfg.BoundsY1 > cfg.BoundsY0 {
		margin := cfg.BoundsMargin
		inBounds := pos.X >= cfg.BoundsX0+margin && pos.X <= cfg.BoundsX1-margin &&
			pos.Y >= cfg.BoundsY0+margin && pos.Y <= cfg.BoundsY1-margin
		if !inBounds {
			toCenter := resolv.Vector{
				X: (cfg.BoundsX0+cfg.BoundsX1)/2 - pos.X,
				Y: (cfg.BoundsY0+cfg.BoundsY1)/2 - pos.Y,
			}
			if l := math.Hypot(toCenter.X, toCenter.Y); l > 0 {
				toCenter.X /= l
				toCenter.Y /= l
			}
			desired.X += toCenter.X * 1.5
			desired.Y += toCenter.Y * 1.5
			if l := math.Hypot(desired.X, desired.Y); l > 0 {
				desired.X /= l
				desired.Y /= l
			}
		}
	}

	// obstacle avoidance: pick the sampled heading with the most clearance
	baseAngle := math.Atan2(desired.Y, desired.X)
	best := desired
	bestScore := math.Inf(-1)

	samples := cfg.HeadingSamples
	if samples < 1 {
		samples = 1
	}

	for i := 0; i < samples; i++ {
		offset := (float64(i) - float64(samples-1)/2) * cfg.HeadingStep
		angle := baseAngle + offset
		dir := resolv.Vector{X: math.Cos(angle), Y: math.Sin(angle)}

		probe := resolv.Vector{X: pos.X + dir.X*cfg.ProbeDist, Y: pos.Y + dir.Y*cfg.ProbeDist}
		clearance := clearance(probe, obstacles)

		// deviation penalty keeps the tank on the desired path when it is free
		score := clearance - math.Abs(offset)*60
		if score > bestScore {
			bestScore = score
			best = dir
		}
	}

	// stuck detection based on actual movement
	moved := math.Hypot(pos.X-st.LastPos.X, pos.Y-st.LastPos.Y)
	st.LastPos = pos
	if moved < cfg.StuckProgress {
		st.StuckFrames++
	} else {
		st.StuckFrames = 0
	}
	if st.StuckFrames > cfg.StuckFramesMax {
		st.StuckFrames = 0
		st.EvasiveFrames = cfg.EvasiveTime
	}

	return best
}

func clearance(p resolv.Vector, obstacles []resolv.Vector) float64 {
	best := math.MaxFloat64
	for _, o := range obstacles {
		d := math.Hypot(o.X-p.X, o.Y-p.Y)
		if d < best {
			best = d
		}
	}
	return best
}
