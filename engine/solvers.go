package engine

import (
	neonMath "Neon/engine/math"
	"Neon/entities"
	"math"
)

// ResolveCollision computes what has to be done during a collision and resolves/calculates all the physics involved with it, given a collision manifold
func (manifold ContactManifold) ResolveCollision() {
	incidentFrame, referenceFrame := manifold.IncidentFrame, manifold.ReferenceFrame

	if !incidentFrame.State.NoKinetic || !referenceFrame.State.NoKinetic {
		// First solve the collision by fetching a solver and invoking it
		manifold.getSolver()(incidentFrame, referenceFrame)

		// Then statically resolve the collision
		referenceFrame.State.ShiftCentroid(manifold.MTV.Scale(-1.0))
		incidentFrame.State.ShiftCentroid(manifold.MTV)
	}
}

// Simple helper function for fetching the appropriate collision solver, will be expanded later to include more specialised solvers
func (manifold *ContactManifold) getSolver() func(incidentFrame *entities.Polygon, referenceFrame *entities.Polygon) {
	switch manifold.ContactCount {
	case 1:
		return manifold.resolvePointCollision
	case 2:
		return manifold.resolvePlanarCollision
	}

	return nil
}

// Resolves collision with one contact point
func (manifold *ContactManifold) resolvePointCollision(incidentFrame, referenceFrame *entities.Polygon) {
	pIncident := manifold.CollisionPoints[0]
	pReference := manifold.CollisionPoints[0]

	impulse := manifold.resolveCollisionWithAppPoints(incidentFrame, referenceFrame, pIncident, pReference)

	// Apply the computed impulse
	incidentFrame.State.ApplyImpulse(impulse, pIncident)
	referenceFrame.State.ApplyImpulse(impulse.Scale(-1.0), pReference)
}

// Resolves collisions with two contact points
func (manifold *ContactManifold) resolvePlanarCollision(incidentFrame, referenceFrame *entities.Polygon) {
	pointA, pointB := manifold.CollisionPoints[0], manifold.CollisionPoints[1]

	// To prevent short circuit evaliation we create two variables as opposed to chaning the invocations with and statements
	impulseOne := manifold.resolveCollisionWithAppPoints(incidentFrame, referenceFrame, pointA, pointA)
	impulseTwo := manifold.resolveCollisionWithAppPoints(incidentFrame, referenceFrame, pointB, pointB)

	// Apply the computed impulses
	incidentFrame.State.ApplyImpulse(impulseOne, pointA)
	referenceFrame.State.ApplyImpulse(impulseOne.Scale(-1.0), pointA)

	incidentFrame.State.ApplyImpulse(impulseTwo, pointB)
	referenceFrame.State.ApplyImpulse(impulseTwo.Scale(-1.0), pointB)
}

// resolveCollisionWithAppPoint resolves a collision given the application point(s)
func (manifold *ContactManifold) resolveCollisionWithAppPoints(incidentFrame, referenceFrame *entities.Polygon, pIncident neonMath.Vector2D, pReference neonMath.Vector2D) neonMath.Vector2D {
	collisionNormal := manifold.MTV.Normalise()

	rI, rR := pIncident.Sub(incidentFrame.State.CentroidPosition).Scale(1.0/neonMath.Metre), pReference.Sub(referenceFrame.State.CentroidPosition).Scale(1.0/neonMath.Metre)
	mR, iR := referenceFrame.State.RetrievePhysicalData()
	mI, iI := incidentFrame.State.RetrievePhysicalData()

	// Compute the velocities at the point of collision
	vPi := incidentFrame.State.Velocity.Sub(rI.CrossUpwardsWithVec(incidentFrame.State.AngularVelocity))
	vPr := referenceFrame.State.Velocity.Add(rR.CrossUpwardsWithVec(referenceFrame.State.AngularVelocity))

	separationVelocity := vPi.Sub(vPr).Dot(collisionNormal)
	if math.IsNaN(separationVelocity) {
		return neonMath.ZeroVec2D
	}

	// weird dampening of the separation velocity, seems to produce more "reasonable" results
	// this number was just arrived at via some weird testing
	restitution := 0.954
	impulse := -(1.0 + restitution) * separationVelocity /
		((1.0/mR + 1.0/mI) +
			math.Pow(rI.CrossMag(collisionNormal), 2)/iI +
			math.Pow(rR.CrossMag(collisionNormal), 2)/iR)
	impulse /= float64(manifold.ContactCount)

	return collisionNormal.Scale(impulse)
}
