package game

import (
	_ "image/png"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	shootCooldown     = time.Millisecond * 250
	rotationPerSecond = math.Pi
	tankSpeed         = 120.0

	bulletSpawnOffset = 40.0
	bulletSpeed       = 10.0
)

type Player struct {
	game *Game

	tank         *Tank
	rotation     float64
	position     Vector
	bulletSprite *ebiten.Image
	bullets      []*Bullet

	shootCooldown *Timer
}

func NewPlayer(game *Game) *Player {
	bulletSprite := game.assets.GetSprite("bulletRed2")

	return &Player{
		game:          game,
		rotation:      0,
		tank:          NewRandomTank(game),
		position:      Vector{X: screenWidth / 2, Y: screenHeight / 2},
		shootCooldown: NewTimer(shootCooldown),
		bulletSprite:  bulletSprite,
	}
}

func (p *Player) Update() {
	rotationSpeed := rotationPerSecond / float64(ebiten.TPS())
	movementSpeed := tankSpeed / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= rotationSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += rotationSpeed
	}

	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.shootCooldown.Reset()

		tankBounds := p.tank.BodySprite.Bounds()
		bulletBounds := p.bulletSprite.Bounds()
		halfWBullet := bulletBounds.Dx() / 2
		halfHBullet := bulletBounds.Dy() / 2
		halfW := float64(tankBounds.Dx()) / 2
		halfH := float64(tankBounds.Dy()) / 2

		spawnPos := Vector{
			p.position.X + halfW - float64(halfWBullet) + math.Sin(p.rotation)*bulletSpawnOffset,
			p.position.Y + halfH - float64(halfHBullet) + math.Cos(p.rotation)*-bulletSpawnOffset,
		}

		p.bullets = append(p.bullets, NewBullet(p.bulletSprite, spawnPos, p.rotation, bulletSpeed))

	}

	// move towards facing
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		// first get the direction the entity is pointed
		dx := math.Sin(p.rotation)
		dy := math.Cos(p.rotation)

		p.position.X += dx * movementSpeed
		p.position.Y -= dy * movementSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		// first get the direction the entity is pointed
		dx := math.Sin(p.rotation)
		dy := math.Cos(p.rotation)
		// then move in that direction
		p.position.X -= dx * movementSpeed
		p.position.Y += dy * movementSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyT) && inpututil.IsKeyJustPressed(ebiten.KeyT) {
		p.tank = NewRandomTank(p.game)
	}

	var visibleBullets []*Bullet
	for _, bullet := range p.bullets {
		bullet.Update()
		if bullet.position.X > 0 && bullet.position.X < float64(p.game.width) && bullet.position.Y > 0 && bullet.position.Y < float64(p.game.height) {
			visibleBullets = append(visibleBullets, bullet)
		}
	}
	p.bullets = visibleBullets
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.tank.Draw(screen, p.position, p.rotation)
	for _, bullet := range p.bullets {
		bullet.Draw(screen)
	}
	// fmt.Println("Bullets", len(p.bullets))
}
