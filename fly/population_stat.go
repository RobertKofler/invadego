package fly

import (
	"fmt"
	"invade/env"
	"sort"
)

/*
Get insertion sites
*/
func (p *Population) GetInsertionSites() []int64 {

	// hash, get the unique sites
	var insertionsites = make(map[int64]bool)
	for _, fly := range p.linearFlies {
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
Count of the silenced individuals
*/
func (p *Population) GetSilenced() int64 {
	c := int64(0)
	for _, f := range p.linearFlies {
		if f.Silenced {
			c++
		}
	}
	return c
}

/*
Frequency of the silenced individuals

*/
func (p *Population) GetSilencedFrequency() float64 {
	return p.Count2Freq(p.GetSilenced())
}

/*
	Get the fixed insertions in a population; the positions are provided
*/
func (p *Population) GetFixedInsertions() []int64 {
	// hash, get the unique sites
	var insertionsites = make(map[int64]int64)
	for _, fly := range p.linearFlies {
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

func (p *Population) medianPosition(numbers []int64) int64 {
	if len(numbers) == 0 {
		return 0
	}

	// Sort in-place (modifies the original slice)
	sort.Slice(numbers, func(i, j int) bool {
		return numbers[i] < numbers[j]
	})

	n := len(numbers)
	if n%2 == 1 {
		// Odd length: middle element
		return numbers[n/2]
	}

	// Even length: average of the two middle elements
	mid1 := numbers[n/2-1]
	mid2 := numbers[n/2]

	// Integer average (truncates toward zero)
	// Common convention for integer median when exact half isn't needed
	return mid1 + (mid2-mid1)/2
}

/*
Get the median x-position of the active invasion wave
*/
func (p *Population) GetMedianXPosActiveWave() int64 {
	//the max-x-position of any insertions having a TE
	maxxpos := make([]int64, len(p.Flies))

	for y, temp := range p.Flies {
		tmaxx := -1
		for x := range temp {
			cf := p.Flies[y][x].FlyStat.CountTotal
			if cf > 0 && x > tmaxx {
				tmaxx = x
			}
			maxxpos[y] = int64(tmaxx)
		}
	}
	return p.medianPosition(maxxpos)
}

/*
Get the median x-position of the silencing wave
*/
func (p *Population) GetMedianXPosSilencingWave() int64 {
	//the max-x-position of any insertions having a TE
	maxxpos := make([]int64, len(p.Flies))

	for y, temp := range p.Flies {
		tmaxx := -1
		for x := range temp {
			cf := p.Flies[y][x].Silenced
			if cf && x > tmaxx {
				tmaxx = x
			}
			maxxpos[y] = int64(tmaxx)
		}
	}
	return p.medianPosition(maxxpos)
}

/*
Get a 2D Matrix with the counts of TE insertions
*/
func (p *Population) GetCountGrid() [][]int64 {
	//ysize, xsize := env.GetYSize(), env.GetXSize()
	popcount := make([][]int64, len(p.Flies))
	for i, _ := range popcount {
		popcount[i] = make([]int64, 0)
	}

	for y, temp := range p.Flies {
		for x := range temp {
			popcount[y] = append(popcount[y], p.Flies[y][x].FlyStat.CountTotal)
		}
	}
	return popcount
}

/*
Get a summary for all the flies including coordinates, te count and silencing status
*/
func (p *Population) GetFlySummaries() []FlySummary {
	//ysize, xsize := env.GetYSize(), env.GetXSize()
	fsums := make([]FlySummary, 0, p.Size())

	for y, temp := range p.Flies {
		for x := range temp {
			cf := p.Flies[y][x]
			fs := FlySummary{CountTE: cf.FlyStat.CountTotal,
				Silenced: cf.Silenced,
				Denovo:   cf.Denovo,
				Ycoord:   int64(y),
				Xcoord:   int64(x)}
			fsums = append(fsums, fs)
		}
	}
	return fsums
}

func (p *Population) GetSilencedGrid() [][]bool {
	//ysize, xsize := env.GetYSize(), env.GetXSize()
	popcount := make([][]bool, len(p.Flies))
	for i, _ := range popcount {
		popcount[i] = make([]bool, 0)
	}

	for y, temp := range p.Flies {
		for x := range temp {
			popcount[y] = append(popcount[y], p.Flies[y][x].Silenced)
		}
	}
	return popcount
}

/*
Count number of flies with a TE
*/
func (p *Population) GetWithTECount() int64 {
	c := int64(0)
	for _, f := range p.linearFlies {
		if f.FlyStat.CountTotal > 0 {
			c++
		}
	}
	return c
}

/*
Count number of flies with a TE and silenced
*/
func (p *Population) GetWithTEAndSilencedCount() int64 {
	c := int64(0)
	for _, f := range p.linearFlies {
		if f.FlyStat.CountTotal > 0 && f.Silenced {
			c++
		}
	}
	return c
}

/*
Count number of flies with a TE and silenced
*/
func (p *Population) GetWithTEAndNotSilencedCount() int64 {
	c := int64(0)
	for _, f := range p.linearFlies {
		if f.FlyStat.CountTotal > 0 && (!f.Silenced) {
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
	toret := float64(count) / float64(p.Size())
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
	for _, f := range p.linearFlies {
		c += f.Fitness
	}
	toret := c / float64(p.Size())
	return toret
}

/*
Get the average number of TE insertions, but only for those having at least one insertion
*/
func (p *Population) GetAverageWithTE() float64 {
	c := float64(0.0)
	positive := 0
	for _, f := range p.linearFlies {
		c += float64(f.FlyStat.CountTotal)
		if f.FlyStat.CountTotal > 0 {
			positive++
		}
	}
	toret := c / float64(positive)
	return toret
}

/*
Get the average number of TE insertions
*/
func (p *Population) GetAverageInsertions() float64 {
	c := float64(0.0)
	for _, f := range p.linearFlies {
		c += float64(f.FlyStat.CountTotal)
	}
	toret := c / float64(p.Size())

	return toret
}

/*
Get the average population frequency of all TE insertions
*/
func (p *Population) GetMHPPopulationFrequency() map[int64]float64 {
	var insertionsites = make(map[int64]int64)
	for _, fly := range p.linearFlies {
		for _, is := range fly.Hap1 {
			insertionsites[is]++
		}
		for _, is := range fly.Hap2 {
			insertionsites[is]++
		}
	}

	insfreq := make(map[int64]float64)
	for pos, val := range insertionsites {
		valfreq := float64(val) / float64(2*p.Size()) // 2 times -> diploids
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

func (p *Population) GetXTransectAvCount() []float64 {
	//ysize, xsize := env.GetYSize(), env.GetXSize()
	transect := make([]float64, env.GetXSize())
	for x := int64(0); x < env.GetXSize(); x++ {
		for y := int64(0); y < env.GetYSize(); y++ {
			transect[x] += float64(p.Flies[y][x].FlyStat.CountTotal)
		}
		transect[x] = transect[x] / float64(env.GetYSize())
	}
	return transect
}
