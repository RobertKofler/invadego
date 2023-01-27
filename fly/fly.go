package fly

import (
	"invade/env"
)

// count the total number of flies
// ToDo problematic; needs to be reset for each replicate;
// maybe make somehow replicate specific
var FLYCOUNTER int64 = 1

type Fly struct {
	FlyNumber int64 // each fly has a number; starting at 1
	Hap1      map[int64]env.TEInsertion
	Hap2      map[int64]env.TEInsertion
	Fitness   float64
	FlyStat   *FlyStatistic
}
type FlyStatistic struct {
	CountTotal   int64
	CountCluster int64
	CountNoc     int64
	TotMap       map[env.TEInsertion]int64
	ClusterMap   map[env.TEInsertion]int64
	NocMap       map[env.TEInsertion]int64
}

func (f *Fly) CountTotalInsertions() int64 {
	return int64(len(f.Hap1) + len(f.Hap2))
}

func (f *Fly) GetClone() *Fly {
	malegam := make(map[int64]env.TEInsertion)
	femgam := make(map[int64]env.TEInsertion)
	for k, v := range f.Hap1 {
		malegam[k] = v
	}
	for k, v := range f.Hap2 {
		femgam[k] = v
	}

	malegam = env.InsertNewTranspositionSites(malegam, f.FlyStat.TotMap, f.FlyStat.ClusterMap, f.FlyStat.CountCluster)
	femgam = env.InsertNewTranspositionSites(femgam, f.FlyStat.TotMap, f.FlyStat.ClusterMap, f.FlyStat.CountCluster)

	malegam = env.IntroduceMutations(malegam)
	femgam = env.IntroduceMutations(femgam)
	return NewFly(malegam, femgam)
}

/*
Get a gamete from the Fly;
First recombination among the two haplotypes will take place;
Second the number of new insertions will be computed based on
i) the TE insertions in the diploid parent
ii) piRNA cluster insertions
iii) the transposition rate.
The number and position of new insertions will be random.
Multiple insertions at the same site will be ignored.
*/
func (f *Fly) GetGamete() map[int64]env.TEInsertion {
	if f.FlyStat == nil {
		panic("Fly statistics not initialized")
	}
	// First get recombined game
	gamete := f.getRecombinedGamete()

	// insert novel transposition sites into the gamete
	gamete = env.InsertNewTranspositionSites(gamete, f.FlyStat.TotMap, f.FlyStat.ClusterMap, f.FlyStat.CountCluster)
	// finally introduce mutations, eg increase or decrease the insertion bias as requested
	gamete = env.IntroduceMutations(gamete)

	return gamete

}

/*
Compute basic statistics for a fly, ie number of cluster insertions, number of reference insertions, total number of insertions etc
*/
func getFlyStat(femgam map[int64]env.TEInsertion, malegam map[int64]env.TEInsertion) FlyStatistic {

	totc, cc, nocc, totmap, cmap, nocmap := env.CountDiploidInsertions(femgam, malegam)
	fs := FlyStatistic{
		CountTotal:   totc,
		CountCluster: cc,
		CountNoc:     nocc,
		ClusterMap:   cmap,
		NocMap:       nocmap,
		TotMap:       totmap,
	}
	return fs
}

/*
Get a unique set of all insertion sites in a diploid fly.
Each heterozygous insertion in hap1 and hap2 is present.
Homozygous insertions are only present once.
Output is sorted.
*/
func (f *Fly) GetInsertionSites() []int64 {
	return env.MergeUniqueSort(f.Hap1, f.Hap2)
}

/*
Recombine the two haplotypes of a fly (f.hap1 and f.hap2) given a sorted list of recombination events;
Implements the RECOMBINATION FIRST principle. Eg if a TE and a rec.event are at site 20, than recombination is done first, and the TE is used second.
(Rec first is important to enable random assortment of the first chromosome!)
*/
func (f *Fly) recombine(recombinationEvents []int64) map[int64]env.TEInsertion {
	hap1 := env.GetSortedKeys(f.Hap1)
	hap2 := env.GetSortedKeys(f.Hap2)
	ihap1 := 0
	ihap2 := 0
	ishap1 := true
	newhap := make(map[int64]env.TEInsertion)

	for _, r := range recombinationEvents {
		for ihap1 < len(hap1) && hap1[ihap1] < r {
			if ishap1 {
				pos := hap1[ihap1]
				te := f.Hap1[pos]
				newhap[pos] = te
			}
			ihap1++
		}
		for ihap2 < len(hap2) && hap2[ihap2] < r {
			if !ishap1 {
				pos := hap2[ihap2]
				te := f.Hap2[pos]
				newhap[pos] = te
			}
			ihap2++
		}
		ishap1 = !ishap1
	}
	// Deal with the last elements
	for _ = ihap1; ihap1 < len(hap1); ihap1++ {
		if ishap1 {
			pos := hap1[ihap1]
			te := f.Hap1[pos]
			newhap[pos] = te
		}
	}
	for _ = ihap2; ihap2 < len(hap2); ihap2++ {
		if !ishap1 {
			pos := hap2[ihap2]
			te := f.Hap2[pos]
			newhap[pos] = te
		}
	}
	return newhap
}

/*
Get a recombined gamete for the two haplotypes of a fly.
Recombination events are random, according to environment settings (i.e. chromosomes, rec.rate)
*/
func (f *Fly) getRecombinedGamete() map[int64]env.TEInsertion {

	recsites := env.GetRecombinationEvents()
	rec := f.recombine(recsites)
	return rec
}

/*
Setup a new Fly; given the gametes, the sex, and the maternal piRNAs;
Will i) merge gametes ii) compute stats iii) determine piRNA status iv) compute fitness v) increase FLYCOUNTER
*/
func NewFly(femgam map[int64]env.TEInsertion, malegam map[int64]env.TEInsertion) *Fly {
	// should give random numbers 0 or 1, ie male female
	fstat := getFlyStat(femgam, malegam)
	// multithreading lock and unlock
	//flylock.Lock()
	currentCounter := FLYCOUNTER
	FLYCOUNTER++
	//flylock.Unlock()

	newFly := Fly{Hap1: malegam, Hap2: femgam, FlyNumber: currentCounter, FlyStat: &fstat}
	newFly.Fitness = GetFitness(&newFly)

	return &newFly
}
