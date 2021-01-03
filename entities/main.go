package entities

import (
	"neon/math"
	"neon/units"
)
import gMath "math"

// Refers the the current state of an entity
type EntityState struct {
	// Motion quantities
	Velocity 			math.Vector2D
	AngularVelocity 	float64 // Angular velocity is of the form: (0, 0, w)
	CentroidPosition	math.Vector2D


	// Inertial stuff
	Mass				float64
	RotationalInertia	float64
	// Toggleable Quantity
	NoKinetic			bool
}




// ApplyImpulse just applies an impulse to the state,
// Note application point is assumed to be outside the actual polygon, as such all vectors corresponding to position are relative to (0, 0) and not the centroid
func (e *EntityState) ApplyImpulse(impulse math.Vector2D, applicationPoint math.Vector2D) {
	if e.NoKinetic {return}

	// Actually apply the impulse
	e.Velocity = e.Velocity.Add(impulse.Scale(1.0/e.Mass))

	if applicationPoint != math.ZeroVec2D {
		e.AngularVelocity += applicationPoint.Sub(e.CentroidPosition).Normalise().CrossMag(impulse.Normalise()) * impulse.Length() / e.RotationalInertia
	}
}



// NextTimeStep computes the next infinitesimal timestamp
func (p *Polygon) NextTimeStep(dt float64) {
	// Update the actual position
	e := &p.State
	if e.NoKinetic {return}

	e.CentroidPosition = e.CentroidPosition.Add(e.Velocity.Scale(units.Metre).Scale(dt))


	// Simple rotational utility
	matrixRotate := func (d math.Vector2D, dTheta float64) math.Vector2D {
		return math.Vector2D{
			X: d.X * gMath.Cos(dTheta) - d.Y * gMath.Sin(dTheta),
			Y: d.X * gMath.Sin(dTheta) + d.Y * gMath.Cos(dTheta),
		}
	}
	// Compute the actual rotation of the entity
	dTheta := dt * e.AngularVelocity
	for i, _ := range p.Vertices {
		p.Vertices[i] = matrixRotate(p.Vertices[i], dTheta)
	}
}



