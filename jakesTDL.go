package channel

import (
	"math"

	"github.com/wiless/vlib"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type rndUniformNTap []distuv.Uniform

// Generator is the simplest fading generator that generates i.i.d samples every time, also optionally supports QuasiStatic
// Default is Gaussian distributed samples..with zero mean, unit variance
type GeneratorTDLJakes struct {
	state          uint64
	quasiDuration  float64
	quasi          bool
	bufferlength   int     // Size of history of samples to be stored
	lastSampletime float64 // time t of the recently sampled
	tInterval      float64
	recentSamples  []complex128
	// a normal distribution generator //
	rndgen   [][]rndUniformNTap /// MIMO
	oldstate uint64
	// Specific params
	N        int // default =20
	fd       float64
	alpham   [][]vlib.VectorF // Each MIMO has vector of initial phases
	basesrc  rand.Source
	nTaps    int
	nTx, nRx int
}

func NewGeneratorTDLJakes(seed uint64, nTx, nRx int) *GeneratorTDLJakes {
	jakes := new(GeneratorTDLJakes)
	jakes.basesrc = rand.NewSource(seed)
	jakes.state = seed
	jakes.nTx, jakes.nRx = nTx, nRx
	jakes.rndgen = make([][]rndUniformNTap, nTx) // distuv.Uniform{Src: rand.NewSource(seed), Min: 0, Max: 2 * math.Pi}
	jakes.alpham = make([][]vlib.VectorF, nTx)

	for i := 0; i < nTx; i++ {
		jakes.rndgen[i] = make([]rndUniformNTap, nRx)
		jakes.alpham[i] = make([]vlib.VectorF, nRx)
	}

	jakes.tInterval = 0
	jakes.lastSampletime = 0
	jakes.N = 18
	jakes.fd = 0
	jakes.nTaps = 0
	return jakes
}
func (g *GeneratorTDLJakes) Init(fdopplerHz float64, tinterval float64) {

	g.fd = fdopplerHz
	g.tInterval = tinterval
}

func (g *GeneratorTDLJakes) Reset(seed uint64) {
	g.basesrc = rand.NewSource((seed))
	g.state = seed

	g.rndgen = make([][]rndUniformNTap, g.nTx) // distuv.Uniform{Src: rand.NewSource(seed), Min: 0, Max: 2 * math.Pi}
	g.alpham = make([][]vlib.VectorF, g.nTx)

	for i := 0; i < g.nTx; i++ {
		g.rndgen[i] = make([]rndUniformNTap, g.nRx)
		g.alpham[i] = make([]vlib.VectorF, g.nRx)
	}

}
func (g *GeneratorTDLJakes) State() uint64 {

	return g.state
}

/// Creation of TAPS and Delays
// CreateTaps creates len(tapPower) taps
func (g *GeneratorTDLJakes) CreateTaps(tapPower []float64) {
	g.nTaps = len(tapPower)
	nTx, nRx := g.nTx, g.nRx
	for i := 0; i < nTx; i++ {
		for j := 0; j < nRx; j++ {
			g.rndgen[i][j] = make([]distuv.Uniform, g.nTaps)
			g.alpham[i][j] = vlib.NewVectorF(g.nTaps)

			for tau := 0; tau < g.nTaps; tau++ {
				s := g.basesrc.Uint64()
				g.rndgen[i][j][tau] = distuv.Uniform{Src: rand.NewSource(s), Min: 0, Max: 2 * math.Pi}
				g.alpham[i][j][tau] = g.rndgen[i][j][tau].Rand() /// generates rv [0,2pi)
			}
		}
	}

}

/// GENERATION OF DATA

// NextSample generates a guassian rv, no depedency on ts
func (g *GeneratorTDLJakes) NextSample(tx, rx int, incr bool) (float64, []complex128) {
	// fmt.Printf("My SEED %v | ", g.seed())
	if incr {
		g.lastSampletime += g.tInterval
	}

	return g.lastSampletime, g.generate(g.lastSampletime, tx, rx)
}

// Generate generates a guassian rv, no depedency on ts
func (g *GeneratorTDLJakes) Generate(ts float64, tx, rx int) []complex128 {
	g.lastSampletime += g.tInterval
	return g.generate(ts, tx, rx)
}

// Generate generates a guassian rv, no depedency on ts
func (g *GeneratorTDLJakes) GenerateN(tstart, tinterval float64, N int, tx, rx int) [][]complex128 {
	g.tInterval = tinterval

	result := make([][]complex128, N)
	t := tstart
	for i, _ := range result {
		result[i] = g.generate(t, tx, rx)
		t += tinterval

	}
	g.lastSampletime = t
	return result
}

// generate implements based on
// https://en.wikipedia.org/wiki/Rayleigh_fading#Jakes's_model
func (g *GeneratorTDLJakes) generate(t float64, tx, rx int) []complex128 {
	g.lastSampletime = t
	twopi := 2 * math.Pi

	// g.alpham = g.rndgen.Rand()
	M := float64(g.N)

	fd := g.fd
	var cos = math.Cos
	var sin = math.Sin
	// var sqrt = math.Sqrt
	cir := vlib.NewVectorC(g.nTaps)
	// Scale := 2.8284 / math.Sqrt(2*M) // 2*sqrt(2)
	Scale := 1.0 / math.Sqrt(2*M) // 2*sqrt(2)

	for tap := 0; tap < g.nTaps; tap++ {
		var r complex128
		theta := g.alpham[tx][rx][tap] // 0.0 // A random phase for each generator
		// fmt.Printf("\n %v [%d,%d] %v", t, tx, rx, theta)
		for m := 1.0; m <= M; m++ {
			betham := math.Pi * m / (M + 1)
			alpha := 0.0

			// // modified Jakes (see wiki)
			// alpha = math.Pi * (m - 0.5) / (2 * M)
			// betham = math.Pi * m / (M)
			// theta = g.alpham[tx][rx][tap] + betham + twopi*(float64(tap-1))/(M+1)
			theta = g.alpham[tx][rx][tap] + twopi*(float64(tap-1))/(M+1)

			alpham := math.Pi * (m - .5) / (2 * M)
			Am := complex(cos(betham), sin(betham))
			Bm := complex(cos(alpha), sin(alpha))

			fn := fd * cos(alpham)
			r += Am*complex(cos(twopi*fn*t+theta), 0) + 0.707*Bm*complex(cos(twopi*fd*t), 0)
		}
		r = r * complex(Scale, 0)
		cir[tap] = r
	}
	return cir
}
