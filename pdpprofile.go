package channel

import (
	"math"
	//"github.com/martinlindhe/unit" for later use
	"github.com/martinlindhe/unit"
	"github.com/wiless/vlib"
)

type PDPprofile struct {
	DelayTaus []float64     // delay Tau in seconds
	Power     []float64     // powerGain in linear
	Ts        unit.Duration // Ts, Sample Interval in uSec
}

/// Normalize the tau with ts interval
func (p *PDPprofile) Reset() {
	p = &PDPprofile{}
}

func (p *PDPprofile) SetDelayUsec(v []float64) {
	p.DelayTaus = v

}
func (p *PDPprofile) GetDelaynSec() []float64 {
	return vlib.VectorF(p.DelayTaus).Scale(1e3) /// Delay is already in u Seconds, x1000, will return NanoSeconds
}

/// Normalize the tau with ts interval
func (p *PDPprofile) AppendTaps(ts unit.Duration, power float64) {
	p.DelayTaus = append(p.DelayTaus, ts.Microseconds())
	p.Power = append(p.Power, power)
}

/// Extrapolate coeffs to NTaps by interpolating at every Ts (µs), set NTaps=-1 to auto select based on PDP.MaxDelay
/// returns samples, and nts=timeInterval in µs
func (p PDPprofile) Extrapolate(Ts unit.Duration, NTaps int, tapcoeff []complex128) (samples []complex128, ntS []float64) {
	maxTau := vlib.Max(p.DelayTaus)
	taps := int(math.Ceil(maxTau/Ts.Microseconds())) + 1
	ntaps := int(math.Max(float64(taps), float64(NTaps)))
	delays := vlib.NewVectorF(ntaps)
	result := vlib.NewVectorC(ntaps)
	for n := 0; n < ntaps; n++ {
		delays[n] = float64(n) * Ts.Microseconds()
	}
	onebyts := 1.0 / Ts.Microseconds()
	for k, v := range p.DelayTaus {
		tt := delays.Sub(v).Scale(onebyts)
		newpdp := vlib.ToVectorC(vlib.SincF(tt)).ScaleC(tapcoeff[k])
		result.PlusEqual(newpdp)
	}

	return result[0:NTaps], delays

}

/// NormalizeCoeff resamples coeffs given at DelayTaus into spacing of Ts
func (p *PDPprofile) NormalizeCoeff(tapcoeff []complex128) []complex128 {
	// find the minimal delta

	maxTau := vlib.Max(p.DelayTaus)
	Ntaps := int(math.Ceil(maxTau/p.Ts.Microseconds())) + 1
	delays := vlib.NewVectorF(Ntaps)
	result := vlib.NewVectorC(Ntaps)
	for n := 0; n < Ntaps; n++ {
		delays[n] = float64(n) * p.Ts.Microseconds()
	}
	onebyts := 1.0 / p.Ts
	for k, v := range p.DelayTaus {
		tt := delays.Sub(v).Scale(onebyts.Microseconds())
		newpdp := vlib.ToVectorC(vlib.SincF(tt)).ScaleC(tapcoeff[k])
		result.PlusEqual(newpdp)
	}

	return result

}

// CreateUPower creates a uniform power nTap channel, with Ts uSeconds spacing
func (p *PDPprofile) CreateUPower(nTaps int, Ts unit.Duration) {
	scale := math.Sqrt(1.0 / float64(nTaps))
	p.DelayTaus = vlib.NewSegmentF(0, Ts.Microseconds(), nTaps)
	p.Power = vlib.NewOnesF(nTaps).Scale(scale)
	p.Ts = unit.Duration(Ts) * unit.Microsecond
}

/// NormalizeInterp normalizes through interpolation ts in Seconds
func (p PDPprofile) InterpolatePDP(tusec unit.Duration) PDPprofile {
	var newpdp PDPprofile
	newpdp.Ts = tusec
	ts := tusec
	maxTau := vlib.Max(p.DelayTaus)
	Ntaps := int(math.Ceil(maxTau/ts.Microseconds())) + 1
	delays := vlib.NewVectorF(Ntaps)
	powers := vlib.NewVectorF(Ntaps)
	for n := 0; n < Ntaps; n++ {
		delays[n] = float64(n) * ts.Microseconds()
	}
	onebyts := 1.0 / ts.Microseconds()
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

}
