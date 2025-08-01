package game

import (
	"fmt"
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

func NewTank(game *Game, bodySprite, barrelSprite, bulletSprite *ebiten.Image, x, y, rotation float64) *Tank {
	tank := &Tank{
		bodySprite: bodySprite,
		bodyWidth:  float64(bodySprite.Bounds().Dx()),
		bodyHeight: float64(bodySprite.Bounds().Dy()),
		position:   Vector{X: x, Y: y},
		rotation:   rotation,
		barrel:     nil,
	}
	shootAnimationSprites := []*ebiten.Image{
		game.assets.GetSprite("shotThin"),
		game.assets.GetSprite("shotLarge"),
		game.assets.GetSprite("shotOrange"),
		game.assets.GetSprite("shotRed"),
		game.assets.GetSprite("shotOrange"),
		game.assets.GetSprite("shotLarge"),
	}

	explosionAnimationSprites := []*ebiten.Image{
		game.assets.GetSprite("explosionSmoke1"),
		game.assets.GetSprite("explosionSmoke2"),
		game.assets.GetSprite("explosionSmoke3"),
		game.assets.GetSprite("explosionSmoke4"),
		game.assets.GetSprite("explosionSmoke5"),
	}

	tank.barrel = NewBarrel(barrelSprite, bulletSprite, tank, shootAnimationSprites, explosionAnimationSprites)
	return tank
}

func NewRandomTank(game *Game, x, y, rotation float64) *Tank {
	bodies := []string{"tankBody_red_outline", "tankBody_blue_outline", "tankBody_dark_outline", "tankBody_green_outline", "tankBody_dark_outline", "tankBody_green_outline", "tankBody_sand_outline"}
	barrels := []string{"tankDark_barrel1_outline", "tankDark_barrel2_outline", "tankDark_barrel3_outline", "tankGreen_barrel1", "tankGreen_barrel1_outline", "tankGreen_barrel2", "tankGreen_barrel2_outline", "tankGreen_barrel3", "tankGreen_barrel3_outline", "tankRed_barrel1", "tankRed_barrel1_outline", "tankRed_barrel2_outline", "tankRed_barrel3_outline", "tankSand_barrel2_outline", "tankSand_barrel3_outline"}
	bullets := []string{"bulletSand3_outline", "bulletGreen3_outline", "bulletBlue3_outline"}

	randomBodyName := bodies[rand.Intn(len(bodies))]
	randomBarrelName := barrels[rand.Intn(len(barrels))]
	randomBulletName := bullets[rand.Intn(len(bullets))]

	fmt.Println("Random tank:", randomBodyName, randomBarrelName, randomBulletName)

	return NewTank(game, game.assets.GetSprite(randomBodyName), game.assets.GetSprite(randomBarrelName), game.assets.GetSprite(randomBulletName), x, y, rotation)
}

func (t *Tank) Fire() *Bullet {
	return t.barrel.Fire()
}

func (t *Tank) Update(tps float64) {
	t.barrel.Update(tps)
}

func (t *Tank) Draw(screen *ebiten.Image) {
	// Draw the tank

	// body
	bodyHalfW := t.bodyWidth / 2
	bodyHalfH := t.bodyHeight / 2
	op_body := &ebiten.DrawImageOptions{}
	op_body.GeoM.Translate(-bodyHalfW, -bodyHalfH)
	op_body.GeoM.Rotate(t.rotation)
	op_body.GeoM.Translate(bodyHalfW, bodyHalfH)
	op_body.GeoM.Translate(t.position.X, t.position.Y)

	screen.DrawImage(t.bodySprite, op_body)

	// barrel
	t.barrel.Draw(screen)
}
