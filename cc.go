package cc

import (
	"sync"

	"github.com/fluxio/multierror"
)

// Pool manages a pool of concurrent workers. It works a bit like a Waitgroup, but with error reporting and concurrency limits
// You create one with New, and run functions with Run. Then you wait on it like a regular WaitGroup.
type Pool struct {
	errors multierror.ConcurrentAccumulator

	semaphore chan bool
	wg        *sync.WaitGroup
}

// New returns a new pool where a limited number (concurrency) of goroutine can work at the same time
func New(concurrency int) *Pool {
	wg := sync.WaitGroup{}
	p := Pool{
		errors:    multierror.ConcurrentAccumulator{},
		semaphore: make(chan bool, concurrency),
		wg:        &wg,
	}
	return &p
}

// Wait blocks and ensures that the channels are closed when all the goroutines end.
// It returns a list of all the errors returned by the goroutine
func (p *Pool) Wait() error {
	p.wg.Wait()
	close(p.semaphore)

	return p.errors.Error()
}

// Run wraps the given function into a goroutine and ensure that the concurrency limits are respected.
// The error returned by the function is stored into the error list returned by Wait
func (p *Pool) Run(fn func() error) {
	p.wg.Add(1)
	go func() {
		p.semaphore <- true
		p.errors.Push(fn())
		<-p.semaphore
		p.wg.Done()
	}()
}

// Stoppable is a function that can be stopped with the method Stop. You can also listen on the Stopped channel to see when it has been stopped.
// Stoppable is different from a context cancelation because it waits until the function has cleaned up before broadcasting on the Stopped channel
type Stoppable struct {
	Stopped chan struct{}
	stop    chan struct{}
	once    sync.Once
}

// Stop signals the provided function that it needs to stop
func (s *Stoppable) Stop() {
	s.once.Do(func() {
		close(s.stop)
	})
}

// Run creates a new stoppable function from the provided func. When you call the Stop method on the returned Stoppable the stop channel fed to the provided func is closed,
// signaling the need to stop. When the provided func returns the Stopped channel on the
// returned Stoppable is closed as well, broadcasting the message that it has finished
func Run(fn func(stop chan struct{})) (s *Stoppable) {
	s = &Stoppable{
		Stopped: make(chan struct{}),
		stop:    make(chan struct{}),
	}

	go func() {
		fn(s.stop)
		close(s.Stopped)
	}()
	return s
}
