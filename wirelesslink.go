package channel

import (
	"log"

	"github.com/wiless/vlib"
	"golang.org/x/exp/rand"
)

type BaseParam struct {
	FcGHz float64 // Carrier Frequency
	BWMHz float64 // Bandwidth of channel used for modelling
}

type WirelessLink struct {
	ID       int
	baseSeed []uint64
	BaseParam
	TxID       int
	RxID       int
	NTx, NRx   int
	singlegen  *SingleTapChannel
	cirgen     *TDLGenerator // always returns a vector for each Tx-Rx pair..
	lastTs     float64       // recent Timesamples
	flatFading bool
}

func (w WirelessLink) IsFlatFading() bool {
	return w.flatFading
}

func (w *WirelessLink) SetFlatFading(stc *SingleTapChannel) {
	if stc == nil {
		w.flatFading = false
		w.singlegen = nil
	}
	w.singlegen = stc
}

func (w *WirelessLink) NextSample() complex128 {

	if !w.IsFlatFading() {
		log.Fatal("MIMO Link:Call NextMIMOSample() instead")
		return complex(0, 0)
	}
	var coeff complex128
	w.lastTs, coeff = w.singlegen.NextSample()
	return coeff

}

func (w *WirelessLink) State() []uint64 {
	return w.baseSeed
}

// func (w *WirelessLink) ResetGenerator(id uint64) {
// 	w.generator.Reset(id)
// }
func (w *WirelessLink) SetMIMO(tx, rx int) {
	w.NTx = tx
	w.NRx = rx
	if w.IsFlatFading() && w.singlegen != nil {
		w.singlegen.SetMIMO(tx, rx)
	}

	if !w.IsFlatFading() && w.cirgen != nil {
		w.cirgen.SetMIMO(tx, rx)
	}

}

func (w WirelessLink) IsMIMO() bool {
	if w.NTx > 1 || w.NRx > 1 {
		return true
	} else {
		return false
	}

}
func (w WirelessLink) Dims() (tx, rx int) {

	return w.singlegen.Dims()

	if w.IsFlatFading() {
		return w.singlegen.Dims()
	} else {
		// return w.cirgen.IsMIMO()
		return w.cirgen.Dims()
	}

}

func (w *WirelessLink) H(t float64) vlib.MatrixC {
	return w.NextMIMOSample()
}
func (w *WirelessLink) NextMIMOSample() vlib.MatrixC {

	var H vlib.MatrixC
	w.lastTs, H = w.singlegen.NextMIMOSample()
	return H
}

func (w *WirelessLink) LastTsample() float64 {
	return w.lastTs
}

func (w *WirelessLink) SetupSingleTapIID() {

	w.singlegen = new(SingleTapChannel)
	w.flatFading = true
	w.singlegen.SetMIMO(w.NTx, w.NRx)
	M, N := w.Dims()
	w.baseSeed = make([]uint64, M*N)

	if w.IsMIMO() {

		for m := 0; m < M; m++ {
			for n := 0; n < N; n++ {
				state := rand.Uint64()
				w.baseSeed[m*N+n] = state
				w.singlegen.genMIMO[m][n] = NewGeneratorIID(state)
			}
		}

	} else {
		state := rand.Uint64()
		iid := NewGeneratorIID(state)
		w.singlegen.generator = iid // NewGeneratorIID(state)
	}

}

// AttachGenerator attaches the fading generator fg,
// if clone=true all fading generator has same seed
func (w *WirelessLink) SetupSingleTapJakes(fd, Ts float64) {

	w.singlegen = new(SingleTapChannel)
	w.flatFading = true
	w.singlegen.SetMIMO(w.NTx, w.NRx)
	M, N := w.Dims()
	w.baseSeed = make([]uint64, M*N)

	if w.IsMIMO() {
		for m := 0; m < M; m++ {
			for n := 0; n < N; n++ {
				state := rand.Uint64()
				w.baseSeed[m*N+n] = state
				jakes := NewGeneratorJakes(state)
				jakes.Init(fd, Ts)
				w.singlegen.genMIMO[m][n] = jakes
			}
		}

	} else {
		state := rand.Uint64()
		w.baseSeed[0] = state
		jakes := NewGeneratorJakes(state)
		jakes.Init(fd, Ts)
		w.singlegen.generator = jakes
	}

}
