package channel

import (
	"fmt"
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

	g.fd = fdopplerHz
	g.tInterval = tinterval
}
func (g *GeneratorJakes) Reset(seed uint64) {
	// fmt.Println("Setting State ", seed)
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

func (g *GeneratorJakes) generate1(t float64) complex128 {
	g.lastSampletime = t
	twopi := 2 * math.Pi
	var re, im float64
	// g.alpham = g.rndgen.Rand()
	M := float64(g.N)
	M = 30
	for m := 1.0; m <= 30; m++ {
		am := g.rndgen.Rand()
		bm := g.rndgen.Rand()
		theta := g.rndgen.Rand()

		term1 := (math.Pi*(2*m-1) + theta) / (4 * M)
		fmt.Printf("\n %v [ %v %v %v] %v", m, theta, am, bm, math.Cos(term1))
		re += math.Cos(twopi*g.fd*math.Cos(term1)*t + am)
		im += math.Sin(twopi*g.fd*math.Cos(term1)*t + bm)
	}
	Scale := 1.0 / math.Sqrt(M)
	re = re * Scale
	im = im * Scale
	return complex(re, im)
}

func (g *GeneratorJakes) generate(t float64) complex128 {
	g.lastSampletime = t
	twopi := 2 * math.Pi

	// g.alpham = g.rndgen.Rand()
	M := float64(g.N)
	M = 30
	fd := g.fd
	var cos = math.Cos
	var sin = math.Sin
	// var sqrt = math.Sqrt
	var r complex128
	theta := g.alpham // 0.0 // A random phase for each generator
	for m := 1.0; m <= 30; m++ {

		betham := math.Pi * m / (M + 1)
		alpha := 0.0
		alpham := math.Pi * (m - .5) / (2 * M)
		Am := complex(cos(betham), sin(betham))
		Bm := complex(cos(alpha), sin(alpha))
		fn := fd * cos(alpham)
		r += Am*complex(cos(twopi*fn*t+theta), 0) + 0.707*Bm*complex(cos(twopi*fd*t), 0)
	}
	Scale := 2.8284 // 2*sqrt(2)
	r = r * complex(Scale, 0)
	return r
}
