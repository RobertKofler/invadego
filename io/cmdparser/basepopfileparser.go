package cmdparser

import (
	"bufio"
	"fmt"
	"invade/env"
	"invade/fly"
	"invade/util"
	"os"
	"strconv"
	"strings"
)

type hapTuple struct {
	hap1   []int64
	hap2   []int64
	silent bool
}

type rangeHapTuple struct {
	coordinates []env.XY
	hap1        []int64
	hap2        []int64
	silent      bool
}

/*
Example file
Y;X;silenced;hap1;hap2
1;2;0;1 100 200 400;0 5 5000
1-10;2-100;1;1 100 200 400;0 5 5000
1-10;1-50;1;1 100 200 400;0 5 5000
*/
func parseBasePopFile(basepopfile string) *fly.Population {
	countgrid := getHapCountGrid(basepopfile)
	ysize, xsize := env.GetYSize(), env.GetXSize()
	pop := make([][]*fly.Fly, ysize)
	for i, _ := range pop {
		pop[i] = make([]*fly.Fly, xsize)
	}
	for y := 0; y < int(ysize); y++ {
		for x := 0; x < int(xsize); x++ {
			silent := countgrid[y][x].silent

			// deep copies
			h1 := append([]int64(nil), countgrid[y][x].hap1...)
			h2 := append([]int64(nil), countgrid[y][x].hap2...)

			nf := fly.NewFly(h1, h2, silent) // the inital flys, per default do not receive parental silencing
			pop[y][x] = nf
		}
	}
	return fly.NewPopulation(pop, fly.INVASION)

}

func getHapCountGrid(basepopfile string) [][]hapTuple {

	pop := getEmptyHapTuples()
	// design decision; latest userspecificiation overwrites previous ones
	hapRanges := loadHapRanges(basepopfile)
	ysize, xsize := env.GetYSize(), env.GetXSize()
	for _, hapr := range hapRanges {
		for _, cord := range hapr.coordinates {
			ht := hapTuple{
				hap1:   hapr.hap1,
				hap2:   hapr.hap2,
				silent: hapr.silent}
			if cord.Y < 1 {
				panic(fmt.Sprintf("y coordinate must not be smaller than 1; got %d", cord.Y))
			}
			if cord.X < 1 {
				panic(fmt.Sprintf("x must not be smaller than 1; got %d", cord.X))
			}
			if cord.Y > ysize {
				panic(fmt.Sprintf("y  must not be larger than ysize; got %d; size %d", cord.Y, ysize))
			}
			if cord.X > xsize {
				panic(fmt.Sprintf("y must not be larger than xsize; got %d; size %d", cord.X, xsize))

			}
			pop[cord.Y-1][cord.X-1] = ht

		}

	}

	return pop

}

/*
Example file - every line is a rangeHapTuple
Y;X;silenced;hap1;hap2
1;2;0;1 100 200 400;0 5 5000
1-10;2-100;1;1 100 200 400;0 5 5000
1-10;1-50;1;1 100 200 400;0 5 5000
-> Returns the parsed entries of the file where each line corresponds to a rangeHapTuple
*/
func loadHapRanges(basepopfile string) []rangeHapTuple {
	readFile, err := os.Open(basepopfile)
	if err != nil {
		panic(err)
	}
	toret := make([]rangeHapTuple, 0)
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line == "" {
			continue // skip emtpy line
		}
		rht := parseSingleLine(line)
		toret = append(toret, rht)

	}
	readFile.Close()

	return toret

}

func parseSingleLine(line string) rangeHapTuple {
	// 1-10;1-50;1;1 100 200 400;0 5 5000
	// split; should give 5 columns
	tmp := strings.Split(line, ";")
	if len(tmp) != 5 {
		panic(fmt.Sprintf("Invalid base population entry %s", line))
	}
	try, trx, silenced, hap1, hap2 := tmp[0], tmp[1], tmp[2], tmp[3], tmp[4]
	cord := getCoordinates(try, trx)
	sil := stringToBool(silenced)

	femhap := parseHaplotpye(hap1)
	malehap := parseHaplotpye(hap2)
	return rangeHapTuple{
		coordinates: cord,
		hap1:        femhap,
		hap2:        malehap,
		silent:      sil}

}

func stringToBool(s string) bool {
	s = strings.TrimSpace(s)
	switch s {
	case "1":
		return true
	case "0":
		return false
	default:
		panic(fmt.Errorf("invalid boolean string: %q", s))
	}
}

/*
eg input '1 100 200 400' should give a slice with these integers
*/
func parseHaplotpye(hap string) []int64 {
	toret := []int64{}
	if hap != "" {
		split := strings.Split(strings.TrimSpace(hap), " ")
		tmp := make([]int64, 0)
		for _, s := range split {
			si, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("Invalid base population character %s", s))
			}
			tmp = append(tmp, si)
		}

		toret = util.UniqueSort(tmp)
	}

	return toret

}
func getEmptyHapTuples() [][]hapTuple {
	// initialize the grid for the next generation
	ysize, xsize := env.GetYSize(), env.GetXSize()
	emptyPop := make([][]hapTuple, ysize)
	for i, _ := range emptyPop {
		//make x
		emptyPop[i] = make([]hapTuple, xsize)
	}
	for y := 0; y < int(ysize); y++ {
		for x := 0; x < int(xsize); x++ {

			h := hapTuple{
				hap1:   []int64{},
				hap2:   []int64{},
				silent: false}
			emptyPop[y][x] = h

		}
	}
	return emptyPop

}
