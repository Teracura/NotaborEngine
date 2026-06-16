package main

import (
	"NotaborEngine/internal/notatask"
	"time"
)

type Task struct {
	handle *notatask.Task
}

type Loop struct {
	handle *notatask.Loop
}

// NewLoop creates and initialises a new Loop at a specified frequency (Hz).
func NewLoop(hz float32) *Loop {
	return &Loop{
		handle: notatask.CreateLoop(hz),
	}
}

func (l *Loop) Start() {
	l.handle.Start()
}

func (l *Loop) Stop() {
	l.handle.Stop()
}

// Alpha gets the frame interpolation fraction between 0.0 and 1.0.
func (l *Loop) Alpha() float32 {
	return l.handle.Alpha(time.Now())
}

func (l *Loop) TickCount() uint64 {
	return l.handle.TickCount()
}

func (l *Loop) Remove(t *Task) {
	if t != nil && t.handle != nil {
		l.handle.Remove(t.handle)
	}
}

func (l *Loop) Do(fn func()) *Task {
	return &Task{handle: l.handle.Do(fn)}
}

func (t *Task) Every(d time.Duration) *Task {
	t.handle.Every(d)
	return t
}

func (t *Task) Delay(d time.Duration) *Task {
	t.handle.Delay(d)
	return t
}

func (t *Task) AfterTicks(count uint32) *Task {
	t.handle.AfterTicks(count)
	return t
}

func (t *Task) Times(n uint32) *Task {
	t.handle.Times(n)
	return t
}

func (t *Task) Once() *Task {
	t.handle.Once()
	return t
}

func (t *Task) Until(cond func() bool) *Task {
	t.handle.Until(cond)
	return t
}

func (t *Task) FinishAfter(d time.Duration) *Task {
	t.handle.FinishAfter(d)
	return t
}
