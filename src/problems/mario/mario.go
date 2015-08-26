package mario

import (
	"fmt"
	neural "github.com/poseidon4o/go-neural/src/neural"
	util "github.com/poseidon4o/go-neural/src/util"
	"math"
	"math/rand"
	"sort"
)

type NeuronName int

const (
	posX      NeuronName = iota
	posY      NeuronName = iota
	velY      NeuronName = iota
	velX      NeuronName = iota
	H1        NeuronName = iota
	H2        NeuronName = iota
	H3        NeuronName = iota
	H4        NeuronName = iota
	H5        NeuronName = iota
	H6        NeuronName = iota
	H7        NeuronName = iota
	H8        NeuronName = iota
	R1        NeuronName = iota
	R2        NeuronName = iota
	R3        NeuronName = iota
	R4        NeuronName = iota
	R5        NeuronName = iota
	R6        NeuronName = iota
	R7        NeuronName = iota
	R8        NeuronName = iota
	jump      NeuronName = iota
	xMove     NeuronName = iota
	NRN_COUNT int        = iota
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
		if len(m.figures) > 1 {
			m.thnikStep(r)
		}
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

	nets := make([]*neural.Net, figCount, figCount)
	for c := range nets {
		nets[c] = neural.NewNet(NRN_COUNT)

		for r := 0; r < (nrn(H8) - nrn(H1)); r++ {
			// input to H
			*nets[c].Synapse(nrn(posX), r+nrn(H1)) = 0.0
			*nets[c].Synapse(nrn(posY), r+nrn(H1)) = 0.0
			*nets[c].Synapse(nrn(velX), r+nrn(H1)) = 0.0
			*nets[c].Synapse(nrn(velY), r+nrn(H1)) = 0.0

			// R to output
			*nets[c].Synapse(r+nrn(R1), nrn(jump)) = 0.0
			*nets[c].Synapse(r+nrn(R1), nrn(xMove)) = 0.0
		}

		for r := 0; r < (nrn(H8) - nrn(H1)); r++ {
			for q := 0; q < (nrn(H8) - nrn(H1)); q++ {
				*nets[c].Synapse(r+nrn(H1), q+nrn(R1)) = 0.0
			}
		}

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

	block := m.lvl.FloorAt(&fig.pos)

	if block == nil || fig.nextPos.Y < block.Y {
		fig.pos.Y = fig.nextPos.Y
	} else {
		// m.drawCb(block, util.NewVector(float64(BLOCK_SIZE), float64(BLOCK_SIZE)), 0xff00ffff)
		// land on block
		fig.vel.Y = 0
		fig.pos.Y = block.Y - 0.1
		fig.Land()
	}

	if fig.pos.X != fig.nextPos.X {
		fig.nextPos.Y = fig.pos.Y
		colide := m.lvl.CubeAt(&fig.nextPos)
		if colide != nil {
			// m.drawCb(colide, util.NewVector(float64(BLOCK_SIZE), float64(BLOCK_SIZE)), 0xff00ffff)
			if fig.pos.X < fig.nextPos.X {
				// collide right
				fig.pos.X = colide.X - 0.1
			} else {
				// colide left
				fig.pos.X = colide.X + float64(BLOCK_SIZE) + 0.1
			}
		} else {
			fig.pos.X = fig.nextPos.X
		}
	}

}

func (m *Mario) thnikStep(c int) {
	discreteX := float64(int(m.figures[c].fig.pos.X / float64(OBSTACLE_SPACING*BLOCK_SIZE)))
	m.figures[c].brain.Stimulate(nrn(posX), discreteX)

	m.figures[c].brain.Stimulate(nrn(posY), m.figures[c].fig.pos.Y)
	m.figures[c].brain.Stimulate(nrn(velX), m.figures[c].fig.vel.X)
	m.figures[c].brain.Stimulate(nrn(velY), m.figures[c].fig.vel.Y)

	m.figures[c].brain.Step()

	if m.figures[c].brain.ValueOf(nrn(jump)) > 0.9 {
		m.figures[c].fig.Jump()
	}

	xMoveValue := m.figures[c].brain.ValueOf(nrn(xMove))
	if math.Abs(xMoveValue) > 0.9 {
		m.figures[c].fig.Move(int(xMoveValue * 10))
	}

	m.figures[c].brain.Clear()
}

const idleThreshold uint32 = 2000

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
			swapChance := (float64(c) / float64(len(m.figures))) * 2.0
			if neural.Chance(swapChance) {
				*m.figures[c].brain = *m.randNet()
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
		} else {
			m.figures[c].idleFrames++
			if m.figures[c].idleFrames >= idleThreshold {
				m.figures[c].dead = true
				// c--
			}
		}

	}
}
