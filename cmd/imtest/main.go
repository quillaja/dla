package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	WIDTH         = 1200
	HEIGHT        = 900
	POINT_MAXSIZE = 10
	POINT_SPEED   = 2
)

var (
	POINT_COLOR      = colornames.Cornflowerblue
	FROZEN_COLOR     = colornames.Ivory
	BACKGROUND_COLOR = colornames.Black
)

type camera struct {
	Position pixel.Vec
	Speed    float64
	Zoom     float64
	ZSpeed   float64
}

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

func addRandPoints(points []*Point, num int) []*Point {
	for i := 0; i < num; i++ {
		points = append(points, NewPoint(
			randFloat(0, WIDTH),
			randFloat(0, HEIGHT),
			randFloat(3, POINT_MAXSIZE)))
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
	seed.C = FROZEN_COLOR
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

	// camera
	cam := camera{Position: pixel.V(WIDTH/2, HEIGHT/2), Speed: 250.0, Zoom: 1.0, ZSpeed: 1.1}
	last := time.Now()

	// options
	showPartitions := false
	paused := false

	// performance
	frames := 0
	iterations := 0
	second := time.Tick(1 * time.Second)

	// run logic
	go func() {

		for !win.Closed() {
			if !paused {
				for _, p := range points {
					p.UpdatePosition()
				}

				// separate into quadrants
				for _, part := range partitions {
					part.ClearPoints()
					part.AddPoints(points, showPartitions)
				}

				// collide within partitions
				for _, part := range partitions {
					part.CollideWithin(
						func(p, other *Point) bool {
							return p.Collides(other)
						},
						func(p *Point) {
							p.SetColor(FROZEN_COLOR)
							p.Frozen = true
						})
				}

				iterations++
			}
		}
	}()

	for !win.Closed() {

		dt := time.Since(last).Seconds()
		last = time.Now()

		camMatrix := pixel.IM.
			Scaled(cam.Position, cam.Zoom).
			Moved(win.Bounds().Center().Sub(cam.Position))
		win.SetMatrix(camMatrix)

		// update user controlled things
		if win.Pressed(pixelgl.KeyLeft) {
			cam.Position.X -= cam.Speed * dt
		}
		if win.Pressed(pixelgl.KeyRight) {
			cam.Position.X += cam.Speed * dt
		}
		if win.Pressed(pixelgl.KeyDown) {
			cam.Position.Y -= cam.Speed * dt
		}
		if win.Pressed(pixelgl.KeyUp) {
			cam.Position.Y += cam.Speed * dt
		}
		// allow toggle paused
		if win.JustPressed(pixelgl.KeyP) {
			paused = !paused
		}
		// allow toggle of coloring to show partitions
		if win.JustPressed(pixelgl.KeyB) {
			showPartitions = !showPartitions
		}
		cam.Zoom *= math.Pow(cam.ZSpeed, win.MouseScroll().Y)

		if !paused {
			batch.Clear()

			// draw to batch
			for _, p := range points {
				// if .Left-p.R <= p.X && p.X <= .Right+p.R &&
				// .Bottom-p.R <= p.Y && p.Y <= .Top+p.R {
				p.Draw()
				p.Visual().Draw(batch)
			}

		}
		// draw batch to window
		win.Clear(BACKGROUND_COLOR)
		batch.Draw(win)

		win.Update()

		//add more points if space bar pressed
		if win.JustPressed(pixelgl.KeySpace) {
			points = addRandPoints(points, 50)
			fmt.Println("num points:", len(points))
		}

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%d fps, %d iter/s",
				frames, iterations))
			frames = 0
			iterations = 0
		default:
		}
	}
}

func main() {
	pixelgl.Run(run)
}
