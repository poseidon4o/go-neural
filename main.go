package main

import (
	"fmt"
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
	Completed() float64
	Done() bool
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Controls:")
	fmt.Println("end:\tfurthest action in the level")
	fmt.Println("home:\tmove back to level begining")
	fmt.Println("left:\tmove screen to the left")
	fmt.Println("right:\tmove screen to the right")
	fmt.Println("esc:\teixt")
	fmt.Println("")

	doDraw := true

	W := 1500
	H := 800
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

	g := flappy.NewFlappy(10000, util.NewVector(float64(LVL_W), float64(H)))
	g.Completed()
	game := mario.NewMario(10000, util.NewVector(float64(LVL_W), float64(H)))

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

	game.SetDrawRectCb(func(pos, size *util.Vector, color uint32) {
		if visible(pos, size) {
			surface.FillRect(toScreen(pos, size), color)
		}
	})

	step := 65

	frame := 0
	var averageFrameTime float64 = FRAME_TIME_MS * 1000000 // in nanosec
	start := time.Now()
	for {
		start = time.Now()

		game.LogicTick(1 / FPS)

		if doDraw {
			window.UpdateSurface()
			surface.FillRect(&clearRect, 0xffffffff)
			game.DrawTick()
		} else if frame%10 == 0 {
			// update only 10% of the frames
			window.UpdateSurface()
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
					window.UpdateSurface()
				case sdl.K_ESCAPE:
					stop = true
				case sdl.K_END:
					offset = int(math.Max(math.Min(float64(LVL_W-W), game.Completed()*float64(LVL_W)-float64(W)/2), 0))
				case sdl.K_HOME:
					offset = 0
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
			fmt.Printf("ftime last: %f\tftime average %f\tcompletion %f%%\n", frameMs, averageFrameTime/1000000, game.Completed()*100)
		}

		// sleep only if drawing and there is time to sleep more than 3ms
		if doDraw && frameMs < FRAME_TIME_MS && FRAME_TIME_MS-frameMs > 3.0 {
			time.Sleep(time.Millisecond * time.Duration(FRAME_TIME_MS-frameMs))
		}
	}

	sdl.Quit()
}