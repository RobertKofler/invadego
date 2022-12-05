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
)

type Fly struct {
	FlyNumber int64 // each fly has a number; starting at 1
	Hap1      []int64
	Hap2      []int64
	Matpirna  int64 // number of the fly that triggered the maternal piRNAs; allows to identify soft sweeps from recurrent mutations!
	Sex       Sex
	Fitness   float64
	FlyStat   *FlyStatistic
}
type FlyStatistic struct {
	CountTotal     int64
	CountCluster   int64
	CountReference int64
	CountPara      int64
	CountTrigger   int64
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
iii) maternal piRNAs and paramutable sites
iv) the transposition rate.
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
	// if f.matpirna > 0 than we have piRNAs and thus no novel insertions (zero is default)
	newsites := env.GetNewTranspositionSites(counttotal, f.Matpirna > 0)

	// merge old and new insertion sites, make them unique and sort
	return util.MergeUniqueSort(gamete, newsites)

}

/*
Compute basic statistics for a fly, ie number of cluster insertions, number of reference insertions, total number of insertions etc
*/
func getFlyStat(femgam []int64, malegam []int64) FlyStatistic {
	totcount := int64(len(femgam) + len(malegam))
	cluster, reference, para, trigger, noe := env.CountDiploidInsertions(femgam, malegam)
	fs := FlyStatistic{
		CountTotal:     totcount,
		CountCluster:   cluster,
		CountReference: reference,
		CountPara:      para,
		CountTrigger:   trigger,
		CountNOE:       noe,
	}
	return fs
}

func getMaternalPirnaStatus(fstat FlyStatistic, matpirna int64, fc int64) int64 {
	// Gain maternal piRNAs
	if matpirna == 0 {
		// check for new trigger events

		if fstat.CountCluster > 0 {
			return fc // if there are cluster insertions -> we gained maternal piRNAs
		} else if fstat.CountTrigger > 0 && fstat.CountPara > 0 {
			return fc // if there are paramutable sites and trigger sites -> we gained maternal piRNAs
		} else {
			return 0 // ok still no maternal piRNAs
		}
	} else {
		// ok there were maternal piRNAs: wuhu
		// if there is a cluster insertion or a paramutable site -> maternal piRNAs are preserved
		// otherwise maternal piRNAs are LOST!

		if fstat.CountCluster > 0 {
			return matpirna // cluster insertion -> preserve maternal piRNAs
		} else if fstat.CountPara > 0 {
			return matpirna // paramutable site -> preserve maternal piRNAs
		} else {
			return 0 // LOSS of maternal PIRNAS
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
*/
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

/*
Get random sex
*/
func GetRandomSex() Sex {
	s := Sex(int(2.0 * rand.Float64()))
	return s
}

/*
Setup a new Fly; given the gametes, the sex, and the maternal piRNAs;
Will i) merge gametes ii) compute stats iii) determine piRNA status iv) compute fitness v) increase FLYCOUNTER
*/
func NewFly(femgam []int64, malegam []int64, sex Sex, matpirna int64) *Fly {
	// should give random numbers 0 or 1, ie male female
	fstat := getFlyStat(femgam, malegam)
	// multithreading lock and unlock
	//flylock.Lock()
	currentCounter := FLYCOUNTER
	FLYCOUNTER++
	//flylock.Unlock()

	matpi := getMaternalPirnaStatus(fstat, matpirna, currentCounter)
	newFly := Fly{Hap1: malegam, Hap2: femgam, FlyNumber: currentCounter, Matpirna: matpi, Sex: Sex(sex), FlyStat: &fstat}
	newFly.Fitness = GetFitness(&newFly)

	return &newFly
}
