package mario

import (
	"fmt"
	neural "github.com/poseidon4o/go-neural/src/neural"
	util "github.com/poseidon4o/go-neural/src/util"
	"math"
	"math/rand"
	"sort"
)

const idleThreshold uint32 = 100

type NeuronName int

const (
	I0  NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	_   NeuronName = iota
	I48 NeuronName = iota

	H1 NeuronName = iota
	H2 NeuronName = iota
	H3 NeuronName = iota
	H4 NeuronName = iota
	H5 NeuronName = iota
	H6 NeuronName = iota
	H7 NeuronName = iota
	H8 NeuronName = iota
	H9 NeuronName = iota

	jump  NeuronName = iota
	xMove NeuronName = iota

	NRN_COUNT int = iota
)

func nrn(name NeuronName) int {
	return int(name)
}

type MarioNode struct {
	fig        *Figure
	brain      *neural.Net
	bestX      float64
	idleX      float64
	dead       bool
	idleFrames uint32
}

type MarioCol []MarioNode

func (figs MarioCol) Len() int {
	return len(figs)
}

func (figs MarioCol) Less(c, r int) bool {
	return figs[c].bestX > figs[r].bestX
}

func (figs MarioCol) Swap(c, r int) {
	figs[c], figs[r] = figs[r], figs[c]
}

type Mario struct {
	figures  MarioCol
	lvl      Level
	drawCb   func(pos, size *util.Vector, color uint32)
	drawSize int
}

func (m *Mario) Complete() float64 {
	return m.figures[0].bestX / m.lvl.size.X
}

func (m *Mario) Done() bool {
	return false
}

func (m *Mario) SetDrawRectCb(cb func(pos, size *util.Vector, color uint32)) {
	m.drawCb = cb
}

func (m *Mario) LogicTick(dt float64) {
	m.lvl.Step(dt)
	sort.Sort(m.figures)

	wg := make(chan struct{}, len(m.figures))

	stepC := func(r int) {
		m.checkStep(r)
		m.mutateStep(r)
		m.thnikStep(r)
		wg <- struct{}{}
	}

	for c := range m.figures {
		go stepC(c)
	}

	for c := 0; c < len(m.figures); c++ {
		<-wg
	}
}

func (m *Mario) Jump() {
	m.figures[0].fig.Jump()
}

func (m *Mario) Move(dir int) {
	m.figures[0].fig.Move(dir)
}

func (m *Mario) Figs() MarioCol {
	return m.figures
}

