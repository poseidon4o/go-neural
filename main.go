package main

import (
	"flag"
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
	StatsReportTick()
	Complete() float64
	Done() bool
	Jump()
	Move(int)
	SaveNetsToFile()
	LoadNetsFromFile()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	agentsCountPtr := flag.Int("agents", 10000, "number of agents alive at a time")
	levelWidthPtr := flag.Int("size", 20, "x in level_size = 1300 * x")

	flag.Parse()

	fmt.Println("Controls:")
	fmt.Println("end:\tfurthest action in the level")
	fmt.Println("home:\tmove back to level begining")
	fmt.Println("left:\tmove screen to the left")
	fmt.Println("right:\tmove screen to the right")
	fmt.Println("1:\tswitch to flappy")
	fmt.Println("2:\tswitch to mario")
	fmt.Println("s:\tsave all nets to file")
	fmt.Println("f:\tturn on fast-forward mode")
	fmt.Println("p:\ttake screenshot in format \"SS-#d.bmp\", where #d is [0,inf)")
	fmt.Println("enter:\tcycle trough mario/flappy")
	fmt.Println("esc:\teixt")
	fmt.Println("")

	doDraw := true
	doDev := true
	doFastForward := false

	W := 1300
	H := 700
	LVL_W := W * *levelWidthPtr

	var FPS float64 = 60.0
	FRAME_TIME_MS := 1000 / FPS
	FRAME_TIME := time.Millisecond * time.Duration(FRAME_TIME_MS)

	sdl.Init(sdl.INIT_EVERYTHING)

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(W), int32(H), sdl.WINDOW_SHOWN)
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

	figCount := *agentsCountPtr
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
	mr := mario.NewMario(figCount, util.NewVector(float64(LVL_W), float64(250)))

	mr.LoadNetsFromFile()
	fl.LoadNetsFromFile()

	fl.SetDrawRectCb(drawCb)
	mr.SetDrawRectCb(drawCb)

	imageCounter := 0

	var game DrawableProblem = fl

	step := 65

	frame := 0

	toggleKeys := map[sdl.Keycode]bool {
		sdl.K_f: true,
		sdl.K_s: true,
		sdl.K_p: true,
		sdl.K_RETURN: true,
	}

	var averageFrameTime float64 = FRAME_TIME_MS * 1000000 // in nanosec
	start := time.Now()
	lastDrawTime := time.Now()
	lastReportTime := time.Now()
	loops := 0
	for {
		start = time.Now()
		loops++

		if doDraw {
			if doFastForward {
				if !lastDrawTime.Add(FRAME_TIME).After(start) {
					lastDrawTime = start
					window.UpdateSurface()
					frame++
					surface.FillRect(&clearRect, 0xffffffff)
					game.DrawTick()
				}
			} else {
				lastDrawTime = start
				window.UpdateSurface()
				frame++
				surface.FillRect(&clearRect, 0xffffffff)
				game.DrawTick()
			}

		} else if loops%10 == 0 {
			// update only 10% of the frames
			window.UpdateSurface()
			frame++
		}

		game.LogicTick(1 / FPS)

		stop := false

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				stop = true
			case *sdl.KeyboardEvent:
				code := t.Keysym.Sym
				if (t.Repeat != 0 || t.Type != sdl.KEYUP) && toggleKeys[code] {
					continue
				}

				switch code {
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
					window.SetSize(int32(W), 700)
					surface, _ = window.GetSurface()
				case sdl.K_2:
					game = mr
					window.SetSize(int32(W), 250)
					surface, _ = window.GetSurface()
				case sdl.K_RETURN:
					if game == fl {
						game = mr
						window.SetSize(int32(W), 250)
						surface, _ = window.GetSurface()
					} else {
						game = fl
						window.SetSize(int32(W), 700)
						surface, _ = window.GetSurface()
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
				case sdl.K_p:
					h := int32(H)
					if game == mr {
						h = 250
					}
					w := int32(float64(LVL_W) * (game.Complete() + 0.001))

					_, rmask, gmask, bmask, amask, err := sdl.PixelFormatEnumToMasks(sdl.PIXELFORMAT_RGB24)
					image, err := sdl.CreateRGBSurface(0, w, h, 24, rmask, gmask, bmask, amask)

					if err == nil {
						name := fmt.Sprintf("SS-%d.bmp", imageCounter)
						imageCounter++

						rect := sdl.Rect{0, 0, w, h}
						image.FillRect(&rect, 0xffffffff)
						game.SetDrawRectCb(func(pos, size *util.Vector, color uint32) {
							rect.X = int32(pos.X)
							rect.Y = int32(pos.Y)
							rect.W = int32(size.X)
							rect.H = int32(size.Y)
							image.FillRect(&rect, color)
						})
						game.DrawTick()
						e := image.SaveBMP(name)
						if e == nil {
							fmt.Println("Saving image: ", name, "Size: ", w, "x", h)
						} else {
							fmt.Println("Failed to save image!", e)
						}
						game.SetDrawRectCb(drawCb)
					}
				}
			}
		}

		if stop {
			break
		}

		elapsed := time.Since(start)
		frameMs := float64(elapsed) / 1000000

		averageFrameTime = averageFrameTime*0.9 + float64(elapsed.Nanoseconds())*0.1

		if game.Done() {
			fmt.Println("Done")
			break
		}

		if !lastReportTime.Add(time.Second).After(start) {
			fmt.Println("")
			game.StatsReportTick()
			fmt.Printf("Last FrameTime: %f\tAverage FrameTime %f\tCompletion %f%%\n", frameMs, averageFrameTime/1000000, game.Complete()*100)
			fmt.Printf("FastForward %t\t Rand buffer refils %f\tFrames for last second %d\n", doFastForward, neural.BUFFER_REFILS, frame)

			lastReportTime = start
			frame = 0
		}

		// sleep only if drawing and there is time to sleep more than 3ms
		if !doFastForward && doDraw && frameMs < FRAME_TIME_MS && FRAME_TIME_MS-frameMs > 3.0 {
			time.Sleep(time.Millisecond * time.Duration(FRAME_TIME_MS-frameMs))
		}
	}

	sdl.Quit()
}
