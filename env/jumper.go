package env

import "invade/util"

type Jumper struct {
	u   float64 // activity without cluster insertion
	uc  float64 // residual activity with cluster insertion; typically 0.0
	sci bool    // selective cluins

}

var jump Jumper

func SetJumper(u float64, uc float64, selectiveCluins bool) {
	jump = Jumper{
		u:   u,
		uc:  uc,
		sci: selectiveCluins}
}

/*
Get the average number of transposition events for a DIPLOID.
For HAPLOIDS divide by two.
Considers u and uc, where u is the regular transposition rate and uc the transposition rate with a cluster insertion
*/
func (j *Jumper) getNovelInsertionCount(totalCount int64, scc int64, tcc int64) float64 {
	activeu := j.u
	if j.sci { // Selective cluster insertion is on
		if scc > 0 { // is the number of specific cluster insertions sufficient?
			activeu = j.uc
		}
	} else { // if selective cluster insertions is OFF,
		if tcc > 0 { // in this case the number of total cluster insertions sufficient?
			activeu = j.uc
		}
	}

	lambda := activeu * float64(totalCount)
	return lambda
}

/*
Get the positions of novel insertions for a haploid gamete; Input parameters are for a DIPLOID fly!
Number of required insertions is then divided by two to obtain estimates for haploid genomes.
returns a list of novel insertion sites; not unique, may contain same site twice
*/
func InsertNewTranspositionSites(gamete map[int64]TEInsertion, totalmap map[TEInsertion]int64, clustermap map[TEInsertion]int64, totclucount int64) map[int64]TEInsertion {

	// for each insertion bias
	// count the number of sites
	// count the number of cluster insertions
	//
	for te, totalcount := range totalmap {
		scc := clustermap[te]                                                               // specific cluster insertions
		newcountAverageDiploid := jump.getNovelInsertionCount(totalcount, scc, totclucount) // specific cluster count, and total cluster count
		newcountAverageHaploid := float64(newcountAverageDiploid) / 2.0                     // is this valid? see below
		newcountHaploid := util.Poisson(newcountAverageHaploid)

		// given an insertion bias (te) -> give me a list of novel insertion sites
		novelInsertions := GetSitesForBias(newcountHaploid, te.BiasFraction())
		// introduce the novel insertion sites into the gamete
		for _, pos := range novelInsertions {
			gamete[pos] = te
		}
	}
	return gamete
	/*
			Info from Internet
			Is the sum of two Poisson processes a Poisson process? (because of haplotype insertions)
		If they're independent of each other, yes. Indeed a more general result is that if there are k independent Poisson processes with rate λi,i=1,2,…,k,
		 then the combined process (the superposition of the component processes) is a Poisson process with rate ∑iλi
	*/
}
