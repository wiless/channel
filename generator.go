package channel

import "gonum.org/v1/gonum/stat/distuv"

// FadeGenerator generic interface that is expected to generate a single-tap complex fading coefficient sample for a given time
// Generator is expected to model a time-variation model of the fading
type FadeGenerator interface {
	Generate(ts float64) complex128                          // Generate a sample for the time t (in sec)
	GenerateN(tstart, tinterval float64, N int) []complex128 // Generate N samples starting from tstart with duration of Tinterval
}

// Generator is the simplest fading generator that generates i.i.d samples every time, also optionally supports QuasiStatic
// Default is Gaussian distributed samples..with zero mean, unit variance
type GeneratorIID struct {
	quasiDuration  float64
	quasi          bool
	bufferlength   int     // Size of history of samples to be stored
	lastSampletime float64 // time t of the recently sampled
	recentSamples  []complex128
	// a normal distribution generator //
	rndgen distuv.Normal
}

func NewGeneratorIID() *GeneratorIID {
	iid := new(GeneratorIID)
	iid.rndgen = distuv.UnitNormal

	// distuv.Normal{
	// 	Mu:    0,
	// 	Sigma: 1.0,
	// }
	return iid
}
