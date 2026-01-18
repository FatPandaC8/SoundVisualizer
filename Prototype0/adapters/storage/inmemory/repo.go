package inmemory

import "shazam/core/model"

type Repo struct {
	Index map[uint32][]model.Hit
}

func (r *Repo) Lookup(hash uint32) []model.Hit {
	return r.Index[hash]
}

func New() *Repo {
	return &Repo{
		Index: make(map[uint32][]model.Hit),
	}
}
