package fly

import (
	"math"
)

// https://blog.knoldus.com/how-to-use-interfaces-in-golang/#:~:text=An%20interface%20in%20Go%20is,a%20similar%20type%20of%20object.

type FitnessFunctionLinear struct {
	x float64
	t float64
}

type FitnessFunctionHeterozygote struct {
	s float64
	h float64
}

type FitnessFunctionMultiplicative struct {
	x float64
}

type IFitnessFunction interface {
	ComputeFitness(*Fly) float64
}

var iff IFitnessFunction

//var minimumFitness float64
//var maximumInsertions float64

func (f FitnessFunctionMultiplicative) ComputeFitness(fly *Fly) float64 {
	frc := fly.CountTotalInsertions()

	// equation is (1-x)^n
	fit := math.Pow(1-f.x, float64(frc))
	return fit
}

func (f FitnessFunctionHeterozygote) ComputeFitness(fly *Fly) float64 {
	homo, hetero := fly.CountHomozygousHeterozygous()
	fit := 1.0
	for i := int64(0); i < homo; i++ {
		fit *= (1 - f.s)
	}
	for i := int64(0); i < hetero; i++ {
		fit *= (1 - f.s*f.h)
	}
	return fit
}

func (f FitnessFunctionLinear) ComputeFitness(fly *Fly) float64 {
	frc := fly.CountTotalInsertions() // fitness relevant count

	// equation is 1.0- n*x^t
	negeffect := f.x * math.Pow(float64(frc), f.t)
	fit := 1.0 - negeffect
	if fit < 0 {
		fit = 0.0
	}
	return fit
}

func SetupFitness(s float64, h float64) {
	iff = FitnessFunctionHeterozygote{s: s, h: h}
}

func GetFitness(f *Fly) float64 {
	return iff.ComputeFitness(f)
}
