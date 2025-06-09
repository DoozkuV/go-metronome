package audio

import (
	"bytes"
	_ "embed"
	"log"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

var (
	//go:embed sounds/metronome-hit.wav
	audioFile []byte
	fileBpm   float64
)

// A simple wrapper around the metronome
type metronome struct {
	*beep.Ctrl
	Vol       *effects.Volume
	resampler *beep.Resampler
}

// Creates a metronome at the specified BPM
func NewMetronome(bpm float64) *metronome {
	// Initialize the streamer from the embedded metronome sound
	streamer, format, err := wav.Decode(bytes.NewReader(audioFile))
	if err != nil {
		log.Fatal(err)
	}
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	sr := format.SampleRate
	// Calculate the length and bpm of the audio file
	if sr > 0 {
		numFrames := buffer.Len() // Total number of frames in the stream
		audioFileLength := sr.D(numFrames)
		fileBpm = (1 / audioFileLength.Minutes())
	} else {
		panic("Invalid audio file")
	}

	// Initialize the speaker
	speaker.Init(sr, sr.N(time.Second/10))

	// Set up the effects chain for the streamer
	s := buffer.Streamer(0, buffer.Len())
	loop := beep.Loop(-1, s)
	ctrl := &beep.Ctrl{Streamer: loop, Paused: true}
	volume := &effects.Volume{Streamer: ctrl, Base: 2, Volume: 0, Silent: false}
	resampler := beep.ResampleRatio(4, bpm/fileBpm, volume)

	speaker.Play(resampler)
	return &metronome{Ctrl: ctrl, resampler: resampler, Vol: volume}
}

func (m *metronome) SetBpm(bpm float64) {
	speaker.Lock()
	defer speaker.Unlock()
	m.resampler.SetRatio(bpm / fileBpm)
}

func (m *metronome) Bpm() float64 {
	return fileBpm * m.resampler.Ratio()
}
