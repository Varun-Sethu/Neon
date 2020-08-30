package engine

import (
	"neon/entities"
	"neon/math"
)


// TODO: https://www.toptal.com/game/video-game-physics-part-ii-collision-detection-for-solid-objects

// This is completely responsible for determining if two objects collide and computing a collision manifold for them
type CollisionManifold struct {
	IncidentFrame	*entities.Polygon
	ReferenceFrame	*entities.Polygon


	MTV				math.Vector2D
	CollisionPoints	[2]math.Vector2D
}





// ComputeCollisionManifold calculates the collision manifold for two intersecting polygons given a specified MTV
func ComputeCollisionManifold(polyA, polyB *entities.Polygon, mtv math.Vector2D) CollisionManifold {
	// Determine the significant faces of the collision, also need to determine the incident and reference frame
	significantA, projA := entities.DetermineSupportingEdge(*polyA, mtv.Normalise())
	significantB, projB := entities.DetermineSupportingEdge(*polyB, mtv.Normalise())

	// Calculate the incident and reference faces/polygons
	reference, referenceFace, incident := polyA, significantA, polyB; if projA < projB {
		reference, referenceFace, incident = polyB, significantB, polyA
	}
	referencePoly := *reference
	incidentPoly  := *incident



	// iterate over all edges in the reference polygon and trim the incident polygon
	for _, edge := range referencePoly.Edges {
		if equalSet(edge, referenceFace) {continue}
		// Compute all information regarding the clip
		line := [2]math.Vector2D{
			referencePoly.Vertices[edge[0]].Add(referencePoly.State.CentroidPosition),
			referencePoly.Vertices[edge[1]].Add(referencePoly.State.CentroidPosition),
		}
		normal := math.ComputeOutwardsNormal(line[0], line[1], referencePoly.State.CentroidPosition).Scale(-1.0)

		incidentPoly = entities.Clip(incidentPoly, line, normal)
	}
	// for the reference face, the manifold is just all points that don't lie in the clipping region
	line := [2]math.Vector2D{
		referencePoly.Vertices[referenceFace[0]].Add(referencePoly.State.CentroidPosition),
		referencePoly.Vertices[referenceFace[1]].Add(referencePoly.State.CentroidPosition)}
	normal := math.ComputeOutwardsNormal(line[0], line[1], referencePoly.State.CentroidPosition).Scale(-1.0)
	internalManifold := entities.PolyVerticesOutside(referencePoly, line, normal)

	// Finally dump the manifold into everything
	return CollisionManifold{
		IncidentFrame: reference,
		ReferenceFrame: incident,
		MTV: mtv,
		CollisionPoints: [2]math.Vector2D{
			incidentPoly.Vertices[internalManifold[0]].Sub(incidentPoly.State.CentroidPosition),
			incidentPoly.Vertices[internalManifold[1]].Sub(incidentPoly.State.CentroidPosition),
		},
	}
}









func DetermineCollision(polyA entities.Polygon, polyB entities.Polygon) bool {
	mtv := entities.SAT(polyA, polyB)
	if mtv != math.ZeroVec2D {return true}
	return false
}