func NewMario(figCount int, size *util.Vector) *Mario {
	fmt.Println("")
	level := NewLevel(int(size.X), int(size.Y))
	level.AddFigures(figCount)

	finp := func(id int) int {
		return nrn(I0 + NeuronName(id))
	}

	nets := make([]*neural.Net, figCount, figCount)
	for c := range nets {
		nets[c] = neural.NewNet(NRN_COUNT)

		for r := 0; r < 6; r++ {
			*nets[c].Synapse(nrn(H1)+r, nrn(jump)) = 0.0
			*nets[c].Synapse(nrn(H1)+r, nrn(xMove)) = 0.0
		}

		*nets[c].Synapse(finp(0), nrn(H1)) = 0.0
		*nets[c].Synapse(finp(1), nrn(H1)) = 0.0
		*nets[c].Synapse(finp(2), nrn(H1)) = 0.0
		*nets[c].Synapse(finp(7), nrn(H1)) = 0.0
		*nets[c].Synapse(finp(8), nrn(H1)) = 0.0
		*nets[c].Synapse(finp(9), nrn(H1)) = 0.0
		*nets[c].Synapse(finp(14), nrn(H1)) = 0.0
		*nets[c].Synapse(finp(15), nrn(H1)) = 0.0
		*nets[c].Synapse(finp(16), nrn(H1)) = 0.0

		*nets[c].Synapse(finp(2), nrn(H2)) = 0.0
		*nets[c].Synapse(finp(3), nrn(H2)) = 0.0
		*nets[c].Synapse(finp(4), nrn(H2)) = 0.0
		*nets[c].Synapse(finp(9), nrn(H2)) = 0.0
		*nets[c].Synapse(finp(10), nrn(H2)) = 0.0
		*nets[c].Synapse(finp(11), nrn(H2)) = 0.0
		*nets[c].Synapse(finp(16), nrn(H2)) = 0.0
		*nets[c].Synapse(finp(17), nrn(H2)) = 0.0
		*nets[c].Synapse(finp(18), nrn(H2)) = 0.0

		*nets[c].Synapse(finp(4), nrn(H3)) = 0.0
		*nets[c].Synapse(finp(5), nrn(H3)) = 0.0
		*nets[c].Synapse(finp(6), nrn(H3)) = 0.0
		*nets[c].Synapse(finp(11), nrn(H3)) = 0.0
		*nets[c].Synapse(finp(12), nrn(H3)) = 0.0
		*nets[c].Synapse(finp(13), nrn(H3)) = 0.0
		*nets[c].Synapse(finp(18), nrn(H3)) = 0.0
		*nets[c].Synapse(finp(19), nrn(H3)) = 0.0
		*nets[c].Synapse(finp(20), nrn(H3)) = 0.0

		*nets[c].Synapse(finp(14), nrn(H4)) = 0.0
		*nets[c].Synapse(finp(15), nrn(H4)) = 0.0
		*nets[c].Synapse(finp(16), nrn(H4)) = 0.0
		*nets[c].Synapse(finp(21), nrn(H4)) = 0.0
		*nets[c].Synapse(finp(22), nrn(H4)) = 0.0
		*nets[c].Synapse(finp(23), nrn(H4)) = 0.0
		*nets[c].Synapse(finp(28), nrn(H4)) = 0.0
		*nets[c].Synapse(finp(29), nrn(H4)) = 0.0
		*nets[c].Synapse(finp(30), nrn(H4)) = 0.0

		*nets[c].Synapse(finp(17), nrn(H5)) = 0.0
		*nets[c].Synapse(finp(17), nrn(H5)) = 0.0
		*nets[c].Synapse(finp(18), nrn(H5)) = 0.0
		*nets[c].Synapse(finp(23), nrn(H5)) = 0.0
		*nets[c].Synapse(finp(24), nrn(H5)) = 0.0
		*nets[c].Synapse(finp(25), nrn(H5)) = 0.0
		*nets[c].Synapse(finp(30), nrn(H5)) = 0.0
		*nets[c].Synapse(finp(31), nrn(H5)) = 0.0
		*nets[c].Synapse(finp(32), nrn(H5)) = 0.0

		*nets[c].Synapse(finp(18), nrn(H6)) = 0.0
		*nets[c].Synapse(finp(19), nrn(H6)) = 0.0
		*nets[c].Synapse(finp(20), nrn(H6)) = 0.0
		*nets[c].Synapse(finp(25), nrn(H6)) = 0.0
		*nets[c].Synapse(finp(26), nrn(H6)) = 0.0
		*nets[c].Synapse(finp(27), nrn(H6)) = 0.0
		*nets[c].Synapse(finp(32), nrn(H6)) = 0.0
		*nets[c].Synapse(finp(33), nrn(H6)) = 0.0
		*nets[c].Synapse(finp(34), nrn(H6)) = 0.0

		*nets[c].Synapse(finp(28), nrn(H7)) = 0.0
		*nets[c].Synapse(finp(29), nrn(H7)) = 0.0
		*nets[c].Synapse(finp(30), nrn(H7)) = 0.0
		*nets[c].Synapse(finp(35), nrn(H7)) = 0.0
		*nets[c].Synapse(finp(36), nrn(H7)) = 0.0
		*nets[c].Synapse(finp(37), nrn(H7)) = 0.0
		*nets[c].Synapse(finp(42), nrn(H7)) = 0.0
		*nets[c].Synapse(finp(43), nrn(H7)) = 0.0
		*nets[c].Synapse(finp(44), nrn(H7)) = 0.0

		*nets[c].Synapse(finp(30), nrn(H8)) = 0.0
		*nets[c].Synapse(finp(31), nrn(H8)) = 0.0
		*nets[c].Synapse(finp(32), nrn(H8)) = 0.0
		*nets[c].Synapse(finp(37), nrn(H8)) = 0.0
		*nets[c].Synapse(finp(38), nrn(H8)) = 0.0
		*nets[c].Synapse(finp(39), nrn(H8)) = 0.0
		*nets[c].Synapse(finp(44), nrn(H8)) = 0.0
		*nets[c].Synapse(finp(45), nrn(H8)) = 0.0
		*nets[c].Synapse(finp(46), nrn(H8)) = 0.0

		*nets[c].Synapse(finp(32), nrn(H9)) = 0.0
		*nets[c].Synapse(finp(33), nrn(H9)) = 0.0
		*nets[c].Synapse(finp(34), nrn(H9)) = 0.0
		*nets[c].Synapse(finp(39), nrn(H9)) = 0.0
		*nets[c].Synapse(finp(40), nrn(H9)) = 0.0
		*nets[c].Synapse(finp(41), nrn(H9)) = 0.0
		*nets[c].Synapse(finp(46), nrn(H9)) = 0.0
		*nets[c].Synapse(finp(47), nrn(H9)) = 0.0
		*nets[c].Synapse(finp(48), nrn(H9)) = 0.0

		nets[c].Randomize()
	}

	figs := make(MarioCol, figCount, figCount)
	for c := range figs {
		figs[c].brain = nets[c]
		figs[c].dead = false
		figs[c].bestX = 0
		figs[c].fig = level.figures[c]
	}

	return &Mario{
		figures:  figs,
		lvl:      *level,
		drawCb:   func(pos, size *util.Vector, color uint32) {},
		drawSize: 5,
	}
}

