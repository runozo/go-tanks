package game

import (
	_ "image/png"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	b2 "github.com/oliverbestmann/box2d-go"
)

const (
	shootCooldown  = time.Millisecond * 250
	rotationSpeed  = 1.0 / 2.0 * math.Pi / 3000 // rad/s
	tankSpeed      = 0.03                       // m/s
	barrelMaxSlope = math.Pi / 4
)

type Player struct {
	game          *Game
	Tank          *Tank
	Bullets       []*Bullet
	shootCooldown *Timer
	netClient     *NetClient
}

func NewPlayer(game *Game) *Player {

	position := b2.Vec2{
		X: screenWidth/2 + screenWidth/4,
		Y: screenHeight / 2,
	}

	newTank := NewRandomTank(game, position, 0.0)

	p := &Player{
		game:          game,
		Tank:          newTank,
		shootCooldown: NewTimer(shootCooldown),
	}

	// TODO
	if game.serverAddress != "" {
		p.netClient = NewNetClient(game.serverAddress)
	}

	return p
}

func (p *Player) Update(tps float64) {
	slopeSpeed := barrelMaxSlope / tps
	p.shootCooldown.Update()

	// rotate tank
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.Tank.tankBody.SetAngularVelocity(-rotationSpeed)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.Tank.tankBody.SetAngularVelocity(rotationSpeed)
	} else {
		p.Tank.tankBody.SetAngularVelocity(0)
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
	// speed is box2d units so meters/second
	movementX := 0.0
	movementY := 0.0

	radians := float64(p.Tank.tankBody.GetRotation().Angle())

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		movementX += math.Sin(radians) * tankSpeed
		movementY -= math.Cos(radians) * tankSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		movementX -= math.Sin(radians) * tankSpeed
		movementY += math.Cos(radians) * tankSpeed
	}

	p.Tank.tankBody.SetLinearVelocity(b2.Vec2{X: float32(movementX), Y: float32(movementY)})
	// fmt.Println("Linear velocity:", float32(movementX), float32(movementY), movementX, movementY, p.Tank.tankBody.GetLinearVelocity())

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
		p.Tank = NewRandomTank(p.game, p.Tank.tankBody.GetPosition(), float64(p.Tank.tankBody.GetRotation().Angle()))
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

	// // send data to server
	// if p.netClient != nil {
	// 	buf := make([]byte, 8) // float64 is 8 bytes long
	// 	binary.LittleEndian.PutUint64(buf, math.Float64bits(p.Tank.Position.X))
	// 	p.netClient.SendData(buf)
	// }
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.Tank.Draw(screen)
}
