package notacore

import (
	"NotaborEngine/notatomic"
	"sync"
)

type InputContext struct {
	mu sync.Mutex

	currHardware  map[StateInput]bool
	keyDownEvents map[StateInput]bool
	keyUpEvents   map[StateInput]bool
	frame         uint64

	// Snapshot taken at the beginning of the frame (safe for readers)
	snapshot *notatomic.Pointer[inputSnapshot]

	signalRegistry *sync.Map
}

// inputSnapshot holds a consistent view of input state for one frame
type inputSnapshot struct {
	Frame uint64

	PrevHardware map[StateInput]bool
	CurrHardware map[StateInput]bool

	KeyDownEvents map[StateInput]bool
	KeyUpEvents   map[StateInput]bool
}

func NewInputContext() *InputContext {
	ctx := &InputContext{
		currHardware:  make(map[StateInput]bool),
		keyDownEvents: make(map[StateInput]bool),
		keyUpEvents:   make(map[StateInput]bool),
		snapshot:      notatomic.NewPointer(&inputSnapshot{}),
	}

	// Initialize snapshot
	initialSnap := &inputSnapshot{
		PrevHardware:  make(map[StateInput]bool),
		CurrHardware:  make(map[StateInput]bool),
		KeyDownEvents: make(map[StateInput]bool),
		KeyUpEvents:   make(map[StateInput]bool),
	}
	ctx.snapshot.Set(initialSnap)

	return ctx
}

// beginFrame should be called at the start of each frame to update state
// This is the only place that does significant work and should be called from the input loop.
func (c *InputContext) beginFrame() {
	oldSnap := c.snapshot.Get()

	c.mu.Lock()
	defer c.mu.Unlock()

	newSnap := &inputSnapshot{
		Frame:         c.frame + 1,
		PrevHardware:  make(map[StateInput]bool),
		CurrHardware:  make(map[StateInput]bool, len(c.currHardware)),
		KeyDownEvents: make(map[StateInput]bool, len(c.keyDownEvents)),
		KeyUpEvents:   make(map[StateInput]bool, len(c.keyUpEvents)),
	}

	if oldSnap != nil {
		for k, v := range oldSnap.CurrHardware {
			newSnap.PrevHardware[k] = v
		}
	}

	for k, v := range c.currHardware {
		newSnap.CurrHardware[k] = v
	}
	for k := range c.keyDownEvents {
		newSnap.KeyDownEvents[k] = true
	}
	for k := range c.keyUpEvents {
		newSnap.KeyUpEvents[k] = true
	}

	c.keyDownEvents = make(map[StateInput]bool)
	c.keyUpEvents = make(map[StateInput]bool)
	c.frame = newSnap.Frame
	c.snapshot.Set(newSnap)
}

// recordKeyDown is called when a hardware key down event occurs
func (c *InputContext) recordKeyDown(input StateInput) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.currHardware[input] {
		c.keyDownEvents[input] = true
	}
	c.currHardware[input] = true
}

// recordKeyUp is called when a hardware key up event occurs
func (c *InputContext) recordKeyUp(input StateInput) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.currHardware[input] {
		c.keyUpEvents[input] = true
	}
	c.currHardware[input] = false
}

func (c *InputContext) currentFrame() uint64 {
	snap := c.snapshot.Get()
	if snap == nil {
		return 0
	}
	return snap.Frame
}

func (c *InputContext) registerSignal(name string, signal *InputSignal) {
	if c == nil || c.signalRegistry == nil || name == "" || signal == nil {
		return
	}
	c.signalRegistry.Store(name, signal)
}

// isKeyHeldThisFrame returns true if the key is physically pressed this frame
func (c *InputContext) isKeyHeldThisFrame(input StateInput) bool {
	snap := c.snapshot.Get()
	if snap == nil {
		return false
	}
	return snap.CurrHardware[input]
}

// wasKeyHeldLastFrame returns true if the key was physically pressed last frame
func (c *InputContext) wasKeyHeldLastFrame(input StateInput) bool {
	snap := c.snapshot.Get()
	if snap == nil {
		return false
	}
	return snap.PrevHardware[input]
}

// isKeyDownThisFrame returns true if a KeyDown event occurred for this input this frame
func (c *InputContext) isKeyDownThisFrame(input StateInput) bool {
	snap := c.snapshot.Get()
	if snap == nil {
		return false
	}
	return snap.KeyDownEvents[input]
}

// isKeyUpThisFrame returns true if a KeyUp event occurred for this input this frame
func (c *InputContext) isKeyUpThisFrame(input StateInput) bool {
	snap := c.snapshot.Get()
	if snap == nil {
		return false
	}
	return snap.KeyUpEvents[input]
}

// GetState returns the current input state for debugging/inspection
func (c *InputContext) GetState(input StateInput) InputState {
	snap := c.snapshot.Get()
	if snap == nil {
		return StateIdle
	}

	curr := snap.CurrHardware[input]
	prev := snap.PrevHardware[input]

	switch {
	case curr && !prev:
		return StatePressed
	case curr && prev:
		return StateHeld
	case !curr && prev:
		return StateReleased
	default:
		return StateIdle
	}
}
