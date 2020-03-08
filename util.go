package channel

import "gonum.org/v1/gonum/unit/constant"

// DopplerHz returns the doppler frequency for the velocity v
// fd=v/Lamda
func DopplerHz(vKmph float64, fcGHz float64) (fHz float64) {
	Lamda := constant.LightSpeedInVacuum.Unit().Value() / (fcGHz * 1e9)
	fHz = (vKmph * 1000 / 3600) / Lamda
	return fHz
}
