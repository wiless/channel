package main

import (
	"fmt"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/wiless/channel"
	"github.com/wiless/vlib"
)

func main() {
	// Reset the seed
	rand.Seed(time.Now().UnixNano())
	fcGHz := 0.700

	speedkmph := 30.0 // walking speed
	fd := channel.DopplerHz(speedkmph, fcGHz)
	Ts := 1e-4
	log.Infof("Doppler (Hz) %4.3fHz", fd)
	env := channel.NewSimpleEnv()
	env.Setup(fcGHz, 10) // fc=700MHz, Bw=10Mhz

	DoMIMO(env, fd, Ts)
	fmt.Println()

}

func DoMIMO(env *channel.Env, fd, Ts float64) {

	// MIMO example
	env.Create(1, 2, 2)

	// Example to set all links to an i.i.d generator
	var pdp channel.PDPprofile
	pdp.CreateUPower(4, 1e-3)

	env.SetupTDLJakes(fd, Ts, pdp)

	x := make([]complex128, 2)
	x[0] = complex(1, 0)
	x[1] = complex(2, 0)
	N := 4 // 100 samples
	// hh := vlib.NewVectorF(N)
	// hh2 := vlib.NewVectorF(N)
	tt := vlib.NewVectorF(N)

	for _, link := range env.Links {
		for t := 0; t < N; t++ {
			{

				x := vlib.RandQPSK(2, 1)
				_ = x
				H := link.NextTDLSample()
				tt[t] = link.LastTsample()

				hh := channel.MIMOCoeff(H)
				fmt.Printf("\nt=%.2es %v", tt[t], hh)
				// h := link.NextSample()

				// y := RxSamples(H, x)
				// _ = idx
				// _ = y

				// hh[t] = cmplx.Abs(H[0][0])

				// fmt.Printf("\nLink (%d) t=%f ", idx, link.LastTsample())

				// fmt.Printf("\nx=%v", x.MatString())
				// fmt.Printf("\nH=%v", H.MatString())
				// fmt.Printf("\ny=%v", y.MatString())
			}
		}
	}

	fmt.Println("\n t=", tt)

}

// RxSamples Returns y=H*x
func RxSamples(H vlib.MatrixC, x vlib.VectorC) vlib.VectorC {

	// mH := mat.NewCDense(H.NRows(), H.NCols(), H.Data())
	// Initialize two matrices, a and b.
	// b := mat.NewCDense(H.NRows(), 1, x)
	// Take the matrix product of a and b and place the result in c.
	result := vlib.NewVectorC(H.NRows())
	for i := 0; i < H.NRows(); i++ {
		h := H.GetRow(i)
		result[i] = vlib.Dotu(h, x)
	}

	return result
}
