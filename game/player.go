package game

import (
	"encoding/binary"
	_ "image/png"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/resolv"
)

const (
	shootCooldown     = time.Millisecond * 250
	rotationPerSecond = math.Pi
	tankSpeed         = 120.0
	barrelMaxSlope    = math.Pi / 4
)

type Player struct {
	game          *Game
	tank          *Tank
	Bullets       []*Bullet
	shootCooldown *Timer
	netClient     *NetClient
}

func NewPlayer(game *Game) *Player {

	position := resolv.Vector{
		X: screenWidth/2 + screenWidth/4,
		Y: screenHeight / 2,
	}

	p := &Player{
		game:          game,
		tank:          NewRandomTank(game, position, 0),
		shootCooldown: NewTimer(shootCooldown),
	}

	if game.serverAddress != "" {
		p.netClient = NewNetClient(game.serverAddress)
	}

	return p
}

func (p *Player) Update(tps float64) {
	moveVec := resolv.Vector{}
	rotation := 0.0
	rotationSpeed := rotationPerSecond / tps
	movementSpeed := tankSpeed / tps
	slopeSpeed := barrelMaxSlope / tps
	p.shootCooldown.Update()

	// rotate tank
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		rotation += rotationSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		rotation -= rotationSpeed
	}

	// rotate barrel
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		for i := 0; i < len(p.tank.barrels); i++ {
			p.tank.barrels[i].relativeRotation -= rotationSpeed
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		for i := 0; i < len(p.tank.barrels); i++ {
			p.tank.barrels[i].relativeRotation += rotationSpeed
		}
	}

	// move tank
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		s, c := math.Sincos(p.tank.solid.Rotation())
		moveVec.X -= s * movementSpeed
		moveVec.Y -= c * movementSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		s, c := math.Sincos(p.tank.solid.Rotation())
		moveVec.X += s * movementSpeed
		moveVec.Y += c * movementSpeed
	}

	// fmt.Println(moveVec)
	p.tank.solid.Rotate(rotation)
	p.tank.solid.MoveVec(moveVec)

	// charge shoot
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		for i := 0; i < len(p.tank.barrels); i++ {
			p.tank.barrels[i].slope += slopeSpeed
			if p.tank.barrels[i].slope > barrelMaxSlope {
				p.tank.barrels[i].slope = barrelMaxSlope
			}

		}
		// fmt.Println(p.tank.barrel.slope)
	}

	// charge fire increasing barrel slope
	yesFire := false
	for _, barrel := range p.tank.barrels {
		if (barrel.slope > 0.0 && inpututil.IsKeyJustReleased(ebiten.KeySpace) || barrel.slope >= barrelMaxSlope) && p.shootCooldown.IsReady() {
			p.shootCooldown.Reset()
			yesFire = true
			// fmt.Println(len(p.bullets))
		}
	}

	// fire and reset barrels slope
	if yesFire {
		p.Bullets = append(p.Bullets, p.tank.Fire()...)
		for i := 0; i < len(p.tank.barrels); i++ {
			p.tank.barrels[i].slope = 0.0
		}
	}

	// new tank
	if ebiten.IsKeyPressed(ebiten.KeyT) && inpututil.IsKeyJustPressed(ebiten.KeyT) {
		p.tank.Destroy()
		p.tank = NewRandomTank(p.game, p.tank.solid.Position(), rotation)
	}

	// update tank(s)
	p.tank.Update(tps)

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
		binary.LittleEndian.PutUint64(buf, math.Float64bits(p.tank.solid.Position().X))
		p.netClient.SendData(buf)
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	p.tank.Draw(screen)
}
