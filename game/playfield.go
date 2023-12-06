package game

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Playfield struct {
	game *Game

	position Vector
	rotation float64
	sprite   *ebiten.Image

	shootCooldown *Timer
}

func NewPlayfield(game *Game) *Playfield {
	// var PlayfieldTile = mustLoadImage("png/tileGrass2.png")
	return &Playfield{
		game: game,
	}
}

func (p *Playfield) Update() {

}

func (p *Playfield) Draw(screen *ebiten.Image) {
}
