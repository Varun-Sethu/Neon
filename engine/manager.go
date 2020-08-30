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
func (receiver PhysicsManager) BeginTracking(poly *entities.Polygon) PhysicsManager {
	receiver.trackingEntities = append(receiver.trackingEntities, poly)
	return receiver
}


// DetectCollisions identifies if any collisions are present
func (receiver PhysicsManager) DetectCollisions() bool {
	for i := range receiver.trackingEntities {
		for j := range receiver.trackingEntities[i+1:] {
			polyA, polyB := *receiver.trackingEntities[i], *receiver.trackingEntities[j]
			if DetermineCollision(polyA, polyB) {
				return true
			}
		}
	}
	return false
}






// Loop is the main loop for the manager, when this method is invoked then the manager begins to check for collisions and apply any other important physics stuff
func (receiver PhysicsManager) Loop() {
	receiver.DetectCollisions()
}