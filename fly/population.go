package fly

import (
	"invade/env"
)

type Population struct {
	Flies  []Fly
	phase  Phase
	minFit float64
}

type Phase int64

const (
	RAPIDINVASION Phase = 0
	SHOTGUN       Phase = 1
	INACTIVE      Phase = 2
)

type PopStatus int64

const (
	BASEPOP PopStatus = 0
	OK      PopStatus = 1
	FAIL0   PopStatus = 2
	FAILW   PopStatus = 3
)

func (p *Population) Size() int64 {
	return int64(len(p.Flies))
}

func InitializePopulation(flies []Fly) *Population {
	p := Population{Flies: flies}
	p.minFit = p.GetAverageFitness()
	p.phase = updatePhase(&p, RAPIDINVASION)
	return &p
}

/*
Get the next generation i) get mate pairs according to fitness ii) get gametes with random recombination and transposition iii) get random sex iv) generate new flies
v) compute fitness and statistics
*/
func (p *Population) GetNextGeneration() *Population {
	matePairs := getMatePairs(p.Flies, int64(len(p.Flies)))
	nextGen := make([]Fly, len(matePairs))
	for i, mp := range matePairs {
		femgam := mp.mate1.GetGamete()
		malgam := mp.mate2.GetGamete()
		newFly := NewFly(femgam, malgam) // maternal piRNAs; only the female passes them
		nextGen[i] = *newFly
	}
	newPop := Population{Flies: nextGen}
	newPhase := updatePhase(&newPop, p.phase)
	newPop.phase = newPhase
	newMinFit := updateFitness(&newPop, p.minFit)
	newPop.minFit = newMinFit
	return &newPop
}

/*
Find the novel minimum Fitness;
Does the new population have a lower average fitness than the previous one?
*/
func updateFitness(newPop *Population, previousMinFit float64) float64 {
	curFit := newPop.GetAverageFitness()
	if curFit < previousMinFit {
		return curFit
	} else {
		return previousMinFit
	}
}

/*
update the phase of the invasion;
RAPIDINVASION->TRIGGERED -> SHOTGUN -> INACTIVE
*/
func updatePhase(newPop *Population, oldPhase Phase) Phase {
	if oldPhase == INACTIVE {
		return INACTIVE // no escape from a fixed cluster insertion
	}

	freq := newPop.GetWithClusterInsertionFrequency()
	if oldPhase == RAPIDINVASION {
		if freq > 0.99 { // condition for trigger -> at least one with piRNAs
			return SHOTGUN
		}
	} else if oldPhase == SHOTGUN {
		// NOTE that heterogenous sites with mixed insertion biases are considered as fixed!
		fixedIns := newPop.GetFixedInsertions()
		_, fclu, _ := env.JustCountHaploidInsertions(fixedIns)
		if fclu > 0 { // conditon for inactive -> at least one fixed cluster insertion
			return INACTIVE
		}
		// Check if the inactive phase was reached
	}
	return oldPhase //if nothing special happens -> oldphase
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
	for _, f := range p.Flies {
		fitcount += f.Fitness
		tecount += int(f.FlyStat.CountTotal)
	}
	avfit := fitcount / float64(p.Size())
	if tecount == 0 {
		return FAIL0
	} else if avfit < minimumFitness {
		return FAILW
	} else {
		return OK
	}
}

func (p *Population) GetHaplotypes() []map[int64]env.TEInsertion {
	toret := make([]map[int64]env.TEInsertion, 0, p.Size()*2)
	for _, f := range p.Flies {
		toret = append(toret, f.Hap1)
		toret = append(toret, f.Hap2)
	}
	return toret
}
