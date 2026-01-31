package internal

import (
	"math"
	"math/cmplx"
)

// A spectrogram is a visual representation of the spectrum of frequencies of a signal as it varies with time.
// x for time, y for freq, z axis is amplitude of particular freq

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