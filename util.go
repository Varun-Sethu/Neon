package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)


// CenterWindow will... center the window
func CenterWindow(win *pixelgl.Window) {
	x, y := pixelgl.PrimaryMonitor().Size()
	width, height := win.Bounds().Size().XY()
	win.SetPos(
		pixel.V(
			x/2-width/2,
			y/2-height/2,
		),
	)
}
