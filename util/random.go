package util

import (
	"math"
	"math/rand"
	"time"
)

func Mean(data []int64) float64 {
	var sum float64
	for _, i := range data {
		sum += float64(i)
	}
	m := sum / float64(len(data))
	return m
}

/*
Population variance - computed as SOS/n
*/
func Variance(data []int64) float64 {
	mean := Mean(data)
	var sos float64
	for _, i := range data {
		res := float64(i) - mean
		sos += math.Pow(res, 2.0)
	}
	variance := sos / float64(len(data))
	return variance
}

func SetSeed(seed int64) int64 {
	if seed != -1 {
		InvadeLogger.Printf("Will use seed provided by user: %d", seed)
		rand.Seed(seed)
		return seed
	} else {
		rseed := time.Now().UnixNano()
		InvadeLogger.Printf("Will use current time as seed: %d", rseed)
		rand.Seed(rseed)
		return rseed
	}
}

/*
Deprecated
has numerical problem when lambda >>700
func PoissonSimple(lambda float64) int64 {
	// algorithm given by Knuth
	L := math.Pow(math.E, -lambda)
	var k int64 = 0
	var p float64 = 1.0

	for p > L {
		k++
		p *= rand.Float64()
	}
	ret := k - 1
	if ret < 0 {
		panic("invalid Poisson number; must not be negative")
	}
	return (ret)
}
*/

/*
Poisson distributed random numbers; works even when lambda >>700
*/
func Poisson(lambda float64) int64 {
	// perfect implementation; solves numerical problem when lambda >700?
	lleft := lambda
	step := 500.0
	p := 1.0
	var k int64 = 0
	for ok := true; ok; ok = p > 1.0 {
		k++
		r := rand.Float64()
		p = p * r
		for p < 1.0 && lleft > 0.0 {
			if lleft > step {
				p *= math.Exp(step)
				lleft -= step
			} else {
				p *= math.Exp(lleft)
				lleft = 0
			}
		}
	}
	ret := k - 1
	if ret < 0 {
		panic("invalid Poisson number; must not be negative")
	}
	return ret
}
