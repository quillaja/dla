package main

import "image/color"

// MakeColorRamp produces a list of colors by interpolating
// color values between a series of stops.
//
// `stops` must contain at least 2 colors, and `numColors` must be greater
// than the number of stops.
func MakeColorRamp(stops []color.RGBA, numColors int) []color.RGBA {
	if len(stops) < 2 {
		panic("Invalid arguments. len(stops) must be >= 2.")
	}
	if len(stops) > numColors {
		panic("Invalid arguments. More stops than number of desired colors.")
	}
	if len(stops) == numColors {
		return stops // why did the idiot even call the func?
	}

	ramp := []color.RGBA{}

	// num colors to add between each stop, and any remainders to distribute
	betweenStops := (numColors - len(stops)) / (len(stops) - 1)
	remainder := (numColors - len(stops)) % (len(stops) - 1)

	// go through pairs of stops, and interpolate the ramp in between
	for i := 0; i < len(stops)-1; i++ {
		cur, next := stops[i], stops[i+1]

		// add current stop to list
		ramp = append(ramp, cur)

		// create interpolated colors bewtween current and next stop.
		// first, distribute some of the remainder, if any
		extra := 0
		if remainder > 0 {
			extra = 1
			remainder--
		}

		// calculate delta for each RGB component
		dR := (int(next.R) - int(cur.R)) / (1 + betweenStops + extra)
		dG := (int(next.G) - int(cur.G)) / (1 + betweenStops + extra)
		dB := (int(next.B) - int(cur.B)) / (1 + betweenStops + extra)

		// do interpolation
		for j := 1; j <= (betweenStops + extra); j++ {
			c := color.RGBA{
				R: uint8(int(cur.R) + (dR * j)),
				G: uint8(int(cur.G) + (dG * j)),
				B: uint8(int(cur.B) + (dB * j)),
				A: 255}
			ramp = append(ramp, c)
		}

		// if next is the last stop, it has to be added
		if i == len(stops)-2 {
			ramp = append(ramp, next)
		}
	}

	return ramp
}
