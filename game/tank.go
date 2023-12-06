package game

import "github.com/hajimehoshi/ebiten/v2"

type Tank struct {
	body   *ebiten.Image
	barrel *ebiten.Image
}

func NewTank(game *Game) *Tank {
	body := game.assets.GetSprite("tankBody_red.png")
	barrel := game.assets.GetSprite("tankGreen_barrel1.png")

	return &Tank{
		body:   body,
		barrel: barrel,
	}
}
