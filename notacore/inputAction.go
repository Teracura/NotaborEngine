package notacore

import "time"

type ActionBehavior int

const (
	RunWhileHeld ActionBehavior = iota
	RunWhileToggled
	RunOnceWhenPressed
	RunOnceWhenReleased
	RunWhileIdle
	Ignore
)

type Action struct {
	Name         string
	Signal       *InputSignal
	Toggled      bool
	HeldTicks    int64
	LastHeldTime time.Duration
	Behavior     ActionBehavior
	lastHold     time.Time
	lastRelease  time.Time

	Cooldown time.Duration
	lastRun  time.Time
}

func (a *Action) ShouldRun() bool {

	a.shouldToggle()
	a.updateHoldInformation()

	if time.Since(a.lastRun) < a.Cooldown {
		return false
	}

	var result bool
	switch a.Behavior {
	case RunWhileHeld:
		result = a.Signal.Held()
	case RunWhileToggled:
		result = a.Toggled
	case RunOnceWhenPressed:
		result = a.Signal.Pressed()
	case RunOnceWhenReleased:
		result = a.Signal.Released()
	case RunWhileIdle:
		result = a.Signal.Idle()
	case Ignore:
		return false
	}

	if result {
		a.lastRun = time.Now()
	}
	return result
}

func (a *Action) shouldToggle() {
	if a.Signal.Pressed() {
		a.Toggled = !a.Toggled
	}
}

func (a *Action) updateHoldInformation() {
	if a.Signal.Released() {
		a.lastRelease = time.Now()
		a.HeldTicks = 0
	}

	if a.Signal.Held() {
		a.lastHold = time.Now()
		a.HeldTicks++
	}

	a.LastHeldTime = a.lastHold.Sub(a.lastRelease)
}
