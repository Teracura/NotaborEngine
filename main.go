package main

import (
	"NotaborEngine/notacore"
	"NotaborEngine/notagl"
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
	"os"
	"path/filepath"
	"runtime"
	"time"

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
		MaxHz:     60,
		Runnables: nil,
		LastTime:  time.Time{},
	}

	logicLoop := &notacore.FixedHzLoop{
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
		LogicLoops: []*notacore.FixedHzLoop{logicLoop},
		RenderLoop: renderLoop,
	}

	window, err := wm.Create(cfg)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	backend := notagl.GLBackend2D{}
	backend.Init()

	notashader.Shaders["basic2d"] = notashader.CreateProgram(notashader.Vertex2D, notashader.Fragment2D).Type

	renderer := notagl.Renderer2D{}

	addRunnables(logicLoop, renderLoop, &renderer)

	err = notacore.Run(wm, window, cfg, &renderer, &backend)
	if err != nil {
		panic(err)
	}
}

func addRunnables(logicLoop *notacore.FixedHzLoop, renderLoop *notacore.RenderLoop, renderer *notagl.Renderer2D) {

	rect := &notagl.Rect{
		Center:    notamath.Po2{X: 0, Y: 0},
		W:         0.5,
		H:         0.5,
		Transform: notamath.NewTransform2D(),
	}

	// Add logic runnable (no snapshot needed - handled by synchronizer)
	logicLoop.Runnables = append(logicLoop.Runnables, func() error {
		rect.Transform.Snapshot()
		rect.Transform.RotateBy(0.01)
		return nil
	})

	// Add render runnable (no channel waiting - handled by synchronizer)
	renderLoop.Runnables = append(renderLoop.Runnables, func() error {
		gl.UseProgram(notashader.Shaders["basic2d"])

		now := time.Now()
		alpha := logicLoop.Alpha(now)

		renderer.Submit(rect, alpha)
		return nil
	})
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
