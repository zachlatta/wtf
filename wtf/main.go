package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

const (
	Title         string = "Where's the Refrigerator?"
	Width, Height int    = 640, 480
	FPSLimit      int64  = 60
)

var (
	texture  gl.Texture
	delta    time.Duration
	lastTime time.Time
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(Width, Height, Title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	gl.Init()

	if err := initScene(); err != nil {
		fmt.Fprintf(os.Stderr, "init: %s\n", err)
		return
	}
	defer destroyScene()

	for !window.ShouldClose() {
		updateDelta()
		draw(window)
		waitToLimitFps()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func initScene() (err error) {
	gl.Enable(gl.TEXTURE_2D)

	gl.ClearColor(0.0, 0.0, 0.0, 0.0)

	gl.Viewport(0, 0, Width, Height)
	gl.Ortho(0, float64(Width), float64(Height), 0, -10, 10)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()

	doge, err := os.Open("doge-freezer.png")
	if err != nil {
		panic(err)
	}
	defer doge.Close()

	texture, err = createTexture(doge)
	return
}

func destroyScene() {
	texture.Delete()
}

var x, y float32 = 0, 0
var w, h float32 = 526.0, 526.0

func updateDelta() {
	delta = time.Now().Sub(lastTime)
	lastTime = time.Now()
}

func draw(window *glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.PushMatrix()
	gl.Translatef(x, y, 0)

	texture.Bind(gl.TEXTURE_2D)

	gl.Begin(gl.QUADS)
	gl.Color4f(1, 1, 1, 1)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(0, 1)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(0, h)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(w, h)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(w, 0)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(w, 0)
	gl.End()

	gl.PopMatrix()

	increment := float32(1 * delta.Seconds())
	if window.GetKey(glfw.KeyW) == glfw.Press {
		y += increment
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		x -= increment
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		y -= increment
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		x += increment
	}
}

func waitToLimitFps() {
	frameTimeTarget := time.Second / time.Duration(FPSLimit)
	time.Sleep(frameTimeTarget - delta)
}
