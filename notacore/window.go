package notacore

import (
	"NotaborEngine/notagl"
	"errors"
	"sync"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type WindowType int

const (
	Windowed WindowType = iota
	Fullscreen
	Borderless
)

type WindowConfig struct {
	X, Y       int
	W, H       int
	Title      string
	Resizable  bool
	Type       WindowType
	LogicLoops []*FixedHzLoop
	RenderLoop *RenderLoop
}

type Window interface {
	GetConfig() *WindowConfig
	GetRuntime() *WindowBaseRuntime
	MakeContextCurrent()
	SwapBuffers()
	ShouldClose() bool
	RunRenderer()
}

type WindowBaseRuntime struct {
	lastRender time.Time
	targetDt   time.Duration
}

type windowRunTime2D struct {
	WindowBaseRuntime
	backend  *notagl.GLBackend2D
	Renderer *notagl.Renderer2D
}

type windowRuntime3D struct {
	WindowBaseRuntime
	backend  *notagl.GLBackend3D
	Renderer *notagl.Renderer3D
}

type GlfwWindow2D struct {
	ID      int
	Handle  *glfw.Window
	Config  WindowConfig
	RunTime windowRunTime2D
}

func (w *GlfwWindow2D) GetConfig() *WindowConfig       { return &w.Config }
func (w *GlfwWindow2D) GetRuntime() *WindowBaseRuntime { return &w.RunTime.WindowBaseRuntime }
func (w *GlfwWindow2D) RunRenderer() {
	w.RunTime.Renderer.Orders = w.RunTime.Renderer.Orders[:0]
	w.Config.RenderLoop.Render()
	w.RunTime.Renderer.Flush(w.RunTime.backend)
}

type GlfwWindow3D struct {
	ID      int
	Handle  *glfw.Window
	Config  WindowConfig
	RunTime windowRuntime3D
}

func (w *GlfwWindow3D) GetConfig() *WindowConfig       { return &w.Config }
func (w *GlfwWindow3D) GetRuntime() *WindowBaseRuntime { return &w.RunTime.WindowBaseRuntime }
func (w *GlfwWindow3D) RunRenderer() {
	w.RunTime.Renderer.Orders = w.RunTime.Renderer.Orders[:0]
	w.Config.RenderLoop.Render()
	w.RunTime.Renderer.Flush(w.RunTime.backend)
}

type windowManager struct {
	mu        sync.Mutex
	windows2D []*GlfwWindow2D
	windows3D []*GlfwWindow3D
	nextID    int
}

func (wm *windowManager) Create2D(cfg WindowConfig) (*GlfwWindow2D, error) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if cfg.W <= 0 || cfg.H <= 0 {
		return nil, errors.New("invalid window size")
	}

	var monitor *glfw.Monitor
	var width, height int
	var x, y int

	switch cfg.Type {
	case Fullscreen:
		monitor = glfw.GetPrimaryMonitor()
		if monitor == nil {
			return nil, errors.New("no primary monitor found")
		}
		videoMode := monitor.GetVideoMode()
		if videoMode == nil {
			return nil, errors.New("could not get video mode")
		}
		width, height = videoMode.Width, videoMode.Height
		x, y = 0, 0

	case Borderless:
		monitor = glfw.GetPrimaryMonitor()
		if monitor == nil {
			return nil, errors.New("no primary monitor found")
		}
		videoMode := monitor.GetVideoMode()
		if videoMode == nil {
			return nil, errors.New("could not get video mode")
		}
		// Use actual window size or fallback to monitor resolution
		if cfg.W == 0 || cfg.H == 0 {
			width, height = videoMode.Width, videoMode.Height
		} else {
			width, height = cfg.W, cfg.H
		}
		// Center window position
		x = (videoMode.Width - width) / 2
		y = (videoMode.Height - height) / 2

	default: // Windowed
		width, height = cfg.W, cfg.H
		x, y = cfg.X, cfg.Y
	}
	if cfg.Type == Borderless {
		glfw.WindowHint(glfw.Decorated, glfw.False)
		if cfg.Resizable {
			glfw.WindowHint(glfw.Resizable, glfw.True)
		} else {
			glfw.WindowHint(glfw.Resizable, glfw.False)
		}
	} else {
		if cfg.Resizable {
			glfw.WindowHint(glfw.Resizable, glfw.True)
		} else {
			glfw.WindowHint(glfw.Resizable, glfw.False)
		}
		glfw.WindowHint(glfw.Decorated, glfw.True)
	}
	var handle *glfw.Window
	var err error

	if cfg.Type == Fullscreen {
		handle, err = glfw.CreateWindow(width, height, cfg.Title, monitor, nil)
	} else {
		handle, err = glfw.CreateWindow(width, height, cfg.Title, nil, nil)
	}

	if err != nil {
		glfw.DefaultWindowHints()
		return nil, err
	}

	glfw.DefaultWindowHints()

	if cfg.Type == Borderless {
		handle.SetPos(x, y)
	} else if cfg.Type == Windowed {
		handle.SetPos(cfg.X, cfg.Y)
	}

	handle.MakeContextCurrent()
	handle.Show()

	win := &GlfwWindow2D{
		ID:     wm.nextID,
		Handle: handle,
		Config: cfg,
		RunTime: windowRunTime2D{
			WindowBaseRuntime: WindowBaseRuntime{
				lastRender: time.Now(),
				targetDt:   time.Second / time.Duration(cfg.RenderLoop.MaxHz),
			},
			backend:  &notagl.GLBackend2D{},
			Renderer: &notagl.Renderer2D{},
		},
	}

	wm.windows2D = append(wm.windows2D, win)
	wm.nextID++
	return win, nil
}

