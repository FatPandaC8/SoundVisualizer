package internal

import (
	"math"
	"math/cmplx"
)

func Normalize(samples []int16) []float64 {
	out := make([]float64, len(samples))
	for i, s := range samples {
		out[i] = float64(s) / 32768.0
	}
	return out
}

func Downsample(samples []float64) []float64 {
	// The only tricky part is that before
	// downsampling a signal, you need to filter the higher
	// frequencies in the sound to avoid aliasing

	temp := make([]float64, len(samples))
	for i := 1; i < len(samples) - 1; i++ {
		temp[i] = (float64(samples[i-1]) + float64(samples[i]) + float64(samples[i+1])) / 3.0
	}

	// Downsample part: of grouping 4 at a time
	floatSample := make([]float64, 0, len(samples) / 4)
	for i := 0; i < len(samples); i += 4 {
		floatSample = append(floatSample, temp[i])
	}

	return floatSample
}

// FFT(x) = FFT(x_even) + W_N^k * FFT(x_odd)

// Where:
//   x_even = [x[0], x[2], x[4], ...]
//   x_odd  = [x[1], x[3], x[5], ...]
//   W_N^k  = e^(-2πik/N)  (twiddle factor)
//   N      = signal length
//   k      = frequency bin index

func FFT(x []complex128) []complex128 {
	n := len(x)
	if n == 1 {
		return x
	}
	
	even := make([]complex128, n /2)
	odd := make([]complex128, n /2)
	
	for i := 0; i < n / 2; i++ {
		even[i] = x[2*i]
		odd[i] = x[2*i + 1]
	}
	
	evenFFT := FFT(even)
	oddFFT := FFT(odd)
	
	result := make([]complex128, n)
	for k := 0; k < n / 2; k++ {
		angle := -2 * math.Pi * float64(k) / float64(n)
		twiddle := cmplx.Rect(1, angle)
		
		t := twiddle * oddFFT[k]
		result[k] = evenFFT[k] + t
		result[k + n/2] = evenFFT[k] - t
	}
	
	return result
}

// Must FFT with a sliding window of the samples

// func Window(samples []float64) []float64 {
	
// }