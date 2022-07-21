package engine

import (
	"Neon/entities"
)

// Physics manager keeps a list to entities it is currently tracking, it is additionally responsible for detecting collisions between only these objects
// The manager furthermore should resolve these collisions
type PhysicsManager struct {
	trackingEntities   []*entities.Polygon
	collisionCallbacks []func(manifold ContactManifold)
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

// Adds a callback function to the set of collision callback functions if a collision ever does occur
func (receiver *PhysicsManager) AddCallback(callbacks ...func(manifold ContactManifold)) {
	receiver.collisionCallbacks = append(receiver.collisionCallbacks, callbacks...)
}

// ResolveCollisions identifies if any collisions are present and resolves them if they are
func (receiver PhysicsManager) ResolveCollisions() {
	for i, a := range receiver.trackingEntities {
		for _, b := range receiver.trackingEntities[i+1:] {
			if collides, manifold := DetermineCollision(a, b); collides {
				manifold.ResolveCollision()

				// Perform the callback operations
				for _, callback := range receiver.collisionCallbacks {
					callback(manifold)
				}
			}
		}
	}
}

// NextTimeStep just progresses everything to the next timestep, just numerical integration.... Note: Every entitiy already has methods for progressing its state
func (receiver *PhysicsManager) NextTimeStep(dt float64) {
	progressEntities := func() {
		for _, e := range receiver.trackingEntities {
			e.NextTimeStep(dt)
		}
	}

	// Progress entities, resolve entities and progress again :D
	// The second progression is for smoother results
	progressEntities()
	receiver.ResolveCollisions()
	progressEntities()
}
