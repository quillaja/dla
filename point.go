package main

import (
	"image/color"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Point struct {
	X      float64
	Y      float64
	R      float64
	C      color.RGBA
	Frozen bool
	Age    int
	v      *imdraw.IMDraw
	dirty  bool
}

// KILL ALL HUMANS
func NewPoint(x, y, r float64) *Point {
	return &Point{x, y, r, POINT_COLOR, false, 0, imdraw.New(nil), true}
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
		// if p.Frozen {
		// 	p.v.Color = color.RGBA{
		// 		R: clampInt255(int(p.C.R) - 16),
		// 		G: clampInt255(int(p.C.G) - 16),
		// 		B: clampInt255(int(p.C.B) - 16),
		// 		A: p.C.A}
		// 	p.v.Push(pixel.V(p.X, p.Y))
		// 	p.v.Circle(p.R, 0.5)
		// }

		p.dirty = false
	}
}

func (p *Point) UpdatePosition() {
	if !p.Frozen {
		p.X, p.Y = randomMovement(p.X, p.Y) // centerDrift(p.X, p.Y)

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

func randomMovement(x, y float64) (float64, float64) {
	return x + randFloat(-POINT_SPEED, POINT_SPEED),
		y + randFloat(-POINT_SPEED, POINT_SPEED)
}

func centerDrift(x, y float64) (float64, float64) {
	// go directly towards center 2% of the time
	if rand.Float32() < 0.02 {
		pos := pixel.V(x, y)
		center := pixel.V(WIDTH/2, HEIGHT/2)
		delta := pos.To(center).
			Unit().
			Scaled(POINT_SPEED)
		return x + delta.X, y + delta.Y
	}

	return randomMovement(x, y)
}

func clampInt255(i int) uint8 {
	if i < 0 {
		return 0
	}
	if i > 255 {
		return 255
	}
	return uint8(i)
}