func (m *Mario) DrawTick() {
	var (
		red   = uint32(0xffff0000)
		green = uint32(0xff00ff00)
		blue  = uint32(0xff0000ff)
	)

	blSize := util.NewVector(float64(BLOCK_SIZE), float64(BLOCK_SIZE))
	blSizeSmall := blSize.Scale(0.5)

	translate := util.NewVector(6, 6)

	size := util.NewVector(float64(m.drawSize), float64(m.drawSize))

	for c := range m.lvl.blocks {
		for r := range m.lvl.blocks[c] {
			if m.lvl.blocks[c][r] != nil {
				m.drawCb(m.lvl.blocks[c][r], blSize, red)
				m.drawCb(m.lvl.blocks[c][r].Add(translate), blSizeSmall, green)
			}
		}
	}

	for c := range m.figures {
		m.drawCb(m.figures[c].fig.pos.Add(size.Scale(0.5).Neg()), size, blue)
	}
}

func (m *Mario) checkStep(c int) {
	fig := m.figures[c].fig

	if fig.nextPos.Y > m.lvl.size.Y || fig.nextPos.Y < 0 {
		m.figures[c].dead = true
		return
	}

	if fig.nextPos.X < 0 {
		fig.nextPos.X = 0
	} else if fig.nextPos.X > m.lvl.size.X {
		fig.nextPos.X = m.lvl.size.X
	}

	bmap := m.lvl.BoolMapAt(&fig.pos)
	dx, dy := m.lvl.OffsetInLevelGrid(&fig.pos, &fig.nextPos)
	tx, ty := m.lvl.ToLevelCoords(&fig.nextPos)

	// if bmap.GridAt(dx, dy) {
	// 	m.drawCb(util.NewVector(float64(tx*BLOCK_SIZE), float64(ty*BLOCK_SIZE)), util.NewVector(float64(BLOCK_SIZE), float64(BLOCK_SIZE)), 0xff00ffff)
	// }

	checkx, checky := true, true

	if dx != 0 && dy != 0 && bmap.GridAt(dx, dy) {
		slideY := bmap.GridAt(dx, 0) && bmap.GridAt(dx, dy) && !bmap.GridAt(0, dy)
		slideX := bmap.GridAt(0, dy) && bmap.GridAt(dx, dy) && !bmap.GridAt(dx, 0)
		if slideX != slideY {
			if slideX {
				checkx = false
				fig.nextPos.Y = fig.pos.Y
			} else {
				checky = false
				fig.nextPos.X = fig.pos.X
			}
		}
	}

	if dx != 0 && bmap.GridAt(dx, dy) && checkx {
		fig.vel.X = 0
		if dx > 0 {
			fig.pos.X = float64(tx*BLOCK_SIZE) - 0.1
		} else {
			fig.pos.X = float64((tx+1)*BLOCK_SIZE) + 0.1
		}
	} else {
		fig.pos.X = fig.nextPos.X
	}

	if dy != 0 && bmap.GridAt(dx, dy) && checky {
		fig.vel.Y = 0
		if dy > 0 {
			fig.pos.Y = float64(ty*BLOCK_SIZE) - 0.1
			fig.Land()
		} else {
			fig.pos.Y = float64((ty+1)*BLOCK_SIZE) + 0.1
		}
	} else {
		fig.pos.Y = fig.nextPos.Y
	}
}

