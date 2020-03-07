package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/wiless/channel"
	"github.com/wiless/vlib"
)

func main() {
	// Reset the seed
	rand.Seed(time.Now().UnixNano())

	env := channel.NewSimpleEnv()
	env.Setup(0.700, 10) // fc=700MHz, Bw=10Mhz

	// // SISO example
	// env.Create(3, 1, 1)

	// // Example to set all links to an i.i.d generator
	// env.AttachGeneratorIID()
	// for idx, link := range env.Links {
	// 	for t := 0.0; t < 5; t++ {
	// 		{
	// 			fmt.Printf("\nLink (%d) t=%f : %v", idx, t, link.NextSample())
	// 		}
	// 	}
	// }

	fmt.Println()

	// MIMO example
	env.Create(1, 2, 2)

	// Example to set all links to an i.i.d generator
	env.AttachGeneratorJakes(10, 1e-3)

	x := make([]complex128, 2)
	x[0] = complex(1, 0)
	x[1] = complex(2, 0)
	for idx, link := range env.Links {
		for t := 0.0; t < 15; t++ {
			{

				x := vlib.RandQPSK(2, 1)
				H := link.NextMIMOSample()
				y := RxSamples(H, x)
				fmt.Printf("\n\n\nLink (%d) t=%f ", idx, link.LastTsample())
				fmt.Printf("\nx=%v", x.MatString())
				fmt.Printf("\nH=%v", H.MatString())
				fmt.Printf("\ny=%v", y.MatString())
			}
		}
	}
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
