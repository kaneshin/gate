package internal

import (
	"os"
	"os/signal"
	"syscall"
)

// Trap captures a signals in asynchronous.
func Trap(sigc chan os.Signal, sigf map[syscall.Signal]func(os.Signal)) {
	sigs := make([]os.Signal, len(sigf), len(sigf))
	for k := range sigf {
		sigs = append(sigs, k)
	}
	signal.Notify(sigc, sigs...)

	go func() {
		for {
			select {
			case c := <-sigc:
				key, ok := c.(syscall.Signal)
				if !ok {
					continue
				}
				sigf[key](c)
			}
		}
	}()
}
