package game

import (
	"image"

	b2 "github.com/oliverbestmann/box2d-go"
)

func doesIntersect(pos1 Vector, box1 image.Rectangle, pos2 Vector, box2 image.Rectangle) bool {
	rec1dx := pos1.X + float64(box1.Bounds().Dx())
	rec1dy := pos1.Y + float64(box1.Bounds().Dy())
	rec2dx := pos2.X + float64(box2.Bounds().Dx())
	rec2dy := pos2.Y + float64(box2.Bounds().Dy())
	return pos1.X < rec2dx &&
		rec1dx > pos2.X &&
		pos1.Y < rec2dy &&
		rec1dy > pos2.Y
}

func NewWorld() b2.World {
	worldDef := b2.DefaultWorldDef()
	worldDef.Gravity = b2.Vec2{X: 0, Y: 0}
	world := b2.CreateWorld(worldDef)
	return world
}
