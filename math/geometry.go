package math




// Returns the intersection of two vectorially defined lines:
// l = la + mu * (lb - la)
// TODO: Update this method to return a zero vec if the intersection is outside of a specific range, also find a better method than determinants
// LineIntervalIntersection returns the intersection between a line an an interval, if there is none then it returns a zero vector
func LineIntervalIntersection(interval, line [2]Vector2D) Vector2D {
	x1, x2, x3, x4 := interval[0].X, interval[1].X, line[0].X, line[1].X
	y1, y2, y3, y4 := interval[0].Y, interval[1].Y, line[0].Y, line[1].Y

	// This is gonna be really ugly
	intersection := Vector2D{
		X: (((x1 * y2 - y1 * x2) * (x3 - x4)) - ((x1 - x2) * (x3 * y4 - y3 * x4))) /
			(((x1 - x2) * (y3 - y4)) - ((y1 - y2) * (x3 - x4))),

		Y: (((x1 * y2 - y1 * x2) * (y3 - y4)) - ((y1 - y2) * (x3 * y4 - y3 * x4))) /
			(((x1 - x2) * (y3 - y4)) - ((y1 - y2) * (x3 - x4))),
	}
	if interval[0].Sub(intersection).Length() > interval[0].Sub(interval[1]).Length() {
		return ZeroVec2D
	}
	return intersection
}


// Given 3 points: A, B and C; ComputeOutwardsNormal computes the normal vector of (A - B) that points away from C
func ComputeOutwardsNormal(A, B, C Vector2D) Vector2D {
	normalVectorAttempt := A.Sub(B).Normal().Normalise()

	if A.Sub(C).Dot(normalVectorAttempt) < 0 {
		normalVectorAttempt = normalVectorAttempt.Scale(-1.0)
	}


	return normalVectorAttempt
}