package cc_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/codeclysm/cc"
)

func Example() {
	p := cc.New(4)
	p.Run(func() error {
		return errors.New("fail1")
	})
	p.Run(func() error {
		return errors.New("fail2")
	})
	p.Run(func() error {
		return nil
	})

	errs := p.Wait()
	fmt.Println(len(errs))
	// Output: 2
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
