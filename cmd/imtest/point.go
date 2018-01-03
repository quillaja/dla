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
	dirty  bool
}

// KILL ALL HUMANS
func NewPoint(x, y, r float64) *Point {
	return &Point{x, y, r, POINT_COLOR, false, imdraw.New(nil), true}
}

func (p *Point) SetColor(c color.RGBA) {
	p.C = c
	p.dirty = true
}

func (p *Point) Draw() {
	if p.dirty {
		p.v.Reset()
		p.v.Clear()
		p.v.Color = p.C
		p.v.Push(pixel.V(p.X, p.Y))
		p.v.Circle(p.R, 0)

		p.dirty = false
	}
}

func (p *Point) UpdatePosition() {
	if !p.Frozen {
		p.X += randFloat(-POINT_SPEED, POINT_SPEED)
		p.Y += randFloat(-POINT_SPEED, POINT_SPEED)

		p.X = clamp(p.X, 0, WIDTH)
		p.Y = clamp(p.Y, 0, HEIGHT)

		p.dirty = true
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
