package main

import (
	"NotaborEngine/notacore"
	"NotaborEngine/notagl"
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
	"fmt"
	"runtime"
	"time"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	engine := &notacore.Engine{
		Settings: &notacore.Settings{},
	}

	if err := engine.InitPlatform(); err != nil {
		return err
	}
	defer engine.Shutdown()

	renderLoop := &notacore.RenderLoop{MaxHz: 60}
	logicLoop := &notacore.FixedHzLoop{Hz: 240}

	cfg := notacore.WindowConfig{
		X:          350,
		Y:          50,
		W:          800,
		H:          800,
		Title:      "Test Window 1",
		Resizable:  true,
		Type:       notacore.Windowed,
		LogicLoops: []*notacore.FixedHzLoop{logicLoop},
		RenderLoop: renderLoop,
	}

	win, err := engine.CreateWindow2D(cfg)
	if err != nil {
		return err
	}
	win.MakeContextCurrent()

	if err := win.CreateShader(notashader.Shader{
		Name:           "circle2D",
		VertexString:   notashader.Circle2DVertex,
		FragmentString: notashader.Circle2DFragment,
	}); err != nil {
		return fmt.Errorf("create shader circle2D: %w", err)
	}

	if err := win.CreateShader(notashader.Shader{
		Name:           "basic2D",
		VertexString:   notashader.Vertex2D,
		FragmentString: notashader.Fragment2D,
	}); err != nil {
		return fmt.Errorf("create shader basic2D: %w", err)
	}

	addRunnables(win)

	if err := engine.Run(); err != nil {
		return err
	}
	return nil
}

func addRunnables(win *notacore.GlfwWindow2D) {
	rect := notagl.Polygon{
		Vertices: []notamath.Po2{
			{-0.5, -0.5},
			{0.5, -0.5},
			{0.5, 0.5},
			{-0.5, 0.5},
		},
		Transform: notamath.NewTransform2D(),
		Colors: []notashader.Color{
			notashader.White,
			notashader.Red,
			notashader.Purple,
			notashader.White,
		},
	}
	rect.Fixate()

	logicLoop := win.Config.LogicLoops[0]
	renderLoop := win.Config.RenderLoop
	renderer := win.RunTime.Renderer

	logicLoop.Runnables = append(logicLoop.Runnables, func() error {
		rect.Transform.Snapshot()
		rect.Transform.RotateBy(0.01)
		return nil
	})

	renderLoop.Runnables = append(renderLoop.Runnables, func() error {
		if err := win.UseShader("circle2D"); err != nil {
			return err
		}

		alpha := logicLoop.Alpha(time.Now())
		renderer.Submit(rect, alpha)
		return nil
	})
}
