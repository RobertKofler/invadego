// the biological environment of the simulations, e.g. chromosomes, recombination,
// transposition rate, piRNA cluster
package env

import (
	"fmt"
	"math"
	"math/rand"
)

/*
Population Structure;
migration barriers
*/
type popStruct struct {
	xBarriers   []int64
	barStrength float64
}

// coordiante Pair
type XY struct {
	X int64
	Y int64
}

type Environment struct {
	genome               *GenomicLandscape
	popStruct            *popStruct
	recombinationWindows []*RecombinationWindow
	triggerThreshold     int64
	epiMode              string
	psilencingLoss       float64 // propability that silencing can be lost
	minimumFitness       float64
	maximumInsertions    float64
	gridX                int64 // size of X-grid
	gridY                int64 // size of Y-grid
	closedX              bool
	closedY              bool
}

func GetMinimumFitness() float64 {
	return env.minimumFitness
}

/*
Propability that silencing can be lost
*/
func GetSilencingLossProbability() float64 {
	return env.psilencingLoss
}

/*
probabilty that a potential mate will be excluded from mating due to populatin structure
*/
func getExclusionProbabilty(xfirst int64, xsecond int64) float64 {
	return env.popStruct.ExclusionProbability(xfirst, xsecond)
}

func ExcludeMigrationBarrier(tofilter []*XY, xcord int64) []*XY {
	filtered := make([]*XY, 0, len(tofilter))
	for _, cand := range tofilter {
		if rand.Float64() < getExclusionProbabilty(cand.X, xcord) {
			// exclusion nothing happens
		} else {
			// inclusion
			filtered = append(filtered, cand)
		}
	}
	return filtered
}

func GetMaximumInsertions() float64 {
	return env.maximumInsertions
}

func (p *popStruct) ExclusionProbability(xfirst int64, xsecond int64) float64 {
	bc := p.barrierscrossed(xfirst, xsecond)
	// 1.0 total barrier,
	// 0.0 no effect
	// 0.1 weak barrier
	// two weak barriers each with 0.1 than we should have 0.19 (1-(1-b)(1-b))
	// similarly with two very strong barriers 0.9 we should have 0.99
	pone := 1.0 - p.barStrength

	effect := math.Pow(pone, float64(bc))
	toret := 1.0 - effect
	return toret
}

/*
calculate the numbers of barriers crossed by two potential mating partners, given the x coordinates
*/
func (p *popStruct) barrierscrossed(xfirst int64, xsecond int64) int64 {

	// cut the crap if no barriers or identical x-coordinates
	if len(p.xBarriers) == 0 || xfirst == xsecond {
		return 0.0
	}

	// sort by x-coordinate
	xleft, xright := xfirst, xsecond
	if xleft > xright {
		xleft, xright = xright, xleft
	}

	// how many barriers have we crossed?
	// with xsize 10 the coordinates 0...9 are used (0-based coordinates)
	// 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
	//
	// now lets assume a mate barrier at site 5 (1-based) this translates to 4 (0-based)
	//            pos 5
	// 0, 1, 2, 3, 4 || 5, 6, 7, 8, 9
	// valid mate barriers would be 1 and 9 in 1-based coordinates ()
	//    1							9
	// 0 || 1, 2, 3, 4, 5, 6, 7, 8 || 9
	toret := int64(0)
	for _, xbar := range p.xBarriers {
		if xbar > xleft && xbar <= xright {
			toret++
		}
	}

	return toret
}

func IsTriggered(diploidcount int64) bool {
	if env.triggerThreshold < 0 {
		return false // if smaller zero; it cannot be triggered
	} else if diploidcount >= env.triggerThreshold {
		return true // if above threshold triggered
	} else {
		return false // if below threshold not triggered
	}
}

func GetYSize() int64 {
	return env.gridY
}

func GetXSize() int64 {
	return env.gridX
}

func GetClosedY() bool {
	return env.closedY
}

func GetClosedX() bool {
	return env.closedX
}

func getX(x int64, xsize int64, closedx bool) int64 {
	if closedx {
		val := (x + xsize) % xsize
		return val

	} else {
		if x < 0 {
			return -1
		}
		if x >= xsize {
			return -1
		}
		return x
	}
}

