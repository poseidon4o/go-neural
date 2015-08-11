package util

type Vector struct {
	X, Y float64
}

func (v Vector) Add(o *Vector) *Vector {
	v.X += o.X
	v.Y += o.Y
	return &v
}

func (v Vector) Neg() *Vector {
	v.X = -v.X
	v.Y = -v.Y
	return &v
}

func (v Vector) Mul(o *Vector) *Vector {
	v.X *= o.X
	v.Y *= o.Y
	return &v
}

func (v Vector) Scale(scalar float64) *Vector {
	v.X *= scalar
	v.Y *= scalar
	return &v
}

func NewVector(x, y float64) *Vector {
	return &Vector{
		X: x,
		Y: y,
	}
}
