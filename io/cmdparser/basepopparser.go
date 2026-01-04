package cmdparser

import (
	"fmt"
	"invade/env"
	"invade/fly"
	"invade/util"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type countTuple struct {
	count  int64
	silent bool
}

type rangeTuple struct {
	coordinates []env.XY
	count       int64
	silent      bool
}

func ParseBasePop(basepop string) *fly.Population {
	// test if file exists
	_, err := os.Stat(basepop)
	if err == nil {
		// if basepop is a file parse it
		return parseBasePopFile(basepop)
	} else {
		// otherwise just parse the regular basepop-definition
		return parseBasePopString(basepop)
	}
}

func parseBasePopString(basepop string) *fly.Population {

	countgrid := getPopCountGrrid(basepop)

	ysize, xsize := env.GetYSize(), env.GetXSize()
	pop := make([][]*fly.Fly, ysize)
	for i, _ := range pop {
		pop[i] = make([]*fly.Fly, xsize)
	}
	for y := 0; y < int(ysize); y++ {
		for x := 0; x < int(xsize); x++ {
			targetinsertions := countgrid[y][x]
			hap1 := make([]int64, 0)
			hap2 := make([]int64, 0)
			for i := 0; i < int(targetinsertions.count); i++ {
				genpos := env.GetRandomSite()
				if rand.Intn(2) == 1 {
					hap1 = append(hap1, genpos)
				} else {
					hap2 = append(hap2, genpos)
				}
			}
			h1 := util.UniqueSort(hap1)
			h2 := util.UniqueSort(hap2)
			nf := fly.NewFly(h1, h2, targetinsertions.silent) // the inital flys, per default do not receive parental silencing
			pop[y][x] = nf
		}
	}
	return fly.NewPopulation(pop, fly.INVASION)
}

func getCoordinates(yparse string, xparse string) []env.XY {
	var ystart, yend, xstart, xend int64
	// yrange
	if strings.Contains(yparse, "-") {
		ytmp := strings.Split(yparse, "-")
		yst, era := strconv.ParseInt(ytmp[0], 10, 64)
		yet, erb := strconv.ParseInt(ytmp[1], 10, 64)
		if era != nil || erb != nil {
			panic(fmt.Sprintf("Invalid base population position %s", yparse))
		}
		if yet < yst {
			yst, yet = yet, yst
		}
		ystart = yst
		yend = yet
	} else {
		ybt, erc := strconv.ParseInt(yparse, 10, 64)
		if erc != nil {
			panic(fmt.Sprintf("Invalid base population entry %s", yparse))
		}
		ystart = ybt
		yend = ybt
	}

	// xrange
	if strings.Contains(xparse, "-") {
		xtmp := strings.Split(xparse, "-")
		xst, era := strconv.ParseInt(xtmp[0], 10, 64)
		xet, erb := strconv.ParseInt(xtmp[1], 10, 64)
		if era != nil || erb != nil {
			panic(fmt.Sprintf("Invalid base population entry %s", xparse))
		}
		if xet < xst {
			xst, xet = xet, xst
		}
		xstart = xst
		xend = xet
	} else {
		xbt, erc := strconv.ParseInt(xparse, 10, 64)
		if erc != nil {
			panic(fmt.Sprintf("Invalid base population entry %s", xparse))
		}
		xstart = xbt
		xend = xbt
	}

	toret := make([]env.XY, 0)
	for y := ystart; y <= yend; y++ {
		for x := xstart; x <= xend; x++ {

			co := env.XY{
				X: x,
				Y: y}
			toret = append(toret, co)
		}
	}
	return toret

}

func getRangeTuples(basepop string) []rangeTuple {
	// Y,X,N:Y,X,N
	toparse := make([]string, 0)
	if strings.Contains(basepop, ":") {
		toparse = strings.Split(basepop, ":")

	} else {
		toparse = append(toparse, basepop)
	}
	toret := make([]rangeTuple, 0)

	for _, tp := range toparse {
		silentb := false
		tmp := strings.Split(tp, ",")
		if len(tmp) != 3 {
			panic(fmt.Sprintf("invalid base population entry; must have three info: %s", tp))
		}
		yparse, xparse, nparse := tmp[0], tmp[1], tmp[2]

		coordinates := getCoordinates(yparse, xparse)

		// count and silencing status
		if strings.HasSuffix(nparse, "'") {
			silentb = true
			nparse = nparse[:len(nparse)-1]
		}
		countint, erc := strconv.ParseInt(nparse, 10, 64)
		if erc != nil {
			panic(fmt.Sprintf("Invalid base population entry %s in %s", nparse, tp))
		}

		rt := rangeTuple{
			coordinates: coordinates,
			count:       countint,
			silent:      silentb,
		}
		toret = append(toret, rt)

	}
	return toret
}

func getPopCountGrrid(basepop string) [][]countTuple {

	// Y,X,N:Y,X,N
	ysize, xsize := env.GetYSize(), env.GetXSize()
	rtuples := getRangeTuples(basepop)

	popcount := make([][]countTuple, ysize)
	for i, _ := range popcount {
		popcount[i] = make([]countTuple, xsize)
	}

	for _, rt := range rtuples {
		for _, cord := range rt.coordinates {
			// validate and check if valid coordinates
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
			popcount[cord.Y-1][cord.X-1] = countTuple{count: rt.count, silent: rt.silent}
		}
	}

	return popcount
}
