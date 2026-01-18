package model

// write based on mostly the Industrial Strength paper
type Fingerprint struct {
	Hash 			uint32 		`json:"hash"`
	TimeQuery		int32		`json:"time"`
}