package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DoozkuV/go-metronome/audio"
	"golang.org/x/term"
)

type rawTerm struct {
	fd       int
	oldState *term.State
}

func CreateRawTerm() *rawTerm {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	raw := rawTerm{fd: fd, oldState: oldState}
	if err != nil {
		panic(err)
	}

	// Handle Ctrl-C gracefully
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigs
		raw.EndProgram()
	}()

	return &raw
}

func (r *rawTerm) EndProgram() {
	term.Restore(r.fd, r.oldState)
	fmt.Println("\nGoodbye!")
	os.Exit(0)
}

func main() {

	// Initialize our metronome
	m := audio.NewMetronome(60)
	m.Ctrl.Paused = false

	raw := CreateRawTerm()
	// TODO: Find a better solution for ending the prog
	defer raw.EndProgram()
	// Main TUI Loop
	for {
		bpm := m.Bpm()
		// Clear line then print
		fmt.Printf("\r+ BPM: %d - ", int(bpm))

		buf := make([]byte, 1)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			log.Fatalf("\nError reading input: %s", err)
		}

		switch buf[0] {
		case '=', '+':
			m.SetBpm(bpm + 2)
		case '-', '_':
			m.SetBpm(bpm - 2)
		case 'q', byte(3): // Ctrl-C
			raw.EndProgram()
		}
	}
}
