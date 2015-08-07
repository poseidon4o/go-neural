package problems

import (
	neural "../neural"
)

type Vector struct {
	x, y float64
}

func (v Vector) Add(o *Vector) *Vector {
	v.x += o.x
	v.y += o.y
	return &v
}

func (v Vector) Neg() *Vector {
	v.x = -v.x
	v.y = -v.y
	return &v
}

func (v Vector) Mul(o *Vector) *Vector {
	v.x *= o.x
	v.y *= o.y
	return &v
}

func (v Vector) Scale(scalar float64) *Vector {
	v.x *= scalar
	v.y *= scalar
	return &v
}

func NewVector(x, y float64) *Vector {
	return &Vector{
		x: x,
		y: y,
	}
}

type Bird struct {
	pos, vel Vector
}

const pylonSpacing int = 50
const G_CONST float64 = 9.8
const G_FORCE Vector = Vector{
	x: 0,
	y: G_CONST,
}

type Level struct {
	size   Vector
	pylons []Vector
	birds  []Bird
}

func NewLevel(w, h int) *Level {
	lvl := &Level{
		size:   *NewVector(float64(w), float64(h)),
		pylons: make([]Vector, w/pylonSpacing),
		birds:  make([]Bird),
	}

	for off := pylonSpacing; off < w; off += pylonSpacing {
		hole := int(neural.RandMax(h/5) + height/10)
		lvl.pylons = append(lvl.pylons, *NewVector(off, hole))
	}

	return lvl
}

func (l *Level) AddBirds(count int) {
	for c := 0; c < count; ++c {
		vpos := int(neural.RandMax(l.size.y / 5) + l.size.y / 10)

		l.birds = append(l.birds, Bird{
			pos: *NewVector(1, vpos),
			vel: *NewVector(0, 0),
		})
	}
}

func (l *Level) GetBirds() []Bird {
	return l.birds
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
		birds[c].pos = birds[c].pos.add(G_FORCE.scale(dt / 2).add(l.birds[c].vel).scale(dt))

		// velocity += timestep * acceleration;
		birds[c].vel = birds[c].vel.add(G_FORCE.scale(dt))
	}
}
