package channel

import (
	"log"
	"time"

	"github.com/wiless/vlib"
	"golang.org/x/exp/rand"
)

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

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
	cirgen     *TDLChannel // always returns a vector for each Tx-Rx pair..
	lastTs     float64     // recent Timesamples
	flatFading bool
	ready      bool
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
		log.Fatal("Not a Flat Fading Channel")
		return complex(0, 0)
	} else {
		if w.IsMIMO() {
			log.Fatal("Not a SISO Flat Fading Channel, call NextMIMOSample()")
			return complex(0, 0)
		} else {
			var coeff complex128
			w.lastTs, coeff = w.singlegen.NextSample()
			return coeff
		}
	}

}

func (w *WirelessLink) NextTDLSample() [][]vlib.VectorC {

	if w.IsFlatFading() {
		log.Fatal("Not a TDL Fading Channel")
		return make([][]vlib.VectorC, 0)
	} else {
		if w.ready {
			w.lastTs = w.cirgen.NextSampleTime()
			coeff := w.cirgen.Hmimot(w.lastTs)
			return coeff
		} else {
			log.Fatal("Seems No Generator attached to the TDLChannel")
			return make([][]vlib.VectorC, 0)
		}

	}

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
	if w.IsFlatFading() {
		var H vlib.MatrixC
		w.lastTs, H = w.singlegen.NextMIMOSample()
		return H
	}

	log.Panicf("Not a Flat Fading Channel, Call appropriate Generator")

	return vlib.NewMatrixC(0, 0)

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

// AttachGenerator attaches the fading generator fg,
// if clone=true all fading generator has same seed
func (w *WirelessLink) SetupTDLJakes(fd, Ts float64, pdp PDPprofile) {
	state := rand.Uint64()

	w.cirgen = new(TDLChannel)
	w.cirgen.Setup(pdp, w.NTx, w.NRx, Ts)
	w.flatFading = false

	jakestdl := NewGeneratorTDLJakes(state, w.NTx, w.NRx)
	jakestdl.Init(fd, Ts)
	jakestdl.CreateTaps(pdp.Power)
	w.cirgen.genMIMOtdl = jakestdl

}

// AttachGenerator attaches the fading generator fg,
func (w *WirelessLink) AttachM2412(m2412tdl *TDLChannel) {
	state := rand.Uint64()
	w.baseSeed = []uint64{state}
	w.cirgen = m2412tdl
	w.flatFading = false

	if w.cirgen.genMIMOtdl != nil {
		w.ready = true
	}
	// jakestdl := NewGeneratorTDLJakes(state, w.NTx, w.NRx)
	// jakestdl.Init(fd, Ts)
	// jakestdl.CreateTaps(pdp.Power)
	// w.cirgen.genMIMOtdl = m2412gen

}
