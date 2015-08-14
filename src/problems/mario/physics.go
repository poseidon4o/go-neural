package mario

import (
	neural "github.com/poseidon4o/go-neural/src/neural"
	util "github.com/poseidon4o/go-neural/src/util"
	"math"
)

const G_CONST float64 = 9.8 * 150

var G_FORCE util.Vector = util.Vector{
	X: 0,
	Y: G_CONST,
}

const BLOCK_SIZE int = 25

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
	if m.jumps >= 1 {
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
	blocks  [][]*util.Vector
	figures []*Figure
}

func (l *Level) makeHole(c *int) {
	*c += 3 // int(3 + neural.RandMax(2))
}

func (l *Level) makeObstacle(c *int) {
	*c--
}

func (l *Level) makeGround(c *int) {
	r := int(l.size.Y/float64(BLOCK_SIZE)) - 1

	l.blocks[*c][r] = util.NewVector(float64(*c*BLOCK_SIZE), float64(r*BLOCK_SIZE))
}

func NewLevel(w, h int) *Level {
	blockH := h / BLOCK_SIZE
	blockW := w / BLOCK_SIZE

	lvl := &Level{
		size:    *util.NewVector(float64(w), float64(h)),
		blocks:  make([][]*util.Vector, blockW, blockW),
		figures: make([]*Figure, 0),
	}

	for c := 0; c < blockW; c++ {
		lvl.blocks[c] = make([]*util.Vector, blockH, blockH)
		for r := 0; r < blockH; r++ {
			lvl.blocks[c][r] = nil
		}
	}

	for c, obs := 0, 1; c < blockW; c, obs = c+1, obs+1 {

		pr := c
		if obs%10 == 0 {
			if neural.Chance(0.5) {
				lvl.makeHole(&c)
			} else {
				lvl.makeObstacle(&c)
			}
		} else {
			lvl.makeGround(&c)
		}
		obs += c - pr
	}

	return lvl
}

func (l *Level) toLevelCoords(pos *util.Vector) (int, int) {
	return int(pos.X / float64(BLOCK_SIZE)), int(pos.Y / float64(BLOCK_SIZE))
}

func (l *Level) validCoord(w, h int) bool {
	return w >= 0 && w < (int(l.size.X)/BLOCK_SIZE) && h >= 0 && h < (int(l.size.Y)/BLOCK_SIZE)
}

func (l *Level) IsSolid(pos *util.Vector) bool {
	return l.CubeAt(pos) != nil
}

func (l *Level) CubeAt(pos *util.Vector) *util.Vector {
	w, h := l.toLevelCoords(pos)
	if l.validCoord(w, h) {
		return l.blocks[w][h]
	} else {
		return nil
	}
}

func (l *Level) FloorAt(pos *util.Vector) *util.Vector {
	wIdx, hIdx := l.toLevelCoords(pos)

	if !l.validCoord(wIdx, hIdx+1) {
		return nil
	}
	return l.blocks[wIdx][hIdx+1]
}

func (l *Level) NewFigurePos() *util.Vector {
	return util.NewVector(1, 1)
}

func (l *Level) AddFigures(count int) {
	for c := 0; c < count; c++ {

		l.figures = append(l.figures, &Figure{
			pos:     *l.NewFigurePos(),
			vel:     *util.NewVector(0, 0),
			nextPos: *util.NewVector(0, 0),
			jumps:   1,
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
		l.figures[c].vel.X *= (1 - 3*dt)
	}
}
