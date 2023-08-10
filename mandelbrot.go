package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

var (
	bkg                                            sdl.Color
	title                                          string = "SDL2 Window"
	width                                          int32  = 800
	height                                         int32  = 800
	renderer                                       *sdl.Renderer
	window                                         *sdl.Window
	frameCount, timerFPS, lastFrame, fps, lastTime uint32
	setFPS                                         uint32 = 60
	mouse                                          sdl.Point
	mousestate                                     uint32
	keystates                                      = sdl.GetKeyboardState()
	event                                          sdl.Event
	running                                        bool

	ms_limit    int = 16
	ms_infinity int = 1

	ms_zoomSpeed float64 = 0.1

	x_lborder float64 = -2.0
	x_rborder float64 = 1.0
	x_delta   float64 = 3.0

	y_lborder float64 = -1.0
	y_rborder float64 = 1.0
	y_delta   float64 = 2.0
)

func setColor(r, g, b, a uint8) sdl.Color {
	var c sdl.Color
	c.R = r
	c.G = g
	c.B = b
	c.A = a
	return c
}

func start() {
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "0")
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		panic(err)
	}
	running = true
}

func startSet(t string, w int32, h int32) {
	title = t
	width = w
	height = h
	start()
}

func quit() {
	running = false
	window.Destroy()
	renderer.Destroy()
	sdl.Quit()
}

func loop() {
	lastFrame = sdl.GetTicks()
	if lastFrame >= (lastTime + 1000) {
		lastTime = lastFrame
		fps = frameCount
		frameCount = 0
	}
	input()
}

func input() {
	keystates = sdl.GetKeyboardState()
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			running = false
			break
		}
	}
	mouse.X, mouse.Y, mousestate = sdl.GetMouseState()
}

func beginRender() {
	renderer.SetDrawColor(bkg.R, bkg.G, bkg.B, bkg.A)
	renderer.Clear()
	for x := 0; x < int(width); x++ {
		for y := 0; y < int(width); y++ {
			//y {-1; 1}
			//x {-1; 2}

			var x_bit float64 = x_delta / float64(width)
			var y_bit float64 = y_delta / float64(width)

			c := complex(x_lborder+(float64(x)*x_bit), y_lborder+(float64(y)*y_bit))

			z := complex(0, 0) + c

			for j := 0; j < ms_limit; j++ {
				z = (z * z) + c
			}

			realZ := real(z) + imag(c)
			var bwColor uint8 = uint8(math.Ceil((realZ / float64(ms_infinity)) * 256))

			renderer.SetDrawColor(bwColor, bwColor, bwColor, 1)
			renderer.DrawPoint(int32(x), int32(y))

		}
	}
	frameCount++
	timerFPS = sdl.GetTicks() - lastFrame
	if timerFPS < (1000 / setFPS) {
		sdl.Delay((1000 / setFPS) - timerFPS)
	}
	renderer.SetDrawColor(0, 0, 0, 255)
}

func endRender() {
	renderer.Present()
}

func main() {
	startSet(title, width, width)
	bkg = setColor(255, 0, 255, 255)
	for running {
		loop()
		beginRender()
		//renderer.DrawRect()
		endRender()
		if keystates[sdl.SCANCODE_ESCAPE] != 0 {
			running = false
		}

		if keystates[sdl.SCANCODE_LEFT] != 0 {
			x_lborder = x_lborder - 0.25
			x_rborder = x_rborder - 0.25
		}
		if keystates[sdl.SCANCODE_RIGHT] != 0 {
			x_lborder = x_lborder + 0.25
			x_rborder = x_rborder + 0.25
		}
		if keystates[sdl.SCANCODE_UP] != 0 {
			x_lborder = y_lborder - 0.25
			x_rborder = y_rborder - 0.25
		}
		if keystates[sdl.SCANCODE_DOWN] != 0 {
			x_lborder = y_lborder + 0.25
			x_rborder = y_rborder + 0.25
		}

		if keystates[sdl.SCANCODE_Q] != 0 {
			ms_limit = ms_limit / 2
		}
		if keystates[sdl.SCANCODE_W] != 0 {
			ms_limit = ms_limit * 2
		}

		if keystates[sdl.SCANCODE_A] != 0 {
			ms_infinity = ms_infinity / 2
		}
		if keystates[sdl.SCANCODE_S] != 0 {
			ms_infinity = ms_infinity * 2
		}

		if keystates[sdl.SCANCODE_X] != 0 {
			x_lborder = x_lborder + ms_zoomSpeed //(x_delta / 2)
			x_rborder = x_rborder - ms_zoomSpeed //(x_delta / 2)
			x_delta = x_rborder - x_lborder

			y_lborder = y_lborder + ms_zoomSpeed //(y_delta / 2)
			y_rborder = y_rborder - ms_zoomSpeed //(y_delta / 2)
			y_delta = y_rborder - y_lborder

		}
		if keystates[sdl.SCANCODE_Z] != 0 {
			x_lborder = x_lborder - ms_zoomSpeed //(x_delta / 2)
			x_rborder = x_rborder + ms_zoomSpeed //(x_delta / 2)
			x_delta = x_rborder - x_lborder

			y_lborder = y_lborder - ms_zoomSpeed //(y_delta / 2)
			y_rborder = y_rborder + ms_zoomSpeed //(y_delta / 2)
			y_delta = y_rborder - y_lborder
		}
	}
	quit()
}
