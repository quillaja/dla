package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Point struct {
	X      float64
	Y      float64
	R      float64
	C      color.RGBA
	Frozen bool
	v      *imdraw.IMDraw
}

// KILL ALL HUMANS
func NewPoint(x, y, r float64) *Point {
	return &Point{x, y, r, POINT_COLOR, false, imdraw.New(nil)}
}

func (p *Point) Draw() {
	p.v.Reset()
	p.v.Clear()
	p.v.Color = p.C
	p.v.Push(pixel.V(p.X, p.Y))
	p.v.Circle(p.R, 0)
}

func (p *Point) UpdatePosition() {
	if !p.Frozen {
		p.X += randFloat(-5, 5)
		p.Y += randFloat(-5, 5)

		p.X = clamp(p.X, 0, WIDTH)
		p.Y = clamp(p.Y, 0, HEIGHT)
	}
}

func (p *Point) Visual() *imdraw.IMDraw {
	return p.v
}

func (p *Point) Collides(other *Point) bool {
	if p.Frozen != other.Frozen {
		dx := other.X - p.X
		dy := other.Y - p.Y
		r := other.R + p.R
		if r*r > dx*dx+dy*dy {
			return true
		}
	}
	return false
}
