package notacore

import (
	"errors"
	"sync"

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
	LogicLoops []FixedHzLoop
	RenderLoop *RenderLoop
}

type Window interface {
	ID() int
	Size() (int, int)
	Position() (int, int)

	MakeContextCurrent()
	SwapBuffers()

	ShouldClose() bool
	Close()

	SetTitle(string)
}

type WindowManager interface {
	Create(WindowConfig) (Window, error)
	Destroy(Window)

	PollEvents()
	Windows() []Window
}

type glfwWindow struct {
	id     int
	handle *glfw.Window
}

type GLFWWindowManager struct {
	mu      sync.Mutex
	windows []*glfwWindow
	nextID  int
}

func NewGLFWWindowManager() (*GLFWWindowManager, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 6) //GL VERSION 4.6
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True) //FOR DEVICES THAT SUPPORT LOWER END VERSIONS

	return &GLFWWindowManager{
		windows: []*glfwWindow{},
		nextID:  0,
	}, nil
}

func (wm *GLFWWindowManager) PollEvents() {
	glfw.PollEvents()
}

func (wm *GLFWWindowManager) Windows() []Window {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	windows := make([]Window, len(wm.windows))
	for i, w := range wm.windows {
		windows[i] = w
	}
	return windows
}

func (w *glfwWindow) ID() int {
	return w.id
}

func (w *glfwWindow) Size() (int, int) {
	return w.handle.GetSize()
}

func (w *glfwWindow) Position() (int, int) {
	return w.handle.GetPos()
}

func (w *glfwWindow) ShouldClose() bool {
	return w.handle.ShouldClose()
}

func (w *glfwWindow) Close() {
	w.handle.SetShouldClose(true)
}

func (w *glfwWindow) SetTitle(title string) {
	w.handle.SetTitle(title)
}

func (w *glfwWindow) MakeContextCurrent() {
	w.handle.MakeContextCurrent()
}

func (w *glfwWindow) SwapBuffers() {
	w.handle.SwapBuffers()
}

func (wm *GLFWWindowManager) Create(cfg WindowConfig) (Window, error) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if cfg.W <= 0 || cfg.H <= 0 {
		return nil, errors.New("invalid window size")
	}

	var monitor *glfw.Monitor
	if cfg.Type == Fullscreen {
		monitor = glfw.GetPrimaryMonitor()
	} else {
		monitor = nil
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

	win := &glfwWindow{
		id:     wm.nextID,
		handle: handle,
	}

	wm.windows = append(wm.windows, win)
	wm.nextID++

	for i := range cfg.LogicLoops {
		cfg.LogicLoops[i].Start()
	}

	if cfg.RenderLoop != nil {
		cfg.RenderLoop.Start()
	}

	return win, nil
}

func (wm *GLFWWindowManager) Destroy(w Window) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	for i, win := range wm.windows {
		if win == w {
			win.Close()
			last := len(wm.windows) - 1
			wm.windows[i] = wm.windows[last]
			wm.windows[last] = nil
			wm.windows = wm.windows[:last]
			break
		}
	}
}
