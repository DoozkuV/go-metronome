package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DoozkuV/go-metronome/audio"
	"github.com/DoozkuV/go-metronome/internal/tui"
)

func main() {

	// Initialize our metronome
	m := audio.NewMetronome(60)
	m.Ctrl.Paused = false

	tui.MakeRawTerm()
	defer tui.RestoreTerm()
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
			fmt.Print("\r\nHave a beautiful day!\r\n")
			return
		}
	}
}
