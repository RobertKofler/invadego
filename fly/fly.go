package fly

import (
	"invade/env"
	"invade/util"
	"math/rand"
)

type Sex int64

// count the total number of flies
// ToDo problematic; needs to be reset for each replicate;
// maybe make somehow replicate specific
var FLYCOUNTER int64 = 1

// var flylock sync.Mutex

// size 2^64 = 1.844674e+19
// hence even if we do simulations for populations of size 100k for 100kgenerations we could run 1 844 674 407 replicates; should be sufficient ;)

const (
	FEMALE Sex = 0
	MALE   Sex = 1
	HERMA  Sex = 2
)

type Position struct {
	X int64
	Y int64
}
type FlySummary struct {
	CountTE  int64
	Silenced bool
	Ycoord   int64
	Xcoord   int64
}

type Fly struct {
	FlyNumber int64 // each fly has a number; starting at 1
	Hap1      []int64
	Hap2      []int64
	Silenced  bool
	Fitness   float64
	FlyStat   *FlyStatistic
}
type FlyStatistic struct {
	CountTotal int64
}

func (f *Fly) CountTotalInsertions() int64 {
	return int64(len(f.Hap1) + len(f.Hap2))
}

/*
Return the number of homozygous, heterozygous insertions yes
*/
func (f *Fly) CountHomozygousHeterozygous() (int64, int64) {
	// novel
	var insertionsites = make(map[int64]int64)

	homo, hetero := int64(0), int64(0)
	for _, is := range f.Hap1 {
		insertionsites[is]++
	}
	for _, is := range f.Hap2 {
		insertionsites[is]++
	}

	for _, val := range insertionsites {
		if val == 1 {
			hetero++
		} else if val == 2 {
			homo++
		} else {
			panic("invalid count of homozygote or heterozygote")
		}

	}
	return homo, hetero
}

/*
Get a gamete from the Fly;
First recombination among the two haplotypes will take place;
Second the number of new insertions will be computed based on
i) the TE insertions in the diploid parent
ii) epigenetic silencing
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
	// if epigenetically sileneced we have no novel insertions (zero is default)
	newsites := env.GetNewTranspositionSites(counttotal, f.Silenced)

	// merge old and new insertion sites, make them unique and sort
	return util.MergeUniqueSort(gamete, newsites)

}

/*
Compute basic statistics for a fly, ie number of cluster insertions, number of reference insertions, total number of insertions etc
*/
func getFlyStat(femgam []int64, malegam []int64) FlyStatistic {
	totcount := int64(len(femgam) + len(malegam))
	fs := FlyStatistic{
		CountTotal: totcount,
	}
	return fs
}

/*
Check if the TE is silenced
Silenced if a) count of TE is larger than threshold b) if the parent transmitted its epigenetic silencing status and the
individum has at least one insertion
*/
func getSilencingStatus(totalCount int64, silenced bool, fc int64) bool {

	// trigger 'de novo' silencing
	if !silenced {
		// check for new trigger events
		if env.IsTriggered(totalCount) {
			return true // silenced
		} else {
			return false
		}
	} else {
		// ok there is epigenetic silencing inherited wuhu
		// if there is a TE insertion it can be preserved,
		// otherwise the epigenetic silencing is lost
		if totalCount > 0 {
			return true // silenced (id of old fly that triggered it)
		} else {
			return false // no te insertion -> epigenetic silencing is lost
		}
	}
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
	Separate flies into males and females; return value provided in this order
func SeparateSexes(flies []Fly) ([]Fly, []Fly) {
	males := []Fly{}
	females := []Fly{}
	for _, fly := range flies {
		if fly.Sex == MALE {
			males = append(males, fly)
		} else if fly.Sex == FEMALE {
			females = append(females, fly)
		}
	}
	return males, females
}
*/

/*
Get random sex
*/
func GetRandomSex() Sex {
	s := Sex(int(2.0 * rand.Float64()))
	return s
}

/*
Setup a new Fly; given the gametes, the sex, and the epigenetic silencing, i.e. is the TE in the parents silenced
Will i) merge gametes ii) compute stats iii) determine silencing status (could be lost) iv) compute fitness v) increase FLYCOUNTER
*/
func NewFly(femgam []int64, malegam []int64, paternalsilenced bool) *Fly {
	// should give random numbers 0 or 1, ie male female
	fstat := getFlyStat(femgam, malegam)
	// multithreading lock and unlock
	//flylock.Lock()
	//
	currentCounter := FLYCOUNTER
	FLYCOUNTER++

	matpi := getSilencingStatus(fstat.CountTotal, paternalsilenced, currentCounter) // update the silencing status, eg if threshold is reached or if all TE insertions are lost
	newFly := Fly{Hap1: malegam, Hap2: femgam, FlyNumber: currentCounter, Silenced: matpi, FlyStat: &fstat}
	newFly.Fitness = GetFitness(&newFly)

	return &newFly
}
