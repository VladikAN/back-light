package backlight

import (
	"image"
	"image/color"
	"math"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/disintegration/imaging"
	"github.com/kbinani/screenshot"
)

const ImageMinSize int = 320

func (worker *Worker) DrawEmpty() []color.RGBA {
	var rs []color.RGBA
	sz := worker.Opt.Width*2 + worker.Opt.Height*2
	for i := 0; i < sz; i++ {
		rs = append(rs, color.RGBA{})
	}

	return rs
}

func (worker *Worker) DrawDebug() []color.RGBA {
	var rs []color.RGBA

	rs = append(rs, color.RGBA{R: 255})
	for i := 0; i < worker.Opt.Width-1; i++ {
		rs = append(rs, color.RGBA{})
	}

	rs = append(rs, color.RGBA{G: 255})
	for i := 0; i < worker.Opt.Height-1; i++ {
		rs = append(rs, color.RGBA{})
	}

	rs = append(rs, color.RGBA{B: 255})
	for i := 0; i < worker.Opt.Width-1; i++ {
		rs = append(rs, color.RGBA{})
	}

	rs = append(rs, color.RGBA{R: 255, G: 255})
	for i := 0; i < worker.Opt.Width-1; i++ {
		rs = append(rs, color.RGBA{})
	}

	return rs
}

func (worker *Worker) DrawScreen() ([]color.RGBA, error) {
	// Take screenshot
	fl, err := screenshot.CaptureDisplay(0)
	if err != nil {
		return nil, err
	}

	// Scale down screenshot and calc sizes
	downScale := imaging.Resize(fl, ImageMinSize, ImageMinSize, imaging.NearestNeighbor)
	wf := downScale.Rect.Max.X
	ws := wf / (worker.Opt.Width + 1)
	hf := downScale.Rect.Max.Y
	hs := hf / (worker.Opt.Height + 1)

	// Crop image into smaller, dependending on LED width / height
	// Find dominant color for each piece
	var rs []color.RGBA

	// top - from left to right
	for i := 0; i < worker.Opt.Width; i++ {
		pt := imaging.Crop(downScale, image.Rectangle{
			Min: image.Point{X: ws * i, Y: 0},
			Max: image.Point{X: ws * (i + 1), Y: hs},
		})

		rs = worker.getDominant(rs, pt)
	}

	// right - from top to bottom
	for i := 0; i < worker.Opt.Height; i++ {
		pt := imaging.Crop(downScale, image.Rectangle{
			Min: image.Point{X: wf - ws, Y: hs * i},
			Max: image.Point{X: wf, Y: hs * (i + 1)},
		})

		rs = worker.getDominant(rs, pt)
	}

	// bottom - from right to left
	for i := 0; i < worker.Opt.Width; i++ {
		pt := imaging.Crop(downScale, image.Rectangle{
			Min: image.Point{X: wf - ws*(i+1), Y: hf - hs},
			Max: image.Point{X: wf - ws*i, Y: hf},
		})

		rs = worker.getDominant(rs, pt)
	}

	// left - from bottom to top
	for i := 0; i < worker.Opt.Height; i++ {
		pt := imaging.Crop(downScale, image.Rectangle{
			Min: image.Point{X: 0, Y: hf - hs*(i+1)},
			Max: image.Point{X: ws, Y: hf - hs*i},
		})

		rs = worker.getDominant(rs, pt)
	}

	worker.Prev = rs
	return rs, nil
}

func (worker *Worker) getDominant(rs []color.RGBA, pt *image.NRGBA) []color.RGBA {
	c, _ := prominentcolor.KmeansWithAll(
		1,
		pt,
		prominentcolor.ArgumentAverageMean|prominentcolor.ArgumentNoCropping,
		prominentcolor.DefaultSize,
		nil)

	record := color.RGBA{}
	if len(c) != 0 {
		fc := c[0].Color
		record = color.RGBA{R: uint8(fc.R), G: uint8(fc.G), B: uint8(fc.B)}
	}

	index := len(rs)
	if len(worker.Prev) > index {
		prev := worker.Prev[index]
		limit := uint8(worker.Opt.ColorLimit)
		if isSimilar(prev, record, limit) {
			record = prev
		}
	}

	return append(rs, record)
}

func isSimilar(a, b color.RGBA, limit uint8) bool {
	d := math.Sqrt(math.Pow(float64(a.R)-float64(b.R), 2) + math.Pow(float64(a.G)-float64(b.G), 2) + math.Pow(float64(a.B)-float64(b.B), 2))
	return d <= float64(limit)
}
