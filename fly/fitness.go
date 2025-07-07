package fly

import (
	"math"
)

type FitnessFunction struct {
	x           float64
	t           float64
	b           float64 // benefit of each cluster insertion
	noxincluins bool
}

var ff FitnessFunction
var minimumFitness float64

func (f *FitnessFunction) ComputeFitness(counttotal int64, countcluster int64) float64 {
	frc := counttotal // fitness relevant count
	if f.noxincluins {
		frc -= countcluster // if clusterinsertions are not deleterios subtract them
	}

	// equation: w = 1 - x*n^t + b*countcluster
	negeffect := f.x * math.Pow(float64(frc), f.t)
	pos := f.b * float64(countcluster)
	fit := 1.0 - negeffect + pos
	if fit < 0 {
		fit = 0.0
	}
	if fit > 1.0 {
		fit = 1.0
	}
	return fit
}

func SetupFitness(x float64, t float64, b float64, noxincluins bool, minFitness float64) {
	ff = FitnessFunction{x: x, t: t, b: b, noxincluins: noxincluins}
	minimumFitness = minFitness
}

func GetFitness(f *Fly) float64 {
	total := f.CountTotalInsertions()
	return ff.ComputeFitness(total, f.FlyStat.CountCluster)
}
