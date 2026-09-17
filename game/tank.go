package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

const (
	barrelMaxSlope    = math.Pi / 3
	rotationPerSecond = math.Pi / 2
	tankSpeed         = 200

	// Enemy AI behavior (tunable)
	enemySpeedFactor    = 0.85 // enemies move slightly slower than the player
	enemyMinEngageDist  = 140  // closer than this: back off
	enemyMaxEngageDist  = 360  // farther than this: approach; in between: orbit
	enemyProbeDist      = 72   // how far ahead we probe for obstacles
	enemyHeadingSamples = 5    // candidate headings sampled around the desired one
	enemyHeadingStep    = math.Pi / 5
	enemyStuckFramesMax = 40  // frames without real movement before evasive maneuver
	enemyStuckProgress  = 1.5 // px moved per frame below which we count as stuck
	enemyEvasiveTime    = 30  // frames of sideways evasive movement after getting stuck
)

type Tank struct {
	ID            string
	Sprite        *ebiten.Image
	Object        *resolv.ConvexPolygon
	barrels       []*Barrel
	Bullets       []*Bullet
	Width         float64
	Height        float64
	game          *Game
	Hit           bool
	IsEnemy       bool
	ShootCooldown *Timer
	// Enemy AI state
	enemyStuckFrames   int
	enemyEvasiveFrames int
	enemyLastPos       resolv.Vector
	enemyOrbitSign     float64 // strafe/orbit direction (+1 or -1), random per enemy
}

func NewTank(g *Game, bodySpriteName, barrelSpriteName, bulletSpriteName string, position resolv.Vector, rotation float64, isEnemy bool) *Tank {

	bodySprite := g.assets.GetSprite(bodySpriteName)
	spriteWidth := float64(bodySprite.Bounds().Dx())
	spriteHeight := float64(bodySprite.Bounds().Dy())

	tank := &Tank{
		ID:             uuid.New().String(),
		Sprite:         bodySprite,
		Object:         resolv.NewRectangle(position.X, position.Y, spriteWidth, spriteHeight),
		Width:          spriteWidth,
		Height:         spriteHeight,
		barrels:        make([]*Barrel, 0),
		Bullets:        make([]*Bullet, 0),
		game:           g,
		Hit:            false,
		IsEnemy:        isEnemy,
		ShootCooldown:  NewTimer(time.Millisecond*2500 + time.Millisecond*time.Duration(rand.Intn(1000))),
		enemyOrbitSign: 1,
	}

	tank.Object.Rotate(rotation)
	tank.Object.SetPositionVec(position)
	if isEnemy {
		tank.Object.Tags().Set(TagEnemy)
		if rand.Intn(2) == 0 {
			tank.enemyOrbitSign = -1
		}
	} else {
		tank.Object.Tags().Set(TagPlayer)
	}

	if bodySpriteName == "tankBody_huge_outline" {
		// this bodies have 2 barrels each
		tank.barrels = []*Barrel{
			NewBarrel(tank, barrelSpriteName, bulletSpriteName, resolv.Vector{X: 0, Y: -spriteHeight / 2 / 2}),
			NewBarrel(tank, barrelSpriteName, bulletSpriteName, resolv.Vector{X: 0, Y: spriteHeight / 2 / 2}),
		}
	} else {
		tank.barrels = []*Barrel{NewBarrel(tank, barrelSpriteName, bulletSpriteName, resolv.Vector{X: 0, Y: 0})}
	}

	tank.Object.SetData(tank)

	maxAttempts := 50
	validSpawn := false

	for attempt := 0; attempt < maxAttempts; attempt++ {
		collision := false

		tank.Object.IntersectionTest(resolv.IntersectionTestSettings{
			TestAgainst: g.space.Shapes(),
			OnIntersect: func(set resolv.IntersectionSet) bool {
				collision = true
				return false // Stop at first collision
			},
		})

		if !collision {
			validSpawn = true
			break // Space found!
		}

		if len(g.Tanks) > 0 {
			targetTank := g.Tanks[rand.Intn(len(g.Tanks))]

			// Random offset
			offsetX := (rand.Float64() - 0.5) * 500
			offsetY := (rand.Float64() - 0.5) * 500

			// Check if the new position is out of bounds
			newX := math.Max(0, math.Min(targetTank.Object.Position().X+offsetX, float64(screenWidth-int(spriteWidth))))
			newY := math.Max(0, math.Min(targetTank.Object.Position().Y+offsetY, float64(screenHeight-int(spriteHeight))))

			tank.Object.SetPositionVec(resolv.Vector{X: newX, Y: newY})
		} else {
			// Random if there are no tanks
			tank.Object.SetPositionVec(resolv.Vector{
				X: float64(rand.Intn(screenWidth - int(spriteWidth))),
				Y: float64(rand.Intn(screenHeight - int(spriteHeight))),
			})
		}
	}

	if !validSpawn {
		fmt.Println("Warning: tank spawn failed due to collision!")
		return nil
	}

	// Add tank to resolv space only now
	g.space.Add(tank.Object)

	return tank
}

