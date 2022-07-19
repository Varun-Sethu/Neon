package main

import (
	"Neon/engine"
	"Neon/engine/math"
	"Neon/entities"
	"flag"
	"image/color"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

// just defines our polygons
func definePolygons() (*entities.Polygon, *entities.Polygon, *entities.Polygon) {
	poly := entities.NewPolygon([]math.Vector2D{
		{0, 100}, {100, 100}, {100, 0}, {0, 0},
	})
	poly.State.Mass = 10.0
	poly.State.RotationalInertia = 5.0

	specialPoly := entities.NewPolygon([]math.Vector2D{
		{200, 300}, {300, 300}, {300, 200}, {200, 200},
	})
	specialPoly.State.Mass = 5.0
	specialPoly.State.RotationalInertia = 2.0

	poly.State.AngularVelocity = 0.3
	specialPoly.State.AngularVelocity = 0.1

	translationalPoly := entities.NewPolygon([]math.Vector2D{
		{600, 300}, {700, 300}, {700, 200}, {600, 200},
	})

	translationalPoly.State.Mass = 2.0
	translationalPoly.State.RotationalInertia = 1.0
	translationalPoly.State.Velocity = math.Vector2D{X: -0.5}
	poly.State.Velocity = math.Vector2D{X: 1, Y: 1}

	return &poly, &specialPoly, &translationalPoly // be freeeeeee to the heap
}

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, 1000, 600),
		Title:  "Neon Spin",
		VSync:  true,
	})
	if err != nil {
		panic(err)
	}
	CenterWindow(win)
	win.SetSmooth(true)
	viewMatrix := pixel.IM.Moved(win.Bounds().Center())
	intermediateCanvas := pixelgl.NewCanvas(win.Bounds())
	imd := imdraw.New(nil)

	// define our smol polygons
	poly, specialPoly, translationalPoly := definePolygons()

	p := Polygon{internal: poly, colour: color.NRGBA{R: 255, G: 255, B: 255, A: 255}}
	p1 := Polygon{internal: specialPoly, colour: color.NRGBA{R: 255, G: 255, B: 255, A: 255}}
	p2 := Polygon{internal: translationalPoly, colour: color.NRGBA{R: 255, G: 255, B: 255, A: 255}}

	// hook em up to the manager
	physicsManager := engine.NewPhysicsManager()
	physicsManager.BeginTracking(poly, specialPoly, translationalPoly)

	// Callback for just drawing in the collision points
	physicsManager.AddCallback(func(manifold engine.ContactManifold) {
		if manifold.ContactCount != 0 {
			imd.Color = color.NRGBA{255, 0, 255, 255}
			for _, p := range manifold.CollisionPoints {
				imd.Push(pixel.V(p.X, p.Y))
				imd.Circle(5, 0)
			}
		}
	})

	start := time.Now()
	for !win.Closed() {
		dt := time.Since(start).Seconds()
		start = time.Now()

		imd.Color = color.NRGBA{R: 0, G: 13, B: 28, A: 255}
		imd.Push(pixel.V(0, win.Bounds().H()), pixel.V(win.Bounds().W(), 0))
		imd.Rectangle(0)

		// core physics
		physicsManager.NextTimeStep(dt)

		p.Render(imd)
		p1.Render(imd)
		p2.Render(imd)

		imd.Draw(intermediateCanvas)
		intermediateCanvas.Draw(win, viewMatrix)
		win.Update()
	}
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	pixelgl.Run(run)
}
