package entities

import "neon/math"

// Polygon data structure represents a polygon internally using a graph
/* Essentially a very simple vertex-vertex mesh */
type Polygon struct {
	Vertices 	[]math.Vector2D
	Edges 	 	map[int][]int		// adjacency matrix for the vertices

	State 	 	EntityState 		// Refers to the current physical state of the polygon
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
		Vertices: []math.Vector2D{},
		Edges: make(map[int][]int),
		State: EntityState{
			CentroidPosition: centroid,
		},
	}

	// Now we need to connect up the edges within the adjacency matrix
	for index, v := range vertices {
		generatedPolygon.Vertices = append(generatedPolygon.Vertices, v.Sub(centroid))

		nextIndex := (index + 1) % (len(vertices))
		generatedPolygon.Edges[index] 		= append(generatedPolygon.Edges[index], nextIndex)
		generatedPolygon.Edges[nextIndex] 	= append(generatedPolygon.Edges[nextIndex], index)

	}

	return generatedPolygon
}

// addVertex adds a new vertex to the polygon, note, it does not update the centroid position
func (p Polygon) addVertex(vertex math.Vector2D) (Polygon, int) {
	// Average position is the centroid
	rawCentroid := p.State.CentroidPosition
	p.Vertices = append(p.Vertices, vertex.Sub(rawCentroid))
	p.Edges[len(p.Vertices)] = []int{}

	return p, len(p.Vertices)
}

// deleteVertex removes a vertex from a polygon, note it does not update the centroid position
func (p Polygon) deleteVertex(vertexID int) Polygon {
	p.Vertices = unsetVec(p.Vertices, vertexID)
	edges := p.Edges[vertexID]; delete(p.Edges, vertexID)

	for _, edge := range edges {
		p.Edges[edge] = del(p.Edges[edge], vertexID)
	}
	return p
}



// Updating the polygon such as adding edges or removing edges shifts the actual centroid position,
func (p Polygon) refreshPolygon() Polygon {
	centroid := p.State.CentroidPosition
	newCentroid := math.ZeroVec2D

	for id, v := range p.Vertices {
		p.Vertices[id] = v.Add(centroid); newCentroid = newCentroid.Add(p.Vertices[id])
	}

	newCentroid = newCentroid.Scale(1.0/float64(len(p.Vertices)))
	for id, v := range p.Vertices { p.Vertices[id] = v.Sub(newCentroid) }

	return p
}