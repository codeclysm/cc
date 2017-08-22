package cc_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/codeclysm/cc"
)

func Example() {
	p := cc.New(4)
	p.Run(func() error {
		time.Sleep(1 * time.Second)
		return errors.New("fail1")
	})
	p.Run(func() error {
		return errors.New("fail2")
	})
	p.Run(func() error {
		return nil
	})

	errs := p.Wait()
	fmt.Println(errs)
	// Output: 2 errors:
	//:   fail2
	//:   fail1
}

func ExampleStoppable() {
	stoppable := cc.Run(func(stop chan struct{}) {
		i := 0
	L:
		for {
			select {
			case <-stop:
				fmt.Println("receive stop signal")
				break L
			default:
				i++
				time.Sleep(250 * time.Millisecond)
				fmt.Println(i)
			}
		}
		fmt.Println("finished with", i)
	})

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("send stop signal")
		stoppable.Stop()
		stoppable.Stop() // It shouldn't explode even if you attempt to close it multiple times
	}()

	<-stoppable.Stopped
	fmt.Println("stopped finally")
	// Output: 1
	// 2
	// 3
	// send stop signal
	// 4
	// receive stop signal
	// finished with 4
	// stopped finally
}

func TestRace(t *testing.T) {
	p := cc.New(4)

	for i := 0; i < 1000; i++ {
		p.Run(func() error {
			return errors.New("fail")
		})
	}

	p.Wait()
}

func Benchmark(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := cc.New(4)
		for i := 0; i < 1000; i++ {
			p.Run(func() error {
				time.Sleep(1 * time.Millisecond)
				return nil
			})
		}
		p.Wait()
	}
}
