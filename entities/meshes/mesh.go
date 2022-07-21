package meshes

import (
	neonMath "Neon/engine/math"
)

// Meshes are represented with an adjacency list, essentially a cylic graph
// File just defines the base properties that are required for a mesh
type Mesh struct {
	Centroid *neonMath.Vector2D // until I can think of a better method this stays... Essentially the centroid points to the centroid as defined in the entity state
	Vertices map[int]neonMath.Vector2D
	EdgeSet  map[int][]int
	Radius   float64 // Only defined for circle mesh

	MeshType
}

// Table of meshes for which we can reasonably compute "intersections"

// Set of mesh types
type MeshType int

const (
	MeshCircle  = 0
	MeshPolygon = 1
)

// Determines if a mesh intersects another mesh and returns the MTV to statically resolve the collision
func (Mesh) Intersects(mesh Mesh) neonMath.Vector2D {
	return neonMath.Vector2D{}
}
