package env

/*
Initialize the entire environment for the simulations, i.e. the chromosomes, the piRNA clusters, the recombination rate
(fitness? mating?)
*/
func SetupEnvironment(gridx int64, gridy int64, materadius int64, chrSizes []int64, recRate []float64, hostsilencetrigger int64, episilence float64, selfrate float64, minFitness float64, maxInsertions float64) {
	//	env.SetupEnvironment(genome, recrate, clp.TriggerSilence, clp.EpiSilencing,clp.SelfingRate, clp.MinFitness, float64(clp.MaxInsertions))
	genome := newGenomicLandscape(chrSizes) // setup genome

	// compute the recombination windows
	recwins := getRecombinationWindows(genome.intervals, recRate)

	env = Environment{
		genome:               genome,
		triggerThreshold:     hostsilencetrigger,
		epiRate:              episilence,
		selfingRate:          selfrate,
		minimumFitness:       minFitness,
		maximumInsertions:    maxInsertions,
		recombinationWindows: recwins,
		gridX:                gridx,
		gridY:                gridy,
		mateRadius:           materadius,
	}
}
