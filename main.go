package main

import (
	"NotaborEngine/notacore"
	"NotaborEngine/notagl"
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
	"os"
	"path/filepath"
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	addNativeDLLPath()

	wm, err := notacore.NewGLFWWindowManager()

	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	renderLoop := &notacore.RenderLoop{
		MaxHz: 60,
	}

	logicLoop := notacore.FixedHzLoop{
		Hz:        240,
		Runnables: nil,
	}

	cfg := notacore.WindowConfig{
		X:          100,
		Y:          100,
		W:          800,
		H:          600,
		Title:      "Test Window",
		Resizable:  true,
		Type:       notacore.Windowed,
		LogicLoops: nil,
		RenderLoop: renderLoop,
	}

	window, err := wm.Create(cfg)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	notashader.Shaders["basic2d"] = notashader.CreateProgram(notashader.Vertex2D, notashader.Fragment2D).Type

	renderer := notagl.Renderer2D{}

	backend := notagl.GLBackend2D{}
	backend.Init()

	rect := &notagl.Rect{
		Center:    notamath.Po2{X: 0, Y: 0},
		W:         0.5,
		H:         0.5,
		Transform: notamath.NewTransform2D(),
	}

	renderLoop.Runnables = append(
		renderLoop.Runnables,
		func() error {
			gl.UseProgram(notashader.Shaders["basic2d"])
			renderer.Submit(rect)
			return nil
		},
	)

	logicLoop.Runnables = append(logicLoop.Runnables, func() error {
		rect.Transform.RotateBy(0.01)
		return nil
	})

	logicLoop.Start()

	for !window.ShouldClose() {
		wm.PollEvents()

		glfw.WaitEventsTimeout(float64(1.0 / renderLoop.MaxHz))

		renderer.Begin()

		window.MakeContextCurrent()
		renderLoop.Render()

		renderer.Flush(&backend)

		window.SwapBuffers()
	}
}

func addNativeDLLPath() {
	switch runtime.GOOS {
	case "windows":
		exeDir, _ := os.Getwd()
		dllDir := filepath.Join(exeDir, "notacore", "native", "windows")
		_ = os.Setenv("PATH", dllDir+";"+os.Getenv("PATH"))

	case "linux":
		// set linux paths later

	case "darwin":
		// set mac paths later
	}
}
