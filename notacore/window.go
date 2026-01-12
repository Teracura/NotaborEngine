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

	autoStartLoops bool
}

type WindowRuntime2D struct {
	LastRender time.Time
	TargetDt   time.Duration
	Backend    *notagl.GLBackend2D
	Renderer   *notagl.Renderer2D
}

type WindowRuntime3D struct {
	LastRender time.Time
	TargetDt   time.Duration
	Backend    *notagl.GLBackend3D
	Renderer   *notagl.Renderer3D
}

type GlfwWindow2D struct {
	ID      int
	Handle  *glfw.Window
	Config  WindowConfig
	RunTime WindowRuntime2D
}

type GlfwWindow3D struct {
	ID      int
	Handle  *glfw.Window
	Config  WindowConfig
	RunTime WindowRuntime3D
}

type GLFWWindowManager struct {
	mu        sync.Mutex
	windows2D []*GlfwWindow2D
	windows3D []*GlfwWindow3D
	nextID    int
}

func NewGLFWWindowManager() (*GLFWWindowManager, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	return &GLFWWindowManager{
		windows2D: []*GlfwWindow2D{},
		windows3D: []*GlfwWindow3D{},
		nextID:    0,
	}, nil
}

func (wm *GLFWWindowManager) Create2D(cfg WindowConfig) (*GlfwWindow2D, error) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if cfg.W <= 0 || cfg.H <= 0 {
		return nil, errors.New("invalid window size")
	}

	var monitor *glfw.Monitor
	if cfg.Type == Fullscreen {
		monitor = glfw.GetPrimaryMonitor()
	}

	handle, err := glfw.CreateWindow(cfg.W, cfg.H, cfg.Title, monitor, nil)
	if err != nil {
		return nil, err
	}
	if !cfg.Resizable {
		handle.SetAttrib(glfw.Resizable, glfw.False)
	}

	handle.SetPos(cfg.X, cfg.Y)
	handle.MakeContextCurrent()
	handle.Show()

	win := &GlfwWindow2D{
		ID:     wm.nextID,
		Handle: handle,
		Config: cfg,
		RunTime: WindowRuntime2D{
			LastRender: time.Now(),
			TargetDt:   time.Second / time.Duration(cfg.RenderLoop.MaxHz),
			Backend:    &notagl.GLBackend2D{},
			Renderer:   &notagl.Renderer2D{},
		},
	}

	wm.windows2D = append(wm.windows2D, win)
	wm.nextID++
	return win, nil
}

func (wm *GLFWWindowManager) Create3D(cfg WindowConfig) (*GlfwWindow3D, error) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if cfg.W <= 0 || cfg.H <= 0 {
		return nil, errors.New("invalid window size")
	}

	var monitor *glfw.Monitor
	if cfg.Type == Fullscreen {
		monitor = glfw.GetPrimaryMonitor()
	}

	handle, err := glfw.CreateWindow(cfg.W, cfg.H, cfg.Title, monitor, nil)
	if err != nil {
		return nil, err
	}
	if !cfg.Resizable {
		handle.SetAttrib(glfw.Resizable, glfw.False)
	}

	handle.SetPos(cfg.X, cfg.Y)
	handle.MakeContextCurrent()
	handle.Show()

	win := &GlfwWindow3D{
		ID:     wm.nextID,
		Handle: handle,
		Config: cfg,
		RunTime: WindowRuntime3D{
			LastRender: time.Now(),
			TargetDt:   time.Second / time.Duration(cfg.RenderLoop.MaxHz),
			Backend:    &notagl.GLBackend3D{},
			Renderer:   &notagl.Renderer3D{},
		},
	}

	wm.windows3D = append(wm.windows3D, win)
	wm.nextID++
	return win, nil
}

func (wm *GLFWWindowManager) PollEvents() {
	glfw.PollEvents()
}

type Window interface {
	MakeContextCurrent()
	SwapBuffers()
	ShouldClose() bool
	Close()
	Size() (int, int)
	Position() (int, int)
}

func (w *GlfwWindow2D) MakeContextCurrent()  { w.Handle.MakeContextCurrent() }
func (w *GlfwWindow2D) SwapBuffers()         { w.Handle.SwapBuffers() }
func (w *GlfwWindow2D) ShouldClose() bool    { return w.Handle.ShouldClose() }
func (w *GlfwWindow2D) Close()               { w.Handle.SetShouldClose(true) }
func (w *GlfwWindow2D) Size() (int, int)     { return w.Handle.GetSize() }
func (w *GlfwWindow2D) Position() (int, int) { return w.Handle.GetPos() }

func (w *GlfwWindow3D) MakeContextCurrent()  { w.Handle.MakeContextCurrent() }
func (w *GlfwWindow3D) SwapBuffers()         { w.Handle.SwapBuffers() }
func (w *GlfwWindow3D) ShouldClose() bool    { return w.Handle.ShouldClose() }
func (w *GlfwWindow3D) Close()               { w.Handle.SetShouldClose(true) }
func (w *GlfwWindow3D) Size() (int, int)     { return w.Handle.GetSize() }
func (w *GlfwWindow3D) Position() (int, int) { return w.Handle.GetPos() }

func (wm *GLFWWindowManager) Destroy2D(win *GlfwWindow2D) {
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

func (wm *GLFWWindowManager) Destroy3D(win *GlfwWindow3D) {
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

func (wc WindowConfig) AddLogicLoop(loop *FixedHzLoop) {
	wc.LogicLoops = append(wc.LogicLoops, loop)
}
