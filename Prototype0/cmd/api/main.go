package main

import (
	"encoding/json"
	"log"
	"net/http"
	"shazam/adapters/storage/inmemory"
	"shazam/core/identify"
	"shazam/core/model"
)

type IndentifyRequest struct {
	Fingerprints			[]model.Fingerprint 	`json:"fingerprints"`
}

type IndentifyResponse struct {
	SongID					int						`json:"song_id"`
	Match 					bool					`json:"match"`
}

func main() {
	cfg := model.Config{
		MinPeakVotes: 2,
		MinConfidence: 0.0,
		TopKCandidate: 1,
		OffsetQuantization: 1,
	}

	repo := inmemory.New()

	repo.Index[111] = []model.Hit{
		{SongID: 1, TimeDB: 10},
		{SongID: 1, TimeDB: 20},
	}

	repo.Index[222] = []model.Hit{
		{SongID: 1, TimeDB: 30},
	}

	http.HandleFunc("/identify", func(w http.ResponseWriter, r *http.Request) {
		var req IndentifyRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result := identify.Identify(repo, req.Fingerprints, cfg)
		resp := IndentifyResponse{
			SongID: int(result.SongID),
			Match: result.Matched,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}