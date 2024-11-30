package sim

import (
	"invade/fly"
	"invade/io/cmdparser"
	"invade/outman"
)

/*
 perform the simulations;
 multiple replicates and generations
*/
func SimulateInvasions(basepop string, generation int64, om outman.IOutputManager) {
	for om.NeedMoreReplicates() {
		pop := cmdparser.ParseBasePop(basepop)
		status := pop.GetStatus()
		om.RecordPopulation(pop, 0, status)
		if status != fly.OK {
			om.FinaliseReplicate(status)
			continue // skip simulation for invalid base populations
		}

		for i := int64(1); i <= generation; i++ { // needs to start 1; 0 is the base population
			pop = pop.GetNextGeneration()
			status := pop.GetStatus()

			om.RecordPopulation(pop, i, status)

			// if the status is not ok abort!
			if status != fly.OK {
				om.FinaliseReplicate(status)
				break
			} else if i == generation {
				om.FinaliseReplicate(status)
			}
		}
	}
}
