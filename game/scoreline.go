package game

import (
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ScoreLine struct {
	Enemy  int
	Player int
	font   *text.GoTextFace
}

func NewScoreLine(game *Game) *ScoreLine {
	return &ScoreLine{
		Enemy:  0,
		Player: 0,
		font:   game.fontMedium,
	}
}

func (s *ScoreLine) IncrementPlayer() {
	s.Player += 1
}

func (s *ScoreLine) IncrementEnemy() {
	s.Enemy += 1
}

func (s *ScoreLine) Reset() {
	s.Enemy = 0
	s.Player = 0
}

func (s *ScoreLine) Draw(screen *ebiten.Image) {
	str := "Your score: " + strconv.Itoa(s.Player) + " \t Comp score: " + strconv.Itoa(s.Enemy)
	width, height := text.Measure(str, s.font, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64((screen.Bounds().Dx()-int(width))/2), float64(height))
	op.ColorScale.ScaleWithColor(color.Black)
	text.Draw(screen, str, s.font, op)
}
