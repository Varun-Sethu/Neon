package entities

import (
	"fmt"
	"neon/math"
	"testing"
)


// TestClipping tests polygon clipping against a line
func TestClipping(t *testing.T) {
	p := NewPolygon([]math.Vector2D{
		{1, 2}, {2, 2}, {3, 3}, {3, 0}, {2, 1}, {1,1},
	})
	line := [2]math.Vector2D{
		{1.5, 3},
		{1.5, 0},
	}
	normal := math.Vector2D{X: 1}
	p = Clip(p, line, normal)

	for i, v := range p.Vertices {
		fmt.Printf("Vertex %d: %v\n", i, v.Add(p.State.CentroidPosition))
	}
	fmt.Printf("\nEdges: %v", p.Edges)
}