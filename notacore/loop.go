package notacore

import (
	"fmt"
	"sync"
	"time"
)

type Loop interface {
	Start()
	Stop()
	Remove(int)
}

type Runnable func() error

type FixedHzLoop struct {
	Hz float32

	mu        sync.Mutex
	Runnables []Runnable

	stop chan struct{}
	wg   sync.WaitGroup
}

type RenderLoop struct {
	MaxHz     float32
	Runnables []Runnable
	LastTime  time.Time
}

func (l *FixedHzLoop) Start() {
	l.stop = make(chan struct{})
	l.wg.Add(1)

	go func() {
		defer l.wg.Done()

		ticker := time.NewTicker(
			time.Duration(float64(time.Second) / float64(l.Hz)),
		)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				l.mu.Lock()
				rs := append([]Runnable(nil), l.Runnables...)
				l.mu.Unlock()
				i := 0
				for _, r := range rs {
					if err := r(); err != nil {
						l.Remove(i)
						fmt.Println(err)
					}
					i++
				}

			case <-l.stop:
				return
			}
		}
	}()
}

func (l *FixedHzLoop) Stop() {
	close(l.stop)
	l.wg.Wait()
}

func (l *FixedHzLoop) Remove(i int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	last := len(l.Runnables) - 1
	if i < 0 || i > last {
		return
	}

	l.Runnables[i] = l.Runnables[last]
	l.Runnables[last] = nil
	l.Runnables = l.Runnables[:last]
}

func (r *RenderLoop) Start() {
	r.LastTime = time.Now()
}

func (r *RenderLoop) Render() {
	now := time.Now()
	minInterval := time.Duration(float64(time.Second) / float64(r.MaxHz))

	if now.Sub(r.LastTime) < minInterval {
		// Too soon, skip this frame
		return
	}

	for _, runnable := range r.Runnables {
		if err := runnable(); err != nil {
			fmt.Println("Render error:", err)
		}
	}

	r.LastTime = now
}
