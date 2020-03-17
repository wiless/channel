package channel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

//SetDefaults loads the default values for the simulation
func (i *TestEnvironment) DefaultEnv() {

	i.ENV = "RMa"
	i.DS.Mu = []float64{-7.49, 7.43, -7.47}
	i.DS.Sigma = []float64{0.55, 0.48, 0.24}
	i.ASD.Mu = []float64{0.90, 0.95, 0.67}
	i.ASD.Sigma = []float64{0.38, 0.45, 0.18}
	i.ASA.Mu = []float64{1.52, 1.52, 1.66}
	i.ASA.Sigma = []float64{0.24, 0.13, 0.21}
	i.ZSA.Mu = []float64{0.47, 0.58, 0.93}
	i.ZSA.Sigma = []float64{0.40, 0.37, 0.22}
	i.SF = []int{0, 0, 8}
	i.KdB.Mu = []float64{7, 0, 0}
	i.KdB.Sigma = []float64{4, 0, 0}
	i.Crosscorr.SFvsK = []float64{0, 0, 0}
	i.Crosscorr.DSvsSF = []float64{-0.5, -0.5, 0}
	i.Crosscorr.ASDvsSF = []float64{0, 0.6, 0}
	i.Crosscorr.ASAvsSF = []float64{0, 0, 0}
	i.Crosscorr.ZSDvsSF = []float64{0.01, -0.04, 0}
	i.Crosscorr.ZSAvsSF = []float64{-0.17, -0.25, 0}
	i.Crosscorr.DSvsK = []float64{0, 0, 0}
	i.Crosscorr.ASDvsK = []float64{0, 0, 0}
	i.Crosscorr.ASAvsK = []float64{0, 0, 0}
	i.Crosscorr.ZSDvsK = []float64{0, 0, 0}
	i.Crosscorr.ZSAvsK = []float64{0, 0, 0}
	i.Crosscorr.ASDvsDS = []float64{0, -0.4, 0}
	i.Crosscorr.ASAvsDS = []float64{0, 0, 0}
	i.Crosscorr.ZSDvsDS = []float64{-0.05, -0.1, 0}
	i.Crosscorr.ZSAvsDS = []float64{0.27, -0.4, 0}
	i.Crosscorr.ASDvsASA = []float64{0, 0, -0.5}
	i.Crosscorr.ZSDvsASD = []float64{0.73, 0.42, 0.66}
	i.Crosscorr.ZSAvsASD = []float64{-0.14, -0.27, 0.47}
	i.Crosscorr.ZSDvsASA = []float64{-0.20, -0.18, -0.55}
	i.Crosscorr.ZSAvsASA = []float64{0.24, 0.26, -0.22}
	i.Crosscorr.ZSDvsZSA = []float64{-0.07, -0.27, 0}
	i.CorrDist.DS = []int{50, 36, 36}
	i.CorrDist.ASD = []int{25, 25, 30}
	i.CorrDist.ASA = []int{35, 35, 40}
	i.CorrDist.SF = []int{37, 120, 120}
	i.CorrDist.DS = []int{40, 0, 0}
	i.CorrDist.ZSA = []int{15, 50, 50}
	i.CorrDist.ZSD = []int{15, 50, 50}
	i.DelayScalingParameter = []float64{3.8, 1.7, 1.7}
	i.XPR.Mu = []float64{12, 7, 7}
	i.XPR.Sigma = []float64{4, 3, 3}
	i.TotalClusters = []int{11, 10, 10}
	i.RaysInCluster = []int{20, 20, 20}
	i.CDS = []int{0, 0, 0}
	i.CASD = []int{2, 2, 2}
	i.CASA = []int{3, 3, 3}
	i.CZSA = []int{3, 3, 3}
	i.Shadowstd = []int{3, 3, 3}
	i.ZSD.Mu = []float64{0, 0}
	i.ZSD.Offset = []float64{0, 0}

}

// Loads the test environment form a file..
func (i *TestEnvironment) Load(fname string) {

	LoadJson(fname, i)

}

// ReadITUConfig reads all the configuration for the app
func NewEnvironment(configname string) (TestEnvironment, error) {
	var cfg TestEnvironment
	// fmt.Println(InDIR)
	cfg.Load(configname)
	cfg.Initialize()
	return cfg, nil
}

