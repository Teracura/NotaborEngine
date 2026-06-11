package notasdl

import (
	"NotaborEngine/notaentity"
	"NotaborEngine/notarender"
	"NotaborEngine/notashader"
	"NotaborEngine/notatask"
	"NotaborEngine/notatexture"
	"fmt"
	"sync"
	"time"

	"github.com/Zyko0/go-sdl3/sdl"
)

type WindowType int

const (
	Windowed WindowType = iota
	Fullscreen
	Borderless
)

type WindowConfig struct {
	X, Y      int
	W, H      int
	Title     string
	Resizable bool
	Type      WindowType

	TargetFPS float32
	Loops     []*notatask.Loop
}

type WindowRuntime struct {
	LastFrame time.Time
	TargetDt  time.Duration

	Backend    *notarender.Backend
	Renderer   *notarender.Renderer
	TextureMgr *notatexture.TextureManager
	SpriteMgr  *notatexture.SpriteManager
	ShaderMgr  *notashader.Manager

	Cameras []*Camera2D
}

type Window struct {
	ID            WindowID
	Handle        *sdl.Window
	Config        *WindowConfig
	Runtime       *WindowRuntime
	DefaultCamera *Camera2D

	ShouldClose bool
	Hidden      bool
	Minimized   bool
	Occluded    bool

	positionMu         sync.Mutex
	pendingSetPosition bool
	pendingX           int
	pendingY           int
	pendingMoveX       int
	pendingMoveY       int

	sizeMu sync.RWMutex
	width  int
	height int
}

func (w *Window) RenderFrame() {
	now := time.Now()
	dt := float32(now.Sub(w.Runtime.LastFrame).Seconds())
	w.Runtime.LastFrame = now

	// Get safe dimensions
	width, height := w.GetSize()

	if !w.canRender() || width <= 0 || height <= 0 {
		w.Runtime.Renderer.Clear()
		return
	}

	w.applyPendingWindowPosition()

	for _, cam := range w.Runtime.Cameras {
		cam.Update(dt)
	}

	w.Runtime.Renderer.FrameID.Inc()

	cmdBuf, err := w.Runtime.Backend.BeginFrame()
	if err != nil {
		return
	}

	if err := w.Runtime.Renderer.Flush(w.Runtime.Backend, cmdBuf, w.Handle); err != nil {
		return
	}
}

func (w *Window) canRender() bool {
	width, height := w.GetSize()
	return !w.ShouldClose &&
		!w.Hidden &&
		!w.Minimized &&
		!w.Occluded &&
		width > 0 &&
		height > 0
}

func (w *Window) MakeCurrent() {
	// SDL3 GPU doesn't require making a context current
	// This is a no-op for GPU rendering
}

func (w *Window) GetConfig() *WindowConfig {
	return w.Config
}

func (w *Window) SetVSync(enabled bool) {
	if w == nil || w.Handle == nil || w.Runtime == nil || w.Runtime.Backend == nil || w.Runtime.Backend.Device == nil {
		return
	}

	mode := sdl.GPU_PRESENTMODE_VSYNC
	if !enabled {
		mode = sdl.GPU_PRESENTMODE_IMMEDIATE
		if !w.Runtime.Backend.Device.WindowSupportsPresentMode(w.Handle, mode) {
			mode = sdl.GPU_PRESENTMODE_MAILBOX
		}
		if !w.Runtime.Backend.Device.WindowSupportsPresentMode(w.Handle, mode) {
			mode = sdl.GPU_PRESENTMODE_VSYNC
		}
	}

	_ = w.Runtime.Backend.Device.SetSwapchainParameters(w.Handle, sdl.GPU_SWAPCHAINCOMPOSITION_SDR, mode)
}

func (w *Window) Destroy() {
	if w == nil {
		return
	}

	w.ShouldClose = true

	if w.Runtime != nil {
		if w.Runtime.SpriteMgr != nil {
			w.Runtime.SpriteMgr.Clear()
		}
		if w.Runtime.TextureMgr != nil {
			w.Runtime.TextureMgr.Clear()
		}
		if w.Runtime.ShaderMgr != nil {
			w.Runtime.ShaderMgr.Clear()
		}
		if w.Runtime.Backend != nil {
			if w.Runtime.Backend.Device != nil && w.Handle != nil {
				w.Runtime.Backend.Device.ReleaseWindow(w.Handle)
			}
			w.Runtime.Backend.Shutdown()
			w.Runtime.Backend = nil
		}
	}

	if w.Handle != nil {
		w.Handle.Destroy()
		w.Handle = nil
	}
}

