package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"neon/engine"
	"neon/entities"
	"neon/math"
	"time"
)




func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, 1000, 600),
		Title: "Neon Spin",
		VSync:  false,
	})
	if err != nil {
		panic(err)
	}
	CenterWindow(win)
	win.SetSmooth(true)
	viewMatrix := pixel.IM.Moved(win.Bounds().Center())
	intermediateCanvas := pixelgl.NewCanvas(win.Bounds())




	poly := entities.NewPolygon([]math.Vector2D{
		{0, 100}, {100, 100}, {100, 0}, {0,0},
	})
	specialPoly := entities.NewPolygon([]math.Vector2D{
		{200, 300}, {300, 300}, {300, 200}, {200,200},
	})
	poly.State.AngularVelocity = math.Vector3D{Z: 1}



	p := Polygon{internal: &poly, colour: color.NRGBA{R: 255, G: 255, B: 255, A: 255}}
	p1 := Polygon{internal: &specialPoly, colour: color.NRGBA{R: 255, G: 255, B: 255, A: 255}}


	poly.State.Velocity = math.Vector2D{X: 50, Y: 50}
	physicsManager := engine.NewPhysicsManager()

	physicsManager = physicsManager.BeginTracking(&poly)
	physicsManager = physicsManager.BeginTracking(&specialPoly)


	start := time.Now()
	for !win.Closed() {
		imd := imdraw.New(nil)
		dt := time.Now().Sub(start).Seconds()
		start = time.Now()


		imd.Color = color.NRGBA{R: 0, G: 13, B: 28, A: 255}
		imd.Push(pixel.V(0, win.Bounds().H()), pixel.V(win.Bounds().W(), 0))
		imd.Rectangle(0)

		fmt.Print(entities.SAT(poly, specialPoly), "\n")


		p.Update(dt)
		p1.Update(dt)
		p.Render(imd)
		p1.Render(imd)


		imd.Draw(intermediateCanvas)
		intermediateCanvas.Draw(win, viewMatrix)
		win.Update()
	}
}








func main() {
	pixelgl.Run(run)
}

