package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	WIDTH  = 1200
	HEIGHT = 800
)

func randFloat(low, high float64) float64 {
	return rand.Float64()*(high-low) + low
}

func clamp(num, min, max float64) float64 {
	if num < min {
		return min
	}
	if num > max {
		return max
	}
	return num
}

type Point struct {
	X float64
	Y float64
	R float64
	C color.RGBA
	v *imdraw.IMDraw
}

// KILL ALL HUMANS
func NewPoint(x, y, r float64) *Point {
	return &Point{x, y, r, colornames.Red, imdraw.New(nil)}
}

func (p *Point) Draw() {
	p.v.Reset()
	p.v.Clear()
	p.v.Color = p.C
	p.v.Push(pixel.V(p.X, p.Y))
	p.v.Circle(p.R, 0)
}

func (p *Point) UpdatePosition() {
	p.X += randFloat(-5, 5)
	p.Y += randFloat(-5, 5)
	p.X = clamp(p.X, 0, WIDTH)
	p.Y = clamp(p.Y, 0, HEIGHT)
}

func (p *Point) Visual() *imdraw.IMDraw {
	return p.v
}

func (p *Point) CollidesAny(others map[*Point]bool) bool {
	for other := range others {
		dx := other.X - p.X
		dy := other.Y - p.Y
		r := other.R + p.R
		if r*r > dx*dx+dy*dy {
			return true
		}
	}
	return false
}

func run() {
	rand.Seed(time.Now().UnixNano())

	cfg := pixelgl.WindowConfig{
		Title:  "Hacked DLA",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// create points in the center
	numPoints := 1000
	points := make([]*Point, 0, numPoints)
	for i := 0; i < numPoints; i++ {
		points = append(points, NewPoint(
			randFloat(200, 1000),
			randFloat(100, 700),
			5))
	}

	// create seed and list of frozens
	frozen := map[*Point]bool{}
	seed := NewPoint(WIDTH/2, HEIGHT/2, 2)
	seed.C = colornames.Blue
	frozen[seed] = true
	points = append(points, seed)

	// the batch
	batch := pixel.NewBatch(&pixel.TrianglesData{}, nil)

	// performance
	frames := 0
	second := time.Tick(1 * time.Second)

	for !win.Closed() {

		batch.Clear()

		for _, p := range points {
			if !frozen[p] {
				p.UpdatePosition()
			}
			p.Draw()
			p.Visual().Draw(batch)

			//collide and change state if necessary
			if _, in := frozen[p]; !in {
				if p.CollidesAny(frozen) {
					// freeze p
					p.C = colornames.Blue
					frozen[p] = true
				}
			}
		}

		win.Clear(colornames.White)
		batch.Draw(win)
		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%d fps", frames))
			frames = 0
			fmt.Println(len(points), len(frozen))

			// add another point every second
			// points = append(points, NewPoint(
			// 	randFloat(200, 1000),
			// 	randFloat(100, 700),
			// 	5))
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
