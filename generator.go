package channel

// FadeGenerator generic interface that is expected to generate a single-tap complex fading coefficient sample for a given time
// Generator is expected to model a time-variation model of the fading
type FadeGenerator interface {
	Reset(uint64)
	State() uint64
	NextSample() (float64, complex128)
	Generate(ts float64) complex128                          // Generate a sample for the time t (in sec)
	GenerateN(tstart, tinterval float64, N int) []complex128 // Generate N samples starting from tstart with duration of Tinterval
}
