package notacore

// InputSignal tracks a binary input across frames/ticks.
// Update State during the frame, then call Snapshot() once per tick to advance LastState.
type InputSignal struct {
	State     bool
	LastState bool
}

// Set updates the current state for this tick/frame.
func (s *InputSignal) Set(state bool) {
	s.State = state
}

// Snapshot advances LastState to match the current State.
func (s *InputSignal) Snapshot() {
	s.LastState = s.State
}

// Down reports whether the input is currently down.
func (s *InputSignal) Down() bool {
	return s.State
}

func (s *InputSignal) Changed() bool {
	return s.State != s.LastState
}

func (s *InputSignal) Held() bool {
	return s.State && s.LastState
}

func (s *InputSignal) Released() bool {
	return !s.State && s.LastState
}

func (s *InputSignal) Pressed() bool {
	return s.State && !s.LastState
}

func (s *InputSignal) Idle() bool {
	return !s.State && !s.LastState
}

func (s *InputSignal) Clone() InputSignal {
	return InputSignal{
		State:     s.State,
		LastState: s.LastState,
	}
}
