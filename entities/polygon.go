package entities

import (
	"neon/math"
)

// Polygon data structure represents a polygon internally using a graph
/* Essentially a very simple vertex-vertex mesh */
type Polygon struct {
	Vertices 	map[int]math.Vector2D
	Edges 	 	map[int][]int		// adjacency matrix for the vertices

	State 	 	EntityState 		// Refers to the current physical state of the polygon


	// internal var for tracking vertex IDs
	prevID		int
}


// Simple method to generate a new polygon
func NewPolygon(vertices []math.Vector2D) Polygon {
	// First compute the centroid
	centroid := math.ZeroVec2D
	for _, vertex := range vertices {
		centroid = centroid.Add(vertex)
	}
	centroid = centroid.Scale(1.0/float64(len(vertices)))

	// The polygon we want to generate
	generatedPolygon := Polygon{
		Vertices: make(map[int]math.Vector2D, len(vertices)),
		Edges: make(map[int][]int),
		State: EntityState{
			CentroidPosition: centroid,
		},
	}

	// Now we need to connect up the edges within the adjacency matrix
	for index, v := range vertices {
		generatedPolygon.Vertices[generatedPolygon.prevID] = v.Sub(centroid)

		nextIndex 										:= (generatedPolygon.prevID + 1) % (len(vertices))
		generatedPolygon.Edges[generatedPolygon.prevID] = append(generatedPolygon.Edges[index], nextIndex)
		generatedPolygon.Edges[nextIndex] 				= append(generatedPolygon.Edges[nextIndex], index)

		generatedPolygon.prevID++
	}

	return generatedPolygon
}




// NewPolygonCopy creates a new polygon via a DEEP COPY of an existing one
func NewPolygonCopy(p Polygon) Polygon {
	new_poly := p
	new_poly.Vertices = make(map[int]math.Vector2D)
	new_poly.Edges = make(map[int][]int)

	for k, v := range p.Vertices {
		new_poly.Vertices[k] = v
	}

	for k, v := range p.Edges {
		new_poly.Edges[k] = v
	}

	return new_poly
}











// addVertex adds a new vertex to the polygon, note, it does not update the centroid position
// the function also assumes that the vertices are in world form
func (p Polygon) addVertex(vertex math.Vector2D) (Polygon, int) {
	// Average position is the centroid
	p.Vertices[p.prevID] = vertex
	p.Edges[p.prevID] = []int{}

	p.prevID++

	return p, p.prevID - 1
}





// deleteVertices deletes a set of vertices from a polygon, note the operation is rather expensive due to the graph representation of a polygon that is used
// also assumes there are no outgoing or incoming connections into the vertices
func (p Polygon) deleteVertices(vertices []int) Polygon {
	for _, vertexID := range vertices {
		delete(p.Vertices, vertexID)
		delete(p.Edges, vertexID)
	}
	return p
}




// Updating the polygon such as adding edges or removing edges shifts the actual centroid position
// Assumes that the polygon is in world-vertex form
func (p Polygon) RefreshPolygon() Polygon {
	newCentroid := math.ZeroVec2D

	for _ , v := range p.Vertices {
		newCentroid = newCentroid.Add(v)
	}

	newCentroid = newCentroid.Scale(1.0/float64(len(p.Vertices)))
	for id, v := range p.Vertices { p.Vertices[id] = v.Sub(newCentroid) }
	p.State.CentroidPosition = newCentroid

	return p
}