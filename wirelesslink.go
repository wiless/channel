package channel

type BaseParam struct {
	Fc float64 // Carrier Frequency
	BW float64 // Bandwidth of channel used for modelling
}

type WirelessLink struct {
	BaseParam
	TxID      int
	RxID      int
	NTx       int //number of Tx ports // number of transmitters to be modelled
	RTx       int //number of Rx ports  //  number of receiving ports to be modelled
	generator FadeGenerator
}
