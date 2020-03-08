package channel

import (
	"math"

	"github.com/wiless/vlib"
)

type PDPprofile struct {
	DelayTaus []float64 // delay Tau in seconds
	Power     []float64 // power in linear
	Ts        float64   // Tau, if delays are normalized index
}

/// Normalize the tau with ts interval
func (p *PDPprofile) Set(ts float64, power float64) {
	p.DelayTaus = append(p.DelayTaus, ts)
	p.Power = append(p.Power, power)
}

/// Normalize the tau with ts interval
func (p *PDPprofile) Normalize(ts float64) {
	// find the minimal delta

}

/// NormalizeInterp normalizes through interpolation
func (p *PDPprofile) NormalizeInterp(ts float64) {

	maxTau := vlib.Max(p.DelayTaus)
	Ntaps := int(math.Ceil(maxTau / ts))
	delays := vlib.NewVectorF(Ntaps)
	powers := vlib.NewVectorF(Ntaps)
	for n := 0; n < Ntaps; n++ {
		delays[n] = float64(n) * ts
	}
	onebyts := 1.0 / ts
	for k, v := range p.DelayTaus {
		tt := delays.Sub(v).Scale(onebyts)
		newpdp := vlib.SincF(tt).Scale(math.Sqrt(p.Power[k]))
		powers.PlusEqual(newpdp)
	}
	for n := 0; n < Ntaps; n++ {
		powers[n] = math.Pow(powers[n], 2.0)
	}
	//  for k=1:length(x)
	// plot(newtt,x(k)*sinc((newtt-tt(k))/ts))
	// newpdp=newpdp+x(k)*sinc((newtt-tt(k))/ts);
	// end

}

type TDLGenerator struct {
	genMIMO [][]FadeGenerator
	NTx     int //number of Tx ports // number of transmitters to be modelled
	NRx     int //number of Rx ports  //  number of receiving ports to be modelled
	fres    float64
	tres    float64
	profile PDPprofile
}

//Setup sets up the generator with given PDP profile, for a ntx x nrx MIMO system
func (tdl *TDLGenerator) Setup(pdp PDPprofile, ntx, nrx int) {

}

func (w TDLGenerator) Dims() (tx, rx int) {

	return w.NTx, w.NRx

}

func (w TDLGenerator) IsMIMO() bool {
	if w.NTx > 1 || w.NRx > 1 {
		return true
	} else {
		return false
	}
}

func (w *TDLGenerator) SetMIMO(tx, rx int) {

	w.NTx, w.NRx = tx, rx
	if tx == 1 && rx == 1 {
		w.genMIMO = nil
		return
	}
	w.genMIMO = make([][]FadeGenerator, tx*rx)

	for i := 0; i < w.NTx; i++ {
		w.genMIMO[i] = make([]FadeGenerator, rx)
	}

}
