package channel

import (
	"math"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// Generator is the simplest fading generator that generates i.i.d samples every time, also optionally supports QuasiStatic
// Default is Gaussian distributed samples..with zero mean, unit variance
type GeneratorJakes struct {
	state          uint64
	quasiDuration  float64
	quasi          bool
	bufferlength   int     // Size of history of samples to be stored
	lastSampletime float64 // time t of the recently sampled
	tInterval      float64
	recentSamples  []complex128
	// a normal distribution generator //
	rndgen   distuv.Uniform
	oldstate uint64
	// Specific params
	N      int // default =20
	fd     float64
	alpham float64
}

func NewGeneratorJakes(seed uint64) *GeneratorJakes {

	iid := new(GeneratorJakes)
	iid.rndgen = distuv.Uniform{Src: rand.NewSource(seed), Min: 0, Max: 2 * math.Pi}
	iid.state = seed
	iid.tInterval = 1.0
	iid.lastSampletime = 0
	iid.N = 20
	iid.fd = 0
	iid.alpham = iid.rndgen.Rand()
	// logrus.Infof("Created Seed %v |  %v", iid.rndgen.Src.Uint64(), seed)
	// distuv.Normal{
	// 	Mu:    0,
	// 	Sigma: 1.0,
	// }
	return iid
}
func (g *GeneratorJakes) Init(fdopplerHz float64, tinterval float64) {
	g.N = 20
	g.fd = fdopplerHz
	g.tInterval = tinterval
}
func (g *GeneratorJakes) Reset(seed uint64) {
	// fmt.Println("Setting State ", seed)
	g.N = 20
	g.state = seed
	g.rndgen.Src.Seed(seed)
	g.alpham = g.rndgen.Rand()
	// fmt.Println("Test.. ", g.rndgen.Rand())

}
func (g *GeneratorJakes) State() uint64 {

	return g.state
}

// NextSample generates a guassian rv, no depedency on ts
func (g *GeneratorJakes) NextSample() (float64, complex128) {
	// fmt.Printf("My SEED %v | ", g.seed())
	g.lastSampletime += g.tInterval
	return g.lastSampletime, g.generate(g.lastSampletime)
}

// Generate generates a guassian rv, no depedency on ts
func (g *GeneratorJakes) Generate(ts float64) complex128 {

	return g.generate(ts)
}

// Generate generates a guassian rv, no depedency on ts
func (g *GeneratorJakes) GenerateN(tstart, tinterval float64, N int) []complex128 {
	g.tInterval = tinterval
	result := make([]complex128, N)
	t := tstart
	for i, _ := range result {
		result[i] = g.generate(t)
		t += tinterval
	}
	return result
}

func (g *GeneratorJakes) generate(t float64) complex128 {
	g.lastSampletime = t
	twopi := 2 * math.Pi
	var re, im float64
	for n := 0; n < g.N; n++ {
		am := g.rndgen.Rand()
		bm := g.rndgen.Rand()
		re += math.Cos(twopi*math.Cos(g.alpham)*g.fd*t + am)
		im += math.Sin(twopi*math.Cos(g.alpham)*g.fd*t + bm)
	}
	return complex(re, im)
}
