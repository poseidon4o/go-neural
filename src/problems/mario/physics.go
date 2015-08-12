package mario

import (
	util "github.com/poseidon4o/go-neural/src/util"
	"math"
)

const G_CONST float64 = 9.8 * 150

var G_FORCE util.Vector = util.Vector{
	X: 0,
	Y: G_CONST,
}

const BLOCK_SIZE int = 24

var JUMP_FORCE util.Vector = util.Vector{
	X: 0,
	Y: -1000,
}

var X_ACCELERATION util.Vector = util.Vector{
	X: 20,
	Y: 0,
}

type Figure struct {
	pos     util.Vector
	vel     util.Vector
	nextPos util.Vector
	jumps   int
}

func (m *Figure) Jump() {
	if m.jumps >= 3 {
		return
	}
	m.jumps++
	m.vel = *m.vel.Add(&JUMP_FORCE)
	m.vel.Y = math.Max(m.vel.Y, JUMP_FORCE.Y)
}

func (f *Figure) Land() {
	f.jumps = 0
}

func (f *Figure) Move(dir int) {
	acc := X_ACCELERATION
	if dir < 0 {
		acc = *acc.Neg()
	}
	f.vel = *f.vel.Add(&acc)
}

type Level struct {
	size    util.Vector
	blocks  [][]util.Vector
	figures []*Figure
}

func NewLevel(w, h int) *Level {
	lvl := &Level{
		size:    *util.NewVector(float64(w), float64(h)),
		blocks:  make([][]util.Vector, 0),
		figures: make([]*Figure, 0),
	}

	// TODO generate level
	for c := 0; c < w; c += BLOCK_SIZE {
		col := make([]util.Vector, 1, 1)
		col[0] = *util.NewVector(float64(c), float64(h-BLOCK_SIZE))
		lvl.blocks = append(lvl.blocks, col)
	}

	return lvl
}

func (l *Level) FloorAt(pos *util.Vector) *util.Vector {
	idx := int(pos.X / float64(BLOCK_SIZE))

	if idx < 0 || idx >= len(l.blocks) {
		return util.NewVector(0, 0)
	}

	for c := len(l.blocks[idx]) - 1; c >= 0; c-- {
		if pos.Y < l.blocks[idx][c].Y {
			// we are above the block
			if c-1 >= 0 && pos.Y > l.blocks[idx][c-1].Y {
				// has next block and we are below it
				return &l.blocks[idx][c]
			} else if c-1 < 0 {
				// dont have next block
				return &l.blocks[idx][c]
			}
		}
	}

	return util.NewVector(0, 0)
}

func (l *Level) NewFigurePos() *util.Vector {
	return util.NewVector(1, l.FloorAt(util.NewVector(1, 1)).Y-float64(BLOCK_SIZE))
}

func (l *Level) AddFigures(count int) {
	for c := 0; c < count; c++ {

		l.figures = append(l.figures, &Figure{
			pos:     *l.NewFigurePos(),
			vel:     *util.NewVector(0, 0),
			nextPos: *util.NewVector(0, 0),
		})
	}
}

func (l *Level) Step(dt float64) {
	for c := range l.figures {
		// position += timestep * (velocity + timestep * acceleration / 2);
		// TODO not use go
		l.figures[c].nextPos = *l.figures[c].pos.Add(G_FORCE.Scale(dt / 2).Add(&l.figures[c].vel).Scale(dt))

		// velocity += timestep * acceleration;
		l.figures[c].vel = *l.figures[c].vel.Add(G_FORCE.Scale(dt))
		l.figures[c].vel.X *= 0.9
	}
}
