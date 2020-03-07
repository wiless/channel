// Package channel provides mechanisms to model wireless medium used for simulating wireless links
// Supports M.2412 based environments and generation of fading generation for each links (SISO and MIMO)
package channel

import (
	"log"
	"math/rand"
)

// Stores the wireless environment related paramters
// Each environment can have multiple wireless link-pairs (SRC-DEST)
type Env struct {
	EnvParams TestEnvironment //TestEnvironment Parameters based on in M.2412
	NLinks    int             // Number of wireless links
	Links     []*WirelessLink
	base      BaseParam
}

// NewChannel creates a Default Envinronment
func NewEnv() *Env {
	env := new(Env)
	env.EnvParams = DefaultEnv()
	env.EnvParams.Initialize()
	return env
}

// NewSimpleEnv creates a Default Rural based Wireless Channel Env, with single Tx-Rx links, with IID generator
func NewSimpleEnv() *Env {
	env := new(Env)
	env.EnvParams = DefaultEnv()
	env.EnvParams.Initialize()
	return env
}

func (e *Env) Setup(fGHz float64, bwMHz float64) {
	e.base.FcGHz = fGHz
	e.base.BWMHz = bwMHz
}

// Create creates nlinks with each link having NxM mimo configuration where mimo=[N,M], mimo=[] then 1x1 system is assumed
func (e *Env) Create(nlinks int, ntx, nrx int) {
	e.NLinks = nlinks
	e.Links = make([]*WirelessLink, nlinks)
	log.Printf("\n Creating %d x %d Links ", ntx, nrx)
	for i, _ := range e.Links {
		e.Links[i] = new(WirelessLink)
		e.Links[i].ID = i
		e.Links[i].BaseParam = e.base
		e.Links[i].SetMIMO(ntx, nrx)
	}
}

// AttachGenerator attaches the fading generator fg,
// if clone=true all fading generator has same seed
func (e *Env) AttachGeneratorIID() {
	if len(e.Links) == 0 {
		return
	}

	for i := 0; i < len(e.Links); i++ {
		if e.Links[i].IsMIMO() {
			M, N := e.Links[i].Dims()

			for m := 0; m < M; m++ {
				for n := 0; n < N; n++ {
					state := rand.Uint64()
					e.Links[i].genMIMO[m][n] = NewGeneratorIID(state)
				}
			}

		} else {
			state := rand.Uint64()
			iid := NewGeneratorIID(state)
			e.Links[i].SetGenerator(iid)
		}

	}
}

// AttachGenerator attaches the fading generator fg,
// if clone=true all fading generator has same seed
func (e *Env) AttachGeneratorJakes(fd, Ts float64) {
	if len(e.Links) == 0 {
		return
	}

	for i := 0; i < len(e.Links); i++ {
		if e.Links[i].IsMIMO() {
			M, N := e.Links[i].Dims()

			for m := 0; m < M; m++ {
				for n := 0; n < N; n++ {
					state := rand.Uint64()
					jakes := NewGeneratorJakes(state)
					Ts := 1e-3 // 1ms sampling
					jakes.Init(7.0, Ts)
					e.Links[i].genMIMO[m][n] = jakes
				}
			}

		} else {
			state := rand.Uint64()
			jakes := NewGeneratorJakes(state)
			e.Links[i].SetGenerator(jakes)
		}

	}
}
