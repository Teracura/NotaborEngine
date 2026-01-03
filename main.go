package main

import (
	"NotaborEngine/notacore"
	"os"
	"path/filepath"
	"runtime"

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

	for !window.ShouldClose() {
		wm.PollEvents()

		window.MakeContextCurrent()
		renderLoop.Render()
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
