package math

import "math"

// Simple package that describes a 2D-3D vector

// 2D vector
type Vector2D struct {
	X, Y float64
}

// 3D Vector
type Vector3D struct {
	X, Y, Z float64
}

var ZeroVec2D = Vector2D{0, 0}
var ZeroVec3D = Vector3D{0, 0, 0}
var BigVec2D = Vector2D{math.Inf(1), math.Inf(1)}

// Conversion functions
func (v Vector2D) ConvertTo3D() Vector3D {
	return Vector3D{X: v.X, Y: v.Y, Z: 0}
}
func (v Vector3D) ConvertTo2D() Vector2D {
	return Vector2D{X: v.X, Y: v.Y}
}

// Length functions
func (v Vector2D) Length() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}
func (v Vector3D) Length() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2) + math.Pow(v.Z, 2))
}

// Scale functions
func (v Vector2D) Scale(x float64) Vector2D {
	v.X *= x
	v.Y *= x
	return v
}
func (v Vector3D) Scale(x float64) Vector3D {
	v.X *= x
	v.Y *= x
	v.Z *= x
	return v
}

// Addition functions
func (v Vector2D) Add(k Vector2D) Vector2D {
	v.X += k.X
	v.Y += k.Y
	return v
}
func (v Vector3D) Add(k Vector3D) Vector3D {
	v.X += k.X
	v.Y += k.Y
	v.Z += k.Z
	return v
}

// Subtraction functions
func (v Vector2D) Sub(k Vector2D) Vector2D {
	return v.Add(k.Scale(-1.0))
}
func (v Vector3D) Sub(k Vector3D) Vector3D {
	return v.Add(k.Scale(-1.0))
}

// General utility functions

// Normalise normalises a vector
func (v Vector2D) Normalise() Vector2D {
	return v.Scale(1.0 / v.Length())
}
func (v Vector3D) Normalise() Vector3D {
	return v.Scale(1.0 / v.Length())
}

// Normal determines the normal of a vector
func (v Vector2D) Normal() Vector2D {
	return Vector2D{
		X: -v.Y,
		Y: v.X,
	}
}

// Dot computes the dot product of two vectors
func (v Vector2D) Dot(k Vector2D) float64 {
	return (v.X * k.X) + (v.Y * k.Y)
}
func (v Vector3D) Dot(k Vector3D) float64 {
	return v.X*k.X + v.Y*k.Y + v.Z*k.Z
}

// Cross product only outputs a 3D vector
func (v Vector2D) Cross(k Vector2D) Vector3D {
	return Vector3D{
		X: 0,
		Y: 0,
		Z: v.X*k.Y - v.Y*k.X,
	}
}
func (v Vector3D) Cross(k Vector3D) Vector3D {
	return Vector3D{
		X: v.Y*k.Z - v.Z*k.Y,
		Y: -(v.X*k.Z - v.Z*k.X),
		Z: v.X*k.Y - v.Y*k.X,
	}
}

// CrossMag computes the cross product but only the magnitude
func (v Vector2D) CrossMag(k Vector2D) float64 {
	return v.X*k.Y - v.Y*k.X
}

// Extension of cross product to the special case of a purely 3d vector cross a 2d one,  purely 3d vec points only in the z direction
func (v Vector2D) CrossWithUpwards(upwards float64) Vector2D {
	return Vector2D{
		X: upwards * v.Y,
		Y: -upwards * v.X,
	}
}

// Extension of cross product to the special case of a purely 3d vector cross a 2d one, purely 3d vec points only in the z direction
func (v Vector2D) CrossUpwardsWithVec(upwards float64) Vector2D {
	return Vector2D{
		X: -upwards * v.Y,
		Y: upwards * v.X,
	}
}

// Projection functions project vector v onto vector k
func (v Vector2D) Project(k Vector2D) Vector2D {
	return k.Scale((v.Dot(k)) /
		k.Dot(k))
}

// Scalar projection of two vectors
func (v Vector2D) ScalarProject(k Vector2D) float64 {
	return (v.Dot(k)) /
		(k.Dot(k))
}

// Min determines the minimum of two vectors
func Min(a Vector2D, b Vector2D) Vector2D {
	min := a
	if b.Length() < a.Length() {
		min = b
	}
	return min
}

// Max determines the maximum of two vectors
func Max(a Vector2D, b Vector2D) Vector2D {
	max := a
	if b.Length() > a.Length() {
		max = b
	}
	return max
}
