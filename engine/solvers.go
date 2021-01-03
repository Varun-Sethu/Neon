package engine

import (
	"math"
	"neon/entities"
	vmath "neon/math"
	"neon/units"
)

// Contains all the collision solvers

// ResolveCollision computes what has to be done during a collision and resolves/calculates all the physics involved with it, given a collision manifold
func (manifold ContactManifold) ResolveCollision(incidentFrame, referenceFrame *entities.Polygon) {
	manifold = manifold.processManifold()

	// First statically resolve the collision
	incidentFrame.State.CentroidPosition = incidentFrame.State.CentroidPosition.Add(manifold.MTV)


	// Then compute impulses
	if manifold.ContactCount == 1 {
		manifold.ResolvePointCollision(incidentFrame, referenceFrame)
	} else if manifold.ContactCount == 2 {
		manifold.ResolvePlanarCollision(incidentFrame, referenceFrame)
	}
}



// Resolves collision with one contact point
func (manifold ContactManifold) ResolvePointCollision(incidentFrame, referenceFrame *entities.Polygon) {
	collision_normal := manifold.MTV.Normalise()
	reference_collision_points := manifold.findContactPointsFor(manifold.ReferenceFrame, manifold.ReferenceFace)


	p_incident := manifold.CollisionPoints[0]
	p_referece := reference_collision_points[0]


	r_i, r_r := p_incident.Sub(incidentFrame.State.CentroidPosition).Scale(1.0/units.Metre), p_referece.Sub(referenceFrame.State.CentroidPosition).Scale(1.0/units.Metre)
	m_r, i_r := retrievePhysicalData(referenceFrame)
	m_i, i_i := retrievePhysicalData(incidentFrame)


	// Compute the velocities at the point of collision
	v_pi := incidentFrame.State.Velocity.Add(r_i.CrossUpwardsWithVec(incidentFrame.State.AngularVelocity))
	v_pr := referenceFrame.State.Velocity.Add(r_r.CrossUpwardsWithVec(incidentFrame.State.AngularVelocity))

	separational_velocity := v_pi.Sub(v_pr).Dot(collision_normal)
	if separational_velocity > 0 {
		return
	}

	impulse := -(2.0) * separational_velocity /
		((1.0/m_r + 1.0/m_i) +
			math.Pow(r_i.CrossMag(collision_normal), 2)/i_i +
			math.Pow(r_r.CrossMag(collision_normal), 2)/i_r)


	incidentFrame.State.ApplyImpulse(collision_normal.Scale(impulse), p_incident)
	referenceFrame.State.ApplyImpulse(collision_normal.Scale(-impulse), p_referece)
}







// Resolves collisions with two contact points
func (manifold ContactManifold) ResolvePlanarCollision(incidentFrame, referenceFrame *entities.Polygon) {

	// Compute the collision normal for simple resolution
	collision_normal := manifold.MTV.Normalise()
	v_pi := incidentFrame.State.Velocity
	v_pr := referenceFrame.State.Velocity
	m_i, _ := retrievePhysicalData(incidentFrame)
	m_r, _ := retrievePhysicalData(referenceFrame)


	impulse := -(2.0) * (v_pi.Sub(v_pr)).Dot(collision_normal) /
		(collision_normal.Dot(collision_normal) * (1.0 / m_i + 1.0 / m_r))

	incidentFrame.State.ApplyImpulse(collision_normal.Scale(impulse), vmath.ZeroVec2D)
	referenceFrame.State.ApplyImpulse(collision_normal.Scale(-impulse), vmath.ZeroVec2D)
}








// Retrieves the physical data of the polygon, note that if the polygon is not kinetic then we say it has "infinite mass"
func retrievePhysicalData(poly *entities.Polygon) (float64, float64) {
	if poly.State.NoKinetic {
		return math.Inf(1), math.Inf(1)
	}
	return poly.State.Mass, poly.State.RotationalInertia
}


// function that just processes the manifold
func (manifold ContactManifold)  processManifold() ContactManifold {
	if manifold.ContactCount == 1 || math.Abs(manifold.ContactDepths[0] - manifold.ContactDepths[1]) < 0.84 { // magic numbers :)
		return manifold
	}

	manifold.ContactCount -= 1
	invalid_contact_point := 0; if manifold.ContactDepths[1] < manifold.ContactDepths[0] {invalid_contact_point = 1}
	manifold.CollisionPoints = append(manifold.CollisionPoints[:invalid_contact_point], manifold.CollisionPoints[invalid_contact_point + 1:]...)

	return manifold
}


// For more accurate collision resolution the incident points are actually projected onto the reference face, thus our "incident" points are really the projections onto these shapes
func (manifold ContactManifold) findContactPointsFor(p entities.Polygon, face []int) []vmath.Vector2D {
	line := [2]vmath.Vector2D{
		p.Vertices[face[0]].Add(p.State.CentroidPosition),
		p.Vertices[face[1]].Add(p.State.CentroidPosition),
	}

	contactPoints := []vmath.Vector2D{}
	for _, v := range manifold.CollisionPoints {
		contactPoints = append(contactPoints, vmath.ProjectPointOntoLine(v, line))
	}

	return contactPoints
}


