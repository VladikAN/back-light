package backlight

import (
	"image"
	"image/color"

	"github.com/EdlinOrg/prominentcolor"
	"github.com/disintegration/imaging"
	"github.com/kbinani/screenshot"
)

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
	downScale := imaging.Resize(fl, 320, 0, imaging.NearestNeighbor)
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

		c, _ := prominentcolor.KmeansWithAll(1, pt, prominentcolor.ArgumentDefault, prominentcolor.DefaultSize, nil)
		rs = append(rs, color.RGBA{R: uint8(c[0].Color.R), G: uint8(c[0].Color.G), B: uint8(c[0].Color.B)})
	}

	// right - from top to bottom
	for i := 0; i < worker.Opt.Height; i++ {
		pt := imaging.Crop(downScale, image.Rectangle{
			Min: image.Point{X: wf - ws, Y: hs * i},
			Max: image.Point{X: wf, Y: hs * (i + 1)},
		})

		c, _ := prominentcolor.KmeansWithAll(1, pt, prominentcolor.ArgumentDefault, prominentcolor.DefaultSize, nil)
		rs = append(rs, color.RGBA{R: uint8(c[0].Color.R), G: uint8(c[0].Color.G), B: uint8(c[0].Color.B)})
	}

	// bottom - from right to left
	for i := 0; i < worker.Opt.Width; i++ {
		pt := imaging.Crop(downScale, image.Rectangle{
			Min: image.Point{X: wf - ws*(i+1), Y: hf - hs},
			Max: image.Point{X: wf - ws*i, Y: hf},
		})

		c, _ := prominentcolor.KmeansWithAll(1, pt, prominentcolor.ArgumentDefault, prominentcolor.DefaultSize, nil)
		rs = append(rs, color.RGBA{R: uint8(c[0].Color.R), G: uint8(c[0].Color.G), B: uint8(c[0].Color.B)})
	}

	// left - from bottom to top
	for i := 0; i < worker.Opt.Height; i++ {
		pt := imaging.Crop(downScale, image.Rectangle{
			Min: image.Point{X: 0, Y: hf - hs*(i+1)},
			Max: image.Point{X: ws, Y: hf - hs*i},
		})

		c, _ := prominentcolor.KmeansWithAll(1, pt, prominentcolor.ArgumentDefault, prominentcolor.DefaultSize, nil)
		rs = append(rs, color.RGBA{R: uint8(c[0].Color.R), G: uint8(c[0].Color.G), B: uint8(c[0].Color.B)})
	}

	return rs, nil
}
