package outman

import "invade/fly"

type IOutputManager interface {
	WriteInfo(userargs string, seed int64, version string)
	FinaliseReplicate(popstat fly.PopStatus)
	NeedMoreReplicates() bool
	End()
	RecordPopulation(p *fly.Population, generation int64, popstat fly.PopStatus) // the bool is success if the invasion was recorded
}

type popentry struct {
	population *fly.Population
	generation int64
	popstat    fly.PopStatus
}

func SetupOutputConditionalDecorator(omclassical *OutputManager) *OutputManagerCondInvasionDecorator {
	ocd := OutputManagerCondInvasionDecorator{
		omi:                   omclassical,
		successfullReplicates: 0,
		candidateRecord:       make([]popentry, 0),
	}
	return &ocd
}

type OutputManagerCondInvasionDecorator struct {
	omi                   *OutputManager
	candidateRecord       []popentry
	successfullReplicates int64
}

func (om *OutputManagerCondInvasionDecorator) WriteInfo(userargs string, seed int64, version string) {
	om.omi.WriteInfo(userargs, seed, version)
}
func (om *OutputManagerCondInvasionDecorator) End() {
	om.omi.End()
}
func (om *OutputManagerCondInvasionDecorator) NeedMoreReplicates() bool {
	if om.successfullReplicates < om.omi.replicatesNeeded {
		return true
	} else {
		return false
	}
}

func (om *OutputManagerCondInvasionDecorator) FinaliseReplicate(popstat fly.PopStatus) {

	if popstat == fly.OK {
		// if we have a successful invasion then a) print all stored entries, b) reset the buffer c) increase the count of successful invasions
		for _, k := range om.candidateRecord {
			om.omi.RecordPopulation(k.population, k.generation, k.popstat)
		}
		om.candidateRecord = make([]popentry, 0)
		om.successfullReplicates++
	}
	// in any case increase the count of invasions
	om.omi.currentReplicate = om.omi.currentReplicate + 1
}

func (om *OutputManagerCondInvasionDecorator) RecordPopulation(p *fly.Population, generation int64, popstat fly.PopStatus) {
	// we are only buffering  successful invasions with the requested step size
	if popstat == fly.OK && generation%om.omi.steps == 0 {
		pe := popentry{
			population: p,
			generation: generation,
			popstat:    popstat,
		}
		om.candidateRecord = append(om.candidateRecord, pe)

	}
}
