package fly

import (
	"math"
)

// https://blog.knoldus.com/how-to-use-interfaces-in-golang/#:~:text=An%20interface%20in%20Go%20is,a%20similar%20type%20of%20object.

type FitnessFunctionLinear struct {
	x           float64
	t           float64
	noxincluins bool
}

type FitnessFunctionMultiplicative struct {
	x           float64
	noxincluins bool
}

type IFitnessFunction interface {
	ComputeFitness(int64, int64, int64) float64
}

var iff IFitnessFunction

//var minimumFitness float64
//var maximumInsertions float64

func (f FitnessFunctionMultiplicative) ComputeFitness(counttotal int64, countcluster int64, countreference int64) float64 {
	frc := counttotal // fitness relevant count
	if f.noxincluins {
		frc -= countcluster   // if clusterinsertions are not deleterios subtract them
		frc -= countreference // reference insertions need to be subtracted as well
	}
	// equation is (1-x)^n
	fit := math.Pow(1-f.x, float64(frc))
	return fit
}

func (f FitnessFunctionLinear) ComputeFitness(counttotal int64, countcluster int64, countreference int64) float64 {
	frc := counttotal // fitness relevant count
	if f.noxincluins {
		frc -= countcluster   // if clusterinsertions are not deleterios subtract them
		frc -= countreference // reference insertions need to be subtracted as well
	}

	// equation is 1.0- n*x^t
	negeffect := f.x * math.Pow(float64(frc), f.t)
	fit := 1.0 - negeffect
	if fit < 0 {
		fit = 0.0
	}
	return fit
}

func SetupFitness(x float64, t float64, noxincluins bool, multiplicative bool) {
	if multiplicative {
		if math.Abs(t-1.0) > 0.0001 {
			panic("epistatic effects not supported for multiplicative fitness, i.e. --t must be 1.0")
		}
		iff = FitnessFunctionMultiplicative{x: x, noxincluins: noxincluins}
	} else {
		iff = FitnessFunctionLinear{x: x, t: t, noxincluins: noxincluins}
	}
}

func GetFitness(f *Fly) float64 {
	total := f.CountTotalInsertions()
	return iff.ComputeFitness(total, f.FlyStat.CountCluster, f.FlyStat.CountReference)
}
