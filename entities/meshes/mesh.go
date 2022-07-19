package meshes

import (
	"Neon/math"
)

// Meshes are represented with an adjacency list, essentially a cylic graph
type Mesh struct {
	Centroid *math.Vector2D // until I can think of a better method this stays... Essentially the centroid points to the centroid as defined in the entity state
	Vertices map[int]math.Vector2D
	EdgeSet  map[int][]int
	Radius   float64 // Only defined for circle mesh

	MeshType

	prevID int // internal var for tracking vertex IDs
}

// Table of meshes for which we can reasonably compute "intersections"

// Set of mesh types
type MeshType int

const (
	MeshCircle  = 0
	MeshPolygon = 1
)

// Determines if a mesh intersects another mesh and returns the MTV to statically resolve the collision
func (Mesh) Intersects(mesh Mesh) math.Vector2D {
	return math.Vector2D{}
}
