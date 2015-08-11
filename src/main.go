package main

import (
	neural "./neural"
	problems "./problems"
	"fmt"
	sdl "github.com/veandco/go-sdl2/sdl"
	"math"
	"runtime"
	"sort"
	"time"
)

type FBird struct {
	bird  *problems.Bird
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

// will check if going from pos to next will collide
func checkFlock(birds Flock, lvl *problems.Level) {

	collide := func(aX, bX, cX float64) bool {
		// c.X == d.X
		return aX-1 <= cX && bX+1 >= cX
	}

	hSize := float64(problems.PylonHole / 2)

	for c := range birds {
		if birds[c].bird.Pos.Y >= lvl.GetSize().Y || birds[c].bird.Pos.Y < 1 {
			// hit ceeling or floor
			birds[c].dead = true
			continue
		}

		pylon := lvl.ClosestPylon(&birds[c].bird.Pos)
		if birds[c].bird.Pos.Y >= pylon.Y-hSize && birds[c].bird.Pos.Y <= pylon.Y+hSize {
			// in the pylon hole
			continue
		}

		if birds[c].bird.Pos.Y > pylon.Y {
			// bottom pylon segment
			birds[c].dead = collide(birds[c].bird.Pos.X, birds[c].bird.NextPos.X, pylon.X)
		} else {
			// top pylon segment
			birds[c].dead = collide(birds[c].bird.Pos.X, birds[c].bird.NextPos.X, pylon.X)
		}
	}

}

func thnikFlock(birds Flock, lvl *problems.Level) {
	wg := make(chan struct{}, len(birds))

	thinkBird := func(c int) {
		birds[c].bestX = math.Max(birds[c].bird.Pos.X, birds[c].bestX)
		next := lvl.FirstPylonAfter(&birds[c].bird.Pos)
		diffY := next.Y - birds[c].bird.Pos.Y
		diffX := next.X - birds[c].bird.Pos.X

		birds[c].brain.Stimulate(0, diffY)
		birds[c].brain.Stimulate(1, diffX)
		birds[c].brain.Stimulate(2, birds[c].bird.Vel.Y)

		birds[c].brain.Step()
		if birds[c].brain.ValueOf(7) > 0.75 {
			birds[c].bird.Vel.Y = -0.4
		}

		birds[c].brain.Clear()
		wg <- struct{}{}
	}

	for c := 0; c < len(birds); c++ {
		go thinkBird(c)
	}

	for c := 0; c < len(birds); c++ {
		<-wg
	}
}

func mutateFlock(birds Flock, lvl *problems.Level) {
	sort.Sort(birds)

	randNet := func() *neural.Net {
		return birds[int(neural.RandMax(float64(len(birds))))].brain
	}

	best := birds[0].brain

	// TODO move dead check out of this loop
	// TODO check if the bird jumps trough the pylon - kill
	for c := range birds {
		if birds[c].dead {
			birds[c].dead = false
			birds[c].bird.Pos = *lvl.NewBirdPos()
			birds[c].bird.Vel = *problems.NewVector(0.1, 0)

			birds[c].brain = neural.Cross(best, randNet())

			if neural.Chance(0.1) {
				birds[c].brain.Mutate(0.33)
			}
		} else {
			birds[c].bird.Pos = birds[c].bird.NextPos
		}
	}

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	doDraw := true

	W := 1500
	H := 800
	LVL_W := W * 50
	fmt.Println(W, H)
	var FPS float64 = 60.0

	lvl := problems.NewLevel(LVL_W, H)

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
	surface.FillRect(&clearRect, 0xffffffff)

	bcount := 10000

	nets := make([]*neural.Net, bcount, bcount)
	for c := range nets {
		nets[c] = neural.NewNet(8)

		// diffY- to hidden
		*nets[c].Synapse(0, 3) = 0.0
		*nets[c].Synapse(0, 4) = 0.0
		*nets[c].Synapse(0, 5) = 0.0
		*nets[c].Synapse(0, 6) = 0.0

		// diffX- to hidden
		*nets[c].Synapse(1, 3) = 0.0
		*nets[c].Synapse(1, 4) = 0.0
		*nets[c].Synapse(1, 5) = 0.0
		*nets[c].Synapse(1, 6) = 0.0

		// velY - to hidden
		*nets[c].Synapse(2, 3) = 0.0
		*nets[c].Synapse(2, 4) = 0.0
		*nets[c].Synapse(2, 5) = 0.0
		*nets[c].Synapse(2, 6) = 0.0

		// hidden to output
		*nets[c].Synapse(3, 7) = 0.0
		*nets[c].Synapse(4, 7) = 0.0
		*nets[c].Synapse(5, 7) = 0.0
		*nets[c].Synapse(6, 7) = 0.0

		nets[c].Randomize()
	}

	lvl.AddBirds(bcount)
	flock := make(Flock, bcount)
	for c := 0; c < bcount; c++ {
		flock[c].bird = (*lvl.GetBirds())[c]
		flock[c].brain = nets[c]
		flock[c].bestX = 0
	}

	offset := 0
	step := 65

	frame := 0
	var frameTime float64 = 1000 / FPS
	start := time.Now()
	for {

		frame++
		thnikFlock(flock, lvl)

		visible := func(x float64) bool {
			return x >= float64(offset) && x < float64(offset+W)
		}

		toScreen := func(x float64) float64 {
			return float64(x - float64(offset))
		}

		brds := *lvl.GetBirds()
		for _, brd := range brds {
			if !visible(brd.Pos.X) {
				continue
			}

			rect.X = int32(toScreen(brd.Pos.X))
			rect.Y = int32(brd.Pos.Y)
			rect.W = 5
			rect.H = 5
			if doDraw {
				surface.FillRect(&rect, 0xffff0000)
			}
		}

		hSize := float64(problems.PylonHole) / 2.0
		for _, pylon := range lvl.GetPylons() {
			if !visible(pylon.X) {
				continue
			}

			rect.X = int32(toScreen(pylon.X))
			rect.Y = int32(0)
			rect.W = 5

			// top part
			rect.H = int32(pylon.Y - hSize)
			if doDraw {
				surface.FillRect(&rect, 0xff00ff00)
			}

			// bottom part
			rect.Y = int32(pylon.Y + hSize)
			rect.H = int32(float64(H) - (pylon.Y + hSize))
			if doDraw {
				surface.FillRect(&rect, 0xff00ff00)
			}

			rect.Y = int32(pylon.Y)
			rect.W = 3
			rect.H = 3
			if doDraw {
				surface.FillRect(&rect, 0xff0000ff)
			}
		}

		elapsed := time.Since(start)

		if doDraw && frameTime < 1000.0/FPS {
			time.Sleep(time.Millisecond * time.Duration(1000.0/FPS))
		}

		start = time.Now()

		frameTime = frameTime*0.9 + float64(elapsed.Nanoseconds())*0.1

		if frame > 60 {
			frame = 0
			fmt.Printf("fps last: %s\tfps average %f\tcompletion %f%%\n", elapsed, frameTime/1000000.0, flock[0].bestX/float64(LVL_W)*100.0)
		}

		window.UpdateSurface()
		lvl.Step(float64(elapsed.Nanoseconds()) / 1000000.0)
		checkFlock(flock, lvl)
		// lvl.DoStep()

		if doDraw {
			surface.FillRect(&clearRect, 0xffffffff)
		}

		stop := false
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				stop = true
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case sdl.K_LEFT:
					offset = int(math.Max(0, float64(offset-step)))
				case sdl.K_RIGHT:
					offset = int(math.Min(float64(LVL_W-W), float64(offset+step)))
				case sdl.K_SPACE:
					doDraw = !doDraw
					surface.FillRect(&clearRect, 0xffaaaaaa)
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
