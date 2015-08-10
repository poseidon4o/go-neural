package problems

import (
	neural "../neural"
)

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

type Bird struct {
	pos, vel Vector
}

const pylonSpacing int = 150
const PylonHole int = 150
const G_CONST float64 = 0.0001

var G_FORCE Vector = Vector{
	X: 0,
	Y: G_CONST,
}

type Level struct {
	size   Vector
	pylons []Vector
	birds  []*Bird
}

func (b *Bird) Pos() *Vector {
	return &b.pos
}

func (b *Bird) Vel() *Vector {
	return &b.vel
}

func NewLevel(w, h int) *Level {
	lvl := &Level{
		size:   *NewVector(float64(w), float64(h)),
		pylons: make([]Vector, 0),
		birds:  make([]*Bird, 0),
	}

	for off := pylonSpacing * 2; off < w; off += pylonSpacing {
		// min offset from top && bottom
		yOffset := float64(PylonHole)/2.0 + float64(h)*0.8

		hole := neural.RandMax(float64(h)-yOffset*2.0) + yOffset

		lvl.pylons = append(lvl.pylons, *NewVector(float64(off), hole))
	}
	return lvl
}

func (l *Level) AddBirds(count int) {
	for c := 0; c < count; c++ {

		l.birds = append(l.birds, &Bird{
			pos: *NewVector(1, l.size.Y/2),
			vel: *NewVector(0.1, 0),
		})
	}
}

func (l *Level) GetBirds() *[]*Bird {
	return &l.birds
}

func (l *Level) GetPylons() []Vector {
	return l.pylons
}

func (l *Level) GetSize() Vector {
	return l.size
}

func (l *Level) Step(dt float64) {
	for c := range l.birds {
		// position += timestep * (velocity + timestep * acceleration / 2);
		l.birds[c].pos = *l.birds[c].pos.Add(G_FORCE.Scale(dt / 2).Add(&l.birds[c].vel).Scale(dt))

		// velocity += timestep * acceleration;
		l.birds[c].vel = *l.birds[c].vel.Add(G_FORCE.Scale(dt))
	}
}
