package engine

import (
	"neon/entities"
	"neon/math"
	"testing"
)



func TestManifoldGeneration(t *testing.T) {
	poly_a := entities.NewPolygon([]math.Vector2D{{0, 100}, {100, 100}, {100, 0}, {0,0}})
	poly_b := entities.NewPolygon([]math.Vector2D{{50, 101}, {120, 120}, {150, 70}, {110,50}})


	mtv, _ := entities.SAT(poly_a, poly_b)
	t.Log(mtv)
	manifold := ComputeContactManifold(poly_a, poly_b)

	t.Logf("Contact Points: %v", manifold.CollisionPoints)
}




func TestClipping(t *testing.T) {

}

