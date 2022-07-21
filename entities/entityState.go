package entities

import (
	neonMath "Neon/engine/math"
	"Neon/entities/meshes"
	"math"
)

// Entity is simply just a 2d object that the physics engine can act on
type Entity struct {
	State EntityState
	Mesh  meshes.Mesh
}

// Refers the the current state of an entity, important for physical calculations
type EntityState struct {
	// Motion quantities
	Velocity         neonMath.Vector2D
	AngularVelocity  float64 // Angular velocity is of the form: (0, 0, w)
	CentroidPosition neonMath.Vector2D

	// Inertial stuff
	Mass              float64
	RotationalInertia float64
	// Toggleable Quantity
	NoKinetic bool
}

// NewEntity creates a completely brand new entity given a meshType and the set of information that defines that mesh
func (entity *Entity) NewEntity(meshType meshes.MeshType, meshDefinition []neonMath.Vector2D) {
	switch meshType {
	case meshes.MeshPolygon:
		break
	case meshes.MeshCircle:
		break
	}
}

// ApplyImpulse just applies an impulse to the state,
// Note application point is assumed to be outside the actual polygon, as such all vectors corresponding to position are relative to (0, 0) and not the centroid
func (e *EntityState) ApplyImpulse(impulse neonMath.Vector2D, applicationPoint neonMath.Vector2D) {
	if e.NoKinetic {
		return
	}

	// Actually apply the impulse
	e.Velocity = e.Velocity.Add(impulse.Scale(1.0 / e.Mass))

	if applicationPoint != neonMath.ZeroVec2D {
		e.AngularVelocity += applicationPoint.Sub(e.CentroidPosition).Normalise().CrossMag(impulse.Normalise()) * impulse.Length() / e.RotationalInertia
	}
}

// ShiftOffset moves the centroid of the entityState by offset, returns true if the centroid was shifted
func (e *EntityState) ShiftCentroid(offset neonMath.Vector2D) bool {
	if e.NoKinetic {
		return false
	}

	e.CentroidPosition = e.CentroidPosition.Add(offset)
	return true
}

// RetrievePhysicalData fetches the physical data of the polygon, note that if the polygon is not kinetic then we say it has "infinite mass" and "infinite moment of inertia"
func (e *EntityState) RetrievePhysicalData() (float64, float64) {
	if e.NoKinetic {
		return math.Inf(1), math.Inf(1)
	}

	return e.Mass, e.RotationalInertia
}

// NextTimeStep computes the next infinitesimal timestamp
func (polygon *Polygon) NextTimeStep(dt float64) {
	// Update the actual position
	e := &polygon.State
	if e.NoKinetic {
		return
	}

	e.CentroidPosition = e.CentroidPosition.Add(e.Velocity.Scale(neonMath.Metre).Scale(dt))

	// Simple rotational utility
	matrixRotate := func(d neonMath.Vector2D, dTheta float64) neonMath.Vector2D {
		return neonMath.Vector2D{
			X: d.X*math.Cos(dTheta) - d.Y*math.Sin(dTheta),
			Y: d.X*math.Sin(dTheta) + d.Y*math.Cos(dTheta),
		}
	}
	// Compute the actual rotation of the entity
	dTheta := dt * e.AngularVelocity
	for i, _ := range polygon.Vertices {
		polygon.Vertices[i] = matrixRotate(polygon.Vertices[i], dTheta)
	}
}
