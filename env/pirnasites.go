package env

import (
	"fmt"
)

type RegionCollection []GenomicInterval

/*
Total size of some regions
*/
func (r RegionCollection) Size() int64 {
	var i int64
	for _, k := range r {
		i += k.Length()
	}
	return i
}

/*
Is a certain position within a region collection
*/
func (r RegionCollection) IsInRegion(position int64) bool {
	// must be sorted by start position of genomic interval; panic if not
	var lastpos int64
	for _, gi := range r {
		if gi.Start < lastpos {
			panic(fmt.Sprintf("Error while searching position of %d; Genomic intervals are not sorted  %d !< %d", position, gi.Start, lastpos))
		}
		if position < gi.Start {
			break
		}
		if position >= gi.Start && position <= gi.End {
			return true
		}
		lastpos = gi.Start
	}
	return false
}

/*
for a haploid genome, count the number of the haploid insertions
*/
func CountHaploidInsertions(positions []int64) int64 {
	var hapins int64
	// Two steps
	// Step 1: separate cluster, reference and no-cluster/no-reference insertions
	for _, p := range positions {
		if p >= env.genome.totalGenome {
			panic(fmt.Sprintf("position outside genome %d", p))
		}
		hapins++

	}
	return hapins
}

/*
for a diploid genome, count the number of the following insertions "cluster, reference, paramutable, trigger and noe (non of either)";
return in the given order
*/
func CountDiploidInsertions(hap1 []int64, hap2 []int64) int64 {
	h1 := CountHaploidInsertions(hap1)
	h2 := CountHaploidInsertions(hap2)
	return h1 + h2
}
