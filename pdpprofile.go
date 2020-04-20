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
func (p *PDPprofile) Reset() {
	p = &PDPprofile{}
}

/// Normalize the tau with ts interval
func (p *PDPprofile) AppendTaps(ts float64, power float64) {
	p.DelayTaus = append(p.DelayTaus, ts)
	p.Power = append(p.Power, power)
}

/// NormalizeCoeff resamples coeffs given at DelayTaus into spacing of Ts
func (p *PDPprofile) NormalizeCoeff(tapcoeff []complex128) []complex128 {
	// find the minimal delta

	maxTau := vlib.Max(p.DelayTaus)
	Ntaps := int(math.Ceil(maxTau / p.Ts))
	delays := vlib.NewVectorF(Ntaps)
	result := vlib.NewVectorC(Ntaps)
	for n := 0; n < Ntaps; n++ {
		delays[n] = float64(n) * p.Ts
	}
	onebyts := 1.0 / p.Ts
	for k, v := range p.DelayTaus {
		tt := delays.Sub(v).Scale(onebyts)
		newpdp := vlib.ToVectorC(vlib.SincF(tt)).ScaleC(tapcoeff[k])
		result.PlusEqual(newpdp)
	}
	return result
	//  for k=1:length(x)
	// plot(newtt,x(k)*sinc((newtt-tt(k))/ts))
	// newpdp=newpdp+x(k)*sinc((newtt-tt(k))/ts);
	// end
}

// CreateUPower creates a uniform power nTap channel
func (p *PDPprofile) CreateUPower(nTaps int, Ts float64) {
	scale := math.Sqrt(1.0 / float64(nTaps))
	p.DelayTaus = vlib.NewSegmentF(0, Ts, nTaps)
	p.Power = vlib.NewOnesF(nTaps).Scale(scale)
	p.Ts = Ts
}

/// NormalizeInterp normalizes through interpolation ts in Seconds
func (p PDPprofile) NormalizeInterp(ts float64) PDPprofile {
	var newpdp PDPprofile
	newpdp.Ts = ts

	maxTau := vlib.Max(p.DelayTaus)
	Ntaps := int(math.Ceil(maxTau/ts)) + 1
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
	newpdp.DelayTaus = delays
	newpdp.Power = powers
	return newpdp
	//  for k=1:length(x)
	// plot(newtt,x(k)*sinc((newtt-tt(k))/ts))
	// newpdp=newpdp+x(k)*sinc((newtt-tt(k))/ts);
	// end

}
