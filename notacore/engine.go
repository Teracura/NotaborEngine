package notacore

import (
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Settings struct {
	Vsync bool
}

type Engine struct {
	Windows2D     []*GlfwWindow2D
	Windows3D     []*GlfwWindow3D
	Settings      *Settings
	WindowManager *GLFWWindowManager
	running       bool
}

func (e *Engine) Run() error {
	e.running = true

	// Start all logic loops
	for _, w := range e.Windows2D {
		for _, loop := range w.Config.LogicLoops {
			loop.Start()
		}
		w.RunTime.LastRender = time.Now()
	}
	for _, w := range e.Windows3D {
		for _, loop := range w.Config.LogicLoops {
			loop.Start()
		}
		w.RunTime.LastRender = time.Now()
	}

	for e.running && !e.AllWindowsClosed() {
		e.WindowManager.PollEvents()
		now := time.Now()

		// Render 2D windows
		for _, win := range e.Windows2D {
			if win.ShouldClose() {
				continue
			}
			elapsed := now.Sub(win.RunTime.LastRender)
			if elapsed < win.RunTime.TargetDt {
				continue
			}
			win.RunTime.LastRender = now

			win.MakeContextCurrent()

			win.RunTime.Renderer.Begin()
			win.Config.RenderLoop.Render()
			win.RunTime.Renderer.Flush(win.RunTime.Backend)

			win.SwapBuffers()
		}

		// Render 3D windows
		for _, win := range e.Windows3D {
			if win.ShouldClose() {
				continue
			}
			elapsed := now.Sub(win.RunTime.LastRender)
			if elapsed < win.RunTime.TargetDt {
				continue
			}
			win.RunTime.LastRender = now

			win.MakeContextCurrent()

			win.RunTime.Renderer.Begin()
			win.Config.RenderLoop.Render()
			win.RunTime.Renderer.Flush(win.RunTime.Backend)

			win.SwapBuffers()
		}
	}

	// Stop logic loops
	for _, w := range e.Windows2D {
		for _, loop := range w.Config.LogicLoops {
			loop.Stop()
		}
	}
	for _, w := range e.Windows3D {
		for _, loop := range w.Config.LogicLoops {
			loop.Stop()
		}
	}

	return nil
}

func (e *Engine) AllWindowsClosed() bool {
	for _, w := range e.Windows2D {
		if !w.ShouldClose() {
			return false
		}
	}
	for _, w := range e.Windows3D {
		if !w.ShouldClose() {
			return false
		}
	}
	return true
}

func (e *Engine) Shutdown() {
	glfw.Terminate()
}

func (e *Engine) InitPlatform() error {
	runtime.LockOSThread()

	if err := addNativeDLLPath(); err != nil {
		return err
	}

	wm, err := NewGLFWWindowManager()
	if err != nil {
		return err
	}

	e.WindowManager = wm
	return nil
}

var glInitialized bool

func (e *Engine) CreateWindow2D(cfg WindowConfig) (*GlfwWindow2D, error) {
	win, err := e.WindowManager.Create2D(cfg)
	if err != nil {
		return nil, err
	}
	win.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		return nil, err
	}
	win.RunTime.Backend.Init()
	e.Windows2D = append(e.Windows2D, win)
	return win, nil
}

func (e *Engine) CreateWindow3D(cfg WindowConfig) (*GlfwWindow3D, error) {
	win, err := e.WindowManager.Create3D(cfg)
	if err != nil {
		return nil, err
	}
	win.MakeContextCurrent()
	if err := gl.Init(); err != nil {
		return nil, err
	}
	win.RunTime.Backend.Init()
	e.Windows3D = append(e.Windows3D, win)
	return win, nil
}

func addNativeDLLPath() error {
	switch runtime.GOOS {
	case "windows":
		exeDir, err := os.Getwd()
		if err != nil {
			return err
		}
		dllDir := filepath.Join(exeDir, "notacore", "native", "windows")
		err = os.Setenv("PATH", dllDir+";"+os.Getenv("PATH"))
		if err != nil {
			return err
		}

	case "linux":
		// set linux paths later

	case "darwin":
	// set mac paths later
	default:
		// return unsupported platform error
	}
	return nil
}
