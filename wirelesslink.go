package channel

import (
	"log"
)

type BaseParam struct {
	FcGHz float64 // Carrier Frequency
	BWMHz float64 // Bandwidth of channel used for modelling
}

type WirelessLink struct {
	BaseParam
	TxID      int
	RxID      int
	NTx       int //number of Tx ports // number of transmitters to be modelled
	NRx       int //number of Rx ports  //  number of receiving ports to be modelled
	generator FadeGenerator
	lastTs    float64 // recent Timesamples

}

func (w *WirelessLink) NextSample() complex128 {
	if w.generator == nil {
		log.Fatal("Link:No Generator Associated..")
	}
	var coeff complex128
	w.lastTs, coeff = w.generator.NextSample()
	return coeff
}

func (w *WirelessLink) State() uint64 {
	return w.generator.State()
}
