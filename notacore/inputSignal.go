package notacore

type InputSignal struct {
	State     bool
	LastState bool
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

func (s *InputSignal) Snapshot() {
	s.LastState = s.State
}

func (s *InputSignal) Clone() InputSignal {
	return InputSignal{s.State, s.LastState}
}
