package fly

import (
	"fmt"
	"math/rand"
	"sort"
)

type matePair struct {
	mate1 *Fly
	mate2 *Fly
}

/*
helper construct for generating the mate pairs
*/
type cumFitFly struct {
	fly    *Fly
	cumFit float64 //cumulative fitness up to this fly

}

/*
Get mate pairs; has random component
*/
func getMatePairs(flies []Fly, n int64) []matePair {

	// cumulative fitness
	hermaphrodites := generateCumFitness(flies)
	merryCouples := make([]matePair, n)
	for i := int64(0); i < n; i++ {
		riherma1 := rand.Float64()
		riherma2 := rand.Float64()
		herma1 := getFlyForRandomNumber(hermaphrodites, riherma1)
		herma2 := getFlyForRandomNumber(hermaphrodites, riherma2)
		// design decision: I'm not preventing self-mating, as this could lead to an endless loop, eg when all flies except one has fitness 0.0
		merryCouples[i] = matePair{mate1: herma1.fly, mate2: herma2.fly}
	}
	return merryCouples
}

func generateCumFitness(flies []Fly) []cumFitFly {
	// Here major go confusion arose with pointers I guess
	// Video
	//https://www.youtube.com/watch?v=sTFJtxJXkaY
	// & ambersand; read it as address off!
	// strange behavior of range, which kind of overwrote the variable always with the same pointer
	// so that all flies in cumFitFly were pointing to the same fly in the end!!
	// This is a very strange behaviour -> major shit feature of GO
	// it may speed up things if sorted; largest fitness first; and cum sum between zero and 1
	sort.Slice(flies, func(i, j int) bool { return flies[i].Fitness > flies[j].Fitness })

	// now get the sum of all fitnesses
	var fitsum float64 = 0.0
	for _, f := range flies {
		if f.Fitness < 0.0 {
			panic("Fitness must be larger than zero")
		}
		fitsum += f.Fitness
	}
	// generate the cumFitFlies

	cumflies := make([]cumFitFly, 0, len(flies))
	var runningsum float64 = 0.0

	for i, f := range flies {
		fi := &flies[i] // Woa that solves my POINTER BUG!
		// I Could not use &f as this was always referring to the same address
		// BE super careful with range and pointers
		//print(fi)
		w := f.Fitness / fitsum // fitness scaled by the total fitness such that the Sum of all is 1.0
		runningsum += w

		c := cumFitFly{fly: fi, cumFit: runningsum}
		cumflies = append(cumflies, c)
	}
	return cumflies
}

/*
given a cumulative fitness array, pick a fly based on the random number;
fast binary search
*/
func getFlyForRandomNumber(cf []cumFitFly, randomIndex float64) *cumFitFly {
	if randomIndex < 0.0 || randomIndex >= 1.0 {
		panic(fmt.Sprintf("invaldi random index; must be between 0-1; got %f", randomIndex))
	}
	var hi int64 = int64(len(cf))
	var lo int64 = 0
	var mid int64 = 0
	for lo < hi {
		mid = int64((float64(lo+hi) / 2.0))
		midval := cf[mid].cumFit
		if midval < randomIndex {
			lo = mid + 1
		} else if midval > randomIndex {
			hi = mid
		} else {
			return &cf[mid+1]
		}
	}
	toret := &cf[lo]
	return toret
}
