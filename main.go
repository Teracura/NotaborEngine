// cmd/test_window_texture/main.go
package main

import (
	"log"
	"path/filepath"

	"NotaborEngine/notacore"
	"NotaborEngine/notagl"
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
)

func main() {
	// Create engine
	engine := &notacore.Engine{
		Settings: &notacore.Settings{Vsync: true},
	}

	// Initialize
	if err := engine.InitPlatform(); err != nil {
		log.Fatal("Init failed:", err)
	}
	defer engine.Shutdown()

	// Create render loop
	renderLoop := &notacore.RenderLoop{
		MaxHz: 60,
	}

	// Create window config
	cfg := notacore.WindowConfig{
		X:          50,
		Y:          50,
		W:          800,
		H:          600,
		Title:      "Texture Loading Test",
		Type:       notacore.Windowed,
		RenderLoop: renderLoop,
	}

	// Create window (this creates OpenGL context and calls gl.Init())
	win, err := engine.CreateWindow2D(cfg)
	if err != nil {
		log.Fatal("Window creation failed:", err)
	}

	log.Println("Window created, OpenGL context is ready")

	// NOW we can load textures (we have OpenGL context)
	currentDir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal("Failed to get current directory:", err)
	}

	texturePath := filepath.Join(currentDir, "resources", "hahaha.jpg")
	log.Printf("Attempting to load texture: %s", texturePath)

	// Load texture using window's method
	texture, err := win.LoadTexture("test_texture", texturePath)
	if err != nil {
		log.Printf("Failed to load texture: %v", err)
		log.Println("Will use colored rectangle instead")

		// Set up render loop for fallback
		renderLoop.Runnables = []notacore.Runnable{
			func() error {
				// Draw colored rectangle
				rect := notagl.CreateRectangle(notamath.Po2{X: 0, Y: 0}, 200, 200)
				rect.SetColor(notashader.Color{R: 1, G: 0, B: 0, A: 1})
				win.RunTime.Renderer.Submit(rect, 1.0)
				return nil
			},
		}
	} else {
		// Success!
		log.Printf("Success! Texture loaded: %dx%d", texture.Width, texture.Height)

		// Set up render loop with texture
		renderLoop.Runnables = []notacore.Runnable{
			func() error {
				// For now, draw a rectangle where texture should be
				// We'll need to implement proper textured rendering later
				rect := notagl.CreateRectangle(
					notamath.Po2{X: 0, Y: 0},
					float32(texture.Width),
					float32(texture.Height),
				)
				rect.SetColor(notashader.Color{R: 1, G: 1, B: 1, A: 1})
				win.RunTime.Renderer.Submit(rect, 1.0)
				return nil
			},
		}
	}

	// Run engine
	log.Println("Starting engine...")
	if err := engine.Run(); err != nil {
		log.Fatal("Engine run failed:", err)
	}

	log.Println("Test completed")
}
