package entities

import (
	gmath "math"
	"neon/math"
	"sort"
)




/*
	This entire file handles elementary geometry for polygons, algorithms such as SAT and Polygon Clipping are implemented here
	Additionally, manifold generation is also implemented here
 */



// GetSupportingPoint returns the supporting point of a polygon along a specific axis
func (polygon *Polygon) GetSupportingPoint(axis math.Vector2D) (math.Vector2D, int) {
	currentMaxProj := gmath.Inf(-1)
	bestVertex	   := math.ZeroVec2D
	vertexID	   := 0

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
func (poly *Polygon) DetermineSupportingEdge(axis math.Vector2D) ([]int, float64) {

	_, vertexID := poly.GetSupportingPoint(axis)
	v			:= poly.Vertices[vertexID]
	A			:= poly.Vertices[poly.Edges[vertexID][0]]; normalA := math.ComputeOutwardsNormal(A, v, poly.State.CentroidPosition)
	B			:= poly.Vertices[poly.Edges[vertexID][1]]; normalB := math.ComputeOutwardsNormal(B, v, poly.State.CentroidPosition)

	// There are two possible other vertices that can connect to this one, we determine which one is significant by comparing dot products
	if gmath.Abs(normalA.Dot(axis)) >= gmath.Abs(normalB.Dot(axis)) {
		return []int{vertexID, poly.Edges[vertexID][0]}, gmath.Abs(normalA.Dot(axis))
	}
	return []int{vertexID, poly.Edges[vertexID][1]}, gmath.Abs(normalB.Dot(axis))
}




// PolyVerticesOutside takes a "line" (2 points) and a normal and returns all points of the polygon that lie outside this line, in the direction anti-parallel to the normal
// Note: the function assumes that the vertices are in the "world frame"
func (poly *Polygon) PolyVerticesOutside(line [2]math.Vector2D, normal math.Vector2D) []int {
	var outside []int

	for i, v := range poly.Vertices {
		// Just ensure that you determine the world position of the vertex
		if v.Sub(line[0]).Dot(normal) < 0 {
			outside = append(outside, i)
		}
	}

	return outside
}



// AxisProjection returns the projection interval of a polygon onto an axis
// Note that the projections are in "world coordinates"
func (polygon *Polygon) AxisProjection(axis math.Vector2D) []float64 {
	// We need to first determine the supporting points along the axis perpendicular to the provided axis
	axis = axis.Normalise()
	sup_left, _  := polygon.GetSupportingPoint(axis)
	sup_right, _ := polygon.GetSupportingPoint(axis.Scale(-1.0))

	interval := []float64{
		sup_left.ScalarProject(axis),
		sup_right.ScalarProject(axis),
	}
	sort.Float64s(interval)
	return interval
}


// projections_overlap determines if two projections onto an axis overlap as well as the degree of overlap
func projections_overlap(projection_a, projection_b []float64) (bool, float64) {
	// evidently there is no overlap
	if projection_a[1] < projection_b[0] || projection_b[1] < projection_a[0] {
		return false, 0.0
	}

	return true, gmath.Min(projection_a[1], projection_b[1]) - gmath.Max(projection_a[0], projection_b[0])
}








// satSinglePolygon just checks if polyB intersects polyA
func satSinglePolygon(polyA Polygon, polyB Polygon) math.Vector2D {
	mtv := math.BigVec2D

	for vertex, edges := range polyA.Edges {
		for _, edge := range edges {
			worldVertex := polyA.Vertices[vertex].Add(polyA.State.CentroidPosition) // worldVertex refers to the world coordinates of the vertex
			worldEdgeV 	:= polyA.Vertices[edge].Add(polyA.State.CentroidPosition) // worldEdgeV is the world coordinates of the other vertex that defines this edge
			normal		:= math.ComputeOutwardsNormal(worldEdgeV, worldVertex, polyA.State.CentroidPosition) // normal is just the normal vector associated with this edge

			projected_axis_polyB := polyB.AxisProjection(normal)
			projected_axis_polyA := polyA.AxisProjection(normal)

			if intersects, overlap := projections_overlap(projected_axis_polyA, projected_axis_polyB); !intersects {
				return math.ZeroVec2D
			} else {
				mtv = math.Min(mtv, normal.Scale(overlap))
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
func (p *Polygon) MapToWorldSpace()  {
	for vertex_id, _ := range p.Vertices {
		p.Vertices[vertex_id] = p.Vertices[vertex_id].Add(p.State.CentroidPosition)
	}
}
// If a polygon is in world space eg. all the coordinates are relative to the global origin then this just maps them all out of it
func (p *Polygon) MapOutofWorldSpace() {
	for vertex_id, _ := range p.Vertices {
		p.Vertices[vertex_id] = p.Vertices[vertex_id].Sub(p.State.CentroidPosition)
	}
}







// SAT determines if two polygons are intersecting and computes the corresponding MTV
// Note that the MTV ALWAYS POINTS FROM A TO B
func SAT(polyA Polygon, polyB Polygon) math.Vector2D {
	// Get both the potential minimum translation vectors
	mtv_for_B := satSinglePolygon(polyA, polyB)
	mtv_for_A := satSinglePolygon(polyB, polyA)


	// Actually figure out which one to return
	// looks like there is a separating axis, hence: no collision
	if mtv_for_B == math.ZeroVec2D || mtv_for_A == math.ZeroVec2D {
		return math.ZeroVec2D
	} else {
		if mtv_for_B.Length() <= mtv_for_A.Length() {
			return mtv_for_B
		} else {
			return mtv_for_A.Scale(-1.0)}
	}
}




