package mario

import (
	neural "github.com/poseidon4o/go-neural/src/neural"
	util "github.com/poseidon4o/go-neural/src/util"
)

type MarioNode struct {
	fig   *Figure
	brain *neural.Net
	bestX float64
	dead  bool
}

type Mario struct {
	figures  []MarioNode
	lvl      Level
	drawCb   func(pos, size *util.Vector, color uint32)
	drawSize int
}

func (m *Mario) Completed() float64 {
	return 0
}

func (m *Mario) Done() bool {
	return false
}

func (m *Mario) SetDrawRectCb(cb func(pos, size *util.Vector, color uint32)) {
	m.drawCb = cb
}

func (m *Mario) LogicTick(dt float64) {

}

func NewMario(figCount int, size *util.Vector) *Mario {
	level := NewLevel(int(size.X), int(size.Y))
	level.AddFigures(1)

	figs := make([]MarioNode, 1)
	figs[0].fig = level.figures[0]
	figs[0].brain = nil

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
			m.drawCb(&m.lvl.blocks[c][r], blSize, red)
			m.drawCb(m.lvl.blocks[c][r].Add(translate), blSizeSmall, green)
		}
	}

	for c := range m.figures {
		m.drawCb(&m.figures[c].fig.Pos, size, blue)
	}
}
