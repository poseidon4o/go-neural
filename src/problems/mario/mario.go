package mario

import (
	"bufio"
	"fmt"
	neural "github.com/poseidon4o/go-neural/src/neural"
	util "github.com/poseidon4o/go-neural/src/util"
	"math"
	"math/rand"
	"os"
	"sort"
)

const idleThreshold uint32 = 300

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

	H1        NeuronName = iota
	H2        NeuronName = iota
	H3        NeuronName = iota
	H4        NeuronName = iota
	R1        NeuronName = iota
	R2        NeuronName = iota
	R3        NeuronName = iota
	R4        NeuronName = iota
	jump      NeuronName = iota
	xMove     NeuronName = iota
	NRN_COUNT int        = iota
)

func nrn(name NeuronName) int {
	return int(name)
}

type MarioOutput struct {
	jump float64
	move float64
}

type MarioNode struct {
	fig        *Figure
	brain      *neural.Net
	cache      map[BoolMap]MarioOutput
	bestX      float64
	idleX      int
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

type MarioStats struct {
	dead      int
	crossed   int
	culled    int
	decisions int
	completed float64
	cacheHits int
	cacheSize int
}

func (m *MarioStats) zero() {
	m.dead = 0
	m.crossed = 0
	m.culled = 0
	m.decisions = 0
	m.cacheHits = 0
	m.cacheSize = 0
}

func (m *MarioStats) print() {
	fmt.Printf("Dead [%d] of them [%d] culled. Crosses [%d], mutations[%d].\nNeural net decisions [%d], from cache (%.3f%%), cached values for last frame: [%d]\n",
		m.dead, m.culled, m.crossed, m.dead-m.crossed, m.decisions, 100*(float64(m.cacheHits)/float64(m.decisions)), m.cacheSize)
}

type DebugMario struct {
	figure *MarioNode
	output MarioOutput
}

type Mario struct {
	stats    MarioStats
	figures  MarioCol
	dbg      DebugMario
	lvl      Level
	drawCb   func(pos, size *util.Vector, color uint32)
	drawSize int
}

const saveName = "mario-save.dat"

func (m *Mario) SaveNetsToFile() {
	file, err := os.Create(saveName)
	if err != nil {
		fmt.Println("Failed to open file ", saveName)
		return
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	fmt.Fprintf(w, "%d\n", len(m.figures))
	for c := range m.figures {
		fmt.Fprintln(w, "")
		m.figures[c].brain.WriteTo(w)
	}
	w.Flush()
	fmt.Println("Saved marion to ", saveName)
}

func (m *Mario) LoadNetsFromFile() {
	file, err := os.Open(saveName)
	if err != nil {
		fmt.Println("Failed to open file ", saveName)
		return
	}
	defer file.Close()

	r := bufio.NewReader(file)
	cnt := 0
	fmt.Fscanf(r, "%d", cnt)
	m = NewMario(cnt, &m.lvl.size)
	for c := 0; c < cnt; c++ {
		m.figures[c].brain.ReadFrom(r)
	}
	fmt.Println("Loaded mario from ", saveName)
}

func (m *Mario) Complete() float64 {
	return m.stats.completed
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
	m.stats.completed = m.figures[0].bestX / m.lvl.size.X
	m.stats.cacheSize = 0

	wg := make(chan struct{}, len(m.figures))

	stepC := func(r int) {
		m.checkStep(r)
		m.mutateStep(r)
		m.thnikStep(r)
		wg <- struct{}{}
	}

	m.dbg.figure = &m.figures[0]
	for c := range m.figures {
		if m.figures[c].fig.pos.X > m.dbg.figure.fig.pos.X {
			m.dbg.figure = &m.figures[c]
		}
	}

	for c := range m.figures {
		m.stats.cacheSize += len(m.figures[c].cache)
		go stepC(c)
	}

	for c := 0; c < len(m.figures); c++ {
		<-wg
	}
}

func (m *Mario) StatsReportTick() {
	m.stats.print()
	m.stats.zero()
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

		for r := 0; r < (nrn(H4) - nrn(H1)); r++ {
			// input to H
			for inp := nrn(I0); inp <= nrn(I48); inp++ {
				*nets[c].Synapse(inp+nrn(I0), r+nrn(H1)) = 0.0
			}

			// R to output
			*nets[c].Synapse(r+nrn(R1), nrn(jump)) = 0.0
			*nets[c].Synapse(r+nrn(R1), nrn(xMove)) = 0.0
		}

		for r := 0; r < (nrn(H4) - nrn(H1)); r++ {
			for q := 0; q < (nrn(H4) - nrn(H1)); q++ {
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
		figs[c].cache = make(map[BoolMap]MarioOutput)
	}

	return &Mario{
		figures: figs,
		lvl:     *level,
		dbg: DebugMario{
			figure: &figs[0],
			output: MarioOutput{
				jump: 0,
				move: 0,
			},
		},
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

	m.drawCb(
		m.dbg.figure.fig.pos.Add(util.NewVector(-3., -20.)),
		util.NewVector(2, 20),
		0xff000000)

	if m.dbg.output.move > 0 {
		m.drawCb(
			m.dbg.figure.fig.pos.Add(util.NewVector(0, -3)),
			util.NewVector(20, 2),
			0xff000000)
	} else {
		m.drawCb(
			m.dbg.figure.fig.pos.Add(util.NewVector(-20, -3)),
			util.NewVector(20, 2),
			0xff000000)
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
	m.stats.decisions++
	bmap := m.lvl.BoolMapAt(&m.figures[c].fig.pos)

	res, ok := m.figures[c].cache[bmap]

	if !ok {
		idx := 0
		for idx = 0; idx < nrn(I48); idx++ {
			if bmap.At(idx) {
				m.figures[c].brain.Stimulate(int(idx)+nrn(I0), -1)
			} else {
				m.figures[c].brain.Stimulate(int(idx)+nrn(I0), 1)
			}
		}
		m.figures[c].brain.Step()

		res = MarioOutput{
			jump: m.figures[c].brain.ValueOf(nrn(jump)),
			move: m.figures[c].brain.ValueOf(nrn(xMove)),
		}
		m.figures[c].cache[bmap] = res
		m.figures[c].brain.Clear()
	} else {
		m.stats.cacheHits++
	}

	if res.jump > 0.75 {
		m.figures[c].fig.Jump()
	}

	if math.Abs(res.move) > 0.75 {
		m.figures[c].fig.Move(int(res.move * 10))
	}

	if m.dbg.figure == &m.figures[c] {
		m.dbg.output = res
	}
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
		m.stats.dead++
		m.figures[c].dead = false
		m.figures[c].cache = make(map[BoolMap]MarioOutput)
		m.figures[c].fig.pos = *m.lvl.NewFigurePos()
		m.figures[c].fig.vel = *util.NewVector(0, 0)

		needsMutation := true

		mutateChance := (float64(c) / float64(len(m.figures))) * 2.0
		forceCross := c >= len(m.figures)/2

		if forceCross || neural.Chance(mutateChance) {
			*m.figures[c].brain = *neural.Cross2(m.randNet(), m.randNet())
			m.stats.crossed++
			needsMutation = false
		}

		if m.figures[c].idleFrames >= idleThreshold {
			m.stats.culled++
			if needsMutation {
				needsMutation = false
				m.figures[c].brain.Mutate(0.75)
				m.figures[c].bestX *= 0.25
			}
		}

		if needsMutation || neural.Chance(0.01) {
			m.figures[c].brain.MutateWithMagnitude(0.01, 0.01)
			m.figures[c].bestX *= 0.975
		}

		m.figures[c].idleFrames = 0
		m.figures[c].idleX = 0
	} else {
		if m.figures[c].fig.pos.X > m.figures[c].bestX {
			m.figures[c].bestX = m.figures[c].fig.pos.X
		}

		if int(m.figures[c].fig.pos.X) > m.figures[c].idleX {
			m.figures[c].idleX = int(m.figures[c].fig.pos.X)
			m.figures[c].idleFrames = 0
		} else {
			m.figures[c].idleFrames++
			if m.figures[c].idleFrames >= idleThreshold {
				m.figures[c].dead = true
			}
		}

	}
}
