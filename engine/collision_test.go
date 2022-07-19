package engine

import (
	"Neon/engine/math"
	"Neon/entities"
	"testing"
)

func TestManifoldGeneration(t *testing.T) {
	polyA := entities.NewPolygon([]math.Vector2D{{X: 0, Y: 100}, {X: 100, Y: 100}, {X: 100, Y: 0}, {X: 0, Y: 0}})
	polyB := entities.NewPolygon([]math.Vector2D{{X: 50, Y: 101}, {X: 120, Y: 120}, {X: 150, Y: 70}, {X: 110, Y: 50}})

	mtv := entities.SAT(polyA, polyB)
	t.Log(mtv)
	manifold := ComputeContactManifold(&polyA, &polyB)

	t.Logf("Contact Points: %v", manifold.CollisionPoints)
}

func TestClipping(t *testing.T) {

}
