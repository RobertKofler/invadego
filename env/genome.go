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
	refRegions           RegionCollection
	paramutables         *RecurrentSites
	triggers             *RecurrentSites
	recombinationWindows []*RecombinationWindow
	minimumFitness       float64
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
Get a random insertio site in the genome;
0-based; ranges from 0 to totalGenome-1
*/
func GetRandomSite() int64 {
	return int64(rand.Intn(int(env.genome.totalGenome)))
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
	/*
	   if(pos>this.genomeSize)throw new IllegalArgumentException("Invalid position outside genome");
	   int chromosome=this.chrSizes.size();
	   for(int i=this.offsets.size()-1; i>=0; i--)
	   {
	       int os=this.offsets.get(i);
	       if(pos>=os){
	           int chrbasedpos=pos-os;
	           return new ChromosomeBasedInsertion(chromosome,chrbasedpos);
	       }
	           chromosome--;
	   }
	   throw new IllegalArgumentException("Invalid position"+pos);
	*/

}
