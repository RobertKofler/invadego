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
func SimulateInvasions(basepop string, popsize int64, replicates int64, generation int64) {
	for k := int64(0); k < replicates; k++ {
		pop := cmdparser.ParseBasePop(basepop, popsize)
		status := pop.GetStatus()
		outman.RecordPopulation(pop, k, 0, status)
		if status != fly.OK {
			continue // skip simulation for invalid base populations
		}

		for i := int64(1); i <= generation; i++ { // needs to start 1; 0 is the base population
			pop = pop.GetNextGeneration() // Fuck multithreading! we loose reproducibitilty with the seeds!
			status := pop.GetStatus()
			outman.RecordPopulation(pop, k, i, status)

			// if the status is not ok abort!
			if status != fly.OK {
				break
			}
		}
	}
}
