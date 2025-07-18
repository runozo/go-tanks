package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tank struct {
	BodySprite   *ebiten.Image
	BarrelSprite *ebiten.Image
	position     Vector
}

func NewTank(game *Game, bodySprite, barrelSprite *ebiten.Image) *Tank {
	return &Tank{
		BodySprite:   bodySprite,
		BarrelSprite: barrelSprite,
		position:     Vector{X: screenWidth / 2, Y: screenHeight / 2},
	}
}

func NewRandomTank(game *Game) *Tank {
	bodies := []string{"tankBody_bigRed", "tankBody_bigRed_outline", "tankBody_blue", "tankBody_blue_outline", "tankBody_dark", "tankBody_darkLarge", "tankBody_darkLarge_outline", "tankBody_dark_outline", "tankBody_green", "tankBody_green_outline", "tankBody_huge", "tankBody_huge_outline", "tankBody_red", "tankBody_red_outline", "tankBody_sand", "tankBody_sand_outline", "tank_bigRed", "tank_blue", "tank_dark", "tank_darkLarge", "tank_green", "tank_huge", "tank_red", "tank_sand"}
	barrels := []string{"tankDark_barrel1", "tankDark_barrel1_outline", "tankDark_barrel2", "tankDark_barrel2_outline", "tankDark_barrel3", "tankDark_barrel3_outline", "tankGreen_barrel1", "tankGreen_barrel1_outline", "tankGreen_barrel2", "tankGreen_barrel2_outline", "tankGreen_barrel3", "tankGreen_barrel3_outline", "tankRed_barrel1", "tankRed_barrel1_outline", "tankRed_barrel2", "tankRed_barrel2_outline", "tankRed_barrel3", "tankRed_barrel3_outline", "tankSand_barrel1", "tankSand_barrel1_outline", "tankSand_barrel2", "tankSand_barrel2_outline", "tankSand_barrel3", "tankSand_barrel3_outline"}

	randomBodyName := bodies[rand.Intn(len(bodies))]
	randomBarrelName := barrels[rand.Intn(len(barrels))]

	return NewTank(game, game.assets.GetSprite(randomBodyName), game.assets.GetSprite(randomBarrelName))
}

func (t *Tank) Fire() *Bullet {
	return nil
}

func (t *Tank) Draw(screen *ebiten.Image, rotation, barrelRotation float64) {
	// Draw the tank

	// body
	bodyBounds := t.BodySprite.Bounds()
	bodyHalfW := float64(bodyBounds.Dx() / 2)
	bodyHalfH := float64(bodyBounds.Dy() / 2)
	op_body := &ebiten.DrawImageOptions{}
	op_body.GeoM.Translate(-bodyHalfW, -bodyHalfH)
	op_body.GeoM.Rotate(rotation)
	op_body.GeoM.Translate(bodyHalfW, bodyHalfH)
	op_body.GeoM.Translate(t.position.X, t.position.Y)

	// barrel
	barrellBounds := t.BarrelSprite.Bounds()
	barrellHalfW := float64(barrellBounds.Dx() / 2)
	barrellHeight := float64(barrellBounds.Dy())
	op_barrell := &ebiten.DrawImageOptions{}
	op_barrell.GeoM.Translate(-barrellHalfW, -barrellHeight)
	op_barrell.GeoM.Rotate(barrelRotation)
	op_barrell.GeoM.Translate(barrellHalfW, barrellHeight)
	op_barrell.GeoM.Translate(
		t.position.X+float64(bodyBounds.Dx())/2-barrellHalfW,
		t.position.Y+float64(bodyBounds.Dy())/2-barrellHeight,
	)

	screen.DrawImage(t.BodySprite, op_body)
	screen.DrawImage(t.BarrelSprite, op_barrell)
}
