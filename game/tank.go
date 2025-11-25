package game

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type Tank struct {
	bodySprite *ebiten.Image
	solid      *resolv.ConvexPolygon
	barrels    []*Barrel
	bodyWidth  float64
	bodyHeight float64
}

func NewTank(g *Game, bodySpriteName, barrelSpriteName, bulletSpriteName string, position resolv.Vector, rotation float64) *Tank {

	bodySprite := g.assets.GetSprite(bodySpriteName)
	spriteWidth := float64(bodySprite.Bounds().Dx())
	spriteHeight := float64(bodySprite.Bounds().Dy())

	tank := &Tank{
		bodySprite: bodySprite,
		solid:      resolv.NewRectangle(position.X, position.Y, spriteWidth, spriteHeight),
		bodyWidth:  spriteWidth,
		bodyHeight: spriteHeight,
		barrels:    make([]*Barrel, 0),
	}

	tank.solid.Rotate(rotation)
	tank.solid.SetPositionVec(position)
	g.space.Add(tank.solid)

	if bodySpriteName == "tankBody_huge_outline" {
		// this bodies have 2 barrels each
		tank.barrels = []*Barrel{
			NewBarrel(g, barrelSpriteName, bulletSpriteName, tank, resolv.Vector{X: 0, Y: tank.bodyHeight / 4}),
			NewBarrel(g, barrelSpriteName, bulletSpriteName, tank, resolv.Vector{X: 0, Y: tank.bodyHeight / 4 * 3}),
		}
	} else {
		tank.barrels = []*Barrel{NewBarrel(g, barrelSpriteName, bulletSpriteName, tank, resolv.Vector{X: 0, Y: 0})}
	}

	return tank
}

func NewRandomTank(game *Game, position resolv.Vector, rotation float64) *Tank {
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
	bodyHalfW := float64(t.bodySprite.Bounds().Dx() / 2.0)
	bodyHalfH := float64(t.bodySprite.Bounds().Dy() / 2.0)
	op_body := &ebiten.DrawImageOptions{}
	op_body.GeoM.Translate(-bodyHalfW, -bodyHalfH)
	op_body.GeoM.Rotate(-t.solid.Rotation())
	op_body.GeoM.Translate(bodyHalfW, bodyHalfH)
	op_body.GeoM.Translate(t.solid.Center().X-bodyHalfW, t.solid.Center().Y-bodyHalfH)

	screen.DrawImage(t.bodySprite, op_body)

	// barrel
	for i := 0; i < len(t.barrels); i++ {
		t.barrels[i].Draw(screen)
	}
}