func (wm *windowManager) Create3D(cfg WindowConfig) (*GlfwWindow3D, error) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if cfg.W <= 0 || cfg.H <= 0 {
		return nil, errors.New("invalid window size")
	}

	var monitor *glfw.Monitor
	var width, height int
	var x, y int

	switch cfg.Type {
	case Fullscreen:
		monitor = glfw.GetPrimaryMonitor()
		if monitor == nil {
			return nil, errors.New("no primary monitor found")
		}
		videoMode := monitor.GetVideoMode()
		if videoMode == nil {
			return nil, errors.New("could not get video mode")
		}
		width, height = videoMode.Width, videoMode.Height
		x, y = 0, 0

	case Borderless:
		monitor = glfw.GetPrimaryMonitor()
		if monitor == nil {
			return nil, errors.New("no primary monitor found")
		}
		videoMode := monitor.GetVideoMode()
		if videoMode == nil {
			return nil, errors.New("could not get video mode")
		}
		if cfg.W == 0 || cfg.H == 0 {
			width, height = videoMode.Width, videoMode.Height
		} else {
			width, height = cfg.W, cfg.H
		}
		x = (videoMode.Width - width) / 2
		y = (videoMode.Height - height) / 2

	default: // Windowed
		width, height = cfg.W, cfg.H
		x, y = cfg.X, cfg.Y
	}

	// Set window hints for borderless
	if cfg.Type == Borderless {
		glfw.WindowHint(glfw.Decorated, glfw.False)
	} else {
		glfw.WindowHint(glfw.Decorated, glfw.True)
	}

	if cfg.Resizable {
		glfw.WindowHint(glfw.Resizable, glfw.True)
	} else {
		glfw.WindowHint(glfw.Resizable, glfw.False)
	}

	var handle *glfw.Window
	var err error

	if cfg.Type == Fullscreen {
		handle, err = glfw.CreateWindow(width, height, cfg.Title, monitor, nil)
	} else {
		handle, err = glfw.CreateWindow(width, height, cfg.Title, nil, nil)
	}

	if err != nil {
		glfw.DefaultWindowHints()
		return nil, err
	}

	glfw.DefaultWindowHints()

	if cfg.Type == Borderless {
		handle.SetPos(x, y)
	} else if cfg.Type == Windowed {
		handle.SetPos(cfg.X, cfg.Y)
	}

	handle.MakeContextCurrent()
	handle.Show()

	win := &GlfwWindow3D{
		ID:     wm.nextID,
		Handle: handle,
		Config: cfg,
		RunTime: windowRuntime3D{
			WindowBaseRuntime: WindowBaseRuntime{
				lastRender: time.Now(),
				targetDt:   time.Second / time.Duration(cfg.RenderLoop.MaxHz),
			},
			backend:  &notagl.GLBackend3D{},
			Renderer: &notagl.Renderer3D{},
		},
	}

	wm.windows3D = append(wm.windows3D, win)
	wm.nextID++
	return win, nil
}

func (wm *windowManager) PollEvents() {
	glfw.PollEvents()
}

func (w *GlfwWindow2D) MakeContextCurrent()  { w.Handle.MakeContextCurrent() }
func (w *GlfwWindow2D) SwapBuffers()         { w.Handle.SwapBuffers() }
func (w *GlfwWindow2D) ShouldClose() bool    { return w.Handle.ShouldClose() }
func (w *GlfwWindow2D) Close()               { w.Handle.SetShouldClose(true) }
func (w *GlfwWindow2D) Size() (int, int)     { return w.Handle.GetSize() }
func (w *GlfwWindow2D) Position() (int, int) { return w.Handle.GetPos() }

