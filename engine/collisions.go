package engine

import (
	"neon/entities"
	"neon/math"
)



// This is completely responsible for determining if two objects collide and computing a collision manifold for them
type ContactManifold struct {
	IncidentFrame	*entities.Polygon
	ReferenceFrame	*entities.Polygon

	IncidentFace	[]int
	ReferenceFace	[]int


	MTV				math.Vector2D
	ContactCount	int
	CollisionPoints	[]math.Vector2D
	ContactDepths	[]float64
}









// ComputeContactManifold computes a contact manifold for two polygon meshes
func ComputeContactManifold(poly_a, poly_b *entities.Polygon) ContactManifold {
	mtv := entities.SAT(*poly_a, *poly_b) // note that the MTV always points from A to B

	// If there is no collision then the MTV is the zero vector which we need to account for
	if mtv.Length() == 0 {
		return ContactManifold{
			ContactCount:  0,
		}
	}
	collisionNormal := mtv.Normalise()


	edgeCandidateA, perpA := poly_a.DetermineSupportingEdge(collisionNormal)
	edgeCandidateB, perpB := poly_b.DetermineSupportingEdge(collisionNormal.Scale(-1.0))

	referencePolygon, referenceFace := poly_a, edgeCandidateA
	incidentPolygon, incidentFace := poly_b, edgeCandidateB
	if perpB >= perpA {
		referencePolygon, referenceFace, incidentPolygon, incidentFace = poly_b, edgeCandidateB, poly_a, edgeCandidateA
		mtv = mtv.Scale(-1.0)
	}



	// Now we need to actually perform the clipping of our reference_polygon onto our incident_polygon
	// After that is done, we simply delete all points of the clipped polygon that are not "behind" the reference face, this is all implemented in the polygon_clip method
	contactPoints, pointDepths := polygonClip(*incidentPolygon, *referencePolygon, incidentFace, referenceFace)

	return ContactManifold{
		IncidentFrame:  incidentPolygon,
		ReferenceFrame: referencePolygon,

		IncidentFace:  incidentFace,
		ReferenceFace: referenceFace,

		MTV:             mtv,
		ContactCount:    len(contactPoints),
		CollisionPoints: contactPoints,
		ContactDepths:   pointDepths,
	}
}




// determineRequiredClippingEdges returns the set of all edges that the incident polygon has to be clipped against
// relatively simple method
func determineRequiredClippingEdges(referencePoly entities.Polygon, referenceFace []int) [][]int {
	var clippingSet [][]int


	for _, vertex := range referenceFace {
		for _, connectedVertex := range referencePoly.Edges[vertex] {
			candidateEdge := []int{vertex, connectedVertex}

			if !equalSet(candidateEdge, referenceFace) {
				clippingSet = append(clippingSet, candidateEdge)
			}
		}
	}
	return clippingSet
}







// Performs the clipping required for manifold computation
func polygonClip(incidentPoly, referencePoly entities.Polygon, incidentFace, referenceFace []int) ([]math.Vector2D, []float64) {

	// get the actual "values" for the incident and reference face
	incidentFaceEdge := incidentPoly.GetEdgeCoordinates(incidentFace)
	referenceFaceEdge := referencePoly.GetEdgeCoordinates(referenceFace)

	// for simplicity we can compute a "set" of edges that the incident polygon need to be clipped against
	requiredClipping := determineRequiredClippingEdges(referencePoly, referenceFace)

	// iterate over all the edges that require clipping
	for _, clippingEdge := range requiredClipping {
		line := referencePoly.GetEdgeCoordinates(clippingEdge)
		orientationNormal := math.ComputeOutwardsNormal(line[0], line[1], referencePoly.State.CentroidPosition).Scale(-1.0)
		incidentFaceEdge = math.IntervalRegionIntersection(incidentFaceEdge, line, orientationNormal)
	}


	// finally the "manifold" is simply points that have actually penetrated the reference_poly, hence they lie below the reference face
	return math.LiesBehindLine(incidentFaceEdge[:], referenceFaceEdge, math.ComputeOutwardsNormal(
		referenceFaceEdge[0], referenceFaceEdge[1], referencePoly.State.CentroidPosition))
}









func DetermineCollision(polyA *entities.Polygon, polyB *entities.Polygon) (bool, ContactManifold) {
	contactManifold := ComputeContactManifold(polyA, polyB)
	return contactManifold.ContactCount != 0, contactManifold
}