// pickTankSprites returns a random (body, barrel, bullet) sprite set for a tank.
func pickTankSprites(isEnemy bool) (string, string, string) {
	enemyBodies := []string{"tankBody_darkLarge", "tankBody_darkLarge_outline", "tankBody_huge_outline"}
	enemyBarrels := []string{"specialBarrel1_outline"}
	enemyBullets := []string{"bulletRed1_outline"}

	playerBodies := []string{"tankBody_red_outline", "tankBody_blue_outline", "tankBody_dark_outline", "tankBody_green_outline", "tankBody_dark_outline", "tankBody_green_outline", "tankBody_sand_outline"}
	playerBarrels := []string{"tankDark_barrel1_outline", "tankDark_barrel2_outline", "tankDark_barrel3_outline", "tankGreen_barrel1", "tankGreen_barrel1_outline", "tankGreen_barrel2", "tankGreen_barrel2_outline", "tankGreen_barrel3", "tankGreen_barrel3_outline", "tankRed_barrel1", "tankRed_barrel1_outline", "tankRed_barrel2_outline", "tankRed_barrel3_outline", "tankSand_barrel2_outline", "tankSand_barrel3_outline"}
	playerBullets := []string{"bulletRed1_outline"}

	if isEnemy {
		return enemyBodies[rand.Intn(len(enemyBodies))],
			enemyBarrels[rand.Intn(len(enemyBarrels))],
			enemyBullets[rand.Intn(len(enemyBullets))]
	}

	return playerBodies[rand.Intn(len(playerBodies))],
		playerBarrels[rand.Intn(len(playerBarrels))],
		playerBullets[rand.Intn(len(playerBullets))]
}

func NewRandomTank(game *Game, rotation float64, isEnemy bool) *Tank {
	randomBodyName, randomBarrelName, randomBulletName := pickTankSprites(isEnemy)

	position := resolv.Vector{X: float64(rand.Intn(screenWidth - tileWidth)), Y: float64(rand.Intn(screenHeight - tileHeight))}

	return NewTank(game, randomBodyName, randomBarrelName, randomBulletName, position, rotation, isEnemy)
}

// NewTankAt creates a tank with random sprites but placed at an authored
// position (used for hand-authored map spawn points).
func NewTankAt(game *Game, position resolv.Vector, rotation float64, isEnemy bool) *Tank {
	randomBodyName, randomBarrelName, randomBulletName := pickTankSprites(isEnemy)
	return NewTank(game, randomBodyName, randomBarrelName, randomBulletName, position, rotation, isEnemy)
}

// ---------------------------------------------------------------------------
// Enemy AI
// ---------------------------------------------------------------------------

// closestTarget returns the nearest non-enemy tank (the local player or a
// remote one), or nil if there is none.
func (t *Tank) closestTarget() *Tank {
	var closest *Tank
	minDist := math.MaxFloat64
	center := t.Object.Center()

	for _, tank := range t.game.Tanks {
		if tank.IsEnemy {
			continue
		}
		other := tank.Object.Center()
		d := math.Hypot(other.X-center.X, other.Y-center.Y)
		if d < minDist {
			minDist = d
			closest = tank
		}
	}

	return closest
}

// obstacleClearance returns the distance from point p to the closest obstacle.
func (t *Tank) obstacleClearance(p resolv.Vector) float64 {
	best := math.MaxFloat64

	for _, sh := range t.game.space.Shapes() {
		if sh.Tags().Has(TagObstacle) {
			c := sh.Bounds().Center()
			d := math.Hypot(c.X-p.X, c.Y-p.Y)
			if d < best {
				best = d
			}
		}
	}

	return best
}

