package audio

import (
	"bytes"
	_ "embed"
	"log"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

//go:embed sounds/metronome-hit.wav
var audioFile []byte

// TODO: Figure out a better solution than hard-coding this
const (
	audioFileLength float64 = 0.417959
	fileBpm         float64 = (1 / audioFileLength) * 60 // BPM
)

// A simple wrapper around the metronome
type Metronome struct {
	Ctrl      *beep.Ctrl
	resampler *beep.Resampler
}

// Creates a metronome at the specified BPM
func NewMetronome(bpm float64) *Metronome {
	// Initialize the streamer from the embedded metronome sound
	// TODO: Load this directly into memory
	streamer, format, err := wav.Decode(bytes.NewReader(audioFile))
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	// Initialize the speaker
	sr := format.SampleRate
	speaker.Init(sr, sr.N(time.Second/10))

	// Set up the effects chain for the streamer
	loop := beep.Loop(-1, streamer)
	ctrl := &beep.Ctrl{Streamer: loop, Paused: true}
	// Speed of the sound
	resampler := beep.ResampleRatio(4, bpm/fileBpm, ctrl)

	speaker.Play(resampler)
	return &Metronome{Ctrl: ctrl, resampler: resampler}
}

func (m *Metronome) SetBpm(bpm float64) {
	speaker.Lock()
	m.resampler.SetRatio(bpm / fileBpm)
	speaker.Unlock()
}

func (m *Metronome) Bpm() (bpm float64) {
	bpm = fileBpm * m.resampler.Ratio()
	return
}
