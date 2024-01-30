package backlight

import (
	"context"
	"image/color"
)

type Options struct {
	// LED line width, how many LED's used for top and bottom parts of the screen. Can be formed in square or circle.
	Width int
	// LED line height, how many LED's used for left and right parts of the screen. Can be formed in square or circle.
	Height int
	// Draw debug with constant colors. Red - top-left, Green - top-right, Blue - bottom-right and Yellow - bottom-left.
	IsDebug bool
	// Time ms between screen captures and data transfer to serial port.
	RefreshRate int
	// Serial data speed, the same as Arduino has
	SerialSpeed int
	// Time ms between serial port searches
	ConnectTimeout int64
	// The limit to consider color value was not changed
	ColorLimit int
	// Invert axis
	Invert bool
}

type Worker struct {
	Opt  *Options
	In   chan string
	Ctx  context.Context
	Prev []color.RGBA
}
