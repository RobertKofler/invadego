package fly

import (
	"invade/env"
	"invade/util"
)

type Sex int64

// count the total number of flies
// ToDo problematic; needs to be reset for each replicate;
// maybe make somehow replicate specific
var FLYCOUNTER int64 = 1

type Fly struct {
	FlyNumber int64 // each fly has a number; starting at 1
	Hap1      []int64
	Hap2      []int64
	Fitness   float64
	FlyStat   *FlyStatistic
}
type FlyStatistic struct {
	CountTotal     int64
	CountCluster   int64
	CountReference int64
	CountNOE       int64
}

func (f *Fly) CountTotalInsertions() int64 {
	return int64(len(f.Hap1) + len(f.Hap2))
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
func (f *Fly) GetGamete() []int64 {
	if f.FlyStat == nil {
		panic("Fly statistics not initialized")
	}
	// First get recombined game
	gamete := f.getRecombinedGamete()

	// Second introduce novel transposition events
	counttotal := int64(len(f.Hap1) + len(f.Hap2))

	// the function generates novel transposition events for a HAPLOID genome, i.e. a gamete
	// if the number of cluster insertions > 0 than we have piRNAs and thus no novel insertions (zero is default)
	newsites := env.GetNewTranspositionSites(counttotal, f.FlyStat.CountCluster)

	// merge old and new insertion sites, make them unique and sort
	return util.MergeUniqueSort(gamete, newsites)

}

/*
Compute basic statistics for a fly, ie number of cluster insertions, number of reference insertions, total number of insertions etc
*/
func getFlyStat(femgam []int64, malegam []int64) FlyStatistic {
	totcount := int64(len(femgam) + len(malegam))
	cluster, reference, noe := env.CountDiploidInsertions(femgam, malegam)
	fs := FlyStatistic{
		CountTotal:     totcount,
		CountCluster:   cluster,
		CountReference: reference,
		CountNOE:       noe,
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
	return util.MergeUniqueSort(f.Hap1, f.Hap2)
}

/*
Recombine the two haplotypes of a fly (f.hap1 and f.hap2) given a sorted list of recombination events;
Implements the RECOMBINATION FIRST principle. Eg if a TE and a rec.event are at site 20, than recombination is done first, and the TE is used second.
(Rec first is important to enable random assortment of the first chromosome!)
*/
func (f *Fly) recombine(recombinationEvents []int64) []int64 {
	hap1 := f.Hap1
	hap2 := f.Hap2
	ihap1 := 0
	ihap2 := 0
	ishap1 := true
	newhap := make([]int64, 0, len(hap1))
	for _, r := range recombinationEvents {
		for ihap1 < len(hap1) && hap1[ihap1] < r {
			if ishap1 {
				newhap = append(newhap, hap1[ihap1])
			}
			ihap1++
		}
		for ihap2 < len(hap2) && hap2[ihap2] < r {
			if !ishap1 {
				newhap = append(newhap, hap2[ihap2])
			}
			ihap2++
		}
		ishap1 = !ishap1
	}
	// Deal with the last elements
	for _ = ihap1; ihap1 < len(hap1); ihap1++ {
		if ishap1 {
			newhap = append(newhap, hap1[ihap1])
		}
	}
	for _ = ihap2; ihap2 < len(hap2); ihap2++ {
		if !ishap1 {
			newhap = append(newhap, hap2[ihap2])
		}
	}
	return newhap
}

/*
Get a recombined gamete for the two haplotypes of a fly.
Recombination events are random, according to environment settings (i.e. chromosomes, rec.rate)
*/
func (f *Fly) getRecombinedGamete() []int64 {

	recsites := env.GetRecombinationEvents()
	rec := f.recombine(recsites)
	return rec
}

/*
Setup a new Fly; given the gametes, the sex, and the maternal piRNAs;
Will i) merge gametes ii) compute stats iii) determine piRNA status iv) compute fitness v) increase FLYCOUNTER
*/
func NewFly(femgam []int64, malegam []int64) *Fly {
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
