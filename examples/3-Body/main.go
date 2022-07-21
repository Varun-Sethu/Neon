package main

import (
	"Neon/engine"
	neonMath "Neon/engine/math"
	"Neon/entities"
	"flag"
	"image/color"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

// just defines our polygons
func definePolygons() []*entities.Polygon {
	poly := entities.NewPolygon([]neonMath.Vector2D{
		{X: 10, Y: 110}, {X: 110, Y: 110}, {X: 110, Y: 10}, {X: 10, Y: 10},
	})
	poly.State.Mass = 2.0
	poly.State.RotationalInertia = 2.0

	specialPoly := entities.NewPolygon([]neonMath.Vector2D{
		{X: 200, Y: 300}, {X: 300, Y: 300}, {X: 300, Y: 200}, {X: 200, Y: 200},
	})
	specialPoly.State.Mass = 2.0
	specialPoly.State.RotationalInertia = 2.0

	poly.State.AngularVelocity = 1.0
	specialPoly.State.AngularVelocity = 0.1

	translationalPoly := entities.NewPolygon([]neonMath.Vector2D{
		{X: 600, Y: 300}, {X: 700, Y: 300}, {X: 700, Y: 200}, {X: 600, Y: 200},
	})

	translationalPoly.State.Mass = 2.1
	translationalPoly.State.RotationalInertia = 2.0
	translationalPoly.State.Velocity = neonMath.Vector2D{X: -0.5}
	poly.State.Velocity = neonMath.Vector2D{X: 0.5, Y: 0.5}

	return []*entities.Polygon{&poly, &specialPoly, &translationalPoly} // be freeeeeee to the heap
}

func defineCornerPolygons(height, width float64) []*entities.Polygon {
	top := entities.NewPolygon([]neonMath.Vector2D{{X: -100, Y: height - 1}, {X: width + 100, Y: height - 1}, {X: width + 100, Y: height + 1}, {X: -100, Y: height + 1}})
	bottom := entities.NewPolygon([]neonMath.Vector2D{{X: -100, Y: 1}, {X: width + 100, Y: 1}, {X: width + 100, Y: -200}, {X: -100, Y: -200}})
	left := entities.NewPolygon([]neonMath.Vector2D{{X: -100, Y: 1}, {X: -100, Y: height - 1}, {X: 1, Y: height - 1}, {X: 1, Y: 0}})
	right := entities.NewPolygon([]neonMath.Vector2D{{X: width + 100, Y: 1}, {X: width + 100, Y: height - 1}, {X: width - 1, Y: height - 1}, {X: width - 1, Y: 0}})

	top.State.NoKinetic = true
	bottom.State.NoKinetic = true
	left.State.NoKinetic = true
	right.State.NoKinetic = true

	return []*entities.Polygon{&top, &bottom, &left, &right}
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

	physicsPolys := append(definePolygons(), defineCornerPolygons(win.Bounds().H(), win.Bounds().W())...)
	drawablePolys := []Polygon{}

	for _, poly := range physicsPolys {
		drawablePolys = append(drawablePolys, Polygon{internal: poly, colour: color.NRGBA{R: 228, G: 233, B: 242, A: 255}})
	}

	// hook em up to the manager
	physicsManager := engine.NewPhysicsManager()
	physicsManager.BeginTracking(physicsPolys...)

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

		// core physics
		intermediateCanvas.Clear(color.NRGBA{R: 0, G: 13, B: 28, A: 255})
		imd.Clear()
		physicsManager.NextTimeStep(dt)

		for _, p := range drawablePolys {
			p.Render(imd)
		}

		imd.Draw(intermediateCanvas)
		intermediateCanvas.Draw(win, viewMatrix)

		win.Update()
	}
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

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

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
