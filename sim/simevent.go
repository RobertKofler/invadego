package sim

import (
	"fmt"
	"invade/env"
	"invade/fly"
)

/*
Reset all events at the beginning of a new replicate
*/
func ResetEvents() {
	if genClusterSwitch > 0 {
		env.SetCluster(defaultCluster)
	}

}

func SetEvents(pop *fly.Population, generation int64) {
	if genClusterSwitch > 0 && generation == genClusterSwitch {
		env.SetCluster(eventCluster)
		pop.ResetStatus()
	}
}

func SetupClusterRemoval(when int64, howmany int64) {
	genClusterSwitch = when
	defaultCluster = env.GetCluster()
	newclustercount := len(defaultCluster) - int(howmany)
	if newclustercount < 0 {
		panic(fmt.Sprintf("Can not remov %d cluster, whent total number of clusters is %d", howmany, len(defaultCluster)))
	}
	newcluster := make([]env.GenomicInterval, newclustercount)
	for i := 0; i < newclustercount; i++ {
		newcluster[i] = defaultCluster[i]
	}
	eventCluster = env.RegionCollection(newcluster)
}

var defaultCluster env.RegionCollection
var eventCluster env.RegionCollection
var genClusterSwitch int64
