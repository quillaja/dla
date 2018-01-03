package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	WIDTH  = 1200
	HEIGHT = 900
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

func addRandPoints(points []*Point, num ...int) []*Point {
	if len(num) < 1 {
		num = append(num, 50)
	}
	if len(num) < 2 {
		num = append(num, 5)
	}

	for i := 0; i < num[0]; i++ {
		points = append(points, NewPoint(
			randFloat(0, WIDTH),
			randFloat(0, HEIGHT),
			float64(num[1])))
	}
	return points
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

	// create points
	numPoints := 200
	points := make([]*Point, 0, numPoints)
	points = addRandPoints(points, numPoints)

	// create seed in center
	seed := NewPoint(WIDTH/2, HEIGHT/2, 2)
	seed.C = colornames.Black
	seed.Frozen = true
	points = append(points, seed)

	// the batch
	batch := pixel.NewBatch(&pixel.TrianglesData{}, nil)

	// collision partitions
	partColors := []color.RGBA{
		colornames.Red, colornames.Green, colornames.Blue,
		colornames.Yellow, colornames.Brown, colornames.Cyan,
		colornames.Darkred}
	numParts := len(partColors) + 1
	partitions := make(map[string]*Partition)
	for w := 0; w < numParts; w++ {
		for h := 0; h < numParts; h++ {
			name := fmt.Sprintf("%d,%d", w, h)
			p := NewPartition()
			p.Left = float64(w * (WIDTH / numParts))
			p.Right = float64((w + 1) * (WIDTH / numParts))
			p.Bottom = float64(h * (HEIGHT / numParts))
			p.Top = float64((h + 1) * (HEIGHT / numParts))
			p.C = partColors[(w+h)%len(partColors)]
			partitions[name] = p
			// fmt.Println(name, p)
		}
	}

	// performance
	frames := 0
	second := time.Tick(1 * time.Second)

	for !win.Closed() {

		batch.Clear()

		// draw
		for _, p := range points {
			p.UpdatePosition()
			p.Draw()
			p.Visual().Draw(batch)
		}

		//collide and change state if necessary

		// separate into quadrants
		for _, part := range partitions {
			part.ClearPoints()
			part.AddPoints(points)
		}

		// collide within partitions
		for _, part := range partitions {
			part.CollideWithin(
				func(p, other *Point) bool {
					return p.Collides(other)
				},
				func(p *Point) {
					p.C = colornames.Black
					p.Frozen = true
				})
		}

		win.Clear(colornames.White)
		batch.Draw(win)
		win.Update()

		if win.JustPressed(pixelgl.KeySpace) {
			//add more points if space bar pressed
			points = addRandPoints(points, 50, rand.Intn(9)+1)
			fmt.Println("num points:", len(points))
		}

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%d fps", frames))
			frames = 0

		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
