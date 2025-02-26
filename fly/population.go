package fly

import (
	"fmt"
	"invade/env"
)

// important testing https://go.dev/tour/flowcontrol/1

type Population struct {
	Flies       [][]*Fly
	linearFlies []*Fly
	size        int64
	activePhase Phase
	triggered   int64
}

type Phase int64

// phases
const (
	INVASION Phase = 0
	DECLINE  Phase = 1
	INACTIVE Phase = 2
)

type PopStatus int64

const (
	BASEPOP PopStatus = 0
	OK      PopStatus = 1
	FAIL0   PopStatus = 2
	FAILW   PopStatus = 3
	FAILMAX PopStatus = 4
)

/*
	Population size;
	total number of individuals in 2D grid
*/
func (p *Population) Size() int64 {
	return p.size
}
func (p *Population) Phase() Phase {
	return p.activePhase
}

func NewPopulation(flies [][]*Fly, prevPhase Phase) *Population {
	p := Population{Flies: flies}
	linear := make([]*Fly, 0, p.Size())
	s := int64(0)

	for _, y := range flies {
		linear = append(linear, y...) // cool append entire vector
		s += int64(len(y))

	}
	p.linearFlies = linear
	p.size = s

	p.activePhase = p.updatePhase(prevPhase)
	return &p
}

func (p *Population) updatePhase(prevPhase Phase) Phase {
	cns := p.GetWithTEAndNotSilencedCount()
	cs := p.GetWithTEAndSilencedCount()

	if prevPhase == INACTIVE {
		return INACTIVE
	} else if prevPhase == INVASION {
		if cs > 0 && cns == 0 {
			return INACTIVE // everyone silenced?
		} else if cs > 0 {
			return DECLINE
		} else {
			return INVASION
		}
	} else if prevPhase == DECLINE {
		if cs > 0 && cns == 0 {
			return INACTIVE
		} else {
			return DECLINE
		}
	} else {
		panic(fmt.Sprintf("invalid phase %d", prevPhase))
	}
}

/*
Get the next generation i) get mate pairs according to fitness ii) get gametes with random recombination and transposition iii) get random sex iv) generate new flies
v) compute fitness and statistics
*/
func (p *Population) GetNextGeneration() *Population {

	// get the merry couples; selfing is considered, in which case male and female are identical
	matePairs := mater.GetMatePairs(p.Flies)

	// initialize the grid for the next generation
	ysize, xsize := env.GetYSize(), env.GetXSize()
	nextGen := make([][]*Fly, ysize)
	for i, _ := range nextGen {
		//make x
		nextGen[i] = make([]*Fly, xsize)
	}
	for y := 0; y < int(ysize); y++ {
		for x := 0; x < int(xsize); x++ {
			mp := matePairs[y][x]
			// selfing considered in choice of mate pair! do not address here
			femgam := mp.female.GetGamete()
			malegam := mp.male.GetGamete()
			issilenced := env.OffspringIsSilenced(mp.female.Silenced, mp.male.Silenced)

			newFly := NewFly(femgam, malegam, issilenced)
			nextGen[y][x] = newFly
		}
	}

	newPop := NewPopulation(nextGen, p.activePhase)
	// update phase would go here

	return newPop
}

/*
Get status, either
OK    		TEs and fitness
fail-0		no TEs left
fail-w		fitness to low
base 		base population
fail-sex 	only males or only females
*/
func (p *Population) GetStatus() PopStatus {
	fitcount := 0.0
	tecount := 0
	for _, f := range p.linearFlies {
		fitcount += f.Fitness
		tecount += int(f.FlyStat.CountTotal)
	}
	avfit := fitcount / float64(p.Size())
	avins := p.GetAverageInsertions()
	if tecount == 0 {
		return FAIL0
	} else if avfit < env.GetMinimumFitness() {
		return FAILW
	} else if avins > env.GetMaximumInsertions() {
		return FAILMAX
	} else {
		return OK
	}
}

func (p *Population) GetHaplotypes() [][]int64 {
	toret := make([][]int64, 0, p.Size()*2)
	for _, f := range p.linearFlies {
		toret = append(toret, f.Hap1)
		toret = append(toret, f.Hap2)
	}
	return toret
}