//sqrtXCorr returns L so that (L*L'=xc) xc should be NxN matrix
func sqrtXCorr(xc *mat.Dense) *mat.Dense {
	N := xc.RawMatrix().Rows
	Rc := mat.NewSymDense(N, xc.RawMatrix().Data)
	//	fmt.Printf("Rc =%v\n", Rc)
	//fmt.Printf("\nRc= %v \n\n", mat.Formatted(Rc, mat.Prefix("    ")))

	/// Create & Save the L matrix, (where  L*L'=Rc )
	var sqRC mat.Cholesky
	var Lmat1 mat.TriDense
	if ok := sqRC.Factorize(Rc); !ok {
		fmt.Println("a matrix is not positive semi-definite... Cannot Factorize")
		return nil
	}
	sqRC.LTo(&Lmat1)
	// fmt.Printf("\nL = %v \n\n", mat.Formatted(&Lmat1, mat.Prefix("    ")))
	// Save(Lmat1.RawTriangular(), output)

	sqx := mat.NewDense(N, N, Lmat1.RawTriangular().Data)
	//	fmt.Printf("\nsqx = %v \n\n", mat.Formatted(sqx, mat.Prefix("      ")))
	return sqx
}

func (e *TestEnvironment) Initialize() {

	// (condition string, testenvironment TestEnvironment, fname string)
	los, nlos, o2i := e.generateXCorr()
	e.SqLOS = sqrtXCorr(los)
	e.SqNLOS = sqrtXCorr(nlos)
	e.SqO2I = sqrtXCorr(o2i)

}

func (e *TestEnvironment) generateXCorr() (los, nlos, o2i *mat.Dense) {

	XCorrtempLOS := mat.NewDense(7, 7, nil)
	XCorrtempNLOS := mat.NewDense(7, 7, nil)
	XCorrtempO2I := mat.NewDense(7, 7, nil)
	rows := 21
	col := 3
	a := mat.NewDense(rows, col, nil)

	for j := 0; j < 3; j++ {
		a.Set(0, j, e.Crosscorr.SFvsK[j])
		a.Set(1, j, e.Crosscorr.DSvsSF[j])
		a.Set(2, j, e.Crosscorr.ASDvsSF[j])
		a.Set(3, j, e.Crosscorr.ASAvsSF[j])
		a.Set(4, j, e.Crosscorr.ZSDvsSF[j])
		a.Set(5, j, e.Crosscorr.ZSAvsSF[j])
		a.Set(6, j, e.Crosscorr.DSvsK[j])
		a.Set(7, j, e.Crosscorr.ASDvsK[j])
		a.Set(8, j, e.Crosscorr.ASAvsK[j])
		a.Set(9, j, e.Crosscorr.ZSDvsK[j])
		a.Set(10, j, e.Crosscorr.ZSAvsK[j])
		a.Set(11, j, e.Crosscorr.ASDvsDS[j])
		a.Set(12, j, e.Crosscorr.ASAvsDS[j])
		a.Set(13, j, e.Crosscorr.ZSDvsDS[j])
		a.Set(14, j, e.Crosscorr.ZSAvsDS[j])
		a.Set(15, j, e.Crosscorr.ASDvsASA[j])
		a.Set(16, j, e.Crosscorr.ZSDvsASD[j])
		a.Set(17, j, e.Crosscorr.ZSAvsASD[j])
		a.Set(18, j, e.Crosscorr.ZSDvsASA[j])
		a.Set(19, j, e.Crosscorr.ZSAvsASA[j])
		a.Set(20, j, e.Crosscorr.ZSDvsZSA[j])

	}
	/// XCorr for LOS

	// "LOS"
	{
		i := 0
		arr0 := []float64{1, a.At(0, i), a.At(1, i), a.At(2, i), a.At(3, i), a.At(4, i), a.At(5, i)}
		arr1 := []float64{a.At(0, i), 1, a.At(6, i), a.At(7, i), a.At(8, i), a.At(9, i), a.At(10, i)}
		arr2 := []float64{a.At(1, i), a.At(6, i), 1, a.At(11, i), a.At(12, i), a.At(13, i), a.At(14, i)}
		arr3 := []float64{a.At(2, i), a.At(7, i), a.At(11, i), 1, a.At(15, i), a.At(16, i), a.At(17, i)}
		arr4 := []float64{a.At(3, i), a.At(8, i), a.At(12, i), a.At(15, i), 1, a.At(18, i), a.At(19, i)}
		arr5 := []float64{a.At(4, i), a.At(9, i), a.At(13, i), a.At(16, i), a.At(18, i), 1, a.At(20, i)}
		arr6 := []float64{a.At(5, i), a.At(10, i), a.At(14, i), a.At(17, i), a.At(19, i), a.At(20, i), 1}

		for k := 0; k < len(arr0); k++ {
			XCorrtempLOS.Set(0, k, arr0[k])
			XCorrtempLOS.Set(1, k, arr1[k])
			XCorrtempLOS.Set(2, k, arr2[k])
			XCorrtempLOS.Set(3, k, arr3[k])
			XCorrtempLOS.Set(4, k, arr4[k])
			XCorrtempLOS.Set(5, k, arr5[k])
			XCorrtempLOS.Set(6, k, arr6[k])
		}

	}

	// "NLOS"
	{
		i := 1
		arr0 := []float64{1, a.At(0, i), a.At(1, i), a.At(2, i), a.At(3, i), a.At(4, i), a.At(5, i)}
		arr1 := []float64{a.At(0, i), 1, a.At(6, i), a.At(7, i), a.At(8, i), a.At(9, i), a.At(10, i)}
		arr2 := []float64{a.At(1, i), a.At(6, i), 1, a.At(11, i), a.At(12, i), a.At(13, i), a.At(14, i)}
		arr3 := []float64{a.At(2, i), a.At(7, i), a.At(11, i), 1, a.At(15, i), a.At(16, i), a.At(17, i)}
		arr4 := []float64{a.At(3, i), a.At(8, i), a.At(12, i), a.At(15, i), 1, a.At(18, i), a.At(19, i)}
		arr5 := []float64{a.At(4, i), a.At(9, i), a.At(13, i), a.At(16, i), a.At(18, i), 1, a.At(20, i)}
		arr6 := []float64{a.At(5, i), a.At(10, i), a.At(14, i), a.At(17, i), a.At(19, i), a.At(20, i), 1}

		for k := 0; k < len(arr0); k++ {
			XCorrtempNLOS.Set(0, k, arr0[k])
			XCorrtempNLOS.Set(1, k, arr1[k])
			XCorrtempNLOS.Set(2, k, arr2[k])
			XCorrtempNLOS.Set(3, k, arr3[k])
			XCorrtempNLOS.Set(4, k, arr4[k])
			XCorrtempNLOS.Set(5, k, arr5[k])
			XCorrtempNLOS.Set(6, k, arr6[k])
		}

	}

	// "O2I"
	{
		i := 2
		arr0 := []float64{1, a.At(0, i), a.At(1, i), a.At(2, i), a.At(3, i), a.At(4, i), a.At(5, i)}
		arr1 := []float64{a.At(0, i), 1, a.At(6, i), a.At(7, i), a.At(8, i), a.At(9, i), a.At(10, i)}
		arr2 := []float64{a.At(1, i), a.At(6, i), 1, a.At(11, i), a.At(12, i), a.At(13, i), a.At(14, i)}
		arr3 := []float64{a.At(2, i), a.At(7, i), a.At(11, i), 1, a.At(15, i), a.At(16, i), a.At(17, i)}
		arr4 := []float64{a.At(3, i), a.At(8, i), a.At(12, i), a.At(15, i), 1, a.At(18, i), a.At(19, i)}
		arr5 := []float64{a.At(4, i), a.At(9, i), a.At(13, i), a.At(16, i), a.At(18, i), 1, a.At(20, i)}
		arr6 := []float64{a.At(5, i), a.At(10, i), a.At(14, i), a.At(17, i), a.At(19, i), a.At(20, i), 1}

		for k := 0; k < len(arr0); k++ {
			XCorrtempO2I.Set(0, k, arr0[k])
			XCorrtempO2I.Set(1, k, arr1[k])
			XCorrtempO2I.Set(2, k, arr2[k])
			XCorrtempO2I.Set(3, k, arr3[k])
			XCorrtempO2I.Set(4, k, arr4[k])
			XCorrtempO2I.Set(5, k, arr5[k])
			XCorrtempO2I.Set(6, k, arr6[k])
		}

	}
	return XCorrtempLOS, XCorrtempNLOS, XCorrtempO2I
}

