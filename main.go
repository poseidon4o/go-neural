package main

import (
	"fmt"
	neural "github.com/poseidon4o/go-neural/src/neural"
	flappy "github.com/poseidon4o/go-neural/src/problems/flappy"
	mario "github.com/poseidon4o/go-neural/src/problems/mario"
	util "github.com/poseidon4o/go-neural/src/util"
	sdl "github.com/veandco/go-sdl2/sdl"
	"math"
	"runtime"
	"time"
)

type DrawableProblem interface {
	SetDrawRectCb(cb func(pos, size *util.Vector, color uint32))
	LogicTick(dt float64)
	DrawTick()
	Complete() float64
	Done() bool
	Jump()
	Move(int)
	SaveNetsToFile()
	LoadNetsFromFile()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Controls:")
	fmt.Println("end:\tfurthest action in the level")
	fmt.Println("home:\tmove back to level begining")
	fmt.Println("left:\tmove screen to the left")
	fmt.Println("right:\tmove screen to the right")
	fmt.Println("1:\tswitch to flappy")
	fmt.Println("2:\tswitch to mario")
	fmt.Println("enter:\tcycle trough mario/flappy")
	fmt.Println("esc:\teixt")
	fmt.Println("")

	doDraw := true
	doDev := false
	doFastForward := false

	W := 1300
	H := 700
	LVL_W := W * 50

	var FPS float64 = 60.0
	FRAME_TIME_MS := 1000 / FPS

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

	clearRect := sdl.Rect{0, 0, int32(W), int32(H)}
	surface.FillRect(&clearRect, 0xffffffff)

	figCount := 100
	if doDev {
		figCount = 1
	}

	offset := 0
	visible := func(pos, size *util.Vector) bool {

		// r1 := sdl.Rect{int32(pos.X), int32(pos.Y), int32(pos.X + size.X), int32(pos.Y + size.Y)}
		// r2 := sdl.Rect{int32(offset), 0, int32(offset + W), int32(H)}
		// return !(r2.X > r1.W || r2.W < r1.X || r2.Y > r1.H || r2.H < r1.Y)

		// so beautiful
		return !(float64(offset) > pos.X+size.X ||
			float64(offset+W) < pos.X ||
			0 > pos.Y+size.Y ||
			float64(H) < pos.Y)
	}

	toScreen := func(pos, size *util.Vector) *sdl.Rect {
		return &sdl.Rect{
			X: int32(pos.X - float64(offset)),
			Y: int32(pos.Y),
			W: int32(size.X),
			H: int32(size.Y),
		}
	}

	drawCb := func(pos, size *util.Vector, color uint32) {
		if visible(pos, size) {
			surface.FillRect(toScreen(pos, size), color)
		}
	}

	fl := flappy.NewFlappy(figCount, util.NewVector(float64(LVL_W), float64(H)))
	mr := mario.NewMario(figCount, util.NewVector(float64(LVL_W), float64(H)))

	mr.LoadNetsFromFile()
	fl.LoadNetsFromFile()

	fl.SetDrawRectCb(drawCb)
	mr.SetDrawRectCb(drawCb)

	var game DrawableProblem = fl

	step := 65

	frame := 0
	var averageFrameTime float64 = FRAME_TIME_MS * 1000000 // in nanosec
	start := time.Now()
	for {
		start = time.Now()

		if doDraw {
			window.UpdateSurface()
			surface.FillRect(&clearRect, 0xffffffff)
			game.DrawTick()
		} else if frame%10 == 0 {
			// update only 10% of the frames
			window.UpdateSurface()
		}

		game.LogicTick(1 / FPS)

		stop := false
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				stop = true
			case *sdl.KeyDownEvent:
				switch t.Keysym.Sym {
				case sdl.K_LEFT:
					if doDev {
						game.Move(-1)
					} else {
						offset = int(math.Max(0, float64(offset-step)))
					}
				case sdl.K_RIGHT:
					if doDev {
						game.Move(1)
					} else {
						offset = int(math.Min(float64(LVL_W-W), float64(offset+step)))
					}
				case sdl.K_SPACE:
					if doDev {
						game.Jump()
					} else {
						doDraw = !doDraw
						surface.FillRect(&clearRect, 0xffaaaaaa)
						window.UpdateSurface()
					}
				case sdl.K_1:
					game = fl
				case sdl.K_2:
					game = mr
				case sdl.K_RETURN:
					if game == fl {
						game = mr
					} else {
						game = fl
					}
				case sdl.K_ESCAPE:
					stop = true
				case sdl.K_END:
					offset = int(math.Max(math.Min(float64(LVL_W-W), game.Complete()*float64(LVL_W)-float64(W)/2), 0))
				case sdl.K_HOME:
					offset = 0
				case sdl.K_f:
					doFastForward = !doFastForward
				case sdl.K_s:
					game.SaveNetsToFile()
				}
			}
		}

		if stop {
			break
		}

		frame++

		elapsed := time.Since(start)
		frameMs := float64(elapsed) / 1000000

		averageFrameTime = averageFrameTime*0.9 + float64(elapsed.Nanoseconds())*0.1

		if game.Done() {
			fmt.Println("Done")
			break
		}

		if frame > int(FPS) {
			frame = 0
			fmt.Printf("CHRand %d\tGRand %d\tG/C %f\n", neural.ChanRand, neural.GlobRand, float64(neural.GlobRand)/float64(neural.ChanRand))
			neural.ChanRand = 0
			neural.GlobRand = 0
			fmt.Printf("ftime last: %f\tftime average %f\tcompletion %f%%\n", frameMs, averageFrameTime/1000000, game.Complete()*100)
		}

		// sleep only if drawing and there is time to sleep more than 3ms
		if !doFastForward && doDraw && frameMs < FRAME_TIME_MS && FRAME_TIME_MS-frameMs > 3.0 {
			time.Sleep(time.Millisecond * time.Duration(FRAME_TIME_MS-frameMs))
		}
	}

	sdl.Quit()
}