// rotationForHeading returns the tank rotation R whose forward direction
// (-sin R, -cos R) equals the given heading vector. The tank moves forward
// along (-sin R, -cos R) (see Tank.Update), hence R = atan2(-h.X, -h.Y).
func rotationForHeading(heading resolv.Vector) float64 {
	return math.Atan2(-heading.X, -heading.Y)
}

// chooseAvoidanceHeading picks, among the headings sampled around the desired
// one, the direction with the most obstacle clearance (slightly preferring the
// desired heading to avoid jittering).
func (t *Tank) chooseAvoidanceHeading(desired resolv.Vector) resolv.Vector {
	baseAngle := math.Atan2(desired.Y, desired.X)
	center := t.Object.Center()

	best := desired
	bestScore := math.Inf(-1)

	for i := 0; i < enemyHeadingSamples; i++ {
		offset := (float64(i) - float64(enemyHeadingSamples-1)/2) * enemyHeadingStep
		angle := baseAngle + offset
		dir := resolv.Vector{X: math.Cos(angle), Y: math.Sin(angle)}

		probe := resolv.Vector{X: center.X + dir.X*enemyProbeDist, Y: center.Y + dir.Y*enemyProbeDist}
		clearance := t.obstacleClearance(probe)

		// deviation penalty keeps the tank on the desired path when it is free
		score := clearance - math.Abs(offset)*60
		if score > bestScore {
			bestScore = score
			best = dir
		}
	}

	return best
}

