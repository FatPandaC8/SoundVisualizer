package model

type Config struct {
	TopKCandidate		int			// candidate pruning
	OffsetQuantization	uint32		// frames per bucket
	MinPeakVotes		int			// reject weak matches
	MinConfidence		float64		// reject noisy matches
}