func (m *Mario) thnikStep(c int) {
	bmap := m.lvl.BoolMapAt(&m.figures[c].fig.pos)

	idx := 0
	for idx = 0; idx < nrn(I48); idx++ {
		if bmap.At(idx) {
			m.figures[c].brain.Stimulate(int(idx)+nrn(I0), -1)
		} else {
			m.figures[c].brain.Stimulate(int(idx)+nrn(I0), 1)
		}
	}

	m.figures[c].brain.Step()

	if m.figures[c].brain.ValueOf(nrn(jump)) > 0.75 {
		m.figures[c].fig.Jump()
	}

	xMoveValue := m.figures[c].brain.ValueOf(nrn(xMove))
	if math.Abs(xMoveValue) > 0.75 {
		m.figures[c].fig.Move(int(xMoveValue * 10))
	}

	m.figures[c].brain.Clear()
}

func (m *Mario) randNet() *neural.Net {
	cutOff := 10.0
	idx := 0
	for {
		r := rand.ExpFloat64()
		if r <= cutOff {
			idx = int((r * float64(len(m.figures))) / cutOff)
			break
		}
	}
	return m.figures[idx].brain
}

func (m *Mario) mutateStep(c int) {

	if m.figures[c].dead {
		m.figures[c].dead = false
		m.figures[c].fig.pos = *m.lvl.NewFigurePos()
		m.figures[c].fig.vel = *util.NewVector(0, 0)

		if m.figures[c].idleFrames >= idleThreshold {
			m.figures[c].brain.Mutate(0.75)
			m.figures[c].bestX *= 0.25
		} else {
			if neural.Chance(0.5) {
				*m.figures[c].brain = *neural.Cross2(m.randNet(), m.randNet())
			}
			m.figures[c].brain.MutateWithMagnitude(0.01, 0.1)
			m.figures[c].bestX *= 0.975
		}

		m.figures[c].idleFrames = 0
		m.figures[c].idleX = 0
	} else {
		if m.figures[c].fig.pos.X > m.figures[c].bestX {
			m.figures[c].bestX = m.figures[c].fig.pos.X
		}

		if m.figures[c].fig.pos.X > m.figures[c].idleX {
			m.figures[c].idleX = m.figures[c].fig.pos.X
			m.figures[c].idleFrames = 0
		} else {
			m.figures[c].idleFrames++
			if m.figures[c].idleFrames >= idleThreshold {
				m.figures[c].dead = true
				// c--
			}
		}

	}
}
