package engine

import (
	vmath "Neon/engine/math"
	"Neon/entities"
	"math"
)

// Contains all the collision solvers

// ResolveCollision computes what has to be done during a collision and resolves/calculates all the physics involved with it, given a collision manifold
func (manifold ContactManifold) ResolveCollision() {
	manifold.processManifold()
	incidentFrame, referenceFrame := manifold.IncidentFrame, manifold.ReferenceFrame

	// First statically resolve the collision
	incidentFrame.State.CentroidPosition = incidentFrame.State.CentroidPosition.Add(manifold.MTV)
	referenceFrame.State.CentroidPosition = referenceFrame.State.CentroidPosition.Add(manifold.MTV.Scale(-1.0))

	// Then compute impulses
	if manifold.ContactCount == 1 {
		manifold.resolvePointCollision(incidentFrame, referenceFrame)
	} else if manifold.ContactCount == 2 {
		manifold.resolvePlanarCollision(incidentFrame, referenceFrame)
	}
}

// Resolves collision with one contact point
func (manifold *ContactManifold) resolvePointCollision(incidentFrame, referenceFrame *entities.Polygon) {
	collisionNormal := manifold.MTV.Normalise()
	referenceCollisionPoints := manifold.findContactPointsFor(manifold.ReferenceFrame, manifold.ReferenceFace)

	pIncident := manifold.CollisionPoints[0]
	pReference := referenceCollisionPoints[0]

	rI, rR := pIncident.Sub(incidentFrame.State.CentroidPosition).Scale(1.0/vmath.Metre), pReference.Sub(referenceFrame.State.CentroidPosition).Scale(1.0/vmath.Metre)
	mR, iR := retrievePhysicalData(referenceFrame)
	mI, iI := retrievePhysicalData(incidentFrame)

	// Compute the velocities at the point of collision
	vPi := incidentFrame.State.Velocity.Add(rI.CrossUpwardsWithVec(incidentFrame.State.AngularVelocity))
	vPr := referenceFrame.State.Velocity.Add(rR.CrossUpwardsWithVec(referenceFrame.State.AngularVelocity))

	separationVelocity := vPi.Sub(vPr).Dot(collisionNormal)
	if separationVelocity > -0.001 {
		return
	}

	impulse := -(2.0) * separationVelocity /
		((1.0/mR + 1.0/mI) +
			math.Pow(rI.CrossMag(collisionNormal), 2)/iI +
			math.Pow(rR.CrossMag(collisionNormal), 2)/iR)

	incidentFrame.State.ApplyImpulse(collisionNormal.Scale(impulse), pIncident)
	referenceFrame.State.ApplyImpulse(collisionNormal.Scale(-impulse), pReference)
}

// Resolves collisions with two contact points
func (manifold *ContactManifold) resolvePlanarCollision(incidentFrame, referenceFrame *entities.Polygon) {

	// Compute the collision normal for simple resolution
	collisionNormal := manifold.MTV.Normalise()
	vPi := incidentFrame.State.Velocity
	vPr := referenceFrame.State.Velocity
	mI, _ := retrievePhysicalData(incidentFrame)
	mR, _ := retrievePhysicalData(referenceFrame)

	separationVelocity := vPi.Sub(vPr).Dot(collisionNormal)
	if separationVelocity > -0.001 {
		return
	}

	impulse := -(2.0) * separationVelocity /
		(collisionNormal.Dot(collisionNormal) * (1.0/mI + 1.0/mR))

	incidentFrame.State.ApplyImpulse(collisionNormal.Scale(impulse), vmath.ZeroVec2D)
	referenceFrame.State.ApplyImpulse(collisionNormal.Scale(-impulse), vmath.ZeroVec2D)
}

// Retrieves the physical data of the polygon, note that if the polygon is not kinetic then we say it has "infinite mass"
func retrievePhysicalData(poly *entities.Polygon) (float64, float64) {
	if poly.State.NoKinetic {
		return math.Inf(1), math.Inf(1)
	}
	return poly.State.Mass, poly.State.RotationalInertia
}

// function that just processes the manifold
func (manifold *ContactManifold) processManifold() {
	if manifold.ContactCount == 1 || math.Abs(manifold.ContactDepths[0]-manifold.ContactDepths[1]) > 0.2 { // magic numbers :)
		return
	}

	manifold.ContactCount -= 1
	invalidContactPoint := 0
	if manifold.ContactDepths[1] < manifold.ContactDepths[0] {
		invalidContactPoint = 1
	}
	manifold.CollisionPoints = append(manifold.CollisionPoints[:invalidContactPoint], manifold.CollisionPoints[invalidContactPoint+1:]...)

	return
}

// For more accurate collision resolution the incident points are actually projected onto the reference face, thus our "incident" points are really the projections onto these shapes
func (manifold ContactManifold) findContactPointsFor(p *entities.Polygon, face []int) []vmath.Vector2D {
	line := p.GetEdgeCoordinates(face)

	var contactPoints []vmath.Vector2D
	for _, v := range manifold.CollisionPoints {
		contactPoints = append(contactPoints, vmath.ProjectPointOntoLine(v, line))
	}

	return contactPoints
}
