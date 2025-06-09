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
	metronome := audio.NewMetronome(60)
	metronome.Paused = false

	err := tui.MakeTermRaw()
	if err != nil {
		log.Fatalf("Failed to set terminal to raw mode: %v. Your terminal might not be interactive.\n", err)
	}
	defer tui.RestoreTerm()

	// Main TUI Loop
	for {
		bpm := metronome.Bpm()
		// Clear line then print
		fmt.Printf("\r+ BPM: %d - ", int(bpm))

		buf := make([]byte, 1)
		_, err := os.Stdin.Read(buf)
		if err != nil {
			log.Fatalf("\nError reading input: %s", err)
		}

		switch buf[0] {
		case '=', '+', 'l':
			metronome.SetBpm(bpm + 2)
		case '-', '_', 'h':
			metronome.SetBpm(bpm - 2)
		case 'j':
			if metronome.Vol.Volume > -4.0 {
				metronome.Vol.Volume -= 0.1
			} else {
				metronome.Vol.Silent = true
			}
		case 'k':
			metronome.Vol.Silent = false
			if metronome.Vol.Volume < 1.5 {
				metronome.Vol.Volume += 0.1
			}
		case ' ', '\r':
			metronome.Paused = !metronome.Paused
		case 'q', byte(3): // Ctrl-C
			fmt.Print("\r\nHave a beautiful day!\r\n")
			return
		}
	}
}
