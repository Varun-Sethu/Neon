package entities

import (
	gmath "math"
	"neon/math"
)




/*
	This entire file handles elementary geometry for polygons, algorithms such as SAT and Polygon Clipping are implemented here
	Additionally, manifold generation is also implemented here
 */



// GetSupportingPoint returns the supporting point of a polygon along a specific axis
func GetSupportingPoint(polygon Polygon, axis math.Vector2D) (math.Vector2D, int) {
	currentMaxProj := gmath.Inf(-1)
	bestVertex	   := math.ZeroVec2D
	vertexID	   := 0

	for id, v := range polygon.Vertices {
		projection := v.Add(polygon.State.CentroidPosition).Dot(axis)
		if projection > currentMaxProj {
			currentMaxProj = projection
			bestVertex = v.Add(polygon.State.CentroidPosition)
			vertexID = id
		}
	}
	return bestVertex, vertexID
}



// DetermineSupportingEdge determines the edge furthest along a specific axis, in essence the specific edge, also returns how "parallel" the edge is with the normal
func DetermineSupportingEdge(poly Polygon, axis math.Vector2D) ([]int, float64) {

	_, vertexID := GetSupportingPoint(poly, axis)
	v			:= poly.Vertices[vertexID]
	A			:= poly.Vertices[poly.Edges[vertexID][0]]; normalA := math.ComputeOutwardsNormal(A, v, math.ZeroVec2D)
	B			:= poly.Vertices[poly.Edges[vertexID][1]]; normalB := math.ComputeOutwardsNormal(B, v, math.ZeroVec2D)

	// There are two possible other vertices that can connect to this one, we determine which one is significant by comparing dot products
	if gmath.Abs(normalA.Dot(axis)) < gmath.Abs(normalB.Dot(axis)) {
		return []int{vertexID, poly.Edges[vertexID][0]}, gmath.Abs(normalA.Dot(axis))
	}
	return []int{vertexID, poly.Edges[vertexID][1]}, gmath.Abs(normalB.Dot(axis))
}




// polyVerticesOutside takes a "line" (2 points) and a normal and returns all points of the polygon that lie outside this line, in the direction anti-parallel to the normal
func PolyVerticesOutside(polygon Polygon, line [2]math.Vector2D, normal math.Vector2D) []int {
	var outside []int

	for i, v := range polygon.Vertices {
		if v.Sub(line[0]).Dot(normal) < 0 {
			outside = append(outside, i)
		}
	}

	return outside
}





// Clip clips a polygon against a line
// Note that the the associated normal vector with the line dictates how we clip the polygon
func Clip(p Polygon, line [2]math.Vector2D, lineNormal math.Vector2D) Polygon {
	newPolygon := p

	// Determine all points that lie outside the line region
	outsidePoints := PolyVerticesOutside(p, line, lineNormal)
	if len(outsidePoints) == 0 {return  p}

	// Since a line can only intersect a polygon twice (if it intersects) we just have to find these intersections and connect them together, this gives us the clipped polygon against this line
	// [2]int refers to the edge that was "intersected" and int refers to the new internal id
	intersections := []int{}

	for _, point := range outsidePoints {
		// Check intersections from outgoing edges, if the intersection returned by the math function is a zero vec that implies there are no intersections
		for _, edge := range p.Edges[point] {
			intersectionEdgeA := math.LineIntervalIntersection(
				[2]math.Vector2D{
					p.Vertices[point].Add(p.State.CentroidPosition), p.Vertices[edge].Add(p.State.CentroidPosition),
				}, line)

			if intersectionEdgeA != math.ZeroVec2D {
				var id int
				newPolygon, id = newPolygon.addVertex(intersectionEdgeA)
				intersections = append(intersections, id)
				newPolygon.Edges[edge] = swap(newPolygon.Edges[edge], point, id)
				newPolygon.Edges[id] = append(newPolygon.Edges[id], edge)
			}
			newPolygon.deleteVertex(point)
		}
	}

	// Now we just need to connect up the two edges
	p.Edges[intersections[0]] = append(p.Edges[intersections[0]], intersections[1])
	p.Edges[intersections[1]] = append(p.Edges[intersections[1]], intersections[0])
	return newPolygon.refreshPolygon()
}




// satSinglePolygon just checks if polyA intersects polyB
func satSinglePolygon(polyA Polygon, polyB Polygon) math.Vector2D {
	mtv := math.BigVec2D

	for vertex, edges := range polyA.Edges {
		for _, edge := range edges {
			worldVertex := polyA.Vertices[vertex].Add(polyA.State.CentroidPosition) // worldVertex refers to the world coordinates of the vertex
			worldEdgeV 	:= polyA.Vertices[edge].Add(polyA.State.CentroidPosition) // worldEdgeV is the world coordinates of the other vertex that defines this edge
			normal		:= math.ComputeOutwardsNormal(worldEdgeV, worldVertex, polyA.State.CentroidPosition) // normal is just the normal vector associated with this edge

			supportingPoint, _ := GetSupportingPoint(polyB, normal.Scale(-1.0))
			separationVector   := normal.Scale(supportingPoint.Sub(worldVertex).Dot(normal))


			if separationVector.Dot(normal) > 0 {
				return math.ZeroVec2D
			} else {
				mtv = math.Min(mtv, separationVector)
			}
		}
	}
	return  mtv
}



// SAT determines if two polygons are intersecting and computes the corresponding MTV
func SAT(polyA Polygon, polyB Polygon) math.Vector2D {
	// Get both the potential minimum translation vectors
	mtvAToB := satSinglePolygon(polyA, polyB)
	mtvBtoA := satSinglePolygon(polyB, polyA)


	// Actually figure out which one to return
	if mtvBtoA == math.ZeroVec2D || mtvAToB == math.ZeroVec2D {
		return mtvAToB
	} else {return math.Min(mtvAToB, mtvBtoA)}
}




