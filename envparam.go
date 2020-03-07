package channel

import "gonum.org/v1/gonum/mat"

// TestEnvironment ....
type TestEnvironment struct {
	ENV string `json:"ENV"`
	DS  struct {
		Mu    []float64 `json:"Mu"`
		Sigma []float64 `json:"Sigma"`
	} `json:"DS"`
	ASD struct {
		Mu    []float64 `json:"Mu"`
		Sigma []float64 `json:"Sigma"`
	} `json:"ASD"`
	ASA struct {
		Mu    []float64 `json:"Mu"`
		Sigma []float64 `json:"Sigma"`
	} `json:"ASA"`
	ZSA struct {
		Mu    []float64 `json:"Mu"`
		Sigma []float64 `json:"Sigma"`
	} `json:"ZSA"`
	SF  []int `json:"SF"`
	KdB struct {
		Mu    []float64 `json:"Mu"`
		Sigma []float64 `json:"Sigma"`
	} `json:"KdB"`
	Crosscorr struct {
		SFvsK    []float64 `json:"SFvsK"`
		DSvsSF   []float64 `json:"DSvsSF"`
		ASDvsSF  []float64 `json:"ASDvsSF"`
		ASAvsSF  []float64 `json:"ASAvsSF"`
		ZSDvsSF  []float64 `json:"ZSDvsSF"`
		ZSAvsSF  []float64 `json:"ZSAvsSF"`
		DSvsK    []float64 `json:"DSvsK"`
		ASDvsK   []float64 `json:"ASDvsK"`
		ASAvsK   []float64 `json:"ASAvsK"`
		ZSDvsK   []float64 `json:"ZSDvsK"`
		ZSAvsK   []float64 `json:"ZSAvsK"`
		ASDvsDS  []float64 `json:"ASDvsDS"`
		ASAvsDS  []float64 `json:"ASAvsDS"`
		ZSDvsDS  []float64 `json:"ZSDvsDS"`
		ZSAvsDS  []float64 `json:"ZSAvsDS"`
		ASDvsASA []float64 `json:"ASDvsASA"`
		ZSDvsASD []float64 `json:"ZSDvsASD"`
		ZSAvsASD []float64 `json:"ZSAvsASD"`
		ZSDvsASA []float64 `json:"ZSDvsASA"`
		ZSAvsASA []float64 `json:"ZSAvsASA"`
		ZSDvsZSA []float64 `json:"ZSDvsZSA"`
	} `json:"Crosscorr"`
	CorrDist struct {
		DS  []int `json:"DS"`
		ASD []int `json:"ASD"`
		ASA []int `json:"ASA"`
		SF  []int `json:"SF"`
		K   []int `json:"K"`
		ZSA []int `json:"ZSA"`
		ZSD []int `json:"ZSD"`
	} `json:"CorrDist"`
	DelayScalingParameter []float64 `json:"DelayScalingParameter"`
	XPR                   struct {
		Mu    []float64 `json:"Mu"`
		Sigma []float64 `json:"Sigma"`
	} `json:"XPR"`
	TotalClusters []int `json:"TotalClusters"`
	RaysInCluster []int `json:"RaysInCluster"`
	CDS           []int `json:"CDS"`
	CASD          []int `json:"CASD"`
	CASA          []int `json:"CASA"`
	CZSA          []int `json:"CZSA"`
	Shadowstd     []int `json:"Shadowstd"`
	ZSD           struct {
		Mu     []float64 `json:"Mu"`
		Offset []float64 `json:"Offset"`
	} `json:"ZSD"`
	fname string `json:"fname"`

	SqLOS, SqNLOS, SqO2I *mat.Dense `json:"-"`
}
