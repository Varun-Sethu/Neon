package engine

import (
	"neon/entities"
)

// Physics manager keeps a list to entities it is currently tracking, it is additionally responsible for detecting collisions between only these objects
// The manager furthermore should resolve these collisions
type PhysicsManager struct {
	trackingEntities	[]*entities.Polygon
}


func NewPhysicsManager() PhysicsManager {
	return PhysicsManager{
		trackingEntities: []*entities.Polygon{},
	}
}



// Adds a polygon to the tracking list
func (receiver *PhysicsManager) BeginTracking(polys ...*entities.Polygon) {
	receiver.trackingEntities = append(receiver.trackingEntities, polys...)
}


// DetectCollisions identifies if any collisions are present and resolves them if they are
func (receiver PhysicsManager) DetectCollisions() {
	for i, a := range receiver.trackingEntities {
		for _, b := range receiver.trackingEntities[i+1:] {
			if collides, manifold := DetermineCollision(a, b); collides {
				manifold.ResolveCollision()
			}
		}
	}
}






// Loop is the main loop for the manager, when this method is invoked then the manager begins to check for collisions and apply any other important physics stuff
func (receiver *PhysicsManager) Loop() {
	receiver.DetectCollisions()
}


// NextTimeStep just progresses everything to the next timestep, just numerical integration.... Note: Every entitiy already has methods for progressing its state
func (receiver *PhysicsManager) NextTimeStep(dt float64) {
	for _, e := range receiver.trackingEntities {
		e.NextTimeStep(dt)
	}
}
