package game

import (
	"encoding/binary"
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
	barrelMaxSlope    = math.Pi / 4
)

type Player struct {
	game          *Game
	Tank          *Tank
	Bullets       []*Bullet
	shootCooldown *Timer
	netClient     *NetClient
}

func NewPlayer(game *Game) *Player {

	position := Vector{
		X: screenWidth/2 + screenWidth/4,
		Y: screenHeight / 2,
	}

	p := &Player{
		game:          game,
		Tank:          NewRandomTank(game, position, 0),
		shootCooldown: NewTimer(shootCooldown),
	}
	if game.serverAddress != "" {
		p.netClient = NewNetClient(game.serverAddress)
	}
	return p
}

func (p *Player) Update(tps float64) {
	rotationSpeed := rotationPerSecond / tps
	movementSpeed := tankSpeed / tps
	slopeSpeed := barrelMaxSlope / tps
	p.shootCooldown.Update()

	// rotate tank
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Tank.Rotation -= rotationSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Tank.Rotation += rotationSpeed
	}

	// rotate barrel
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		for i := 0; i < len(p.Tank.barrels); i++ {
			p.Tank.barrels[i].relativeRotation -= rotationSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		for i := 0; i < len(p.Tank.barrels); i++ {
			p.Tank.barrels[i].relativeRotation += rotationSpeed
		}
	}

	// move
	movementX := 0.0
	movementY := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		movementX += math.Sin(p.Tank.Rotation) * movementSpeed
		movementY -= math.Cos(p.Tank.Rotation) * movementSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		movementX -= math.Sin(p.Tank.Rotation) * movementSpeed
		movementY += math.Cos(p.Tank.Rotation) * movementSpeed
	}
	p.Tank.Position.X += movementX
	p.Tank.Position.Y += movementY

	// charge shoot
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		for i := 0; i < len(p.Tank.barrels); i++ {
			p.Tank.barrels[i].slope += slopeSpeed
			if p.Tank.barrels[i].slope > barrelMaxSlope {
				p.Tank.barrels[i].slope = barrelMaxSlope
			}

		}
		// fmt.Println(p.tank.barrel.slope)
	}

	// fire
	yesFire := false
	for i := 0; i < len(p.Tank.barrels); i++ {
		if (p.Tank.barrels[i].slope > 0.0 && inpututil.IsKeyJustReleased(ebiten.KeySpace) || p.Tank.barrels[i].slope >= barrelMaxSlope) && p.shootCooldown.IsReady() {
			p.shootCooldown.Reset()
			yesFire = true
			// fmt.Println(len(p.bullets))
		}
	}

	// fire and reset barrels slope
	if yesFire {
		p.Bullets = append(p.Bullets, p.Tank.Fire()...)
		for i := 0; i < len(p.Tank.barrels); i++ {
			p.Tank.barrels[i].slope = 0.0
		}
	}

	// new tank
	if ebiten.IsKeyPressed(ebiten.KeyT) && inpututil.IsKeyJustPressed(ebiten.KeyT) {
		p.Tank = NewRandomTank(p.game, p.Tank.Position, p.Tank.Rotation)
	}

	// update tank(s)
	p.Tank.Update(tps)

	// update bullets
	var activeBullets []*Bullet
	for _, bullet := range p.Bullets {
		bullet.Update(tps)
		if !bullet.exploded {
			activeBullets = append(activeBullets, bullet)
		}
	}
	p.Bullets = activeBullets

	// send data to server
	if p.netClient != nil {
		buf := make([]byte, 8) // float64 is 8 bytes long
		binary.LittleEndian.PutUint64(buf, math.Float64bits(p.Tank.Position.X))
		p.netClient.SendData(buf)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.Tank.Draw(screen)
	for _, bullet := range p.Bullets {
		bullet.Draw(screen)
	}
	// fmt.Println("Bullets", len(p.bullets))
}
