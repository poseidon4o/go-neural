package flappy

import (
	neural "github.com/poseidon4o/go-neural/src/neural"
	util "github.com/poseidon4o/go-neural/src/util"
	"math"
	"sort"
)

type NeuronName int

const (
	diffY     NeuronName = iota
	diffX     NeuronName = iota
	velY      NeuronName = iota
	H1        NeuronName = iota
	H2        NeuronName = iota
	H3        NeuronName = iota
	H4        NeuronName = iota
	jump      NeuronName = iota
	NRN_COUNT int        = iota
)

func nrn(name NeuronName) int {
	return int(name)
}

type FBird struct {
	bird  *Bird
	brain *neural.Net
	bestX float64
	dead  bool
}

type Flock []FBird

func (birds Flock) Len() int {
	return len(birds)
}

func (birds Flock) Less(c, r int) bool {
	return birds[c].bestX > birds[r].bestX
}

func (birds Flock) Swap(c, r int) {
	birds[c], birds[r] = birds[r], birds[c]
}

type Flappy struct {
	birds    Flock
	lvl      Level
	drawCb   func(pos, size *util.Vector, color uint32)
	drawSize int
}

func (f *Flappy) Completed() float64 {
	return f.birds[0].bestX / f.lvl.size.X
}

func (f *Flappy) Done() bool {
	return f.birds[0].bestX > f.lvl.pylons[len(f.lvl.pylons)-1].X
}

func (f *Flappy) SetDrawRectCb(cb func(pos, size *util.Vector, color uint32)) {
	f.drawCb = cb
}

func (f *Flappy) LogicTick(dt float64) {
	f.lvl.Step(dt)
	f.checkFlock()
	f.mutateFlock()
	f.thnikFlock()
}

func (f *Flappy) DrawTick() {
	var (
		red   = uint32(0xffff0000)
		green = uint32(0xff00ff00)
		blue  = uint32(0xff0000ff)
	)

	var tl, size util.Vector
	size.X = float64(f.drawSize)
	size.Y = float64(f.drawSize)

	for c := range f.birds {
		f.drawCb(&f.birds[c].bird.Pos, &size, red)
	}

	hSize := float64(PylonHole) / 2.0
	for _, pylon := range f.lvl.pylons {
		// top part
		tl.X = pylon.X
		tl.Y = 0
		size.Y = pylon.Y - hSize
		f.drawCb(&tl, &size, green)

		// bottom part
		tl.Y = pylon.Y + hSize
		size.Y = f.lvl.size.Y - (pylon.Y + hSize)
		f.drawCb(&tl, &size, green)

		// middle point
		tl.Y = pylon.Y
		size.Y = float64(f.drawSize)
		f.drawCb(&tl, &size, blue)
	}
}

func NewFlappy(birdCount int, size *util.Vector) *Flappy {
	level := NewLevel(int(size.X), int(size.Y))

	nets := make([]*neural.Net, birdCount, birdCount)
	for c := range nets {
		nets[c] = neural.NewNet(NRN_COUNT)

		// diffY- to hidden
		*nets[c].Synapse(nrn(diffY), nrn(H1)) = 0.0
		*nets[c].Synapse(nrn(diffY), nrn(H2)) = 0.0
		*nets[c].Synapse(nrn(diffY), nrn(H3)) = 0.0
		*nets[c].Synapse(nrn(diffY), nrn(H4)) = 0.0

		// diffX- to hidden
		*nets[c].Synapse(nrn(diffX), nrn(H1)) = 0.0
		*nets[c].Synapse(nrn(diffX), nrn(H2)) = 0.0
		*nets[c].Synapse(nrn(diffX), nrn(H3)) = 0.0
		*nets[c].Synapse(nrn(diffX), nrn(H4)) = 0.0

		// velY - to hidden
		*nets[c].Synapse(nrn(velY), nrn(H1)) = 0.0
		*nets[c].Synapse(nrn(velY), nrn(H2)) = 0.0
		*nets[c].Synapse(nrn(velY), nrn(H3)) = 0.0
		*nets[c].Synapse(nrn(velY), nrn(H4)) = 0.0

		// hidden to output
		*nets[c].Synapse(nrn(H1), nrn(jump)) = 0.0
		*nets[c].Synapse(nrn(H2), nrn(jump)) = 0.0
		*nets[c].Synapse(nrn(H3), nrn(jump)) = 0.0
		*nets[c].Synapse(nrn(H4), nrn(jump)) = 0.0

		nets[c].Randomize()
	}

	level.AddBirds(birdCount)
	flock := make(Flock, birdCount)
	for c := 0; c < birdCount; c++ {
		flock[c].bird = level.birds[c]
		flock[c].brain = nets[c]
		flock[c].bestX = 0
		flock[c].dead = false
	}

	return &Flappy{
		birds:    flock,
		lvl:      *level,
		drawCb:   func(pos, size *util.Vector, color uint32) {},
		drawSize: 5,
	}
}

