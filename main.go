package main

import (
	"NotaborEngine/notacore"
	"NotaborEngine/notagl"
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
	"fmt"
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

	if err = gl.Init(); err != nil {
		panic(err)
	}

	notashader.Shaders["basic2d"] = notashader.CreateProgram(notashader.Vertex2D, notashader.Fragment2D).Type

	renderer := notagl.Renderer2D{}
	backend := notagl.GLBackend2D{}
	backend.Init()

	addRunnables(logicLoop, cfg.RenderLoop, &renderer)

	for _, loop := range cfg.LogicLoops {
		loop.Start()
	}

	lastRenderTime := time.Now()
	for !window.ShouldClose() {
		now := time.Now()
		elapsed := now.Sub(lastRenderTime)

		targetFrameDuration := time.Second / time.Duration(renderLoop.MaxHz)

		if elapsed < targetFrameDuration {
			time.Sleep(targetFrameDuration - elapsed)
			continue
		}

		lastRenderTime = now
		wm.PollEvents()
		renderer.Begin()
		window.MakeContextCurrent()

		// Remove the timing check inside Render() since we handle it here
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, runnable := range renderLoop.Runnables {
			if err := runnable(); err != nil {
				fmt.Println("Render error:", err)
			}
		}

		renderer.Flush(&backend)
		window.SwapBuffers()
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
