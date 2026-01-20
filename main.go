package main

import (
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
		X:          0,
		Y:          0,
		W:          1000,
		H:          1000,
		Title:      "Test Window 1",
		Resizable:  true,
		Type:       notacore.Borderless,
		LogicLoops: []*notacore.FixedHzLoop{logicLoop1},
		RenderLoop: renderLoop1,
	}
	circle := notashader.Shader{
		Name:           "circle2D",
		VertexString:   notashader.Circle2DVertex,
		FragmentString: notashader.Circle2DFragment,
	}
	rectShader := notashader.Shader{
		Name:           "basic2D",
		VertexString:   notashader.Vertex2D,
		FragmentString: notashader.Fragment2D,
	}

	win1, err := engine.CreateWindow2D(cfg1)
	if err != nil {
		panic(err)
	}
	win1.MakeContextCurrent()
	err = win1.CreateShader(circle)
	err = win1.CreateShader(rectShader)
	if err != nil {
		panic(err)
	}

	addRunnables(win1)

	if err := engine.Run(); err != nil {
		panic(err)
	}
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
			notashader.Red,
			notashader.Yellow,
			notashader.Blue,
			notashader.Green,
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
		err := win.UseShader("circle2D")
		if err != nil {
			panic(err)
		}
		alpha := logicLoop.Alpha(time.Now())

		renderer.Submit(rect, alpha)
		return nil
	})
}
