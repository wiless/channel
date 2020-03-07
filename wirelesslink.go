package channel

import (
	"log"

	"github.com/wiless/vlib"
)

type BaseParam struct {
	FcGHz float64 // Carrier Frequency
	BWMHz float64 // Bandwidth of channel used for modelling
}

type WirelessLink struct {
	ID int
	BaseParam
	TxID      int
	RxID      int
	NTx       int //number of Tx ports // number of transmitters to be modelled
	NRx       int //number of Rx ports  //  number of receiving ports to be modelled
	generator FadeGenerator
	genMIMO   [][]FadeGenerator
	lastTs    float64 // recent Timesamples

}

func (w *WirelessLink) NextSample() complex128 {
	if w.generator == nil {
		if w.IsMIMO() {
			log.Fatal("MIMO Link:Call NextMIMOSample() instead")
			return complex(0, 0)
		}
		log.Fatal("Link:No Generator Associated..")
	} else {
		// fmt.Printf("\n %d My Generator ", w.ID)
	}
	var coeff complex128
	w.lastTs, coeff = w.generator.NextSample()
	return coeff
}

func (w *WirelessLink) State() uint64 {
	return w.generator.State()
}

func (w *WirelessLink) ResetGenerator(id uint64) {
	w.generator.Reset(id)
}
func (w *WirelessLink) SetMIMO(tx, rx int) {

	w.NTx, w.NRx = tx, rx
	if tx == 1 && rx == 1 {
		w.genMIMO = nil
		return
	}
	w.genMIMO = make([][]FadeGenerator, tx*rx)
	w.generator = nil
	for i := 0; i < w.NTx; i++ {
		w.genMIMO[i] = make([]FadeGenerator, rx)
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
	return w.NTx, w.NRx
}
func (w *WirelessLink) SetGenerator(gen FadeGenerator) {

	if w.NTx > 1 || w.NRx > 1 {
		log.Fatal("Link is MIMO, Initialize MIMO Generator")
		return
	}
	w.genMIMO = nil
	w.generator = gen
}

func (w *WirelessLink) H(t float64) vlib.MatrixC {
	return w.NextMIMOSample()
}
func (w *WirelessLink) NextMIMOSample() vlib.MatrixC {
	res := vlib.NewMatrixC(w.NTx, w.NRx)
	for i := 0; i < w.NTx; i++ {
		for j := 0; j < w.NRx; j++ {
			_, res[i][j] = w.genMIMO[i][j].NextSample()
		}
	}
	return res
}
