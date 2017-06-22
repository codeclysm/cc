package cc

import "sync"

// Pool manages a pool of concurrent workers. It works a bit like a Waitgroup, but with error reporting and concurrency limits
// You create one with New, and run functions with Run. Then you wait on it like a regular WaitGroup.
//
// Example:
//
//   p := cc.New(4)
//   p.Run(func() error {
//		 afunction()
//       return nil
//   })
//   errs := p.Wait()
//
//   for err := range errs {
//
//   }
type Pool struct {
	errors chan error

	semaphore chan bool
	wg        *sync.WaitGroup
}

// New returns a new pool where a limited number (concurrency) of goroutine can work at the same time
func New(concurrency int) *Pool {
	wg := sync.WaitGroup{}
	p := Pool{
		errors:    make(chan error),
		semaphore: make(chan bool, concurrency),
		wg:        &wg,
	}
	return &p
}

// Wait blocks and ensures that the channels are closed when all the goroutines end.
// It returns a list of all the errors returned by the goroutine
func (p *Pool) Wait() []error {
	go func() {
		p.wg.Wait()

		close(p.semaphore)
		close(p.errors)
	}()

	errs := []error{}

	for err := range p.errors {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

// Run wraps the given function into a goroutine and ensure that the concurrency limits are respected.
// The error returned by the function is stored into the error list returned by Wait
func (p *Pool) Run(fn func() error) {
	p.wg.Add(1)
	go func() {
		p.semaphore <- true
		p.errors <- fn()
		<-p.semaphore
		p.wg.Done()
	}()
}