// will check if going from pos to next will collide
func (f *Flappy) checkFlock() {

	// just for readability
	collide := func(aX, bX, cX float64) bool {
		// c.X == d.X
		return aX-1 <= cX && bX+1 >= cX
	}

	hSize := float64(PylonHole / 2)

	for c := range f.birds {
		if f.birds[c].bird.Pos.Y >= f.lvl.size.Y || f.birds[c].bird.Pos.Y < 1 {
			// hit ceeling or floor
			f.birds[c].dead = true
			continue
		}

		pylon := f.lvl.ClosestPylon(&f.birds[c].bird.Pos)
		if f.birds[c].bird.Pos.Y >= pylon.Y-hSize && f.birds[c].bird.Pos.Y <= pylon.Y+hSize {
			// in the pylon hole
			continue
		}

		f.birds[c].dead = collide(f.birds[c].bird.Pos.X, f.birds[c].bird.NextPos.X, pylon.X)
	}

}

func (f *Flappy) thnikFlock() {
	wg := make(chan struct{}, len(f.birds))

	thinkBird := func(c int) {
		next := f.lvl.FirstPylonAfter(&f.birds[c].bird.Pos)
		diffYval := next.Y - f.birds[c].bird.Pos.Y
		diffXval := next.X - f.birds[c].bird.Pos.X

		f.birds[c].brain.Stimulate(nrn(diffY), diffYval)
		f.birds[c].brain.Stimulate(nrn(diffX), diffXval)
		f.birds[c].brain.Stimulate(nrn(velY), f.birds[c].bird.Vel.Y)

		f.birds[c].brain.Step()
		if f.birds[c].brain.ValueOf(nrn(jump)) > 0.75 {
			f.birds[c].bird.Vel.Y = -500
		}

		f.birds[c].brain.Clear()
		wg <- struct{}{}
	}

	for c := 0; c < len(f.birds); c++ {
		go thinkBird(c)
	}

	for c := 0; c < len(f.birds); c++ {
		<-wg
	}
}

func (f *Flappy) mutateFlock() {
	sort.Sort(f.birds)

	randNet := func() *neural.Net {
		return f.birds[int(neural.RandMax(float64(len(f.birds))))].brain
	}

	best := f.birds[0].brain

	for c := range f.birds {
		if f.birds[c].dead {
			f.birds[c].dead = false
			f.birds[c].bird.Pos = *f.lvl.NewBirdPos()
			f.birds[c].bird.Vel = *util.NewVector(SCROLL_SPEED, 0)

			f.birds[c].brain = neural.Cross(best, randNet())

			if neural.Chance(0.1) {
				// penalize best achievement due to mutation
				f.birds[c].bestX *= 0.99
				f.birds[c].brain.Mutate(0.33)
			}
		} else {
			f.birds[c].bird.Pos = f.birds[c].bird.NextPos
			f.birds[c].bestX = math.Max(f.birds[c].bird.Pos.X, f.birds[c].bestX)
		}
	}

}
