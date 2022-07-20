package main

import (
	"Neon/engine/math"
	"Neon/entities"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Convert internal vector into a pixel vector
func internalToPixelVec(v math.Vector2D) pixel.Vec {
	return pixel.V(v.X, v.Y)
}

// Drawable represent abstractions over internal entities
type Polygon struct {
	internal *entities.Polygon
	colour   color.NRGBA
}

// Called in the main update loop
func (p Polygon) Update(dt float64) {
	p.internal.NextTimeStep(dt)
}

// Render takes an IMDraw object and draws all the vertices to it for rendering
func (p Polygon) Render(imd *imdraw.IMDraw) {
	imd.Color = p.colour
	for _, vc := range p.internal.Vertices {
		cp := internalToPixelVec(vc.Add(p.internal.State.CentroidPosition))
		imd.Push(cp)
	}
	imd.Polygon(0.0)
}
