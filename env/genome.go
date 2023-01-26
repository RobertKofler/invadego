// the biological environment of the simulations, e.g. chromosomes, recombination,
// transposition rate, piRNA cluster
package env

import (
	"fmt"
	"math/rand"
)

type Environment struct {
	genome               *GenomicLandscape
	clusters             RegionCollection
	nonClusters          RegionCollection
	recombinationWindows []*RecombinationWindow
	minimumFitness       float64
	mubias               float64
}

var env Environment

/*
An interval in a genome;
Start and End are both within the interval;
If Start == End the length of the interval is 1
*/
type GenomicInterval struct {
	Start int64
	End   int64
}

func (gi GenomicInterval) Length() int64 {
	return gi.End - gi.Start + 1
}

//
// ********** Genomic Landscape ****************+
//

// Construct a GenomicLandscape based on a list of chromosome sizes
func newGenomicLandscape(lens []int64) *GenomicLandscape {
	var offsets []int64
	var chrmSizes []int64
	var intervals []GenomicInterval
	var totGenome int64
	var currentOffset int64
	for _, len := range lens {
		totGenome += len
		offsets = append(offsets, int64(currentOffset))
		end := currentOffset + len - 1
		gi := GenomicInterval{Start: currentOffset, End: end}
		intervals = append(intervals, gi)
		chrmSizes = append(chrmSizes, len)
		currentOffset = end + 1
	}
	return &GenomicLandscape{offsets: offsets, chrmSizes: chrmSizes, intervals: intervals, totalGenome: totGenome}
}

/*
A 0-based linear representation of a genome.
Chromosomes are modeled as GenomicIntervals that occupy a range in a linear integer space.
The offset determines the start position of each chromosome in the linear space
*/
type GenomicLandscape struct {
	offsets     []int64
	chrmSizes   []int64
	intervals   []GenomicInterval
	totalGenome int64
}

/*
likely not needed with insertion bias; TODO -> commented out to avoid problems
Get a random insertio site in the genome;
0-based; ranges from 0 to totalGenome-1
*/
func GetRandomSite() int64 {
	return int64(rand.Intn(int(env.genome.totalGenome)))
}

// Get random site WITHIN a piRNA cluster
func GetRandomClusterSite() int64 {
	//
	ci := int64(rand.Intn(int(env.clusters.Size())))
	return env.clusters.Positions[ci]
}

// Get random site OUTSIDE of a piRNA cluster
func GetRandomNonClusterSite() int64 {
	//
	nci := int64(rand.Intn(int(env.nonClusters.Size())))
	return env.nonClusters.Positions[nci]
}

/*
Insertion probability into a piRNA cluster when the insertion bias is considered
with a bias of 0 -> p=fc (size of the cluster eg 3%)
with a bias of -1 -> p=0
with a bias of 1 -> p=1
acts as the threshold for random number between 0 and 1
*/
func GetInsertionProbabilityWithBias(insertionbias float64) float64 {
	clusterfrac := float64(env.clusters.Size()) / float64(env.genome.totalGenome)
	return getProbForBiasAndClusi(insertionbias, clusterfrac)
}

// helper function for insertion bias, allows for more easy debugging
func getProbForBiasAndClusi(insertionbias float64, clufrac float64) float64 {
	if insertionbias < -1.0 || insertionbias > 1.0 {
		panic(fmt.Sprintf("invalid insertion bias, must be between -1.0 and 1.0; got %f", insertionbias))
	}
	if clufrac < 0.0 || clufrac > 1.0 {
		panic(fmt.Sprintf("invalid cluster fraction; must be between 0.0 and 1.0; got %f", clufrac))
	}

	genomefrac := 1 - clufrac
	clusterfit := (insertionbias + 1.0) / 2.0
	genomefit := 1.0 - clusterfit
	totalfit := clufrac*clusterfit + genomefrac*genomefit

	threshold := clufrac * clusterfit / totalfit

	return threshold
}

// get novel insertion sites for a given bias (-1.0 to 1.0)
func GetSitesForBias(numberofsites int64, insertionbias float64) []int64 {
	//

	sites := make([]int64, numberofsites)
	threshold := GetInsertionProbabilityWithBias(insertionbias)
	for i := 0; i < int(numberofsites); i++ {
		if rand.Float64() < threshold {
			// cluster
			sites[i] = GetRandomClusterSite()
		} else {
			// non cluster
			sites[i] = GetRandomNonClusterSite()
		}

	}
	return sites

}

// TODO TEST
func TranslateCoordinates(pos int64) (int64, int64) {
	if pos >= env.genome.totalGenome {
		panic("invalid genomic position; larger than genome")
	}
	for i := len(env.genome.offsets) - 1; i >= 0; i-- {
		curos := env.genome.offsets[i]
		if pos >= curos {
			chrnum := int64(i + 1)
			chrpos := pos - curos + 1
			return chrnum, chrpos
		}
	}
	panic(fmt.Sprintf("invalid index; smaller than allowed %d", pos))

}
