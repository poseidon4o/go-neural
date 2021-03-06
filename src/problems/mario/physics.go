package mario

import (
	"fmt"
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
	Y: -500,
}

var X_ACCELERATION util.Vector = util.Vector{
	X: 15,
	Y: 0,
}

const OBSTACLE_SPACING int = 10

type Figure struct {
	pos     util.Vector
	vel     util.Vector
	nextPos util.Vector
	jumps   int
}

func (f *Figure) Jump() {
	if f.jumps >= 1 {
		return
	}
	f.jumps++
	f.vel = *f.vel.Add(&JUMP_FORCE)
	f.vel.Y = math.Max(f.vel.Y, JUMP_FORCE.Y)
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

var BOOL_MAP_GRID_SIDE int = 9

type BoolMap struct {
	hi uint64
	lo uint64
}

func NewBMap() BoolMap {
	return BoolMap{
		hi: 0,
		lo: 0,
	}
}

func (b *BoolMap) At(idx int) bool {
	if idx > 63 {
		return (b.hi>>uint(idx))&1 == 1
	} else {
		return (b.lo>>uint(idx))&1 == 1
	}
}

func (b *BoolMap) Set(idx int) {
	if idx > 63 {
		b.hi |= 1 << uint(idx)
	} else {
		b.lo |= 1 << uint(idx)
	}
}

func (b *BoolMap) GridAt(x, y int) bool {
	idx := y + x*BOOL_MAP_GRID_SIDE + (int(BOOL_MAP_GRID_SIDE*BOOL_MAP_GRID_SIDE) / 2)
	return b.At(idx)
}

func (b *BoolMap) GridSet(x, y int) {
	idx := y + x*BOOL_MAP_GRID_SIDE + (int(BOOL_MAP_GRID_SIDE*BOOL_MAP_GRID_SIDE) / 2)
	b.Set(idx)
}

type Level struct {
	size    util.Vector
	blocks  [][]*util.Vector
	bmap    [][]BoolMap
	figures []*Figure
}

func (l *Level) makeHole(c *int) {
	size := int(3 + neural.RandMax(2))

	height := int(l.size.Y/float64(BLOCK_SIZE)) - 1 - 3

	skip := OBSTACLE_SPACING * 2

	if neural.Chance(0.75) {
		size = skip
		for iter := 0; iter < skip; iter++ {
			if iter+*c >= len(l.blocks) {
				break
			}
			if iter%5 == 0 {
				x := float64((iter + *c) * BLOCK_SIZE)
				y := float64((height + 3) * BLOCK_SIZE)
				l.blocks[iter+*c][height+3] = util.NewVector(x, y)
			}
		}
	}

	if neural.Chance(0.3) {
		for iter := -2; iter < size+3; iter++ {
			xIdx := iter + *c
			if xIdx < 0 || xIdx >= len(l.blocks) {
				continue
			}

			x := float64(xIdx * BLOCK_SIZE)
			y := float64(height * BLOCK_SIZE)

			l.blocks[xIdx][height] = util.NewVector(x, y)
		}
	}

	*c += size
}

func (l *Level) makeObstacle(c *int) {
	r := int(l.size.Y/float64(BLOCK_SIZE)) - 1
	for q := 0; q < 4; q++ {
		l.blocks[*c][r-q] = util.NewVector(float64(*c*BLOCK_SIZE), float64((r-q)*BLOCK_SIZE))
	}
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
		bmap:    make([][]BoolMap, blockW, blockW),
		figures: make([]*Figure, 0),
	}
	fmt.Println("Mario: generating level...")
	for c := 0; c < blockW; c++ {
		lvl.blocks[c] = make([]*util.Vector, blockH, blockH)
		for r := 0; r < blockH; r++ {
			lvl.blocks[c][r] = nil
		}
	}

	for c, obs := 0, 1; c < blockW; c, obs = c+1, obs+1 {

		pr := c
		if obs%OBSTACLE_SPACING == 0 {
			if neural.Chance(0.7) {
				lvl.makeHole(&c)
			} else {
				lvl.makeObstacle(&c)
			}
		} else {
			lvl.makeGround(&c)
		}
		obs += c - pr
	}

	fmt.Println("Mario: generating bool map...")
	for c := 0; c < blockW; c++ {
		lvl.bmap[c] = make([]BoolMap, blockH, blockH)
		for r := 0; r < blockH; r++ {
			lvl.bmap[c][r] = lvl.boolMapAtIdx(c, r)
		}
	}

	return lvl
}

func (l *Level) ToLevelCoords(pos *util.Vector) (int, int) {
	return int(pos.X / float64(BLOCK_SIZE)), int(pos.Y / float64(BLOCK_SIZE))
}

func (l *Level) OffsetInLevelGrid(from, to *util.Vector) (int, int) {
	fx, fy := l.ToLevelCoords(from)
	tx, ty := l.ToLevelCoords(to)
	return tx - fx, ty - fy
}

func (l *Level) validCoord(w, h int) bool {
	return w >= 0 && w < (int(l.size.X)/BLOCK_SIZE) && h >= 0 && h < (int(l.size.Y)/BLOCK_SIZE)
}

func (l *Level) IsSolid(pos *util.Vector) bool {
	return l.CubeAt(pos) != nil
}

func (l *Level) boolMapAtIdx(cx, cy int) BoolMap {
	var res BoolMap = NewBMap()
	cx -= (BOOL_MAP_GRID_SIDE - 1) / 2
	cy -= (BOOL_MAP_GRID_SIDE - 1) / 2

	off := 0
	for c := 0; c < BOOL_MAP_GRID_SIDE; c++ {
		for r := 0; r < BOOL_MAP_GRID_SIDE; r++ {
			if l.validCoord(cx+c, cy+r) && l.blocks[cx+c][cy+r] != nil {
				res.Set(off)
			}
			off++
		}
	}

	return res
}

func (l *Level) BoolMapAt(pos *util.Vector) BoolMap {
	cx, cy := l.ToLevelCoords(pos)
	if l.validCoord(cx, cy) {
		return l.bmap[cx][cy]
	} else {
		return NewBMap()
	}
}

func (l *Level) CubeAt(pos *util.Vector) *util.Vector {
	w, h := l.ToLevelCoords(pos)
	if l.validCoord(w, h) {
		return l.blocks[w][h]
	} else {
		return nil
	}
}

func (l *Level) FloorAt(pos *util.Vector) *util.Vector {
	wIdx, hIdx := l.ToLevelCoords(pos)

	if !l.validCoord(wIdx, hIdx+1) {
		return nil
	}
	return l.blocks[wIdx][hIdx+1]
}

func (l *Level) NewFigurePos() *util.Vector {
	return util.NewVector(float64(BLOCK_SIZE)/2, l.size.Y-float64(BLOCK_SIZE)*2)
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
		// nextPos = position + timestep * (velocity + timestep * acceleration / 2);
		// TODO not use go
		l.figures[c].nextPos = *l.figures[c].pos.Add(G_FORCE.Scale(dt / 2).Add(&l.figures[c].vel).Scale(dt))

		// velocity += timestep * acceleration;
		l.figures[c].vel = *l.figures[c].vel.Add(G_FORCE.Scale(dt))
		l.figures[c].vel.X *= (1 - 3*dt)
		// l.figures[c].vel.X = 0
	}
}
