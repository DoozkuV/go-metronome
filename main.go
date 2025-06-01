package main

import (
	"fmt"
	"strconv"

	"github.com/DoozkuV/go-metronome/audio"
)

func main() {
	metronome := audio.NewMetronome(60)
	metronome.Ctrl.Paused = false
	for {
		bpm := metronome.Bpm()
		fmt.Printf("Current BPM: %v\n", bpm)
		var input string
		fmt.Scan(&input)
		switch input {
		case "+":
			metronome.SetBpm(bpm + 2)
		case "-":
			metronome.SetBpm(bpm - 2)
		default:
			if inputNum, err := strconv.Atoi(input); err != nil {
				fmt.Println("Bad input")
			} else {
				metronome.SetBpm(float64(inputNum))
			}
		}
	}
}
