package problems

import (
	neural "github.com/poseidon4o/go-neural/src/neural"
	util "github.com/poseidon4o/go-neural/src/util"
	"math"
)

type Bird struct {
	Pos     util.Vector
	Vel     util.Vector
	NextPos util.Vector
}

const pylonSpacing int = 150
const PylonHole int = 150
const G_CONST float64 = 9.8 * 100
const SCROLL_SPEED float64 = 75

var G_FORCE util.Vector = util.Vector{
	X: 0,
	Y: G_CONST,
}

type Level struct {
	size   util.Vector
	pylons []util.Vector
	birds  []*Bird
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
			Pos:     *l.NewBirdPos(),
			Vel:     *util.NewVector(SCROLL_SPEED, 0),
			NextPos: *util.NewVector(0, 0),
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

func (l *Level) GetBirds() *[]*Bird {
	return &l.birds
}

func (l *Level) GetPylons() []util.Vector {
	return l.pylons
}

func (l *Level) GetSize() util.Vector {
	return l.size
}

func (l *Level) Step(dt float64) {
	for c := range l.birds {
		// position += timestep * (velocity + timestep * acceleration / 2);
		// TODO not use go
		l.birds[c].NextPos = *l.birds[c].Pos.Add(G_FORCE.Scale(dt / 2).Add(&l.birds[c].Vel).Scale(dt))

		// velocity += timestep * acceleration;
		l.birds[c].Vel = *l.birds[c].Vel.Add(G_FORCE.Scale(dt))
	}
}
