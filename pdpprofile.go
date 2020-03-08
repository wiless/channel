package channel

import (
	"math"

	"github.com/wiless/vlib"
)

type PDPprofile struct {
	DelayTaus []float64 // delay Tau in seconds
	Power     []float64 // powerGain in linear
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

// CreateUPower creates a uniform power nTap channel
func (p *PDPprofile) CreateUPower(nTaps int, Ts float64) {
	scale := math.Sqrt(1.0 / float64(nTaps))
	p.DelayTaus = vlib.NewSegmentF(0, Ts, nTaps)
	p.Power = vlib.NewOnesF(nTaps).Scale(scale)
	p.Ts = Ts
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
