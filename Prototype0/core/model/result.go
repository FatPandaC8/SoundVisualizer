package model

type MatchResult struct {
	SongID			uint32
	Score			int
	Confidence		float64
	Matched			bool
}