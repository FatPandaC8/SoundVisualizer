package repo

import (
	"shazam/core/model"
)

type FingerprintIndex interface {
	Lookup(hash uint32) []model.Hit 
}