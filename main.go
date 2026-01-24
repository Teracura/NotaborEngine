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
		Input:    &notacore.InputManager{},
	}

	if err := engine.InitPlatform(); err != nil {
		return err
	}
	defer engine.Shutdown()

	renderLoop := &notacore.RenderLoop{MaxHz: 60}
	logicLoop := &notacore.FixedHzLoop{Hz: 1000}
	logicLoop.EnableMonitor(1 * time.Second)

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

	win, err := engine.CreateWindow3D(cfg)
	if err != nil {
		return err
	}
	win.MakeContextCurrent()

	if err := win.CreateShader(notashader.Shader{
		Name:           "basic3D",
		VertexString:   notashader.Vertex3D,
		FragmentString: notashader.Fragment3D,
	}); err != nil {
		return fmt.Errorf("create shader: %w", err)
	}
	addRunnables(engine, win)

	if err := engine.Run(); err != nil {
		return err
	}
	return nil
}

func addRunnables(engine *notacore.Engine, win *notacore.GlfwWindow3D) {
	cube := &notagl.Mesh{
		Vertices: []notamath.Po3{
			// Back face
			{-0.5, -0.5, -0.5}, {0.5, -0.5, -0.5}, {0.5, 0.5, -0.5},
			{-0.5, -0.5, -0.5}, {0.5, 0.5, -0.5}, {-0.5, 0.5, -0.5},

			// Front face
			{-0.5, -0.5, 0.5}, {0.5, 0.5, 0.5}, {0.5, -0.5, 0.5},
			{-0.5, -0.5, 0.5}, {-0.5, 0.5, 0.5}, {0.5, 0.5, 0.5},

			// Left face
			{-0.5, -0.5, -0.5}, {-0.5, 0.5, -0.5}, {-0.5, 0.5, 0.5},
			{-0.5, -0.5, -0.5}, {-0.5, 0.5, 0.5}, {-0.5, -0.5, 0.5},

			// Right face
			{0.5, -0.5, -0.5}, {0.5, 0.5, 0.5}, {0.5, 0.5, -0.5},
			{0.5, -0.5, -0.5}, {0.5, -0.5, 0.5}, {0.5, 0.5, 0.5},

			// Bottom face
			{-0.5, -0.5, -0.5}, {-0.5, -0.5, 0.5}, {0.5, -0.5, 0.5},
			{-0.5, -0.5, -0.5}, {0.5, -0.5, 0.5}, {0.5, -0.5, -0.5},

			// Top face
			{-0.5, 0.5, -0.5}, {0.5, 0.5, 0.5}, {-0.5, 0.5, 0.5},
			{-0.5, 0.5, -0.5}, {0.5, 0.5, -0.5}, {0.5, 0.5, 0.5},
		},
		Transform: notamath.NewTransform3D(),
		Colors: []notashader.Color{
			notashader.Red,
			notashader.Red,
			notashader.Red,
			notashader.Blue,
			notashader.Blue,
			notashader.Blue,
			notashader.Red,
			notashader.Red,
			notashader.Red,
			notashader.Blue,
			notashader.Blue,
			notashader.Blue,
			notashader.Navy,
			notashader.Navy,
			notashader.Navy,
			notashader.Blue,
			notashader.Navy,
			notashader.Navy,
			notashader.Maroon,
			notashader.Maroon,
			notashader.Maroon,
			notashader.Red,
			notashader.Maroon,
			notashader.Maroon,
			notashader.Olive,
			notashader.Olive,
			notashader.Olive,
			notashader.Olive,
			notashader.Olive,
			notashader.Olive,
			notashader.Purple,
			notashader.Purple,
			notashader.Purple,
			notashader.Purple,
			notashader.Purple,
			notashader.Purple,
		},
	}
	cube.Fixate()
	cube.SetDepthGradient(notashader.Blue, notashader.Magenta)

	logicLoop := win.Config.LogicLoops[0]
	renderLoop := win.Config.RenderLoop
	renderer := win.RunTime.Renderer

	var currentAxis = notamath.Vec3{X: 0.2, Y: 1, Z: 0.5}
	cube.Transform.RotationAxis = currentAxis

	aSignal := &notacore.InputSignal{}
	engine.Input.BindInput(notacore.KeyA, aSignal)

	aAction := &notacore.Action{
		Behavior: notacore.RunWhileHeld,
	}
	aAction.BindSignal(aSignal)
	aAction.AddRunnable(func() error {
		cube.Transform.Snapshot()
		cube.Transform.RotateBy(0.01)
		return nil
	})

	logicLoop.Runnables = append(logicLoop.Runnables, func() error {
		cube.Transform.Snapshot()
		if aAction.ShouldRun() {
			return aAction.Run()
		}
		return nil
	})

	renderLoop.Runnables = append(renderLoop.Runnables, func() error {
		if err := win.UseShader("basic3D"); err != nil {
			return err
		}

		alpha := logicLoop.Alpha(time.Now())
		renderer.Submit(cube, alpha)
		return nil
	})
}
