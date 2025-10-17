package game

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	b2 "github.com/oliverbestmann/box2d-go"
)

type Tank struct {
	tankBody   b2.Body
	bodySprite *ebiten.Image
	barrels    []*Barrel
	bodyWidth  float32
	bodyHeight float32
	Rotation   float64
}

func NewTank(game *Game, bodySpriteName, barrelSpriteName, bulletSpriteName string, position b2.Vec2, rotation float64) *Tank {

	bodySprite := game.assets.GetSprite(bodySpriteName)

	bodyDef := b2.DefaultBodyDef()
	bodyDef.Position = b2.Vec2{X: float32(position.X), Y: float32(position.Y)}
	shape := b2.MakeBox(float32(bodySprite.Bounds().Dx())/2, float32(bodySprite.Bounds().Dy())/2)
	body := game.world.CreateBody(bodyDef)
	body.CreatePolygonShape(b2.DefaultShapeDef(), shape)

	tank := &Tank{
		bodySprite: bodySprite,
		bodyWidth:  float32(bodySprite.Bounds().Dx()),
		bodyHeight: float32(bodySprite.Bounds().Dy()),
		Rotation:   rotation,
		barrels:    make([]*Barrel, 0),
		tankBody:   body,
	}

	if bodySpriteName == "tankBody_huge_outline" {
		// this bodies have 2 barrels each
		tank.barrels = []*Barrel{
			NewBarrel(game, barrelSpriteName, bulletSpriteName, tank, b2.Vec2{X: tank.bodyWidth / 2, Y: tank.bodyHeight / 4}),
			NewBarrel(game, barrelSpriteName, bulletSpriteName, tank, b2.Vec2{X: tank.bodyWidth / 2, Y: tank.bodyHeight / 4 * 3}),
		}
	} else {
		tank.barrels = []*Barrel{
			NewBarrel(game, barrelSpriteName, bulletSpriteName, tank, b2.Vec2{X: tank.bodyWidth / 2, Y: tank.bodyHeight / 2}),
		}
	}

	return tank
}

func NewRandomTank(game *Game, position b2.Vec2, rotation float64) *Tank {
	bodies := []string{"tankBody_red_outline", "tankBody_blue_outline", "tankBody_dark_outline", "tankBody_green_outline", "tankBody_dark_outline", "tankBody_green_outline", "tankBody_sand_outline"}
	barrels := []string{"tankDark_barrel1_outline", "tankDark_barrel2_outline", "tankDark_barrel3_outline", "tankGreen_barrel1", "tankGreen_barrel1_outline", "tankGreen_barrel2", "tankGreen_barrel2_outline", "tankGreen_barrel3", "tankGreen_barrel3_outline", "tankRed_barrel1", "tankRed_barrel1_outline", "tankRed_barrel2_outline", "tankRed_barrel3_outline", "tankSand_barrel2_outline", "tankSand_barrel3_outline"}
	bullets := []string{"bulletSand3_outline", "bulletGreen3_outline", "bulletBlue3_outline"}

	randomBodyName := bodies[rand.Intn(len(bodies))]
	randomBarrelName := barrels[rand.Intn(len(barrels))]
	randomBulletName := bullets[rand.Intn(len(bullets))]

	fmt.Println("Random tank:", randomBodyName, randomBarrelName, randomBulletName)

	return NewTank(game, randomBodyName, randomBarrelName, randomBulletName, position, rotation)
}

func (t *Tank) Fire() []*Bullet {
	// Fire all the barrels
	bullets := make([]*Bullet, len(t.barrels))
	for b := 0; b < len(t.barrels); b++ {
		bullets[b] = t.barrels[b].Fire()
	}
	return bullets
}

func (t *Tank) Update(tps float64) {
	for i := 0; i < len(t.barrels); i++ {
		t.barrels[i].Update(tps)
	}
}

func (t *Tank) Draw(screen *ebiten.Image) {
	// Draw the tank

	// body
	bodyHalfW := t.bodyWidth / 2
	bodyHalfH := t.bodyHeight / 2
	op_body := &ebiten.DrawImageOptions{}
	op_body.GeoM.Translate(float64(-bodyHalfW), float64(-bodyHalfH))
	op_body.GeoM.Rotate(t.Rotation)
	op_body.GeoM.Translate(float64(bodyHalfW), float64(bodyHalfH))
	pos := t.tankBody.GetPosition()
	op_body.GeoM.Translate(float64(pos.X), float64(pos.Y))

	screen.DrawImage(t.bodySprite, op_body)

	// barrel
	for i := 0; i < len(t.barrels); i++ {
		t.barrels[i].Draw(screen)
	}
}
