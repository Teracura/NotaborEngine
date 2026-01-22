package notacore

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Runnable func() error

type FixedHzLoop struct {
	Hz float32

	mu               sync.Mutex
	Runnables        []Runnable
	OneTimeRunnables []Runnable

	stop chan struct{}
	wg   sync.WaitGroup

	lastTick time.Time
	delta    time.Duration
}

type RenderLoop struct {
	MaxHz     float32
	Runnables []Runnable
	LastTime  time.Time
}

func (l *FixedHzLoop) Start() {

	l.mu.Lock()
	// Prevent multiple starts
	if l.stop != nil {
		l.mu.Unlock()
		return
	}
	l.stop = make(chan struct{})
	l.mu.Unlock()

	l.stop = make(chan struct{})
	l.wg.Add(1)

	go func() {
		defer l.wg.Done()

		interval := time.Duration(float64(time.Second) / float64(l.Hz))
		l.delta = interval
		l.lastTick = time.Now()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Execute all runnables
				l.mu.Lock()
				l.lastTick = time.Now()
				otr := append([]Runnable(nil), l.OneTimeRunnables...)
				atr := append([]Runnable(nil), l.Runnables...) //all time runnables
				l.mu.Unlock()

				for _, r := range otr {
					if err := r(); err != nil {
						fmt.Println(err)
					}
				}

				newRunnables := atr[:0] // reuse underlying array
				for _, r := range atr {
					if err := r(); err != nil {
						fmt.Println(err)
					} else {
						newRunnables = append(newRunnables, r)
					}
				}

				l.mu.Lock()
				l.OneTimeRunnables = nil
				l.Runnables = newRunnables
				l.mu.Unlock()

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

func (r *RenderLoop) Render() {
	now := time.Now()

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, runnable := range r.Runnables {
		if err := runnable(); err != nil {
			fmt.Println("Render error:", err)
		}
	}
	r.LastTime = now
}

func (l *FixedHzLoop) Alpha(now time.Time) float32 {
	l.mu.Lock()
	last := l.lastTick
	delta := l.delta
	l.mu.Unlock()

	if delta <= 0 {
		return 1
	}

	alpha := float32(now.Sub(last).Seconds() / delta.Seconds())

	if alpha < 0 {
		return 0
	}
	if alpha > 1 {
		return 1
	}
	return alpha
}