// SetFullscreen switches between windowed and fullscreen mode
func (w *GlfwWindow2D) SetFullscreen(fullscreen bool) error {
	monitor := glfw.GetPrimaryMonitor()
	if monitor == nil {
		return errors.New("no primary monitor found")
	}

	videoMode := monitor.GetVideoMode()
	if videoMode == nil {
		return errors.New("could not get video mode")
	}

	if fullscreen {
		w.Handle.SetMonitor(monitor, 0, 0, videoMode.Width, videoMode.Height, videoMode.RefreshRate)
	} else {
		// Restore to windowed mode with previous or default size
		w.Handle.SetMonitor(nil, 100, 100, 800, 600, videoMode.RefreshRate)
	}

	w.Config.Type = Fullscreen

	return nil
}
func (w *GlfwWindow2D) SetBorderless(borderless bool) error {
	monitor := glfw.GetPrimaryMonitor()
	if monitor == nil {
		return errors.New("no primary monitor found")
	}

	videoMode := monitor.GetVideoMode()
	if videoMode == nil {
		return errors.New("could not get video mode")
	}

	if borderless {
		w.Handle.SetAttrib(glfw.Decorated, glfw.False)
		width, height := w.Size()
		x := (videoMode.Width - width) / 2
		y := (videoMode.Height - height) / 2
		w.Handle.SetPos(x, y)
	} else {
		w.Handle.SetAttrib(glfw.Decorated, glfw.True)
	}

	w.Config.Type = Borderless

	return nil
}

func (w *GlfwWindow2D) SetWindowed(x, y, width, height int) error {
	w.Handle.SetMonitor(nil, x, y, width, height, glfw.DontCare)
	w.Handle.SetAttrib(glfw.Decorated, glfw.True)
	w.Config.Type = Windowed

	return nil
}

func (w *GlfwWindow3D) MakeContextCurrent()  { w.Handle.MakeContextCurrent() }
func (w *GlfwWindow3D) SwapBuffers()         { w.Handle.SwapBuffers() }
func (w *GlfwWindow3D) ShouldClose() bool    { return w.Handle.ShouldClose() }
func (w *GlfwWindow3D) Close()               { w.Handle.SetShouldClose(true) }
func (w *GlfwWindow3D) Size() (int, int)     { return w.Handle.GetSize() }
func (w *GlfwWindow3D) Position() (int, int) { return w.Handle.GetPos() }

func (w *GlfwWindow3D) SetFullscreen(fullscreen bool) error {
	monitor := glfw.GetPrimaryMonitor()
	if monitor == nil {
		return errors.New("no primary monitor found")
	}

	videoMode := monitor.GetVideoMode()
	if videoMode == nil {
		return errors.New("could not get video mode")
	}

	if fullscreen {
		w.Handle.SetMonitor(monitor, 0, 0, videoMode.Width, videoMode.Height, videoMode.RefreshRate)
	} else {
		w.Handle.SetMonitor(nil, 100, 100, 800, 600, videoMode.RefreshRate)
	}

	w.Config.Type = Fullscreen

	return nil
}

func (w *GlfwWindow3D) SetBorderless(borderless bool) error {
	monitor := glfw.GetPrimaryMonitor()
	if monitor == nil {
		return errors.New("no primary monitor found")
	}

	videoMode := monitor.GetVideoMode()
	if videoMode == nil {
		return errors.New("could not get video mode")
	}

	if borderless {
		w.Handle.SetAttrib(glfw.Decorated, glfw.False)
		width, height := w.Size()
		x := (videoMode.Width - width) / 2
		y := (videoMode.Height - height) / 2
		w.Handle.SetPos(x, y)
	} else {
		w.Handle.SetAttrib(glfw.Decorated, glfw.True)
	}

	w.Config.Type = Borderless

	return nil
}

func (w *GlfwWindow3D) SetWindowed(x, y, width, height int) error {
	w.Handle.SetMonitor(nil, x, y, width, height, glfw.DontCare)
	w.Handle.SetAttrib(glfw.Decorated, glfw.True)
	w.Config.Type = Windowed

	return nil
}

func (wm *windowManager) Destroy2D(win *GlfwWindow2D) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	for i, w := range wm.windows2D {
		if w == win {
			w.Close()
			wm.windows2D = append(wm.windows2D[:i], wm.windows2D[i+1:]...)
			break
		}
	}
}

func (wm *windowManager) Destroy3D(win *GlfwWindow3D) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	for i, w := range wm.windows3D {
		if w == win {
			w.Close()
			wm.windows3D = append(wm.windows3D[:i], wm.windows3D[i+1:]...)
			break
		}
	}
}
