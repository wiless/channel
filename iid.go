package channel

import (
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// Generator is the simplest fading generator that generates i.i.d samples every time, also optionally supports QuasiStatic
// Default is Gaussian distributed samples..with zero mean, unit variance
type GeneratorIID struct {
	state          uint64
	quasiDuration  float64
	quasi          bool
	bufferlength   int     // Size of history of samples to be stored
	lastSampletime float64 // time t of the recently sampled
	tInterval      float64
	recentSamples  []complex128
	// a normal distribution generator //
	rndgen   distuv.Normal
	oldstate uint64
}

func NewGeneratorIID(seed uint64) *GeneratorIID {

	iid := new(GeneratorIID)
	iid.rndgen = distuv.Normal{Mu: 0, Sigma: 1, Src: rand.NewSource(seed)}
	iid.state = seed
	iid.tInterval = 1.0
	iid.lastSampletime = 0
	// distuv.Normal{
	// 	Mu:    0,
	// 	Sigma: 1.0,
	// }
	return iid
}

func (g *GeneratorIID) Reset(seed uint64) {
	// fmt.Println("Setting State ", seed)
	g.state = seed
	g.rndgen.Src.Seed(seed)
	// fmt.Println("Test.. ", g.rndgen.Rand())

}
func (g *GeneratorIID) State() uint64 {

	return g.state
}

// NextSample generates a guassian rv, no depedency on ts
func (g *GeneratorIID) NextSample() (float64, complex128) {
	// fmt.Printf("My SEED %v | ", g.seed())

	g.lastSampletime += g.tInterval

	return g.lastSampletime, complex(g.rndgen.Rand(), g.rndgen.Rand())
}

// Generate generates a guassian rv, no depedency on ts
func (g *GeneratorIID) Generate(ts float64) complex128 {
	g.lastSampletime = ts
	return complex(g.rndgen.Rand(), g.rndgen.Rand())
}

// Generate generates a guassian rv, no depedency on ts
func (g *GeneratorIID) GenerateN(tstart, tinterval float64, N int) []complex128 {
	g.tInterval = tinterval
	g.lastSampletime = tstart + float64(N-1)*tinterval
	result := make([]complex128, N)
	for i, _ := range result {
		result[i] = complex(g.rndgen.Rand(), g.rndgen.Rand())
	}
	return result
}
