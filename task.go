package main

import (
	"NotaborEngine/internal/notatask"
	"time"
)

// Task represents a schedulable unit of work that can be configured with
// delays, intervals, repeat counts, and completion conditions.
type Task struct {
	handle *notatask.Task
}

// Loop is a fixed-timestep execution loop that runs tasks at a target frequency.
type Loop struct {
	handle *notatask.Loop
}

// NewLoop creates a new Loop running at the given frequency (Hz).
func NewLoop(hz float32) *Loop {
	return &Loop{
		handle: notatask.CreateLoop(hz),
	}
}

// Start begins loop execution (non-blocking; runs in a goroutine).
func (l *Loop) Start() {
	l.handle.Start()
}

// Stop halts loop execution and waits for workers to finish.
func (l *Loop) Stop() {
	l.handle.Stop()
}

// Alpha returns the frame interpolation fraction [0.0, 1.0] for smooth rendering.
func (l *Loop) Alpha() float32 {
	return l.handle.Alpha(time.Now())
}

// TickCount returns the total number of ticks executed by this loop.
func (l *Loop) TickCount() uint64 {
	return l.handle.TickCount()
}

// Remove cancels a task so it will not run again.
func (l *Loop) Remove(t *Task) {
	if t != nil && t.handle != nil {
		l.handle.Remove(t.handle)
	}
}

// Do schedules a function to run every tick.
func (l *Loop) Do(fn func()) *Task {
	return &Task{handle: l.handle.Do(fn)}
}

// Every sets the task to repeat at the given interval.
func (t *Task) Every(d time.Duration) *Task {
	t.handle.Every(d)
	return t
}

// Delay sets an initial delay before the first execution.
func (t *Task) Delay(d time.Duration) *Task {
	t.handle.Delay(d)
	return t
}

// AfterTicks delays execution by a number of loop ticks.
func (t *Task) AfterTicks(count uint32) *Task {
	t.handle.AfterTicks(count)
	return t
}

// Times limits the task to run a fixed number of times.
func (t *Task) Times(n uint32) *Task {
	t.handle.Times(n)
	return t
}

// Once configures the task to run exactly once.
func (t *Task) Once() *Task {
	t.handle.Once()
	return t
}

// Until configures the task to stop when cond returns true.
func (t *Task) Until(cond func() bool) *Task {
	t.handle.Until(cond)
	return t
}

// FinishAfter sets a hard deadline after which the task stops.
func (t *Task) FinishAfter(d time.Duration) *Task {
	t.handle.FinishAfter(d)
	return t
}
