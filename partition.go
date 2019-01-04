package main

import "image/color"

type Partition struct {
	Top, Bottom, Left, Right float64
	Points                   []*Point
	C                        color.RGBA
}

func NewPartition() *Partition {
	return &Partition{Points: make([]*Point, 0, 200)}
}

func (part *Partition) ClearPoints() {
	part.Points = make([]*Point, 0, 200)
}

func (part *Partition) AddPoints(points []*Point, applyColor bool) {
	for _, p := range points {
		if part.ShouldContain(p) {
			part.Points = append(part.Points, p)
			if applyColor {
				p.SetColor(part.C)
			} else if p.C == part.C {
				if p.Frozen {
					p.SetColor(COLOR_RAMP[p.Age])
				} else {
					p.SetColor(POINT_COLOR)
				}
			}
		}
	}
}

func (part *Partition) ShouldContain(p *Point) bool {
	return part.Left-p.R <= p.X && p.X <= part.Right+p.R &&
		part.Bottom-p.R <= p.Y && p.Y <= part.Top+p.R
}

func (part *Partition) CollideWithin(test func(p, other *Point) bool, action func(p *Point)) {
	for i := 0; i < len(part.Points); i++ {
		for j := 1 + 1; j < len(part.Points); j++ {
			if test(part.Points[i], part.Points[j]) {
				if part.Points[i].Frozen {
					action(part.Points[j])
				} else {
					action(part.Points[i])
				}
			}
		}
	}
}
