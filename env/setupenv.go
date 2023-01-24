package env

/*
Initialize the entire environment for the simulations, i.e. the chromosomes, the piRNA clusters, the recombination rate
*/
func SetupEnvironment(chrSizes []int64, cluSizes []int64, recRate []float64, minFitness float64) {
	genome := newGenomicLandscape(chrSizes)                          // setup genome
	clusters, nonClusters := newClusterNonClusters(cluSizes, genome) // setup cluster, they depend on the genome

	// compute the recombination windows
	recwins := getRecombinationWindows(genome.intervals, recRate)

	env = Environment{
		genome:               genome,
		clusters:             *clusters,
		nonClusters:          *nonClusters,
		minimumFitness:       minFitness,
		recombinationWindows: recwins,
	}
}
