package backlight

import (
	"fmt"
	"image/color"
	"math"
	"strings"
)

func (worker *Worker) ToSerial(rs []color.RGBA) string {
	sb := &strings.Builder{}
	for _, c := range rs {
		sb.WriteString(fmt.Sprintf("%02x%02x%02x;", c.R, c.G, c.B))
	}

	return sb.String()
}

func (worker *Worker) FilterOutput(rs []color.RGBA) []color.RGBA {
	// Invert sequence
	if worker.Opt.Invert {
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
		if item.R < 50 && item.G < 50 && item.B < 50 {
			rs[index] = color.RGBA{}
		}
	}

	return rs
}
