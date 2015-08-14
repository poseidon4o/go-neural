package flappy

import (
	neural "github.com/poseidon4o/go-neural/src/neural"
	util "github.com/poseidon4o/go-neural/src/util"
	"math"
)

const pylonSpacing int = 150
const PylonHole int = 150
const G_CONST float64 = 9.8 * 100

const SCROLL_SPEED float64 = 75

var X_ACCELERATION util.Vector = util.Vector{
	X: SCROLL_SPEED,
	Y: 0,
}

var G_FORCE util.Vector = util.Vector{
	X: 0,
	Y: G_CONST,
}

var JUMP_FORCE util.Vector = util.Vector{
	X: 0,
	Y: -500,
}

type Level struct {
	size   util.Vector
	pylons []util.Vector
	birds  []*Bird
}

type Bird struct {
	pos     util.Vector
	vel     util.Vector
	nextPos util.Vector
}

func (b *Bird) Jump() {
	b.vel = *b.vel.Add(&JUMP_FORCE)
	b.vel.Y = math.Max(b.vel.Y, JUMP_FORCE.Y)
}

func (b *Bird) Land() {
}

func (b *Bird) Move(dir int) {
}

func NewLevel(w, h int) *Level {
	lvl := &Level{
		size:   *util.NewVector(float64(w), float64(h)),
		pylons: make([]util.Vector, 0),
		birds:  make([]*Bird, 0),
	}

	// min offset from top && bottom
	yOffset := float64(PylonHole)

	for off := pylonSpacing; off < w; off += pylonSpacing {
		hole := neural.RandMax(float64(h)-yOffset*2.0) + yOffset

		lvl.pylons = append(lvl.pylons, *util.NewVector(float64(off), hole))
	}
	return lvl
}

func (l *Level) NewBirdPos() *util.Vector {
	return util.NewVector(1, l.size.Y/2)
}

func (l *Level) AddBirds(count int) {
	for c := 0; c < count; c++ {

		l.birds = append(l.birds, &Bird{
			pos:     *l.NewBirdPos(),
			vel:     *util.NewVector(0, 0).Add(&X_ACCELERATION),
			nextPos: *util.NewVector(0, 0),
		})
	}
}

func (l *Level) FirstPylonAfterIdx(pos *util.Vector) int {
	// TODO not use GO
	start := int(math.Max(float64(int(pos.X/float64(pylonSpacing))-1), 0))

	for ; start < len(l.pylons); start++ {
		if l.pylons[start].X > pos.X {
			return start
		}
	}
	return -1
}

func (l *Level) ClosestPylon(pos *util.Vector) util.Vector {
	idx := l.FirstPylonAfterIdx(pos)

	if idx == -1 {
		return *util.NewVector(0, 0)
	}

	nextX := l.pylons[idx].X - pos.X
	// TODO srsly?
	//prevX := idx > 0 ? pos.X - l.pylons[idx - 1] : pylonSpacing * 2

	var prevX float64 = 0
	if idx > 0 {
		prevX = pos.X - l.pylons[idx-1].X
	} else {
		// will be bigger than any spacing to next pylon
		prevX = float64(pylonSpacing) * 2
	}

	if prevX < nextX {
		return l.pylons[idx-1]
	} else {
		return l.pylons[idx]
	}
}

func (l *Level) FirstPylonAfter(pos *util.Vector) util.Vector {
	idx := l.FirstPylonAfterIdx(pos)
	if idx >= 0 {
		return l.pylons[idx]
	}
	return *util.NewVector(0, 0)

	// TODO srsly?
	// idx >= 0 ? l.pylons[idx] : *util.NewVector(0, 0)
}

func (l *Level) Step(dt float64) {
	for c := range l.birds {
		// position += timestep * (velocity + timestep * acceleration / 2);
		// TODO not use go
		l.birds[c].nextPos = *l.birds[c].pos.Add(G_FORCE.Scale(dt / 2).Add(&l.birds[c].vel).Scale(dt))

		// velocity += timestep * acceleration;
		l.birds[c].vel = *l.birds[c].vel.Add(G_FORCE.Scale(dt))
	}
}
