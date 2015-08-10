package main

import (
	// neural "./neural"
	problems "./problems"
	sdl "github.com/veandco/go-sdl2/sdl"
	"time"
)

func main() {
	W := 900
	H := 500
	var FPS float64 = 60.0

	lvl := problems.NewLevel(W, H)
	lvl.AddBirds(1)

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

	for {
		rect.X = int32(lvl.GetBirds()[0].Pos().X)
		rect.Y = int32(lvl.GetBirds()[0].Pos().Y)
		surface.FillRect(&rect, 0xffff0000)

		time.Sleep(time.Millisecond * time.Duration(1000.0/FPS))
		window.UpdateSurface()
		lvl.Step(1000.0 / FPS)
		surface.FillRect(&clearRect, 0xffffffff)

		stop := false
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				stop = true
			case *sdl.KeyUpEvent:
				if t.Keysym.Sym == sdl.K_SPACE {
					lvl.GetBirds()[0].Vel().Y = -0.3
				}
			}
		}

		if stop {
			break
		}

	}

	sdl.Quit()
}
