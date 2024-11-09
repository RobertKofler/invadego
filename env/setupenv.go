package env

import "fmt"

/*
Initialize the entire environment for the simulations, i.e. the chromosomes, the piRNA clusters, the recombination rate
(fitness? mating?)
*/
func SetupEnvironment(gridx int64, gridy int64, chrSizes []int64, recRate []float64, hostsilencetrigger int64, episilence string, minFitness float64, maxInsertions float64) {
	//	env.SetupEnvironment(genome, recrate, clp.TriggerSilence, clp.EpiSilencing,clp.SelfingRate, clp.MinFitness, float64(clp.MaxInsertions))
	genome := newGenomicLandscape(chrSizes) // setup genome

	// compute the recombination windows
	recwins := getRecombinationWindows(genome.intervals, recRate)
	epimode := NONE
	if episilence == "none" {
		epimode = NONE
	} else if episilence == "ara" {
		epimode = ARABIDOPSIS
	} else if episilence == "dros" {
		epimode = DROSOPHILA
	} else {
		panic(fmt.Sprintf("nknown epigenetic mode %s", episilence))
	}

	env = Environment{
		genome:               genome,
		triggerThreshold:     hostsilencetrigger,
		epiMode:              epimode,
		minimumFitness:       minFitness,
		maximumInsertions:    maxInsertions,
		recombinationWindows: recwins,
		gridX:                gridx,
		gridY:                gridy,
	}
}

func SetupEpigeneticSilencing(epis string) {

}
