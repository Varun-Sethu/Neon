package entities

import (
	neonMath "Neon/engine/math"
	"math"
	"sort"
)

/*
	This entire file handles elementary geometry for polygons, algorithms such as SAT and Polygon Clipping are implemented here
	Additionally, manifold generation is also implemented here
*/

// GetSupportingPoint returns the supporting point of a polygon along a specific axis
func (polygon *Polygon) GetSupportingPoint(axis neonMath.Vector2D) (neonMath.Vector2D, int) {
	currentMaxProj := math.Inf(-1)
	bestVertex := neonMath.ZeroVec2D
	vertexID := 0

	for id, v := range polygon.Vertices {
		projection := v.Dot(axis)
		if projection > currentMaxProj {
			currentMaxProj = projection
			bestVertex = v.Add(polygon.State.CentroidPosition)
			vertexID = id
		}
	}
	return bestVertex, vertexID
}

// DetermineSupportingEdge determines the edge furthest along a specific axis, in essence the specific edge, also returns how "parallel" the edge's normal is with the provided normal
func (polygon *Polygon) DetermineSupportingEdge(axis neonMath.Vector2D) ([]int, float64) {

	_, vertexID := polygon.GetSupportingPoint(axis)
	v := polygon.Vertices[vertexID]
	A := polygon.Vertices[polygon.Edges[vertexID][0]]
	normalA := neonMath.ComputeOutwardsNormal(A, v, polygon.State.CentroidPosition)
	B := polygon.Vertices[polygon.Edges[vertexID][1]]
	normalB := neonMath.ComputeOutwardsNormal(B, v, polygon.State.CentroidPosition)

	// There are two possible other vertices that can connect to this one, we determine which one is significant by comparing dot products
	if math.Abs(normalA.Dot(axis)) >= math.Abs(normalB.Dot(axis)) {
		return []int{vertexID, polygon.Edges[vertexID][0]}, math.Abs(normalA.Dot(axis))
	}
	return []int{vertexID, polygon.Edges[vertexID][1]}, math.Abs(normalB.Dot(axis))
}

// PolyVerticesOutside takes a "line" (2 points) and a normal and returns all points of the polygon that lie outside this line, in the direction anti-parallel to the normal
// Note: the function assumes that the vertices are in the "world frame"
func (polygon *Polygon) PolyVerticesOutside(line [2]neonMath.Vector2D, normal neonMath.Vector2D) []int {
	var outside []int

	for i, v := range polygon.Vertices {
		// Just ensure that you determine the world position of the vertex
		if v.Sub(line[0]).Dot(normal) < 0 {
			outside = append(outside, i)
		}
	}

	return outside
}

// AxisProjection returns the projection interval of a polygon onto an axis
// Note that the projections are in "world coordinates"
func (polygon *Polygon) AxisProjection(axis neonMath.Vector2D) []float64 {
	// We need to first determine the supporting points along the axis perpendicular to the provided axis
	axis = axis.Normalise()
	supLeft, _ := polygon.GetSupportingPoint(axis)
	supRight, _ := polygon.GetSupportingPoint(axis.Scale(-1.0))

	interval := []float64{
		supLeft.ScalarProject(axis),
		supRight.ScalarProject(axis),
	}
	sort.Float64s(interval)
	return interval
}

// projectionsOverlap determines if two projections onto an axis overlap as well as the degree of overlap
func projectionsOverlap(projectionA, projectionB []float64) (bool, float64) {
	// evidently there is no overlap
	if projectionA[1] < projectionB[0] || projectionB[1] < projectionA[0] {
		return false, 0.0
	}

	return true, math.Min(projectionA[1], projectionB[1]) - math.Max(projectionA[0], projectionB[0])
}

// satSinglePolygon just checks if polyB intersects polyA
func satSinglePolygon(polyA Polygon, polyB Polygon) neonMath.Vector2D {
	mtv := neonMath.BigVec2D

	for vertex, edges := range polyA.Edges {
		for _, edge := range edges {
			worldVertex := polyA.Vertices[vertex].Add(polyA.State.CentroidPosition)                         // worldVertex refers to the world coordinates of the vertex
			worldEdgeV := polyA.Vertices[edge].Add(polyA.State.CentroidPosition)                            // worldEdgeV is the world coordinates of the other vertex that defines this edge
			normal := neonMath.ComputeOutwardsNormal(worldEdgeV, worldVertex, polyA.State.CentroidPosition) // normal is just the normal vector associated with this edge

			projectedAxisPolyb := polyB.AxisProjection(normal)
			projectedAxisPolya := polyA.AxisProjection(normal)

			if intersects, overlap := projectionsOverlap(projectedAxisPolya, projectedAxisPolyb); !intersects {
				return neonMath.ZeroVec2D
			} else {
				mtv = neonMath.Min(mtv, normal.Scale(overlap))
			}
		}
	}
	// Polygons with parallel edges have an issue with computing the same normal axis multiple times, the inevitable consequence of this is that sometimes the MTV faces TOWARDS polygon A instead of away
	// to resolve this we just need to check that the MTV points away and flip it if it doesnt, if the dot product of the normal and the separation vector is positive then it points towards it
	flip := mtv.Dot(polyB.State.CentroidPosition.Sub(polyA.State.CentroidPosition)) < 0
	if flip {
		return mtv.Scale(-1.0)
	} else {
		return mtv
	}
}

// MapToWorldSpace takes a polygon whose vertices internally are in the COM frame and converts them to the "world frame"
func (polygon *Polygon) MapToWorldSpace() {
	for vertexId, _ := range polygon.Vertices {
		polygon.Vertices[vertexId] = polygon.Vertices[vertexId].Add(polygon.State.CentroidPosition)
	}
}

// If a polygon is in world space eg. all the coordinates are relative to the global origin then this just maps them all out of it
func (polygon *Polygon) MapOutofWorldSpace() {
	for vertexId, _ := range polygon.Vertices {
		polygon.Vertices[vertexId] = polygon.Vertices[vertexId].Sub(polygon.State.CentroidPosition)
	}
}

// SAT determines if two polygons are intersecting and computes the corresponding MTV
// Note that the MTV ALWAYS POINTS FROM A TO B
func SAT(polyA Polygon, polyB Polygon) neonMath.Vector2D {
	// Get both the potential minimum translation vectors
	mtvForB := satSinglePolygon(polyA, polyB)
	mtvForA := satSinglePolygon(polyB, polyA)

	// Actually figure out which one to return
	// looks like there is a separating axis, hence: no collision
	if mtvForB == neonMath.ZeroVec2D || mtvForA == neonMath.ZeroVec2D {
		return neonMath.ZeroVec2D
	} else {
		if mtvForB.Length() <= mtvForA.Length() {
			return mtvForB
		} else {
			return mtvForA.Scale(-1.0)
		}
	}
}
