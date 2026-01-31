package internal
// https://hasan-hasanov.com/post/2023/10/how_to_parse_wav_file/

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func ReadWav(filepath string) (*AudioData, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var riff [12]byte
	if _, err := io.ReadFull(file, riff[:]); err != nil {
		return nil, err
	}

	if string(riff[0:4]) != "RIFF" || string(riff[8:12]) != "WAVE" {
		return nil, fmt.Errorf("not a WAV file")
	}

	var (
		audioFormat   	uint16
		numChannels   	uint16
		sampleRate    	uint32
		byteRate      	uint32
		blockAlign    	uint16
		bitsPerSample 	uint16
		data          	[]byte
		samples 		[]int16
	)

	for {
		var subchunkID [4]byte
		if _, err := io.ReadFull(file, subchunkID[:]); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("cannot read subchunk id")
		}

		var subchunkSize uint32
		if err := binary.Read(file, binary.LittleEndian, &subchunkSize); err != nil {
			return nil, err
		}

		switch string(subchunkID[:]) {

		case "fmt ":
			start, _ := file.Seek(0, io.SeekCurrent)

			binary.Read(file, binary.LittleEndian, &audioFormat)
			binary.Read(file, binary.LittleEndian, &numChannels)
			binary.Read(file, binary.LittleEndian, &sampleRate)
			binary.Read(file, binary.LittleEndian, &byteRate)
			binary.Read(file, binary.LittleEndian, &blockAlign)
			binary.Read(file, binary.LittleEndian, &bitsPerSample)

			fmt.Println("AudioFormat:", audioFormat)
			fmt.Println("Channels:", numChannels)
			fmt.Println("SampleRate:", sampleRate)
			fmt.Println("BitsPerSample:", bitsPerSample)
			if bitsPerSample != 16 {
				return nil, fmt.Errorf("only 16-bit PCM supported (got %d-bit)", bitsPerSample)
			}

			readNow, _ := file.Seek(0, io.SeekCurrent)
			remaining := int64(subchunkSize) - (readNow - start)
			if remaining > 0 {
				io.CopyN(io.Discard, file, remaining)
			}

		case "data":
			data = make([]byte, subchunkSize)
			if _, err := io.ReadFull(file, data); err != nil {
				return nil, err
			}
			if len(data) == 0 {
				return nil, fmt.Errorf("no data chunk found")
			}
			fmt.Println("Read data bytes:", len(data))

		default:
			fmt.Println("JUNK")
			io.CopyN(io.Discard, file, int64(subchunkSize)) // discard this subchunksize as its subchunkID is junk
		}
	}

	if audioFormat != 1 {
		return nil, fmt.Errorf("compressed WAV not supported (format=%d)", audioFormat)
	}

	numSamples := len(data) / 2
	samples = make([]int16, numSamples*int(numChannels))

	for i := 0; i < numSamples; i++ {
		samples[i] = int16(binary.LittleEndian.Uint16(data[i*2 : (i+1)*2]))
	}

	if numChannels == 2 {
		mono := make([]int16, numSamples)
		for i := 0; i < numSamples; i++ {
			mono[i] = (samples[i*2] + samples[i*2 + 1]) / 2
		}
		samples = mono
		numChannels = 1
	}
	duration := float64(len(data)) / float64(byteRate)

	return &AudioData{
		Samples:    samples,
		SampleRate: int(sampleRate),
		Channels:   int(numChannels),
		Duration:   duration,
	}, nil
}

