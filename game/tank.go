package game

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	b2 "github.com/oliverbestmann/box2d-go"
)

type Tank struct {
	tankBody   b2.Body
	bodySprite *ebiten.Image
	barrels    []*Barrel
	bodyWidth  float64
	bodyHeight float64
}

func NewTank(game *Game, bodySpriteName, barrelSpriteName, bulletSpriteName string, position b2.Vec2, rotation float64) *Tank {

	bodySprite := game.assets.GetSprite(bodySpriteName)

	// create box2d body
	bodyDef := b2.DefaultBodyDef()
	bodyDef.Position = position

	s, c := math.Sincos(rotation)
	bodyDef.Rotation = b2.Rot{C: float32(c), S: float32(s)}
	bodyDef.Type1 = b2.DynamicBody
	shape := b2.MakeBox(float32(bodySprite.Bounds().Dx())/2.0, float32(bodySprite.Bounds().Dy())/2.0)
	body := game.world.CreateBody(bodyDef)
	// attach shape for collisions
	shapeDef := b2.DefaultShapeDef()
	shapeDef.Filter.CategoryBits = uint64(FB_TANK)
	shapeDef.Filter.MaskBits = uint64(FB_BULLET)
	body.CreatePolygonShape(shapeDef, shape)

	tank := &Tank{
		bodySprite: bodySprite,
		bodyWidth:  float64(bodySprite.Bounds().Dx()),
		bodyHeight: float64(bodySprite.Bounds().Dy()),
		barrels:    make([]*Barrel, 0),
		tankBody:   body,
	}

	if bodySpriteName == "tankBody_huge_outline" {
		// this bodies have 2 barrels each
		tank.barrels = []*Barrel{
			NewBarrel(game, barrelSpriteName, bulletSpriteName, tank, b2.Vec2{X: float32(tank.bodyWidth / 2.0), Y: float32(tank.bodyHeight / 4.0)}),
			NewBarrel(game, barrelSpriteName, bulletSpriteName, tank, b2.Vec2{X: float32(tank.bodyWidth / 2.0), Y: float32(tank.bodyHeight / 4.0 * 3.0)}),
		}
	} else {
		tank.barrels = []*Barrel{
			NewBarrel(game, barrelSpriteName, bulletSpriteName, tank, b2.Vec2{X: float32(tank.bodyWidth / 2.0), Y: float32(tank.bodyHeight / 2.0)}),
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

	// fmt.Println("Random tank:", randomBodyName, randomBarrelName, randomBulletName)

	return NewTank(game, randomBodyName, randomBarrelName, randomBulletName, position, rotation)
}

func (t *Tank) Fire() []*Bullet {
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
	pos := t.tankBody.GetPosition()
	bodyHalfW := t.bodyWidth / 2.0
	bodyHalfH := t.bodyHeight / 2.0
	op_body := &ebiten.DrawImageOptions{}
	op_body.GeoM.Translate(-bodyHalfW, -bodyHalfH)
	op_body.GeoM.Rotate(float64(t.tankBody.GetRotation().Angle()))
	op_body.GeoM.Translate(bodyHalfW, bodyHalfH)
	op_body.GeoM.Translate(float64(pos.X), float64(pos.Y))

	screen.DrawImage(t.bodySprite, op_body)

	// barrel
	for i := 0; i < len(t.barrels); i++ {
		t.barrels[i].Draw(screen)
	}
}
