package env

import (
	"fmt"
	"invade/util"
)

/*
Count the number of regions
*/
func (r RegionCollection) Count() int64 {
	return int64(len(r.Intervals))
}

/*
Total size of the regions
*/
func (r RegionCollection) Size() int64 {
	return int64(len(r.Positions))
}

// compute the size of some genomic intervals
func sizeOfGenomicIntervals(gis []GenomicInterval) int64 {
	var i int64
	for _, k := range gis {
		i += k.Length()
	}
	return i
}

type RegionCollection struct {
	Intervals []GenomicInterval
	Positions []int64 // translation from position array to index in the GenomicInteval (TODO make unit test testing each position if is in the interval)
}

func newRegionCollection(gis []GenomicInterval) *RegionCollection {

	// make the translation array for genomic Interval positions; the translation array helps finding random position in an interval collection;
	// pos is an array having all the indici of the genomic intervals; hence we need to find a random number with pos and thand
	// just get the attached number to get the random number in the interval
	// eg if intervals are 5..10 and 15..20
	// the pos will be [5 6 7 8 9 10 15 16 17 18 19 20] pos has length 12
	// now given a random number between 0..11 eg 7, the random position in the interval is thus pos[7]=17
	size := sizeOfGenomicIntervals(gis)
	pos := make([]int64, size)
	ri := 0 // running index
	for _, gi := range gis {
		for i := gi.Start; i <= gi.End; i++ {
			pos[ri] = i
			ri++
		}
	}
	return &RegionCollection{
		Intervals: gis,
		Positions: pos}
}

/*
initialize the genomic regions, i.e. what is a piRNA cluster and what is not a piRNA cluster
*/
func newClusterNonClusters(cl []int64, genome *GenomicLandscape) (*RegionCollection, *RegionCollection) {
	if cl == nil {
		// if user did not provide clusters (nil)
		// return an empty slice for the clusters and the entire genome for the nonClusters
		return newRegionCollection(make([]GenomicInterval, 0)), newRegionCollection(genome.intervals)
	}
	if len(genome.chrmSizes) != len(cl) {
		panic("Invalid number of piRNA clusters; must match number of chromosomes")
	}
	gis := genome.intervals
	clusters := make([]GenomicInterval, 0)
	nonClusters := make([]GenomicInterval, 0)
	for i, gi := range gis {
		clusterLength := cl[i]
		if clusterLength > gi.Length() {
			panic(fmt.Sprintf("Invalid size of piRNA cluster (%d), must not be larger than chromosome (%d)", clusterLength, gi.Length()))
		}
		nonclusterLength := gi.Length() - clusterLength
		clusterStart := gi.Start
		clusterEnd := clusterStart + clusterLength - 1
		nonclusterStart := clusterEnd + 1
		nonclusterEnd := gi.End
		if clusterLength > 0 {
			clui := GenomicInterval{Start: clusterStart, End: clusterEnd}
			clusters = append(clusters, clui)
		}
		if nonclusterLength > 0 {
			nonclui := GenomicInterval{Start: nonclusterStart, End: nonclusterEnd}
			nonClusters = append(nonClusters, nonclui)
		}

	}
	util.InvadeLogger.Printf("Will use piRNA clusters %v", clusters)
	util.InvadeLogger.Printf("Will use non-piRNA clusters %v", nonClusters)
	return newRegionCollection(clusters), newRegionCollection(nonClusters)

}

