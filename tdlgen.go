package channel

import (
	"fmt"

	"github.com/wiless/vlib"
)

type MIMOCoeff [][]vlib.VectorC

type TDLChannel struct {
	genMIMOtdl TDLFadeGenerator
	NTx        int //number of Tx ports // number of transmitters to be modelled
	NRx        int //number of Rx ports  //  number of receiving ports to be modelled
	fres       float64
	tres       float64
	profile    PDPprofile
	tInterval  float64
	t          float64
}

func (tdl *TDLChannel) SetGenerator(tdlgen TDLFadeGenerator) {
	tdl.genMIMOtdl = tdlgen
}

//Setup sets up the generator with given PDP profile, for a ntx x nrx MIMO system
func (tdl *TDLChannel) Setup(pdp PDPprofile, ntx, nrx int, Ts float64) {
	tdl.profile = pdp
	tdl.tInterval = Ts
	tdl.SetMIMO(ntx, nrx)
}

func (tdl *TDLChannel) SetMIMO(ntx, nrx int) {
	tdl.NTx = ntx
	tdl.NRx = nrx
}

func (w TDLChannel) Dims() (tx, rx int) {

	return w.NTx, w.NRx

}

func (w TDLChannel) IsMIMO() bool {
	if w.NTx > 1 || w.NRx > 1 {
		return true
	} else {
		return false
	}
}

//Ht returns a time-domain impulse response of the channel
func (w *TDLChannel) NextSampleTime() float64 {
	return w.t + w.tInterval
}

//Ht returns a time-domain impulse response of the channel
func (w *TDLChannel) Ht(t float64, tx, rx int) []complex128 {
	w.t = t
	return w.genMIMOtdl.Generate(t, tx, rx)
}

//Ht returns a time-domain impulse response of the channel
func (w *TDLChannel) Hmimot(t float64) [][]vlib.VectorC {
	M, N := w.Dims()
	w.t = t
	var result [][]vlib.VectorC
	result = make([][]vlib.VectorC, M)
	for m := 0; m < M; m++ {
		result[m] = make([]vlib.VectorC, N)
		for n := 0; n < N; n++ {
			result[m][n] = w.genMIMOtdl.Generate(t, m, n)
		}
	}
	return result
}

func (H MIMOCoeff) String() string {
	M := len(H)
	N := len(H[0])
	var str string
	for m := 0; m < M; m++ {
		for n := 0; n < N; n++ {
			str += fmt.Sprintf("\n[%d,%d]=%v", m, n, H[m][n])
		}
	}
	return str
}
