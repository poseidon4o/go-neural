package problems

import (
	neural "../neural"
)

type Vector struct {
	x, y float64
}

const ZeroVector Vector = &Vector{
	x: 0.0,
	y: 0.0,
}

func (v *Vector) Add(o *Vector) *Vector {
	v.x += o.x
	v.y += o.y
	return v
}

func (v *Vector) Neg() *Vector {
	v.x = -v.x
	v.y = -v.y
	return v
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

}
