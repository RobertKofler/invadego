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
func SimulateInvasions(basepop string, popsize int64, replicates int64, generation int64, clonal bool) {
	for k := int64(0); k < replicates; k++ {
		pop := cmdparser.ParseBasePop(basepop, popsize) // reload base population; can not just reuse old one, needs to be novel random insertion sites!
		status := pop.GetStatus()
		outman.RecordPopulation(pop, k, 0, status)
		if status != fly.OK {
			continue // skip simulation for invalid base populations
		}

		for i := int64(1); i <= generation; i++ { // needs to start 1; 0 is the base population
			if clonal {
				pop = pop.GetNextGenerationClonal() // clonal propagation
			} else {
				pop = pop.GetNextGeneration() // sexual mating propagation
			}
			status := pop.GetStatus()
			outman.RecordPopulation(pop, k, i, status)

			// if the status is not ok abort!
			if status != fly.OK {
				break
			}
		}
	}
}
