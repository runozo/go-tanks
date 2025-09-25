package game

import "image"

func DoesIntersect(pos1 Vector, box1 image.Rectangle, pos2 Vector, box2 image.Rectangle) bool {
	rec1dx := pos1.X + float64(box1.Bounds().Dx())
	rec1dy := pos1.Y + float64(box1.Bounds().Dy())
	rec2dx := pos2.X + float64(box2.Bounds().Dx())
	rec2dy := pos2.Y + float64(box2.Bounds().Dy())
	return pos1.X < rec2dx &&
		rec1dx > pos2.X &&
		pos1.Y < rec2dy &&
		rec1dy > pos2.Y
}
