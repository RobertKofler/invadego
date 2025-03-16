package env

/*
Initialize the entire environment for the simulations, i.e. the chromosomes, the piRNA clusters, the recombination rate
(fitness? mating?)
*/
func SetupEnvironment(gridx int64, gridy int64, chrSizes []int64, recRate []float64,
	popstruct []int64, barrierStrength float64, hostsilencetrigger int64, episilence string,
	epsiloss float64, minFitness float64, maxInsertions float64) {
	//	env.SetupEnvironment(clp.GridX, clp.GridY, genome, recrate,
	// popstruct, clp.BarrierEffect, clp.TriggerSilence, clp.EpiSilencing,
	// clp.PSiLoss, clp.MinFitness, float64(clp.MaxInsertions))
	genome := newGenomicLandscape(chrSizes) // setup genome

	// compute the recombination windows
	recwins := getRecombinationWindows(genome.intervals, recRate)

	ps := popStruct{
		xBarriers:   popstruct,
		barStrength: barrierStrength,
	}

	env = Environment{
		genome:               genome,
		triggerThreshold:     hostsilencetrigger,
		epiMode:              episilence,
		minimumFitness:       minFitness,
		popStruct:            &ps,
		psilencingLoss:       epsiloss,
		maximumInsertions:    maxInsertions,
		recombinationWindows: recwins,
		gridX:                gridx,
		gridY:                gridy,
	}
}
