package engine

import (
	"neon/entities"
)

// Contains all the collision solvers

// ResolveCollision computes what has to be done during a collision and resolves/calculates all the physics involved with it, given a collision manifold
func (manifold ContactManifold) ResolveCollision(incidentFrame, referenceFrame *entities.Polygon, restitution float64) {
	if incidentFrame.State.NoKinetic {
		incidentFrame = referenceFrame
		manifold.MTV = manifold.MTV.Scale(-1.0)
	}


	incidentFrame.State.CentroidPosition = incidentFrame.State.CentroidPosition.Add(manifold.MTV)

	negatable_velocity := incidentFrame.State.Velocity.Project(manifold.MTV.Normalise()).Scale(-1.0)
	incidentFrame.State.Velocity = incidentFrame.State.Velocity.Add(negatable_velocity)

	// Currently can only handle a single collision point
	//if len(manifold.CollisionPoints) > 1 {return}

	// yes . . .
	//collision_point := manifold.CollisionPoints[0]
	//relative_vector_incident, relative_vector_reference := collision_point.Sub(incidentFrame.State.CentroidPosition), collision_point.Sub(referenceFrame.State.CentroidPosition)


	// Compute the collision velocities at the point of contact
	//v_incident  := incidentFrame.State.Velocity.Add(incidentFrame.State.AngularVelocity.Cross(relative_vector_incident.ConvertTo3D()).ConvertTo2D())
	//v_reference := referenceFrame.State.Velocity.Add(referenceFrame.State.AngularVelocity.Cross(relative_vector_reference.ConvertTo3D()).ConvertTo2D())
	//v_relative 	:= v_incident.Sub(v_reference).Scale(-(1 + restitution)).Dot(manifold.MTV.Normalise())

	//inertial_component := 1.0/incidentFrame.State.Mass + 1.0/referenceFrame.State.Mass +
	//	(relative_vector_incident.ConvertTo3D().Cross(manifold.MTV.ConvertTo3D().Normalise()).Cross(relative_vector_incident.ConvertTo3D())).Length() / incidentFrame.State.RotationalInertia +
	//	(relative_vector_reference.ConvertTo3D().Cross(manifold.MTV.ConvertTo3D().Normalise()).Cross(relative_vector_reference.ConvertTo3D())).Length() / referenceFrame.State.RotationalInertia

	//impulse_mag := v_relative / inertial_component

	//incidentFrame.State.ApplyImpulse(manifold.MTV.Normalise().Scale(-impulse_mag), collision_point)
	//referenceFrame.State.ApplyImpulse(manifold.MTV.Normalise().Scale(impulse_mag), collision_point)
}



// Resolves collisions with two contact points
func (manifold ContactManifold) ResolvePlanarCollision(incidentFrame, referenceFrame *entities.Polygon, restitution float64) {

	// Compute the collision normal for simple resolution
	//collision_normal = math.ComputeOutwardsNormal(manifold.CollisionPoints[0], manifold.CollisionPoints[1], referenceFrame.State.CentroidPosition)



}

