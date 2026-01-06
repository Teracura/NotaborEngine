package notacore

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
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
		fmt.Printf("Starting FixedHzLoop at %.2f Hz (interval: %v)\n", l.Hz, interval)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// Execute all runnables
				l.mu.Lock()
				rs := append([]Runnable(nil), l.Runnables...)
				l.mu.Unlock()

				for i, r := range rs {
					if err := r(); err != nil {
						l.Remove(i)
						fmt.Println(err)
					}
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
	} //too soon, skip

	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, runnable := range r.Runnables {
		if err := runnable(); err != nil {
			fmt.Println("Render error:", err)
		}
	}

	r.LastTime = now
}
