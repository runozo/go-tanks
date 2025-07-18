package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tank struct {
	bodySprite *ebiten.Image
	barrel     *Barrel
	bodyWidth  float64
	bodyHeight float64
	position   Vector
	rotation   float64
}

func NewTank(game *Game, bodySprite, barrelSprite, bulletSprite *ebiten.Image) *Tank {
	tank := &Tank{
		bodySprite: bodySprite,
		bodyWidth:  float64(bodySprite.Bounds().Dx()),
		bodyHeight: float64(bodySprite.Bounds().Dy()),
		position:   Vector{X: screenWidth / 2, Y: screenHeight / 2},
		rotation:   0.0,
		barrel:     nil,
	}
	tank.barrel = NewBarrel(barrelSprite, bulletSprite, tank)
	return tank
}

func NewRandomTank(game *Game) *Tank {
	bodies := []string{"tankBody_bigRed", "tankBody_bigRed_outline", "tankBody_blue", "tankBody_blue_outline", "tankBody_dark", "tankBody_darkLarge", "tankBody_darkLarge_outline", "tankBody_dark_outline", "tankBody_green", "tankBody_green_outline", "tankBody_huge", "tankBody_huge_outline", "tankBody_red", "tankBody_red_outline", "tankBody_sand", "tankBody_sand_outline", "tank_bigRed", "tank_blue", "tank_dark", "tank_darkLarge", "tank_green", "tank_huge", "tank_red", "tank_sand"}
	barrels := []string{"tankDark_barrel1", "tankDark_barrel1_outline", "tankDark_barrel2", "tankDark_barrel2_outline", "tankDark_barrel3", "tankDark_barrel3_outline", "tankGreen_barrel1", "tankGreen_barrel1_outline", "tankGreen_barrel2", "tankGreen_barrel2_outline", "tankGreen_barrel3", "tankGreen_barrel3_outline", "tankRed_barrel1", "tankRed_barrel1_outline", "tankRed_barrel2", "tankRed_barrel2_outline", "tankRed_barrel3", "tankRed_barrel3_outline", "tankSand_barrel1", "tankSand_barrel1_outline", "tankSand_barrel2", "tankSand_barrel2_outline", "tankSand_barrel3", "tankSand_barrel3_outline"}
	bullets := []string{"bulletSand3_outline", "bulletGreen3_outline", "bulletBlue3_outline"}

	randomBodyName := bodies[rand.Intn(len(bodies))]
	randomBarrelName := barrels[rand.Intn(len(barrels))]
	randomBulletName := bullets[rand.Intn(len(bullets))]

	return NewTank(game, game.assets.GetSprite(randomBodyName), game.assets.GetSprite(randomBarrelName), game.assets.GetSprite(randomBulletName))
}

func (t *Tank) Fire() *Bullet {
	return t.barrel.Fire()
}

func (t *Tank) Update() {
	t.barrel.Update()
}

func (t *Tank) Draw(screen *ebiten.Image, rotation float64) {
	// Draw the tank

	// body
	bodyHalfW := t.bodyWidth / 2
	bodyHalfH := t.bodyHeight / 2
	op_body := &ebiten.DrawImageOptions{}
	op_body.GeoM.Translate(-bodyHalfW, -bodyHalfH)
	op_body.GeoM.Rotate(rotation)
	op_body.GeoM.Translate(bodyHalfW, bodyHalfH)
	op_body.GeoM.Translate(t.position.X, t.position.Y)

	screen.DrawImage(t.bodySprite, op_body)

	// barrel
	t.barrel.Draw(screen)
}
