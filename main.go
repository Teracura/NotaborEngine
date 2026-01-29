// cmd/test_textured/main.go
package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"NotaborEngine/notacore"
	"NotaborEngine/notagl"
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
)

func main() {
	engine := &notacore.Engine{
		Settings: &notacore.Settings{Vsync: true},
	}

	if err := engine.InitPlatform(); err != nil {
		log.Fatal("Init failed:", err)
	}
	defer engine.Shutdown()

	renderLoop := &notacore.RenderLoop{
		MaxHz: 60,
	}

	logicLoop := &notacore.FixedHzLoop{Hz: 500}

	cfg := notacore.WindowConfig{
		X:          50,
		Y:          50,
		W:          800,
		H:          600,
		Title:      "Textured Rendering Test",
		Type:       notacore.Windowed,
		Resizable:  true,
		RenderLoop: renderLoop,
		LogicLoops: []*notacore.FixedHzLoop{logicLoop},
	}

	win, err := engine.CreateWindow2D(cfg)
	if err != nil {
		log.Fatal("Window creation failed:", err)
	}

	currentDir, _ := filepath.Abs(".")
	texturePath := filepath.Join(currentDir, "resources", "hahaha.jpg")
	texture, err := win.LoadTexture("test", texturePath)
	if err != nil {
		log.Fatal("Failed to load texture:", err)
	}

	texturedShader := notashader.Shader{
		Name:           "textured",
		VertexString:   notashader.TexturedVertex2D,
		FragmentString: notashader.TexturedFragment2D,
	}

	if err := win.CreateShader(texturedShader); err != nil {
		log.Fatal("Failed to create shader:", err)
	}

	if err := win.UseShader("textured"); err != nil {
		log.Fatal("Failed to use shader:", err)
	}
	quad := notagl.CreateTextureQuad(notamath.Po2{X: 0, Y: 0}, 2, 1)

	renderLoop.Runnables = []notacore.Runnable{
		func() error {
			texture.Bind(0)
			alpha := logicLoop.Alpha(time.Now())
			win.RunTime.Renderer.Submit(quad, alpha)
			return nil
		},
	}

	logicLoop.Runnables = []notacore.Runnable{
		func() error {
			quad.Transform.Snapshot()
			quad.Transform.RotateBy(0.01)
			return nil
		},
	}

	fmt.Println("Running textured rendering test...")
	if err := engine.Run(); err != nil {
		log.Fatal("Engine run failed:", err)
	}
}
