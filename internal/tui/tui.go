package tui

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/term"
)

type rawTerm struct {
	fd       int
	oldState *term.State
	restored bool
	mu       sync.Mutex
}

// Holds the rawTerm state
var raw *rawTerm

// EnableRawTerm puts the terminal into raw mode and returns a RawTerm
// that can be used to restore the original terminal state.
//
// Panics if the terminal cannot be put into raw mode.
// Sets up automatic cleanup on SIGINT and SIGTERM signals.
func EnableRawTerm() {
	if raw != nil && !raw.restored {
		panic("raw term already enabled")
	}

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	raw = &rawTerm{fd: fd, oldState: oldState, restored: false}
	if err != nil {
		panic(err)
	}

	// Handle Ctrl-C gracefully
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		CleanupRawTerm()
		os.Exit(130)
	}()
}

func CleanupRawTerm() {
	raw.mu.Lock()
	defer raw.mu.Unlock()

	if !raw.restored && raw.oldState != nil {
		term.Restore(raw.fd, raw.oldState)
		raw.restored = true
	}
}
