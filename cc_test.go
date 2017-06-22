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

func TestRace(t *testing.T) {
	p := cc.New(4)

	for i := 0; i < 1000; i++ {
		p.Run(func() error {
			return errors.New("fail")
		})
	}

	p.Wait()
}