func getY(y int64, ysize int64, closedy bool) int64 {
	if closedy {
		val := (y + ysize) % ysize
		return val

	} else {
		if y < 0 {
			return -1
		}
		if y >= ysize {
			return -1
		}
		return y
	}
}

/*
for given coordinates, delinitate the neighborhood in the grid; considering the boundaries
*/
func GetNeighborhoodCoordinates(ycord int64, xcord int64, radius int64) []*XY {

	// are coordinates within bounds?
	ysize, xsize := env.gridY, env.gridX
	closedx, closedy := env.closedX, env.closedY
	if ycord < 0 || ycord >= ysize {
		panic(fmt.Sprintf("invalid y-coordinate %d", ycord))
	}
	if xcord < 0 || xcord >= xsize {
		panic(fmt.Sprintf("invalid x-coordinate %d", xcord))
	}
	// find the coordinates of the neighborhood
	ystart, yend := ycord-radius, ycord+radius
	xstart, xend := xcord-radius, xcord+radius
	coords := make([]*XY, 0)

	for y := ystart; y <= yend; y++ {
		for x := xstart; x <= xend; x++ {
			xt := getX(x, xsize, closedx)
			yt := getY(y, ysize, closedy)
			if xt != -1 && yt != -1 {
				nc := XY{
					X: xt,
					Y: yt}
				coords = append(coords, &nc)
			}

		}
	}
	return coords

}

var env Environment

/*
An interval in a genome;
Start and End are both within the interval;
If Start == End the length of the interval is 1
*/
type GenomicInterval struct {
	Start int64
	End   int64
}

func (gi GenomicInterval) Length() int64 {
	return gi.End - gi.Start + 1
}

//
// ********** Genomic Landscape ****************+
//

// Construct a GenomicLandscape based on a list of chromosome sizes
func newGenomicLandscape(lens []int64) *GenomicLandscape {
	var offsets []int64
	var chrmSizes []int64
	var intervals []GenomicInterval
	var totGenome int64
	var currentOffset int64
	for _, len := range lens {
		totGenome += len
		offsets = append(offsets, int64(currentOffset))
		end := currentOffset + len - 1
		gi := GenomicInterval{Start: currentOffset, End: end}
		intervals = append(intervals, gi)
		chrmSizes = append(chrmSizes, len)
		currentOffset = end + 1
	}
	return &GenomicLandscape{offsets: offsets, chrmSizes: chrmSizes, intervals: intervals, totalGenome: totGenome}
}

/*
A 0-based linear representation of a genome.
Chromosomes are modeled as GenomicIntervals that occupy a range in a linear integer space.
The offset determines the start position of each chromosome in the linear space
*/
type GenomicLandscape struct {
	offsets     []int64
	chrmSizes   []int64
	intervals   []GenomicInterval
	totalGenome int64
}

/*
Get a random insertio site in the genome;
0-based; ranges from 0 to totalGenome-1
*/
func GetRandomSite() int64 {
	return int64(rand.Intn(int(env.genome.totalGenome)))
}

// TODO TEST
func TranslateCoordinates(pos int64) (int64, int64) {
	if pos >= env.genome.totalGenome {
		panic("invalid genomic position; larger than genome")
	}
	for i := len(env.genome.offsets) - 1; i >= 0; i-- {
		curos := env.genome.offsets[i]
		if pos >= curos {
			chrnum := int64(i + 1)
			chrpos := pos - curos + 1
			return chrnum, chrpos
		}
	}
	panic(fmt.Sprintf("invalid index; smaller than allowed %d", pos))
	/*
	   if(pos>this.genomeSize)throw new IllegalArgumentException("Invalid position outside genome");
	   int chromosome=this.chrSizes.size();
	   for(int i=this.offsets.size()-1; i>=0; i--)
	   {
	       int os=this.offsets.get(i);
	       if(pos>=os){
	           int chrbasedpos=pos-os;
	           return new ChromosomeBasedInsertion(chromosome,chrbasedpos);
	       }
	           chromosome--;
	   }
	   throw new IllegalArgumentException("Invalid position"+pos);
	*/

}
