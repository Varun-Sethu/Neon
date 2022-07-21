package neonMath

// Matrix2 is a 2x2 square matrix
type Matrix2 struct {
	M [2][2]float64
}

// Identity matrix
var Identity2 = Matrix2{
	M: [2][2]float64{{1.0, 0}, {0, 1.0}},
}

// Multiply two 2x2 matrices
func (A Matrix2) Multiply(B Matrix2) Matrix2 {
	a := A.M
	b := B.M

	return Matrix2{M: [2][2]float64{
		{a[0][0]*b[0][0] + a[0][1]*b[1][0], a[0][0]*b[0][1] + a[0][1]*b[1][1]},
		{a[0][1]*b[0][0] + a[1][1]*b[1][0], a[1][0]*b[0][1] + a[1][0]*b[1][1]},
	}}
}

// VectorMultiply multiplies a vector by a matrix
func (A Matrix2) VectorMultiply(v Vector2D) Vector2D {
	return Vector2D{
		X: A.M[0][0]*v.X + A.M[0][1]*v.Y,
		Y: A.M[1][0]*v.X + A.M[1][1]*v.Y,
	}
}

// Scale just scales an entire matrix by x
func (A Matrix2) Scale(x float64) Matrix2 {
	return Matrix2{M: [2][2]float64{
		{x * A.M[0][0], x * A.M[0][1]},
		{x * A.M[1][0], x * A.M[1][1]},
	}}
}

// Add two 2x2 Matrices
func (A Matrix2) Add(B Matrix2) Matrix2 {
	a := A.M
	b := B.M

	return Matrix2{M: [2][2]float64{
		{a[0][0] + b[0][0], a[0][1] + b[0][1]},
		{a[1][0] + b[1][0], a[1][1] + b[1][1]},
	}}
}

// Computes the Determinant of a 2x2 Matrix
func (A Matrix2) Determinant() float64 {
	return A.M[0][0]*A.M[1][1] - A.M[0][1]*A.M[1][0]
}

// Computes the inverse of a 2x2 Matrix
func (A Matrix2) Inverse() Matrix2 {
	return Matrix2{M: [2][2]float64{
		{A.M[1][1], -A.M[0][1]},
		{-A.M[1][0], A.M[0][0]},
	}}.Scale(1.0 / A.Determinant())
}
