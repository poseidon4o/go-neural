package mario

import (
	util "github.com/poseidon4o/go-neural/src/util"
)

const G_CONST float64 = 9.8

var G_FORCE util.Vector = util.Vector{
	X: 0,
	Y: G_CONST,
}

const BLOCK_SIZE int = 24

var JUMP_FORCE util.Vector = util.Vector{
	X: 0,
	Y: -500,
}

var X_ACCELERATION util.Vector = util.Vector{
	X: 100,
	Y: 0,
}

type Figure struct {
	Pos     util.Vector
	Vel     util.Vector
	NextPos util.Vector
	jumps   int
}

func (m *Figure) Jump() {
	if m.jumps > 2 {
		return
	}
	m.jumps++
	m.Vel = *m.Vel.Add(&JUMP_FORCE)
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
	return &l.blocks[idx][0]
}

func (l *Level) NewFigurePos() *util.Vector {
	return util.NewVector(1, l.FloorAt(util.NewVector(1, 1)).Y-float64(BLOCK_SIZE))
}

func (l *Level) AddFigures(count int) {
	for c := 0; c < count; c++ {

		l.figures = append(l.figures, &Figure{
			Pos:     *l.NewFigurePos(),
			Vel:     *util.NewVector(0, 0),
			NextPos: *util.NewVector(0, 0),
		})
	}
}

func (l *Level) Step(dt float64) {
	for c := range l.figures {
		// position += timestep * (velocity + timestep * acceleration / 2);
		// TODO not use go
		l.figures[c].NextPos = *l.figures[c].Pos.Add(G_FORCE.Scale(dt / 2).Add(&l.figures[c].Vel).Scale(dt))

		// velocity += timestep * acceleration;
		l.figures[c].Vel = *l.figures[c].Vel.Add(G_FORCE.Scale(dt))
	}
}