/*
Is a certain position within a region collection
*/
func (r RegionCollection) IsInRegion(position int64) bool {
	// must be sorted by start position of genomic interval; panic if not
	var lastpos int64
	for _, gi := range r.Intervals {
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
is a given TE insertion (int64) a cluster insertion
*/
func IsClusterInsertion(position int64) bool {
	return env.clusters.IsInRegion(position)
}

/*
TE insertions separated into categories.
There may be overlap between Paramutated and Trigger
*/
type InsertionCollection struct {
	Cluster []int64
	NOE     []int64 // non of either
}

/*
Environmental function;
Separate TE insertions into distinct categories
*/
func SeparateInsertions(positions []int64) InsertionCollection {
	var cluster, noe []int64
	for _, p := range positions {
		if p >= env.genome.totalGenome {
			panic(fmt.Sprintf("position outside genome %d", p))
		}
		if IsClusterInsertion(p) {
			cluster = append(cluster, p)
		} else {
			noe = append(noe, p)
		}
	}

	return InsertionCollection{
		Cluster: cluster,
		NOE:     noe}
}

/*
for a haploid genome, count the number of the following insertions "cluster and non-cluster";
return in the given order

func CountHaploidInsertions(positions []int64) (int64, int64) {
	var cluster, noe int64
	// Two steps
	// Step 1: separate cluster, reference and no-cluster/no-reference insertions
	for _, p := range positions {
		if p >= env.genome.totalGenome {
			panic(fmt.Sprintf("position outside genome %d", p))
		}
		if IsClusterInsertion(p) {
			cluster++
		} else {

			noe++
		}
	}
	return cluster, noe
}
*/

/*
for a diploid genome, count the number of the following insertions "cluster, non-cluster ";
return in the given order
*/
func CountDiploidInsertions(hap1 map[int64]TEInsertion, hap2 map[int64]TEInsertion) (int64, int64, int64, map[TEInsertion]int64, map[TEInsertion]int64, map[TEInsertion]int64) {
	totc1, cc1, nocc1, totmap1, cmap1, nocmap1 := CountHaploidInsertions(hap1)
	totc2, cc2, nocc2, totmap2, cmap2, nocmap2 := CountHaploidInsertions(hap2)
	return totc1 + totc2, cc1 + cc2, nocc1 + nocc2,
		mergeTEInsertionCountmap(totmap1, totmap2),
		mergeTEInsertionCountmap(cmap1, cmap2),
		mergeTEInsertionCountmap(nocmap1, nocmap2)
}

func mergeTEInsertionCountmap(m1 map[TEInsertion]int64, m2 map[TEInsertion]int64) map[TEInsertion]int64 {
	toret := make(map[TEInsertion]int64)
	for k, v := range m1 {
		toret[k] += v
	}
	for k, v := range m2 {
		toret[k] += v
	}
	return toret

}

/*
totalcount, clustercount, no-cluster count, total map, cluster map, no-cluster map
*/
func CountHaploidInsertions(hap1 map[int64]TEInsertion) (int64, int64, int64, map[TEInsertion]int64, map[TEInsertion]int64, map[TEInsertion]int64) {
	totc, cc, nocc := int64(0), int64(0), int64(0)
	totmap, cmap, nocmap := make(map[TEInsertion]int64), make(map[TEInsertion]int64), make(map[TEInsertion]int64)
	for pos, te := range hap1 {
		if pos >= env.genome.totalGenome {
			panic(fmt.Sprintf("position outside genome %d", pos))
		}
		totc++
		totmap[te]++
		if IsClusterInsertion(pos) {
			cc++
			cmap[te]++
		} else {

			nocc++
			nocmap[te]++
		}

	}
	return totc, cc, nocc, totmap, cmap, nocmap
}

func JustCountHaploidInsertions(hap1 []int64) (int64, int64, int64) {
	totc, cc, nocc := int64(0), int64(0), int64(0)
	for _, pos := range hap1 {
		if pos >= env.genome.totalGenome {
			panic(fmt.Sprintf("position outside genome %d", pos))
		}
		totc++
		if IsClusterInsertion(pos) {
			cc++
		} else {
			nocc++
		}

	}
	return totc, cc, nocc
}

/*
Given a TE insertion (position), check if the insertion is in a cluster (clu), reference region (ref) or none of either (noe)
*/
func ScoreInsertion(position int64) string {
	// Two steps
	// Step 1: separate cluster, reference and no-cluster/no-reference insertions
	if position >= env.genome.totalGenome {
		panic(fmt.Sprintf("position outside genome %d", position))
	}
	if IsClusterInsertion(position) {
		return "clu"
	} else {
		return "noc"
	}

}
