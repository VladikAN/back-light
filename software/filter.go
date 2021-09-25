package main

import (
	"image/color"
	"math"
)

func (worker *worker) filterOutput(rs []color.RGBA) []color.RGBA {
	// Invert sequence
	if worker.opt.Invert {
		for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
			rs[i], rs[j] = rs[j], rs[i]
		}
	}

	for index, item := range rs {
		// Remove dark gray color cause its brighter then others
		if math.Abs(float64(item.R-item.G)) < 20 && math.Abs(float64(item.R-item.B)) < 20 && item.R < 100 {
			rs[index] = color.RGBA{}
		}

		// Remove weak colors cause it has the same brightness like others
		if item.R < 100 && item.G < 100 && item.B < 100 {
			rs[index] = color.RGBA{}
		}
	}

	return rs
}
