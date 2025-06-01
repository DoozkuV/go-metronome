package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

func main() {
	f, err := os.Open("sounds/metronome-hit.wav")
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	fmt.Println("Playing my cool sound!")
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)
	select {}

}
