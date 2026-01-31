package internal

type AudioData struct {
	// Sample: a finite part of a sound
	Samples []int16
	SampleRate int
	Channels int
	Duration float64
}