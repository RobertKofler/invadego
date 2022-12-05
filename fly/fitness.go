package fly

import (
	"math"
)

type FitnessFunction struct {
	x           float64
	t           float64
	noxincluins bool
}

var ff FitnessFunction
var minimumFitness float64

func (f *FitnessFunction) ComputeFitness(counttotal int64, countcluster int64, countreference int64) float64 {
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

func SetupFitness(x float64, t float64, noxincluins bool, minFitness float64) {
	ff = FitnessFunction{x: x, t: t, noxincluins: noxincluins}
	minimumFitness = minFitness
}

func GetFitness(f *Fly) float64 {
	total := f.CountTotalInsertions()
	return ff.ComputeFitness(total, f.FlyStat.CountCluster, f.FlyStat.CountReference)
}