// Draw queues entities for rendering using SDL-backed renderer.
// If cam is nil, default camera is used.
func (w *Window) Draw(alpha float32, cam *Camera2D, entities ...*notaentity.Entity) error {
	if cam == nil {
		cam = w.DefaultCamera
	}

	view := cam.ViewMatrix()

	for _, e := range entities {
		if e == nil {
			continue
		}
		if err := e.DrawWithView(w.Runtime.Renderer, view, alpha); err != nil {
			return err
		}
	}

	return nil
}

func (w *Window) SetPosition(x, y int) {
	if w.Handle == nil {
		return
	}

	w.positionMu.Lock()
	w.pendingSetPosition = true
	w.pendingX = x
	w.pendingY = y
	w.pendingMoveX = 0
	w.pendingMoveY = 0
	w.positionMu.Unlock()
}

func (w *Window) applyPendingWindowPosition() {
	w.positionMu.Lock()
	if !w.pendingSetPosition && w.pendingMoveX == 0 && w.pendingMoveY == 0 {
		w.positionMu.Unlock()
		return
	}

	x := w.Config.X
	y := w.Config.Y
	if w.pendingSetPosition {
		x = w.pendingX
		y = w.pendingY
	}
	x += w.pendingMoveX
	y += w.pendingMoveY

	w.pendingSetPosition = false
	w.pendingMoveX = 0
	w.pendingMoveY = 0
	w.positionMu.Unlock()

	_ = w.setPositionNow(x, y)
}

func (w *Window) setPositionNow(x, y int) error {
	if w.Handle == nil {
		return nil
	}

	bounds, err := w.getCurrentDisplayBounds()
	if err != nil {
		return err
	}

	minX := -w.Config.W + 50
	minY := -w.Config.H + 50
	maxX := int(bounds.W - 50)
	maxY := int(bounds.H - 50)

	if x < minX {
		x = minX
	}
	if x > maxX {
		x = maxX
	}
	if y < minY {
		y = minY
	}
	if y > maxY {
		y = maxY
	}

	if err := w.Handle.SetPosition(int32(x), int32(y)); err != nil {
		return err
	}

	w.setCachedPosition(x, y)
	return nil
}

func (w *Window) setCachedPosition(x, y int) {
	w.positionMu.Lock()
	defer w.positionMu.Unlock()
	w.Config.X = x
	w.Config.Y = y
}

// Move moves the window by delta with boundary checking
func (w *Window) Move(dx, dy int) {
	if w.Handle == nil {
		return
	}

	w.positionMu.Lock()
	w.pendingMoveX += dx
	w.pendingMoveY += dy
	w.positionMu.Unlock()
}

// Helper to get current display bounds

func (w *Window) getCurrentDisplayBounds() (*sdl.Rect, error) {
	if w.Handle == nil {
		return nil, fmt.Errorf("window handle is nil, initialize window before calling getCurrentDisplayBounds")
	}

	bounds, err := sdl.GetDisplayForWindow(w.Handle).Bounds()
	if err != nil {
		return nil, fmt.Errorf("failed to get display bounds: %w", err)
	}

	return bounds, nil
}

// GetPosition returns current window position
func (w *Window) GetPosition() (x, y int) {
	if w.Handle == nil {
		return w.Config.X, w.Config.Y
	}

	w.positionMu.Lock()
	defer w.positionMu.Unlock()

	x = w.Config.X
	y = w.Config.Y
	if w.pendingSetPosition {
		x = w.pendingX
		y = w.pendingY
	}

	return x + w.pendingMoveX, y + w.pendingMoveY
}

// GetSize returns the current window size safely
func (w *Window) GetSize() (width, height int) {
	w.sizeMu.RLock()
	defer w.sizeMu.RUnlock()

	if w.width > 0 && w.height > 0 {
		return w.width, w.height
	}
	return w.Config.W, w.Config.H
}

// SetSize updates the window size safely
func (w *Window) SetSize(width, height int) {
	w.sizeMu.Lock()
	defer w.sizeMu.Unlock()

	w.width = width
	w.height = height
	w.Config.W = width
	w.Config.H = height
}