func Load(fname string) *mat.Dense {

	f, er := os.Open(fname)

	if er != nil {
		logrus.Error("Error ", er)
		return nil
	} else {
		data, er := ioutil.ReadAll(f)
		var tmp blas64.General

		if er == nil {
			er := json.Unmarshal(data, &tmp)
			if er != nil {
				fmt.Printf("Error json ", er)
			} else {
				a := mat.NewDense(1, 1, nil)
				a.SetRawMatrix(tmp)
				return a
				// a.SetRawMatrix(tmp)
			}
		} else {
			// logrus.Error("Error ", er, n)
			return nil
		}

	}

	return nil
}

func Save(a *mat.Dense, fname string) {
	//fmt.Printf("Saving ..", mat.Formatted(a, mat.Prefix("    ")))

	b, e := json.Marshal(a.RawMatrix())
	if e != nil {
		logrus.Println("Error ", e)
	} else {
		f, er := os.Create(fname)
		_ = er
		f.Write(b)
		f.Close()
	}
}

func LoadJson(fname string, a interface{}) error {
	f, er := os.Open(fname)
	if er != nil {
		logrus.Error("Error ", er)
		return er
	} else {
		data, er := ioutil.ReadAll(f)
		logrus.Info("Read .. ", string(data))
		if er == nil {
			er := json.Unmarshal(data, a)
			if er != nil {
				fmt.Printf("Error json ", er)
			}
		}
		return er
	}
}
