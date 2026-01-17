package main

import (
	"github.com/go-gl/gl/v4.6-core/gl"

	"NotaborEngine/notacore"
	"NotaborEngine/notagl"
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
	"runtime"
	"time"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	engine := &notacore.Engine{
		Settings: &notacore.Settings{},
	}

	if err := engine.InitPlatform(); err != nil {
		panic(err)
	}
	defer engine.Shutdown()

	renderLoop1 := &notacore.RenderLoop{MaxHz: 60}
	logicLoop1 := &notacore.FixedHzLoop{Hz: 240}

	cfg1 := notacore.WindowConfig{
		X:          100,
		Y:          100,
		W:          800,
		H:          600,
		Title:      "Test Window 1",
		Resizable:  true,
		Type:       notacore.Windowed,
		LogicLoops: []*notacore.FixedHzLoop{logicLoop1},
		RenderLoop: renderLoop1,
	}
	win1, err := engine.CreateWindow2D(cfg1)
	if err != nil {
		panic(err)
	}
	win1.MakeContextCurrent()
	shader1 := notashader.CreateProgram(
		notashader.Vertex2D,
		notashader.Fragment2D,
	).Type

	addRunnables(win1, shader1)

	if err := engine.Run(); err != nil {
		panic(err)
	}
}

func addRunnables(win *notacore.GlfwWindow2D, shaderProgram uint32) {
	rect := notagl.Polygon{
		Vertices: []notamath.Po2{
			{0, 0},
			{1, 0},
			{1, 1},
			{0, 1},
			{-1, -1},
		},
		Transform: notamath.NewTransform2D(),
	}
	rect.Fixate()

	logicLoop := win.Config.LogicLoops[0]
	renderLoop := win.Config.RenderLoop
	backend := win.RunTime.Backend
	renderer := win.RunTime.Renderer

	logicLoop.Runnables = append(logicLoop.Runnables, func() error {
		rect.Transform.Snapshot()
		rect.Transform.RotateBy(0.01)
		return nil
	})

	renderLoop.Runnables = append(renderLoop.Runnables, func() error {
		gl.UseProgram(shaderProgram)
		alpha := logicLoop.Alpha(time.Now())

		renderer.Begin()
		renderer.Submit(rect, alpha)
		renderer.Flush(backend)
		return nil
	})
}
