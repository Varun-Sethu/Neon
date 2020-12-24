package engine

import (
	"fmt"
	"neon/entities"
	"neon/math"
)


// TODO: https://www.toptal.com/game/video-game-physics-part-ii-collision-detection-for-solid-objects

// This is completely responsible for determining if two objects collide and computing a collision manifold for them
type ContactManifold struct {
	IncidentFrame	entities.Polygon
	ReferenceFrame	entities.Polygon


	MTV				math.Vector2D
	ContactCount	int
	CollisionPoints	[]math.Vector2D
}









// ComputeContactManifold computes a contact manifold for two polygon meshes
func ComputeContactManifold(poly_a, poly_b entities.Polygon) ContactManifold {
	mtv := entities.SAT(poly_a, poly_b) // note that the MTV always points from A to B

	// If there is no collision then the MTV is the zero vector which we need to account for
	if mtv.Length() == 0 {
		return ContactManifold{
			ContactCount:  0,
		}
	}
	collision_normal := mtv.Normalise()


	edge_candidate_a, perp_a := entities.DetermineSupportingEdge(poly_a, collision_normal) // TODO: Fix this issue: get supporting edge doesn't account for ones in the "negative" direction
	edge_candidate_b, perb_b := entities.DetermineSupportingEdge(poly_b, collision_normal.Scale(-1.0))

	reference_polygon, reference_face := poly_a, edge_candidate_a
	incident_polygon,  incident_face  := poly_b, edge_candidate_b
	if perb_b >= perp_a {
		reference_polygon, reference_face, incident_polygon, incident_face = incident_polygon, incident_face, reference_polygon, reference_face
	}



	// Now we need to actually perform the clipping of our reference_polygon onto our incident_polygon
	// In essence this involves clipping the incident polygon against all edges of the reference polygon EXCEPT the reference_face
	// After that is done, we simply delete all points of the clipped polygon that are not "behind" the reference face, this is all implemented in the polygon_clip method
	contact_points := polygon_clip(incident_polygon, reference_polygon, incident_face, reference_face) // note that contact points are in "world-coordinates" aka. Relative to the world origin
	// TODO: Fix contact point generation

	return ContactManifold{
		IncidentFrame:  incident_polygon,
		ReferenceFrame: reference_polygon,

		MTV: mtv,
		ContactCount: len(contact_points),
		CollisionPoints: contact_points,
	}
}




// determine_required_clipping_edges returns the set of all edges that the incident polygon has to be clipped against
// relatively simple method
func determine_required_clipping_edges(reference_poly entities.Polygon, reference_face []int) [][]int {
	clipping_set := [][]int{}


	for _, vertex := range reference_face {
		for _, connected_vertex := range reference_poly.Edges[vertex] {
			candidate_edge := []int{vertex, connected_vertex}

			if !equalSet(candidate_edge, reference_face) {
				clipping_set = append(clipping_set, candidate_edge)
			}
		}
	}
	return clipping_set
}







// Performs the clipping required for manifold computation
func polygon_clip(incident_poly, reference_poly entities.Polygon, incident_face, reference_face []int) []math.Vector2D {

	reference_poly = entities.MapToWorldSpace(reference_poly)
	incident_poly  = entities.MapToWorldSpace(incident_poly)


	fmt.Printf("Incident Polygon: %v\nReference Polygon: %v\nReference Face: %v\n\n", incident_poly, reference_poly, reference_face)


	// get the actual "values" for the incident and reference face
	//incident_face_edge := [2]math.Vector2D {
	//	incident_poly.Vertices[incident_face[0]],
	//	incident_poly.Vertices[incident_face[1]],
	//}
	reference_face_edge := [2]math.Vector2D {
		reference_poly.Vertices[reference_face[0]],
		reference_poly.Vertices[reference_face[1]],
	}

	// for simplicity we can compute a "set" of edges that the incident polygon need to be clipped against
	//required_clipping := determine_required_clipping_edges(reference_poly, reference_face)

	// iterate over all the edges that require clipping
	/**for _, clipping_edge := range required_clipping {
		line := [2]math.Vector2D {
			reference_poly.Vertices[clipping_edge[0]],
			reference_poly.Vertices[clipping_edge[1]],
		}
		orientation_normal := math.ComputeOutwardsNormal(line[0], line[1], reference_poly.State.CentroidPosition).Scale(-1.0)

		incident_face_edge = math.IntervalRegionIntersection(incident_face_edge, line, orientation_normal)
	}**/

	// finally the "manifold" is simply points that have actually penetrated the reference_poly, hence they lie below the reference face
	return reference_face_edge[:] //math.LiesBehindLine(incident_face_edge[:], reference_face_edge, math.ComputeOutwardsNormal(
		//reference_face_edge[0], reference_face_edge[1], reference_poly.State.CentroidPosition).Scale(-1.0))
}









func DetermineCollision(polyA entities.Polygon, polyB entities.Polygon) (bool, ContactManifold) {
	contact_manifold := ComputeContactManifold(polyA, polyB)
	return contact_manifold.ContactCount != 0, contact_manifold
}