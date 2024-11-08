package fly

import (
	"fmt"
	"invade/env"
	"math/rand"
	"sort"
)

type matePair struct {
	female *Fly
	male   *Fly
}

type neighborhood struct {
	xstart int64
	xend   int64
	ystart int64
	yend   int64
}

/*
helper construct for generating the mate pairs
*/
type cumFitFly struct {
	fly    *Fly
	cumFit float64 //cumulative fitness up to this fly

}

/*
for given coordinates, delinitate the neighborhood in the grid; considering the boundaries
*/
func getNeighborhoodCoordinates(ycord int64, xcord int64, radius int64) neighborhood {

	// are coordinates within bounds?
	ysize, xsize := env.GetYSize(), env.GetXSize()
	if ycord < 0 || ycord >= ysize {
		panic(fmt.Sprintf("invalid y-coordinate %d", ycord))
	}
	if xcord < 0 || xcord >= xsize {
		panic(fmt.Sprintf("invalid x-coordinate %d", xcord))
	}
	// find the coordinates of the neighborhood
	ystart, yend := ycord-radius, ycord+radius
	xstart, xend := xcord-radius, xcord+radius
	if ystart < 0 {
		ystart = 0
	}
	if yend >= ysize {
		yend = ysize - 1
	}
	if xstart < 0 {
		xstart = 0
	}
	if xend >= xsize {
		xend = xsize - 1
	}
	return neighborhood{
		xstart: xstart,
		xend:   xend,
		ystart: ystart,
		yend:   yend,
	}

}

/*
find the neighbors for a given coordinate
*/
func getNeighbors(flies [][]*Fly, ycord int64, xcord int64, radius int64) []*Fly {

	ncoord := getNeighborhoodCoordinates(ycord, xcord, radius)
	ncount := (ncoord.yend - ncoord.ystart + 1) * (ncoord.xend - ncoord.xstart + 1)
	neighbors := make([]*Fly, 0, ncount)

	for y := ncoord.ystart; y <= ncoord.yend; y++ {
		for x := ncoord.xstart; x <= ncoord.xend; x++ {
			neighbors = append(neighbors, flies[y][x])
		}
	}
	return neighbors
}

/*
 Get mate pairs;
 Get neighbors;
 consider selfing (in which case both parents are identical)
*/
func getMatePairs(flies [][]*Fly) [][]matePair {

	// initialize
	ysize, xsize := env.GetYSize(), env.GetXSize()
	merryCouples := make([][]matePair, ysize)
	for i, _ := range merryCouples {
		//make x
		merryCouples[i] = make([]matePair, xsize)
	}

	for y := 0; y < int(ysize); y++ {
		for x := 0; x < int(xsize); x++ {
			neighbors := getNeighbors(flies, int64(y), int64(x), env.GetMateRadius())
			neigcum := generateCumFitness(neighbors)
			fem := getFlyForRandomNumber(neigcum, rand.Float64())
			if rand.Float64() < env.GetSelfingRate() { // eg selfing rate 0.9 and number 0.0 - 0.89999 will result in selfing
				// here goes selfing
				merryCouples[y][x] = matePair{female: fem.fly, male: fem.fly} // for selfing, male and female is the same
			} else {
				// here goes non selfing
				male := getFlyForRandomNumber(neigcum, rand.Float64())
				// avoid selfing by mistake!!
				counter := 0
				for male.fly.FlyNumber == fem.fly.FlyNumber {
					male = getFlyForRandomNumber(neigcum, rand.Float64())
					counter++
					if counter > 5 { // prevent being stuck in an endless loop
						break
					}
				}
				merryCouples[y][x] = matePair{female: fem.fly, male: male.fly} // non-selfing: male and female are different
			}
		}
	}
	return merryCouples

}

func generateCumFitness(flies []*Fly) []cumFitFly {
	// Here major go confusion arose with pointers I guess
	// Video
	//https://www.youtube.com/watch?v=sTFJtxJXkaY
	// & ambersand; read it as address off!
	// strange behavior of range, which kind of overwrote the variable always with the same pointer
	// so that all flies in cumFitFly were pointing to the same fly in the end!!
	// This is a very strange behaviour -> major shit feature of GO
	// it may speed up things if sorted; largest fitness first; and cum sum between zero and 1
	sort.Slice(flies, func(i, j int) bool { return flies[i].Fitness > flies[j].Fitness })

	// now get the sum of all fitnesses
	var fitsum float64 = 0.0
	for _, f := range flies {
		if f.Fitness < 0.0 {
			panic("Fitness must be larger than zero")
		}
		fitsum += f.Fitness
	}
	// generate the cumFitFlies

	cumflies := make([]cumFitFly, 0, len(flies))
	var runningsum float64 = 0.0

	for i, f := range flies { // maybe use for loop if bug persists
		fi := &flies[i] // Woa that solves my POINTER BUG! todo check if still true! pointer pug could have been reintroduced
		// I Could not use &f as this was always referring to the same address
		// BE super careful with range and pointers
		//print(fi)
		w := f.Fitness / fitsum // fitness scaled by the total fitness such that the Sum of all is 1.0
		runningsum += w

		c := cumFitFly{fly: *fi, cumFit: runningsum}
		cumflies = append(cumflies, c)
	}
	return cumflies
}

/*
given a cumulative fitness array, pick a fly based on the random number;
fast binary search
*/
func getFlyForRandomNumber(cf []cumFitFly, randomIndex float64) *cumFitFly {
	if randomIndex < 0.0 || randomIndex >= 1.0 {
		panic(fmt.Sprintf("invaldi random index; must be between 0-1; got %f", randomIndex))
	}
	var hi int64 = int64(len(cf))
	var lo int64 = 0
	var mid int64 = 0
	for lo < hi {
		mid = int64((float64(lo+hi) / 2.0))
		midval := cf[mid].cumFit
		if midval < randomIndex {
			lo = mid + 1
		} else if midval > randomIndex {
			hi = mid
		} else {
			return &cf[mid+1]
		}
	}
	toret := &cf[lo]
	return toret
}
