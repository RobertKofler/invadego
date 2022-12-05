package env

import (
	"fmt"
	"invade/util"
	"math"
	"math/rand"
	"sort"
)

/*
The switch of the haplotype must be made before the TE;
If recombination site is 0 and a TE is at position 0, i must first introduce a recombination event and than take the TE;
RECOMBINATION FIRST
Recombination first!

now with the recombination sites i need two loops

loop 1:
just go over each chromosome and ask if there is random segregation of the chromosome,
providing the very first position of each chromosome if there is a recombination event

loop2:
go over each recombination window and get the recombination events

finalize:
a) merge loop1 and loop2
b) set, must be unique
c) sort and return
*/

type RecombinationWindow struct {
	genint GenomicInterval
	lambda float64
}

func (rw *RecombinationWindow) getRecombinationNumber() int64 {
	return util.Poisson(rw.lambda)
}

// Get a random position in the recombination window;
// the last position is included as possible recombination sites
// NOTE the first site is not included, because the first site recombines on random assortment of chromosomes and should not be used for regular cross-over events
//
// TODO: exclude first site?
func (rw *RecombinationWindow) getRandomPosition() int64 {
	// AT
	// 01
	// 100 = start; 101 = end
	// len = 2
	// randoffest 0/1 = rand.Intn
	// randpos = start + randoffset
	//
	// include first site (random assortment -> should not be used for regular recombination events )
	//gilen := rw.genint.Length()
	//randOffset := rand.Intn(int(gilen))
	//return rw.genint.Start + int64(randOffset)

	// code excluding first site -> first site is reserved for random assortment of chromosomes;
	gilen := rw.genint.Length()
	randOffset := rand.Intn(int(gilen - 1))
	return rw.genint.Start + 1 + int64(randOffset)
}

func GetRecombinationEvents() []int64 {
	// recombination events in set; avoid double events
	var recombinationEvents = make(map[int64]bool)
	// first
	// handle the random assortment of chromosomes
	for _, chromosomeOffset := range env.genome.offsets {
		if rand.Float64() < 0.5 {
			recombinationEvents[chromosomeOffset] = true
		}
	}
	// second
	// handle the recombination events of each Window
	for _, rw := range env.recombinationWindows {
		recEvents := rw.getRecombinationNumber()
		for i := 0; i < int(recEvents); i++ {
			randpos := rw.getRandomPosition()
			recombinationEvents[randpos] = true
		}
	}
	// sort the recombinatin events by position
	keys := make([]int64, len(recombinationEvents))
	i := 0
	for k := range recombinationEvents {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys

}

/*
Translate the recombination rate from cM/Mb to lambda, i.e the mean recombination events in a window;
Rely on haldane-transformation
*/
func cmpmb2lambda(cmpmb float64, winsize int64) float64 {
	/*
	   // Recombination rate is in cM
	    double recFraction=cmpmb/100;
	    // return 0 for small recombination rates
	    double cWin=1000000.0/windowsize;
	    double distance=haldane1919mapFunction(recFraction);
	    double lambdawin=distance/cWin;
	    return lambdawin;
	*/
	recfraction := cmpmb / 100.0
	distance := haldane1919mapFunction(recfraction)

	cWin := 1000000.0 / float64(winsize)
	lambdawin := distance / cWin
	return lambdawin
}

/*
Haldanes map function
*/
func haldane1919mapFunction(rf float64) float64 {
	lambda := -0.5 * math.Log(1.0-2.0*rf)
	return lambda
}

/*
Get the recombination windows for some genomic Intervals and the recombination rate in cM/Mb
*/
func getRecombinationWindows(genIntervals []GenomicInterval, recRate []float64) []*RecombinationWindow {
	if len(genIntervals) != len(recRate) {
		panic(fmt.Sprintf("Different lengths of genomic intervals (%d) and recombination rates (%d)", len(genIntervals), len(recRate)))
	}
	recwins := make([]*RecombinationWindow, 0, len(genIntervals))
	for i, rr := range recRate {
		gi := genIntervals[i]
		lambda := cmpmb2lambda(rr, gi.Length())
		rwin := RecombinationWindow{genint: gi, lambda: lambda}
		recwins = append(recwins, &rwin)
	}
	return recwins
}
