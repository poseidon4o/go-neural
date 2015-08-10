package main

import (
	neural "./neural"
	problems "./problems"
	"fmt"
	sdl "github.com/veandco/go-sdl2/sdl"
	"math"
	"sort"
	"time"
)

type FBird struct {
	bird  *problems.Bird
	brain *neural.Net
	grade float64
	bestX float64
}

type Flock []FBird

func (birds Flock) Len() int {
	return len(birds)
}

func (birds Flock) Less(c, r int) bool {
	return birds[c].grade > birds[r].grade
}

func (birds Flock) Swap(c, r int) {
	birds[c], birds[r] = birds[r], birds[c]
}

func gradeFlock(birds Flock, lvl *problems.Level) {
	maxX := lvl.GetSize().X
	for c := range birds {
		birds[c].grade = birds[c].bestX / maxX
	}
}

func thnikFlock(birds Flock, lvl *problems.Level) {
	pylons := lvl.GetPylons()

	nextPylon := func(from *problems.Vector) problems.Vector {
		for _, p := range pylons {
			if p.X < from.X-1 {
				continue
			} else {
				// first pylon after from
				return p
			}
		}
		return *problems.NewVector(0, 0)
	}

	for c := range birds {
		birds[c].bestX = math.Max(birds[c].bird.Pos().X, birds[c].bestX)
		next := nextPylon(birds[c].bird.Pos())
		diff := next.Y - birds[c].bird.Pos().Y
		birds[c].brain.Stimulate(0, diff)

		birds[c].brain.Step()
		if birds[c].brain.ValueOf(5) > 0.75 {
			birds[c].bird.Vel().Y -= 0.1
		}

		birds[c].brain.Clear()
	}
}

func mutateFlock(birds Flock, lvl *problems.Level) {
	h := lvl.GetSize().Y
	sort.Sort(birds)

	randNet := func() *neural.Net {
		return birds[int(neural.RandMax(float64(len(birds))))].brain
	}
	hSize := float64(problems.PylonHole / 2)
	hitsPylon := func(pos, pyl *problems.Vector) bool {
		hits := true
		// in range of pylon
		hits = hits && (pos.X >= pyl.X-1 && pos.X <= pyl.X+1)
		// not in hole
		hits = hits && (pos.Y < pyl.Y-hSize || pos.Y > pyl.Y+hSize)
		return hits
	}

	best := birds[0].brain

	for c := range birds {
		brd := &birds[c]
		pos := brd.bird.Pos()
		isDead := pos.Y >= h || pos.Y < 10

		if !isDead {
			for _, p := range lvl.GetPylons() {
				if hitsPylon(pos, &p) {
					isDead = true
					break
				}
			}
		}

		if isDead {
			*pos = *lvl.NewBirdPos()
			*brd.bird.Vel() = *problems.NewVector(0.1, 0)

			brd.brain = neural.Cross(best, randNet())

			if neural.Chance(0.1) {
				brd.brain.Mutate(0.33)
			}
		}

	}

}

func main() {
	W := 1500
	H := 800
	fmt.Println(W, H)
	var FPS float64 = 60.0

	lvl := problems.NewLevel(W, H)

	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int(W), int(H), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	rect := sdl.Rect{0, 0, 5, 5}
	clearRect := sdl.Rect{0, 0, int32(W), int32(H)}

	bcount := 10

	nets := make([]*neural.Net, bcount, bcount)
	for c := range nets {
		nets[c] = neural.NewNet(7)

		// input 0 - to hidden
		*nets[c].Synapse(0, 1) = 0.0
		*nets[c].Synapse(0, 2) = 0.0
		*nets[c].Synapse(0, 3) = 0.0
		*nets[c].Synapse(0, 4) = 0.0

		// hidden to output
		*nets[c].Synapse(1, 5) = 0.0
		*nets[c].Synapse(2, 5) = 0.0
		*nets[c].Synapse(3, 5) = 0.0
		*nets[c].Synapse(4, 5) = 0.0

		nets[c].Randomize()
	}

	lvl.AddBirds(bcount)
	flock := make(Flock, bcount)
	for c := 0; c < bcount; c++ {
		flock[c].bird = (*lvl.GetBirds())[c]
		flock[c].brain = nets[c]
		flock[c].grade = 0
	}

	for {
		thnikFlock(flock, lvl)
		gradeFlock(flock, lvl)

		brds := *lvl.GetBirds()
		for _, brd := range brds {
			rect.X = int32(brd.Pos().X)
			rect.Y = int32(brd.Pos().Y)
			rect.W = 5
			rect.H = 5
			surface.FillRect(&rect, 0xffff0000)
		}

		hSize := float64(problems.PylonHole) / 2.0
		for _, pylon := range lvl.GetPylons() {
			rect.X = int32(pylon.X)
			rect.Y = int32(0)
			rect.W = 5

			// top part
			rect.H = int32(pylon.Y - hSize)
			surface.FillRect(&rect, 0xff00ff00)

			// bottom part
			rect.Y = int32(pylon.Y + hSize)
			rect.H = int32(float64(H) - (pylon.Y + hSize))
			surface.FillRect(&rect, 0xff00ff00)
		}

		time.Sleep(time.Millisecond * time.Duration(1000.0/FPS))
		window.UpdateSurface()
		lvl.Step(1000.0 / FPS)
		surface.FillRect(&clearRect, 0xffffffff)

		stop := false
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				stop = true
			case *sdl.KeyDownEvent:
				if t.Keysym.Sym == sdl.K_SPACE {
					// (*lvl.GetBirds())[0].Vel().Y = -0.3
				}
			}
		}

		if stop {
			break
		}

		mutateFlock(flock, lvl)
	}

	sdl.Quit()
}
