package engine

import (
	"neon/entities"
	"reflect"
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
func (receiver PhysicsManager) BeginTracking(poly *entities.Polygon) PhysicsManager {
	receiver.trackingEntities = append(receiver.trackingEntities, poly)
	return receiver
}


// DetectCollisions identifies if any collisions are present
func (receiver PhysicsManager) DetectCollisions() (bool, ContactManifold) {
	for i, a := range receiver.trackingEntities {
		for _, b := range receiver.trackingEntities[i+1:] {
			polyA, polyB := entities.NewPolygonCopy(*a), entities.NewPolygonCopy(*b)

			if collides, manifold := DetermineCollision(polyA, polyB); collides {
				reference_frame, incident_frame := a, b
				if reflect.DeepEqual(polyB, manifold.ReferenceFrame) {reference_frame, incident_frame = b, a} // TODO: you know this shouldn't even be a thing....
				manifold.ResolveCollision(incident_frame, reference_frame, 1.0)
				return true, manifold
			}
		}
	}
	return false, ContactManifold{}
}






// Loop is the main loop for the manager, when this method is invoked then the manager begins to check for collisions and apply any other important physics stuff
func (receiver PhysicsManager) Loop() {
	receiver.DetectCollisions()
}