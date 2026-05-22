package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

const (
	barrelMaxSlope    = math.Pi / 3
	rotationPerSecond = math.Pi / 2
	tankSpeed         = 200
)

type Tank struct {
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
}

func NewTank(g *Game, bodySpriteName, barrelSpriteName, bulletSpriteName string, position resolv.Vector, rotation float64, isEnemy bool) *Tank {

	bodySprite := g.assets.GetSprite(bodySpriteName)
	spriteWidth := float64(bodySprite.Bounds().Dx())
	spriteHeight := float64(bodySprite.Bounds().Dy())

	tank := &Tank{
		Sprite:        bodySprite,
		Object:        resolv.NewRectangle(position.X, position.Y, spriteWidth, spriteHeight),
		Width:         spriteWidth,
		Height:        spriteHeight,
		barrels:       make([]*Barrel, 0),
		Bullets:       make([]*Bullet, 0),
		game:          g,
		Hit:           false,
		IsEnemy:       isEnemy,
		ShootCooldown: NewTimer(time.Millisecond*2500 + time.Millisecond*time.Duration(rand.Intn(1000))),
	}

	tank.Object.Rotate(rotation)
	tank.Object.SetPositionVec(position)
	if isEnemy {
		tank.Object.Tags().Set(TagEnemy)
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

func NewRandomTank(game *Game, rotation float64, isEnemy bool) *Tank {
	enemyBodies := []string{"tankBody_darkLarge", "tankBody_darkLarge_outline", "tankBody_huge_outline"}
	enemyBarrels := []string{"specialBarrel1_outline"}
	enemyBullets := []string{"bulletRed1_outline"}

	playerBodies := []string{"tankBody_red_outline", "tankBody_blue_outline", "tankBody_dark_outline", "tankBody_green_outline", "tankBody_dark_outline", "tankBody_green_outline", "tankBody_sand_outline"}
	playerBarrels := []string{"tankDark_barrel1_outline", "tankDark_barrel2_outline", "tankDark_barrel3_outline", "tankGreen_barrel1", "tankGreen_barrel1_outline", "tankGreen_barrel2", "tankGreen_barrel2_outline", "tankGreen_barrel3", "tankGreen_barrel3_outline", "tankRed_barrel1", "tankRed_barrel1_outline", "tankRed_barrel2_outline", "tankRed_barrel3_outline", "tankSand_barrel2_outline", "tankSand_barrel3_outline"}
	playerBullets := []string{"bulletRed1_outline"}

	randomBodyName := playerBodies[rand.Intn(len(playerBodies))]
	randomBarrelName := playerBarrels[rand.Intn(len(playerBarrels))]
	randomBulletName := playerBullets[rand.Intn(len(playerBullets))]

	if isEnemy {
		randomBodyName = enemyBodies[rand.Intn(len(enemyBodies))]
		randomBarrelName = enemyBarrels[rand.Intn(len(enemyBarrels))]
		randomBulletName = enemyBullets[rand.Intn(len(enemyBullets))]
	}

	position := resolv.Vector{X: float64(rand.Intn(screenWidth - tileWidth)), Y: float64(rand.Intn(screenHeight - tileHeight))}

	return NewTank(game, randomBodyName, randomBarrelName, randomBulletName, position, rotation, isEnemy)
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
		// --- NUOVA IA: RICERCA DEL GIOCATORE PIÙ VICINO ---
		var closestPlayer *Tank
		minDistance := math.MaxFloat64
		enemyCenter := t.Object.Center()

		// Cicliamo su tutti i tank per trovare il giocatore più vicino
		for _, tank := range t.game.Tanks {
			if !tank.IsEnemy {
				playerCenter := tank.Object.Center()

				// Calcoliamo la distanza tra questo nemico e il potenziale bersaglio
				dx := playerCenter.X - enemyCenter.X
				dy := playerCenter.Y - enemyCenter.Y
				distance := math.Hypot(dx, dy)

				// Se è più vicino del precedente trovato, lo memorizziamo
				if distance < minDistance {
					minDistance = distance
					closestPlayer = tank
				}
			}
		}

		// Se abbiamo individuato un bersaglio valido nel raggio d'azione
		if closestPlayer != nil {
			targetPosition := closestPlayer.Object.Center()

			// Ruotiamo la canna (o le canne) verso il giocatore più vicino
			for _, b := range t.barrels {
				b.relativeRotation = -math.Atan2(targetPosition.Y-enemyCenter.Y, targetPosition.X-enemyCenter.X) - math.Pi/2
			}

			// Spara autonomamente seguendo il cooldown del timer
			if t.ShootCooldown.IsReady() {
				t.ShootCooldown.Reset()
				// I nemici scelgono un'inclinazione (gittata) casuale per simulare l'imprecisione
				randomSlope := rand.Float64() * barrelMaxSlope
				for _, b := range t.barrels {
					b.slope = randomSlope
				}
				t.Bullets = append(t.Bullets, t.Fire()...)
			}
		}
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
