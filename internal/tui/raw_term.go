package tui

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/term"
)

type rawTermState struct {
	fd       int
	oldState *term.State
	restored bool
	mu       sync.Mutex
}

var (
	// To setup the interrupt handler only once
	setupOnce sync.Once
	raw       *rawTermState
)

// MakeTermRaw puts the terminal into raw mode.
//
// Returns Err the terminal cannot be put into raw mode or
// if the term was already in raw mode when the func was called.
// Sets up automatic cleanup on SIGINT and SIGTERM signals.
func MakeTermRaw() error {
	if raw != nil && !raw.restored {
		return fmt.Errorf("terminal is already in raw mode by this package")
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to make terminal raw: %w", err)
	}
	raw = &rawTermState{fd: fd, oldState: oldState, restored: false}

	// Handle Ctrl-C gracefully
	setupOnce.Do(func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigs
			RestoreTerm()
			os.Exit(130)
		}()
	})
	return nil
}

// Restores the state of the terminal to canonical mode if in raw mode.
// No-op if terminal is already in canonical mode.
func RestoreTerm() {
	if raw == nil {
		return
	}

	raw.mu.Lock()
	defer raw.mu.Unlock()

	if !raw.restored && raw.oldState != nil {
		term.Restore(raw.fd, raw.oldState)
		raw.restored = true
	}
}
