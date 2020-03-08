package channel

import (
	"log"

	"github.com/wiless/vlib"
)

type SingleTapChannel struct {
	generator FadeGenerator
	genMIMO   [][]FadeGenerator
	NTx       int //number of Tx ports // number of transmitters to be modelled
	NRx       int //number of Rx ports  //  number of receiving ports to be modelled
}

func (s *SingleTapChannel) NextSample() (lastTs float64, coeff complex128) {

	if s.generator == nil {
		log.Fatal("Link:No Generator Associated..")
		return 0, complex(0, 0)

	} else {
		// fmt.Printf("\n %d My Generator ", w.ID)
	}

	return s.generator.NextSample()
}

func (w SingleTapChannel) IsMIMO() bool {
	if w.NTx > 1 || w.NRx > 1 {
		return true
	} else {
		return false
	}
}

func (w SingleTapChannel) Dims() (tx, rx int) {

	return w.NTx, w.NRx

}

func (w *SingleTapChannel) SetGenerator(gen FadeGenerator) {

	if w.NTx > 1 || w.NRx > 1 {
		log.Fatal("Link is MIMO, Initialize MIMO Generator")
		return
	}
	w.genMIMO = nil
	w.generator = gen
}

func (w *SingleTapChannel) SetMIMO(tx, rx int) {

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

func (w *SingleTapChannel) NextMIMOSample() (lastTs float64, H vlib.MatrixC) {

	if !w.IsMIMO() {
		log.Panic("Link is not MIMO ")
	}
	res := vlib.NewMatrixC(w.NTx, w.NRx)
	for i := 0; i < w.NTx; i++ {
		for j := 0; j < w.NRx; j++ {
			lastTs, res[i][j] = w.genMIMO[i][j].NextSample()
		}
	}
	return lastTs, res
}
