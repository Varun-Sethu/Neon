package entities

import (
	neonMath "Neon/engine/math"
)

// Polygon data structure represents a polygon internally using a graph
/* Essentially a very simple vertex-vertex mesh */
type Polygon struct {
	Vertices map[int]neonMath.Vector2D
	Edges    map[int][]int // adjacency matrix for the vertices

	State EntityState // Refers to the current physical state of the polygon

	// internal var for tracking vertex IDs
	prevID int
}

// Simple method to generate a new polygon
func NewPolygon(vertices []neonMath.Vector2D) Polygon {
	// First compute the centroid
	centroid := neonMath.ZeroVec2D
	for _, vertex := range vertices {
		centroid = centroid.Add(vertex)
	}
	centroid = centroid.Scale(1.0 / float64(len(vertices)))

	// The polygon we want to generate
	generatedPolygon := Polygon{
		Vertices: make(map[int]neonMath.Vector2D, len(vertices)),
		Edges:    make(map[int][]int),
		State: EntityState{
			CentroidPosition: centroid,
		},
	}

	// Now we need to connect up the edges within the adjacency matrix
	for index, v := range vertices {
		generatedPolygon.Vertices[generatedPolygon.prevID] = v.Sub(centroid)

		nextIndex := (generatedPolygon.prevID + 1) % (len(vertices))
		generatedPolygon.Edges[generatedPolygon.prevID] = append(generatedPolygon.Edges[index], nextIndex)
		generatedPolygon.Edges[nextIndex] = append(generatedPolygon.Edges[nextIndex], index)

		generatedPolygon.prevID++
	}

	return generatedPolygon
}

// Returns the endpoints of the interval defined by an edge
func (polygon *Polygon) GetEdgeCoordinates(face []int) [2]neonMath.Vector2D {
	return [2]neonMath.Vector2D{
		polygon.Vertices[face[0]].Add(polygon.State.CentroidPosition),
		polygon.Vertices[face[1]].Add(polygon.State.CentroidPosition),
	}
}