// ballisticSlope returns the barrel slope needed to hit a target at horizontal
// distance dist, or -1 when the target is out of range.
// Derived from the bullet physics in bullet.go:
//
//	horizontalDistance(slope) = (bulletSpeed^2 * tps / gravity) * sin(2*slope)
func (t *Tank) ballisticSlope(dist, tps float64) float64 {
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

// barrelWorldRotation returns the world rotation the barrel must have so that
// a bullet fired from `from` flies toward `to`. Bullets travel along the
// forward direction (-sin R, -cos R), hence R = -atan2(dy,dx) - pi/2.
func barrelWorldRotation(from, to resolv.Vector) float64 {
	return -math.Atan2(to.Y-from.Y, to.X-from.X) - math.Pi/2
}

// updateEnemyAI drives an enemy tank: it seeks/orbits the closest target while
// avoiding obstacles, and fires at it with a ballistic aim.
func (t *Tank) updateEnemyAI(tps float64) {
	target := t.closestTarget()
	if target == nil {
		return
	}

	enemyCenter := t.Object.Center()
	targetCenter := target.Object.Center()
	dist := math.Hypot(targetCenter.X-enemyCenter.X, targetCenter.Y-enemyCenter.Y)

	// desired unit direction towards the target
	desired := resolv.Vector{X: targetCenter.X - enemyCenter.X, Y: targetCenter.Y - enemyCenter.Y}
	length := math.Hypot(desired.X, desired.Y)
	if length > 0 {
		desired.X /= length
		desired.Y /= length
	}

	switch {
	case dist < enemyMinEngageDist:
		// too close: back off
		desired.X, desired.Y = -desired.X, -desired.Y
	case dist <= enemyMaxEngageDist:
		// engagement band: orbit the target (mostly sideways, slightly inward)
		perp := resolv.Vector{X: -desired.Y * t.enemyOrbitSign, Y: desired.X * t.enemyOrbitSign}
		desired = resolv.Vector{X: perp.X*0.8 + desired.X*0.2, Y: perp.Y*0.8 + desired.Y*0.2}
	}

	// evasive maneuver: when stuck, briefly move sideways instead
	if t.enemyEvasiveFrames > 0 {
		perp := resolv.Vector{X: -desired.Y * t.enemyOrbitSign, Y: desired.X * t.enemyOrbitSign}
		desired = perp
		t.enemyEvasiveFrames--
	}

	// keep the enemy on the visible playfield: near the screen edge, steer
	// back towards the center so they never disappear off-screen
	edgeMargin := 90.0
	if enemyCenter.X < edgeMargin || enemyCenter.X > float64(screenWidth)-edgeMargin ||
		enemyCenter.Y < edgeMargin || enemyCenter.Y > float64(screenHeight)-edgeMargin {
		toCenter := resolv.Vector{X: float64(screenWidth)/2 - enemyCenter.X, Y: float64(screenHeight)/2 - enemyCenter.Y}
		toCenterLen := math.Hypot(toCenter.X, toCenter.Y)
		if toCenterLen > 0 {
			toCenter.X /= toCenterLen
			toCenter.Y /= toCenterLen
			desired.X += toCenter.X * 1.5
			desired.Y += toCenter.Y * 1.5
			desiredLen := math.Hypot(desired.X, desired.Y)
			if desiredLen > 0 {
				desired.X /= desiredLen
				desired.Y /= desiredLen
			}
		}
	}

	heading := t.chooseAvoidanceHeading(desired)

	// rotate the tank so it faces the chosen heading
	targetRotation := rotationForHeading(heading)
	angleDiff := math.Atan2(math.Sin(targetRotation-t.Object.Rotation()), math.Cos(targetRotation-t.Object.Rotation()))
	if math.Abs(angleDiff) > 0.02 {
		rot := rotationPerSecond / tps
		if angleDiff < 0 {
			rot = -rot
		}
		t.Object.Rotate(rot)
	}

	// move forward at a slightly reduced speed
	s, c := math.Sincos(t.Object.Rotation())
	speed := tankSpeed / tps * enemySpeedFactor
	t.Object.MoveVec(resolv.Vector{X: -s * speed, Y: -c * speed})

	// stuck detection based on actual movement
	moved := math.Hypot(t.Object.Center().X-t.enemyLastPos.X, t.Object.Center().Y-t.enemyLastPos.Y)
	t.enemyLastPos = t.Object.Center()
	if moved < enemyStuckProgress {
		t.enemyStuckFrames++
	} else {
		t.enemyStuckFrames = 0
	}
	if t.enemyStuckFrames > enemyStuckFramesMax {
		t.enemyStuckFrames = 0
		t.enemyEvasiveFrames = enemyEvasiveTime
	}

	// aim the barrels at the target: the barrel solid rotation is
	// tankRotation + relativeRotation, so relativeRotation must encode the
	// world aim minus the current tank rotation, otherwise the turret would
	// rotate together with the tank body.
	worldAim := barrelWorldRotation(enemyCenter, targetCenter)
	for _, b := range t.barrels {
		b.relativeRotation = worldAim - t.Object.Rotation()
	}

	// fire with ballistic slope when the target is in range
	slope := t.ballisticSlope(dist, tps)
	if slope >= 0 && t.ShootCooldown.IsReady() {
		t.ShootCooldown.Reset()
		for _, b := range t.barrels {
			b.slope = slope
		}
		t.Bullets = append(t.Bullets, t.Fire()...)
	}
}

func (t *Tank) Fire() []*Bullet {
	// Fire all the barrels
	bullets := make([]*Bullet, len(t.barrels))
	for b := 0; b < len(t.barrels); b++ {
		bullets[b] = t.barrels[b].Fire()
	}
	return bullets
}

func (t *Tank) Destroy() {
	// barrel
	for _, b := range t.barrels {
		b.Destroy()
	}
	t.game.space.Remove(t.Object)
}

func (t *Tank) Update(tps float64) {
	moveVec := resolv.Vector{}
	rotation := 0.0
	rotationSpeed := rotationPerSecond / tps
	movementSpeed := tankSpeed / tps
	slopeSpeed := barrelMaxSlope / tps

	if t.IsEnemy {
		t.updateEnemyAI(tps)
	} else {
		// player
		// rotate tank
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			rotation += rotationSpeed
		}

		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			rotation -= rotationSpeed
		}

		// rotate barrel
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			for _, b := range t.barrels {
				b.relativeRotation += rotationSpeed
			}
		}

		if ebiten.IsKeyPressed(ebiten.KeyD) {
			for _, b := range t.barrels {
				b.relativeRotation -= rotationSpeed
			}
		}

		// move tank
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			s, c := math.Sincos(t.Object.Rotation())
			moveVec.X -= s * movementSpeed
			moveVec.Y -= c * movementSpeed
		}
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			s, c := math.Sincos(t.Object.Rotation())
			moveVec.X += s * movementSpeed
			moveVec.Y += c * movementSpeed
		}

		// fmt.Println(moveVec)
		t.Object.Rotate(rotation)
		t.Object.MoveVec(moveVec)

		// charge shoot
		if t.ShootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
			for i := 0; i < len(t.barrels); i++ {
				t.barrels[i].slope += slopeSpeed
				if t.barrels[i].slope > barrelMaxSlope {
					t.barrels[i].slope = barrelMaxSlope
				}
			}
			// fmt.Println(t.barrel.slope)
		}

		// charge fire increasing barrel slope
		yesFire := false
		for _, barrel := range t.barrels {
			if (barrel.slope > 0.0 && inpututil.IsKeyJustReleased(ebiten.KeySpace) || barrel.slope >= barrelMaxSlope) && t.ShootCooldown.IsReady() {
				t.ShootCooldown.Reset()
				yesFire = true
				// fmt.Println(len(p.bullets))
			}
		}

		// fire and reset barrels slope
		if yesFire {
			t.Bullets = append(t.Bullets, t.Fire()...)
			for i := 0; i < len(t.barrels); i++ {
				t.barrels[i].slope = 0.0
			}
		}

		// fire and reset barrels slope
		if yesFire {
			t.Bullets = append(t.Bullets, t.Fire()...)
			for i := 0; i < len(t.barrels); i++ {
				t.barrels[i].slope = 0.0
			}
		}

		// --- NETWORKING: Send data to server ---
		if t.game.netClient != nil {
			t.game.netClient.SendTankData(t)
		}
	}

	t.ShootCooldown.Update()

	// update bullets
	var activeBullets []*Bullet
	for _, bullet := range t.Bullets {
		bullet.Update(tps)
		if !bullet.exploded {
			activeBullets = append(activeBullets, bullet)
		}
	}
	t.Bullets = activeBullets

	for _, b := range t.barrels {
		b.Update(tps)
	}

	otherTanks := t.Object.SelectTouchingCells(2).FilterShapes().ByTags(TagEnemy | TagPlayer | TagObstacle)
	// bulletObjects := t.Object.SelectTouchingCells(2).FilterShapes().ByTags(TagBullet)
	t.Object.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: otherTanks,
		OnIntersect: func(set resolv.IntersectionSet) bool {
			t.Object.MoveVec(set.MTV)
			return true
		},
	})
}

