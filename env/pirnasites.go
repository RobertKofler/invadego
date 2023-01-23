package env

import (
	"fmt"
	"invade/util"
)

type RegionCollection []GenomicInterval

func newCluster(cl []int64, genome *GenomicLandscape) RegionCollection {
	if cl == nil {
		// if user did not provide clusters (nil)
		// return an empty slice
		return RegionCollection(make([]GenomicInterval, 0))
	}
	if len(genome.chrmSizes) != len(cl) {
		panic("Invalid number of piRNA clusters; must match number of chromosomes")
	}
	gis := genome.intervals
	clusters := make([]GenomicInterval, len(cl)) // initialize slice in correct size
	for i, gi := range gis {
		clusterLength := cl[i]
		if clusterLength > gi.Length() {
			panic(fmt.Sprintf("Invalid size of piRNA cluster (%d), must not be larger than chromosome (%d)", clusterLength, gi.Length()))
		}
		clusterStart := gi.Start
		clusterEnd := clusterStart + clusterLength - 1
		ngi := GenomicInterval{clusterStart, clusterEnd}
		clusters[i] = ngi
	}
	util.InvadeLogger.Printf("Will use piRNA clusters %v", clusters)
	return RegionCollection(clusters) // can not return a pointer to Cluster, likely because cluster is a slice which is already a reference type

}

func newReferenceRegions(rl []int64, genome *GenomicLandscape) RegionCollection {
	if rl == nil {
		// if user did not provide a reference region (nil)
		// return an empty slice
		return RegionCollection(make([]GenomicInterval, 0))
	}
	if len(genome.chrmSizes) != len(rl) {
		panic("Invalid number of reference regions; must match number of chromosomes")
	}
	gis := genome.intervals
	refs := make([]GenomicInterval, len(rl)) // initialize slice in correct size
	for i, gi := range gis {
		refLength := rl[i]
		if refLength > gi.Length() {
			panic(fmt.Sprintf("Invalid size of reference region (%d), must not be larger than chromosome (%d)", refLength, gi.Length()))
		}
		refEnd := gi.End
		refStart := gi.End - refLength + 1

		ngi := GenomicInterval{refStart, refEnd}
		refs[i] = ngi
	}
	util.InvadeLogger.Printf("Will use reference regions %v", refs)
	return RegionCollection(refs) // can not return a pointer to Cluster, likely because cluster is a slice which is already a reference type
}

/*
check if piRNA clusters are overlapping with reference regions
*/
func isClusterOverlappingReferences(clus RegionCollection, refs RegionCollection) bool {
	// deal with not provided clusters or references
	// command line nil turned into -> slice of size zero
	if len(clus) == 0 || len(refs) == 0 {
		return false
	}
	for i, cl := range clus {
		re := refs[i]
		if cl.End >= re.Start {
			return true
		}
	}
	return false
}

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
is a given TE insertion (int64) a cluster insertion
*/
func IsClusterInsertion(position int64) bool {
	return env.clusters.IsInRegion(position)
}

func IsReferenceInsertion(position int64) bool {
	return env.refRegions.IsInRegion(position)
}

/*
TE insertions separated into categories.
There may be overlap between Paramutated and Trigger
*/
type InsertionCollection struct {
	Cluster   []int64
	RefRegion []int64
	NOE       []int64 // non of either
}

/*
Environmental function;
Separate TE insertions into distinct categories
*/
func SeparateInsertions(positions []int64) InsertionCollection {
	var cluster, noe, refregion []int64
	for _, p := range positions {
		if p >= env.genome.totalGenome {
			panic(fmt.Sprintf("position outside genome %d", p))
		}
		if IsClusterInsertion(p) {
			cluster = append(cluster, p)
		} else if IsReferenceInsertion(p) {
			refregion = append(refregion, p)
		} else {
			noe = append(noe, p)
		}
	}

	return InsertionCollection{
		Cluster:   cluster,
		RefRegion: refregion,
		NOE:       noe}
}

/*
for a haploid genome, count the number of the following insertions "cluster, reference and noe (non of either)";
return in the given order
*/
func CountHaploidInsertions(positions []int64) (int64, int64, int64) {
	var cluster, noe, refregion int64
	// Two steps
	// Step 1: separate cluster, reference and no-cluster/no-reference insertions
	for _, p := range positions {
		if p >= env.genome.totalGenome {
			panic(fmt.Sprintf("position outside genome %d", p))
		}
		if IsClusterInsertion(p) {
			cluster++
		} else if IsReferenceInsertion(p) {
			refregion++
		} else {

			noe++
		}
	}
	return cluster, refregion, noe
}

/*
for a diploid genome, count the number of the following insertions "cluster, reference and noe (non of either)";
return in the given order
*/
func CountDiploidInsertions(hap1 []int64, hap2 []int64) (int64, int64, int64) {
	cluster1, refregion1, noe1 := CountHaploidInsertions(hap1)
	cluster2, refregion2, noe2 := CountHaploidInsertions(hap2)
	return cluster1 + cluster2,
		refregion1 + refregion2,
		noe1 + noe2
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
	} else if IsReferenceInsertion(position) {
		return "ref"
	} else {
		return "noe"
	}

}
