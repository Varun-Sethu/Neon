package math

import (
	"math"
)

// Returns the intersection of two vectorially defined lines:
// l = la + mu * (lb - la)
// TODO: Update this method to return a zero vec if the intersection is outside of a specific range, also find a better method than determinants
// LineIntervalIntersection returns the intersection between a line an an interval, if there is none then it returns a zero vector
func LineIntervalIntersection(interval, line [2]Vector2D) Vector2D {
	x1, x2, x3, x4 := interval[0].X, interval[1].X, line[0].X, line[1].X
	y1, y2, y3, y4 := interval[0].Y, interval[1].Y, line[0].Y, line[1].Y

	// This is gonna be really ugly
	intersection := Vector2D{
		X: (((x1*y2 - y1*x2) * (x3 - x4)) - ((x1 - x2) * (x3*y4 - y3*x4))) /
			(((x1 - x2) * (y3 - y4)) - ((y1 - y2) * (x3 - x4))),

		Y: (((x1*y2 - y1*x2) * (y3 - y4)) - ((y1 - y2) * (x3*y4 - y3*x4))) /
			(((x1 - x2) * (y3 - y4)) - ((y1 - y2) * (x3 - x4))),
	}
	length := interval[0].Sub(intersection).Length()
	max_length := math.Max(interval[1].Sub(intersection).Length(), length)

	if max_length > interval[0].Sub(interval[1]).Length() || math.IsNaN(length) {
		return ZeroVec2D
	}
	return intersection
}

// Given 3 points: A, B and C; ComputeOutwardsNormal computes the normal vector of (A - B) that points away from C
func ComputeOutwardsNormal(A, B, C Vector2D) Vector2D {
	normalVectorAttempt := A.Sub(B).Normal().Normalise()

	if A.Sub(C).Dot(normalVectorAttempt) < 0.0001 {
		normalVectorAttempt = normalVectorAttempt.Scale(-1.0)
	}

	return normalVectorAttempt
}

// LiesBehindLine produces all points that lie behind a certain line given an oriented normal vector
func LiesBehindLine(points []Vector2D, line [2]Vector2D, axis Vector2D) ([]Vector2D, []float64) {
	axis = axis.Normalise()
	valid := []Vector2D{}
	depths := []float64{}

	for _, p := range points {
		depth := p.Sub(line[0]).Dot(axis)
		if depth <= 0 {
			valid = append(valid, p)
			depths = append(depths, math.Abs(depth))
		}
	}

	return valid, depths
}

// IntervalRegionIntersection returns the intersection between a line and an intersection
func IntervalRegionIntersection(interval, region_boundary [2]Vector2D, region_orientation Vector2D) [2]Vector2D {
	// determine the intersection points between the interval and the region_boundary
	region_interval_intersection := LineIntervalIntersection(interval, region_boundary)
	output := [2]Vector2D{}

	// if the interval intersects the region_boundary then we must replace one end of the interval with the intersection
	if region_interval_intersection != ZeroVec2D {
		new_interval, _ := LiesBehindLine(interval[:], region_boundary, region_orientation.Scale(-1.0))
		copy(output[:], append(new_interval, region_interval_intersection))
		return output
	}

	// otherwise just return the interval IFF it is valid
	if e, _ := LiesBehindLine(interval[:], region_boundary, region_orientation.Scale(-1.0)); len(e) == 2 {
		return interval
	} else {
		return [2]Vector2D{ZeroVec2D, ZeroVec2D}
	}
}

// Reasonably simple method, just projects a point onto a line, primarily used within the collision solver
func ProjectPointOntoLine(point Vector2D, line [2]Vector2D) Vector2D {
	projection := point.Project(line[1].Sub(line[0]))
	return line[0].Add(projection)
}
