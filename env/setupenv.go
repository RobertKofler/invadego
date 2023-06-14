package env

/*
Initialize the entire environment for the simulations, i.e. the chromosomes, the piRNA clusters, the recombination rate
(fitness? mating?)
*/
func SetupEnvironment(chrSizes []int64, cluSizes []int64, refSizes []int64,
	trigger []bool, para []bool, recRate []float64, minFitness float64, maxInsertions float64) {
	genome := newGenomicLandscape(chrSizes)             // setup genome
	clusters := newCluster(cluSizes, genome)            // setup cluster, they depend on the genome
	refRegions := newReferenceRegions(refSizes, genome) // setup reference regions
	if isClusterOverlappingReferences(clusters, refRegions) {
		// clusters must not overlap with reference regions
		panic("Invalid definition of clusters and reference regions; must not overlap")
	}
	triggers := newRecurrentSite(trigger)
	paramutables := newRecurrentSite(para)

	// compute the recombination windows
	recwins := getRecombinationWindows(genome.intervals, recRate)

	env = Environment{
		genome:               genome,
		clusters:             clusters,
		refRegions:           refRegions,
		triggers:             triggers,
		paramutables:         paramutables,
		minimumFitness:       minFitness,
		maximumInsertions:    maxInsertions,
		recombinationWindows: recwins,
	}
}
