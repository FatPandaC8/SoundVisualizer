package main

import (
	"log"
	"shazam/internal"
)

func main() {
	_, err := internal.ReadWav("songs/Uoc_Mua.wav");
	if err != nil {
		log.Fatal(err)
	}
}