func (t *Tank) Draw(screen *ebiten.Image) {
	// Draw the tank

	// body
	bodyHalfW := t.Width / 2.0
	bodyHalfH := t.Height / 2.0
	op_body := &ebiten.DrawImageOptions{}
	op_body.GeoM.Translate(-bodyHalfW, -bodyHalfH)
	op_body.GeoM.Rotate(-t.Object.Rotation())
	op_body.GeoM.Translate(bodyHalfW, bodyHalfH)
	op_body.GeoM.Translate(t.Object.Center().X-bodyHalfW, t.Object.Center().Y-bodyHalfH)

	screen.DrawImage(t.Sprite, op_body)

	// barrel
	for i := 0; i < len(t.barrels); i++ {
		t.barrels[i].Draw(screen)
	}

	// --- DISEGNO DEL MIRINO (Solo per il Player) ---
	if !t.IsEnemy {
		for _, b := range t.barrels {
			// Mostriamo il mirino solo se il giocatore sta caricando il colpo
			if b.slope > 0 {
				// Calcoliamo la distanza stimata in base all'inclinazione (slope).
				// barrelMaxSlope è l'inclinazione massima. Mappiamo questo valore
				// su una distanza in pixel sullo schermo (es. 600 pixel di gittata massima).
				maxDistance := 600.0 // Puoi aggiustare questo valore per farlo coincidere con la fisica esatta
				distance := (b.slope / barrelMaxSlope) * maxDistance

				// Calcoliamo la direzione in cui punta la canna
				s, c := math.Sincos(b.solid.Rotation())

				// Troviamo le coordinate del bersaglio
				targetX := b.solid.Center().X - s*distance
				targetY := b.solid.Center().Y - c*distance

				// Colore del mirino: rosso semitrasparente
				crosshairColor := color.RGBA{R: 255, G: 50, B: 50, A: 150}

				// 1. Cerchio esterno del mirino
				vector.StrokeCircle(screen, float32(targetX), float32(targetY), 15, 2, crosshairColor, true)

				// 2. Linea orizzontale
				vector.StrokeLine(screen, float32(targetX-22), float32(targetY), float32(targetX+22), float32(targetY), 2, crosshairColor, true)

				// 3. Linea verticale
				vector.StrokeLine(screen, float32(targetX), float32(targetY-22), float32(targetX), float32(targetY+22), 2, crosshairColor, true)
			}
		}
	}
}
