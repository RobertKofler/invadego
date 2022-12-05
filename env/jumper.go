package env

import "invade/util"

type Jumper struct {
	u  float64 // activity without cluster insertion
	uc float64 // residual activity with cluster insertion; typically 0.0

}

var jump Jumper

func SetJumper(u float64, uc float64) {
	jump = Jumper{
		u:  u,
		uc: uc}
}

/*
	Get the average number of transposition events for a DIPLOID.
	For HAPLOIDS divide by two.
	Considers u and uc, where u is the regular transposition rate and uc the transposition rate with i) a cluster insertion or ii) maternal piRNAs and paramutable site

*/
func (j *Jumper) getNovelInsertionCount(totalCount int64, pirna bool) float64 {
	activeu := j.u
	if pirna {
		activeu = j.uc
	}
	lambda := activeu * float64(totalCount)
	return lambda
}

/*
Get the positions of novel insertions for a haploid gamete; Input parameters are for a DIPLOID fly!
Number of required insertions is then divided by two to obtain estimates for haploid genomes.
returns a list of novel insertion sites; not unique, may contain same site twice
*/
func GetNewTranspositionSites(totalCount int64, pirna bool) []int64 {
	newcountAverageDiploid := jump.getNovelInsertionCount(totalCount, pirna)
	newcountAverageHaploid := float64(newcountAverageDiploid) / 2.0 // is this valid? see below
	newcountHaploid := util.Poisson(newcountAverageHaploid)
	toret := make([]int64, newcountHaploid)
	for i := int64(0); i < newcountHaploid; i++ {
		toret[i] = GetRandomSite()
	}
	return toret
	/*
			Info from Internet
			Is the sum of two Poisson processes a Poisson process? (because of haplotype insertions)
		If they're independent of each other, yes. Indeed a more general result is that if there are k independent Poisson processes with rate λi,i=1,2,…,k,
		 then the combined process (the superposition of the component processes) is a Poisson process with rate ∑iλi
	*/
}
