package engine

import (
	vmath "Neon/engine/math"
	"Neon/entities"
	"math"
)

// Contains all the collision solvers

// ResolveCollision computes what has to be done during a collision and resolves/calculates all the physics involved with it, given a collision manifold
func (manifold ContactManifold) ResolveCollision() {
	incidentFrame, referenceFrame := manifold.IncidentFrame, manifold.ReferenceFrame

	if incidentFrame.State.NoKinetic && referenceFrame.State.NoKinetic {
		return
	}

	// First compute the impulses to apply to the physics objects
	if manifold.ContactCount == 1 {
		manifold.resolvePointCollision(incidentFrame, referenceFrame)
	} else if manifold.ContactCount == 2 {
		manifold.resolvePlanarCollision(incidentFrame, referenceFrame)
	}

	// Then statically resolve the collision
	if !referenceFrame.State.NoKinetic {
		referenceFrame.State.CentroidPosition = referenceFrame.State.CentroidPosition.Add(manifold.MTV.Scale(-1.0))
	} else if !incidentFrame.State.NoKinetic {
		incidentFrame.State.CentroidPosition = incidentFrame.State.CentroidPosition.Add(manifold.MTV)
	}
}

// Resolves collision with one contact point
func (manifold *ContactManifold) resolvePointCollision(incidentFrame, referenceFrame *entities.Polygon) bool {
	referenceCollisionPoints := manifold.findContactPointsFor(manifold.ReferenceFrame, manifold.ReferenceFace)

	pIncident := manifold.CollisionPoints[0]
	pReference := referenceCollisionPoints[0]

	return manifold.resolveCollisionWithAppPoints(incidentFrame, referenceFrame, pIncident, pReference)
}

// Resolves collisions with two contact points
func (manifold *ContactManifold) resolvePlanarCollision(incidentFrame, referenceFrame *entities.Polygon) bool {
	referenceCollisionPoints := manifold.findContactPointsFor(manifold.ReferenceFrame, manifold.ReferenceFace)

	pIncident := (manifold.CollisionPoints[0].Add(manifold.CollisionPoints[1])).Scale(0.5)
	pReference := (referenceCollisionPoints[0].Add(referenceCollisionPoints[1])).Scale(0.5)

	return manifold.resolveCollisionWithAppPoints(incidentFrame, referenceFrame, pIncident, pReference)
}

// resolveCollisionWithAppPoint resolves a collision given the application point(s)
func (manifold *ContactManifold) resolveCollisionWithAppPoints(incidentFrame, referenceFrame *entities.Polygon, pIncident vmath.Vector2D, pReference vmath.Vector2D) bool {
	collisionNormal := manifold.MTV.Normalise()

	rI, rR := pIncident.Sub(incidentFrame.State.CentroidPosition).Scale(1.0/vmath.Metre), pReference.Sub(referenceFrame.State.CentroidPosition).Scale(1.0/vmath.Metre)
	mR, iR := retrievePhysicalData(referenceFrame)
	mI, iI := retrievePhysicalData(incidentFrame)

	// Compute the velocities at the point of collision
	vPi := incidentFrame.State.Velocity.Add(rI.CrossUpwardsWithVec(incidentFrame.State.AngularVelocity))
	vPr := referenceFrame.State.Velocity.Sub(rR.CrossUpwardsWithVec(referenceFrame.State.AngularVelocity))

	separationVelocity := vPi.Sub(vPr).Dot(collisionNormal)
	if math.IsNaN(separationVelocity) {
		return false
	}

	// weird dampening of the separation velocity, seems to produce more "reasonable" results
	// this number was just arrived at via some weird testing
	restitution := 0.73
	if math.Abs(separationVelocity) < 0.1 {
		restitution = 1.0
	}
	impulse := -(1.0 + restitution) * separationVelocity /
		((1.0/mR + 1.0/mI) +
			math.Pow(rI.CrossMag(collisionNormal), 2)/iI +
			math.Pow(rR.CrossMag(collisionNormal), 2)/iR)

	incidentFrame.State.ApplyImpulse(collisionNormal.Scale(impulse), pIncident)
	referenceFrame.State.ApplyImpulse(collisionNormal.Scale(-impulse), pReference)
	return true
}

// Retrieves the physical data of the polygon, note that if the polygon is not kinetic then we say it has "infinite mass" and "infinite moment of inertia"
func retrievePhysicalData(poly *entities.Polygon) (float64, float64) {
	if poly.State.NoKinetic {
		return math.Inf(1), math.Inf(1)
	}
	return poly.State.Mass, poly.State.RotationalInertia
}

// For more accurate collision resolution the incident points are actually projected onto the reference face, thus our "incident" points are really the projections onto these shapes
func (manifold ContactManifold) findContactPointsFor(p *entities.Polygon, face []int) []vmath.Vector2D {
	line := p.GetEdgeCoordinates(face)

	var contactPoints []vmath.Vector2D
	for _, v := range manifold.CollisionPoints {
		contactPoints = append(contactPoints, vmath.ProjectPointOntoLine(v, line).Add(p.State.CentroidPosition))
	}

	return contactPoints
}
