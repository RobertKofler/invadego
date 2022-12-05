package fly

import (
	"fmt"
	"invade/env"
	"sort"
)

func (p *Population) GetInsertionSites() []int64 {

	// hash, get the unique sites
	var insertionsites = make(map[int64]bool)
	for _, fly := range p.Flies {
		for _, is := range fly.Hap1 {
			insertionsites[is] = true
		}
		for _, is := range fly.Hap2 {
			insertionsites[is] = true
		}
	}
	// sort the keys
	keys := make([]int64, len(insertionsites))
	i := 0
	for k := range insertionsites {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

/*
Count of the individuals with piRNAs
*/
func (p *Population) GetWithPirnaCount() int64 {
	c := int64(0)
	for _, f := range p.Flies {
		if f.Matpirna > 0 {
			c++
		}
	}
	return c
}

/*
Frequency of the individuals with piRNAs;
may be due to cluster or paramutation+trigger
*/
func (p *Population) GetWithPirnaFrequency() float64 {
	return p.Count2Freq(p.GetWithPirnaCount())
}

/*
	Get the fixed insertions in a population; the positions are provided
*/
func (p *Population) GetFixedInsertions() []int64 {
	// hash, get the unique sites
	var insertionsites = make(map[int64]int64)
	for _, fly := range p.Flies {
		for _, is := range fly.Hap1 {
			insertionsites[is]++
		}
		for _, is := range fly.Hap2 {
			insertionsites[is]++
		}
	}
	// sort the keys
	keys := make([]int64, 0)
	hapChrom := int64(2 * len(p.Flies)) // 2 times -> diploid
	for key, val := range insertionsites {
		if val == hapChrom {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func (p *Population) GetWithTECount() int64 {
	c := int64(0)
	for _, f := range p.Flies {
		if f.FlyStat.CountTotal > 0 {
			c++
		}
	}
	return c
}

/*
	helper function for a population;
	turns counts into a frequency
*/
func (p *Population) Count2Freq(count int64) float64 {
	toret := float64(count) / float64(len(p.Flies))
	if toret < 0 || toret > 1.0 {
		panic(fmt.Sprintf("Invalid frequency %f", toret))
	}
	return toret
}

func (p *Population) GetWithTEFrequency() float64 {
	return p.Count2Freq(p.GetWithTECount())
}

func (p *Population) GetAverageFitness() float64 {
	c := float64(0.0)
	for _, f := range p.Flies {
		c += f.Fitness
	}
	toret := c / float64(len(p.Flies))
	return toret
}

/*
Get the average number of TE insertions
*/
func (p *Population) GetAverageInsertions() float64 {
	c := float64(0.0)
	for _, f := range p.Flies {
		c += float64(f.FlyStat.CountTotal)
	}
	toret := c / float64(len(p.Flies))
	return toret
}

/*
Get the average population frequency of all TE insertions
*/
func (p *Population) GetMHPPopulationFrequency() map[int64]float64 {
	var insertionsites = make(map[int64]int64)
	for _, fly := range p.Flies {
		for _, is := range fly.Hap1 {
			insertionsites[is]++
		}
		for _, is := range fly.Hap2 {
			insertionsites[is]++
		}
	}

	insfreq := make(map[int64]float64)
	for pos, val := range insertionsites {
		valfreq := float64(val) / float64(2*len(p.Flies)) // 2 times -> diploids
		insfreq[pos] = valfreq

	}
	return insfreq
}

/*
Get the average population frequency of all TE insertions
*/
func (p *Population) GetAveragePopulationFrequency() float64 {
	insfreq := p.GetMHPPopulationFrequency()
	c := float64(0.0)
	keycount := float64(0.0)
	for _, valfreq := range insfreq {
		c += valfreq
		keycount += 1.0
	}
	return c / keycount
}

func (p *Population) GetPhase() Phase {
	return p.phase
}

func (p *Population) GetMinimumFitness() float64 {
	return p.minFit
}

/*
Number of males
*/
func (p *Population) GetMaleCount() int64 {
	c := int64(0)
	for _, f := range p.Flies {
		if f.Sex == MALE {
			c++
		}
	}
	return c
}

func (p *Population) GetMaleFrequency() float64 {
	return p.Count2Freq(p.GetMaleCount())
}

/*
Frequency of individuals with a cluster insertion
*/
func (p *Population) GetWithClusterInsertionFrequency() float64 {
	c := int64(0)
	for _, f := range p.Flies {
		if f.FlyStat.CountCluster > 0 {
			c++
		}
	}
	return p.Count2Freq(c)
}

/*
Frequency of individuals with an insertion into paramutable locus and piRNAs
*/
func (p *Population) GetWithParamutationYesPirnaFrequency() float64 {
	c := int64(0)
	for _, f := range p.Flies {
		if f.FlyStat.CountPara > 0 && f.Matpirna > 0 {
			c++
		}
	}
	return p.Count2Freq(c)
}

/*
Frequency of individuals with an insertion into paramutable locus but no piRNAs
*/
func (p *Population) GetWithParamutationNoPirnaFrequency() float64 {
	c := int64(0)
	for _, f := range p.Flies {
		if f.FlyStat.CountPara > 0 && f.Matpirna == 0 {
			c++
		}
	}
	return p.Count2Freq(c)
}

/*
Get the average number of cluster insertions
*/
func (p *Population) GetAverageClusterInsertions() float64 {
	c := float64(0.0)
	for _, f := range p.Flies {
		c += float64(f.FlyStat.CountCluster)
	}
	toret := c / float64(len(p.Flies))
	return toret
}

/*
Get the average number of paramutable insertions
*/
func (p *Population) GetAverageParaInsertions() float64 {
	c := float64(0.0)
	for _, f := range p.Flies {
		c += float64(f.FlyStat.CountPara)
	}
	toret := c / float64(len(p.Flies))
	return toret
}

func (p *Population) GetFixedClusterInsertionCount() int64 {
	fixedIns := p.GetFixedInsertions()
	fclu, _, _, _, _ := env.CountHaploidInsertions(fixedIns)
	return fclu
}

func (p *Population) GetFixedParaInsertionCount() int64 {
	fixedIns := p.GetFixedInsertions()
	_, _, fpara, _, _ := env.CountHaploidInsertions(fixedIns)
	return fpara
}

func (p *Population) GetPirnaOriginCount() int64 {
	var origins = make(map[int64]bool)
	for _, fly := range p.Flies {
		if fly.Matpirna > 0 {
			origins[fly.Matpirna] = true
		}
	}
	return int64(len(origins))
}

/*
 Get piRNA origin map
*/
func (p *Population) GetPirnaOriginFrequencies() []OriginFreq {
	var origins = make(map[int64]int64)
	for _, fly := range p.Flies {
		if fly.Matpirna > 0 {
			origins[fly.Matpirna]++
		}
	}

	toret := make([]OriginFreq, 0, len(origins))
	for key, value := range origins {
		freq := float64(value) / float64(len(p.Flies))
		toret = append(toret, OriginFreq{Origin: key, Freq: freq})
	}
	return toret
}

type OriginFreq struct {
	Origin int64
	Freq   float64
}

func haplotypeContainsPosition(s []int64, p int64) bool {
	for _, v := range s {
		if v == p {
			return true
		}
	}

	return false
}

/*
Compute linkage disequilibrium (D) between two loci. Use carefully, super slow.
*/
func (p *Population) GetD(locus1 int64, locus2 int64) float64 {
	//https://en.wikipedia.org/wiki/Linkage_disequilibrium
	x11 := 0
	p1 := 0
	p2 := 0
	haps := p.GetHaplotypes()
	hapcount := float64(len(haps))
	for _, hap := range haps {
		valid := true
		if haplotypeContainsPosition(hap, locus1) {
			p1++
		} else {
			valid = false
		}
		if haplotypeContainsPosition(hap, locus2) {
			p2++
		} else {
			valid = false
		}
		if valid {
			x11++
		}

	}
	fx11 := float64(x11) / hapcount
	fp1 := float64(p1) / hapcount
	fp2 := float64(p2) / hapcount
	d := fx11 - fp1*fp2
	return d

}
