package identify

import (
	"sort"
	"shazam/core/model"
	"shazam/core/repo"
)

func Identify(index repo.FingerprintIndex, query []model.Fingerprint, cfg model.Config) model.MatchResult {
	// Phase 1: candidate generation --	[To perform a search, the above fingerprinting step is
	// 									performed on a captured sample sound file to generate a set
	// 									of hash:time offset records]

	songHitCount := make(map[uint32]int)

	for _, fp := range query {
		hits := index.Lookup(fp.Hash)
		for _, h := range hits {
			songHitCount[h.SongID]++
		}
	}

	if len(songHitCount) == 0 {
		return model.MatchResult{Matched: false}
	}

	type candidate struct {
		songID 		uint32
		count 		int
	}

	candidates := make([]candidate, 0, len(songHitCount))
	for id, c := range songHitCount {
		candidates = append(candidates, candidate{id, c})
	}
	
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].count > candidates[j].count
	})

	if len(candidates) > cfg.TopKCandidate {
		candidates = candidates[:cfg.TopKCandidate]
	}

	// Phase 2: offset clustering -- 	[If the files match, matching
	// 									features should occur at similar relative offsets from the
	// 									beginning of the file, i.e. a sequence of hashes in one file
	// 									should also occur in the matching file with the same
	// 									relative time sequence. The problem of deciding whether a
	// 									match has been found reduces to detecting a significant
	// 									cluster of points forming a diagonal line within the
	// 									scatterplot.]

	var bestSong uint32 		= uint32(0)
	var bestPeak int 			= 0
	var bestConfidence float64 	= 0.0

	for _, cand := range candidates {
		offsetHistogram := make(map[int]int)
		totalVotes := 0

		for _, fp := range query {
			hits := index.Lookup(fp.Hash)
			for _, h := range hits {
				if h.SongID != cand.songID {
					continue
				}

				offset := int(
					(int64(h.TimeDB) - int64(fp.TimeQuery)) / int64(cfg.OffsetQuantization),
				)
				offsetHistogram[offset]++
				totalVotes++
			}

			if totalVotes == 0 {
				continue
			}

			peak := 0
			for _, count := range offsetHistogram {
				if count > peak {
					peak = count
				}
			}

			confidence := float64(peak) / float64(totalVotes)

			if peak >= cfg.MinPeakVotes && confidence >= cfg.MinConfidence && peak > bestPeak {
				bestPeak = peak
				bestSong = cand.songID
				bestConfidence = confidence
			}
		}
	}

	if bestPeak == 0 {
		return model.MatchResult{Matched: false}
	}

	return model.MatchResult {
		SongID:     bestSong,
		Score:      bestPeak,
		Confidence: bestConfidence,
		Matched:    true,
	}
}