// Package channel provides mechanisms to model wireless medium used for simulating wireless links
// Supports M.2412 based environments and generation of fading generation for each links (SISO and MIMO)
package channel

// Stores the wireless environment related paramters
// Each environment can have multiple wireless link-pairs (SRC-DEST)
type Env struct {
	EnvParams TestEnvironment //TestEnvironment Parameters based on in M.2412
	NLinks    int             // Number of wireless links
	Links     []WirelessLink
}



