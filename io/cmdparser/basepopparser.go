package cmdparser

import (
	"fmt"
	"invade/env"
	"invade/fly"
	"invade/util"
	"math/rand"
	"strconv"
	"strings"
)

type countTuple struct {
	count  int64
	silent bool
}

func ParseBasePop(basepop string) *fly.Population {
	var countgrid [][]countTuple
	if strings.HasPrefix(basepop, "all") {
		countgrid = getAllCountGrrid(basepop)
	} else {
		countgrid = getPopCountGrrid(basepop)
	}

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

func getPopCountGrrid(basepop string) [][]countTuple {
	ysize, xsize := env.GetYSize(), env.GetXSize()
	popcount := make([][]countTuple, ysize)
	for i, _ := range popcount {
		popcount[i] = make([]countTuple, xsize)
	}

	//Y,X,count;Y,X,count
	toparse := make([]string, 0)
	if strings.Contains(basepop, ":") {
		toparse = strings.Split(basepop, ":")

	} else {
		toparse = append(toparse, basepop)
	}
	for _, tp := range toparse {
		silentb := false
		tmp := strings.Split(tp, ",")
		yco, ery := strconv.ParseInt(tmp[0], 10, 64)
		xco, erx := strconv.ParseInt(tmp[1], 10, 64)
		cstring := tmp[2]

		if strings.HasSuffix(cstring, "'") {
			silentb = true
			cstring = cstring[:len(cstring)-1]
		}
		countint, erc := strconv.ParseInt(cstring, 10, 64)
		if ery != nil || erx != nil || erc != nil {
			panic(fmt.Sprintf("Invalid base population entry %s", tp))
		}
		if yco >= ysize {
			panic(fmt.Sprintf("Y-coordinate of specimen in base popualtion outside of spatial grid: %d", yco))
		}
		if xco >= xsize {
			panic(fmt.Sprintf("Y-coordinate of specimen in base popualtion outside of spatial grid: %d", xco))
		}

		popcount[yco][xco] = countTuple{count: countint, silent: silentb}

	}

	return popcount
}

func getAllCountGrrid(basepop string) [][]countTuple {
	// default
	ysize, xsize := env.GetYSize(), env.GetXSize()
	popcount := make([][]countTuple, ysize)
	for i, _ := range popcount {
		popcount[i] = make([]countTuple, xsize)
	}

	// parse basepop and check validity
	tmp := strings.Split(basepop, ":")
	if len(tmp) != 2 || tmp[0] != "all" {
		panic(fmt.Sprintf("Invalid base population %s", basepop))
	}
	toparse := tmp[1]
	defsilent := false
	if strings.HasSuffix(toparse, "'") {
		defsilent = true
		toparse = toparse[:len(toparse)-1]
	}
	defcount, er := strconv.ParseInt(toparse, 10, 64)
	if er != nil {
		panic(fmt.Sprintf("Invalid base population count entry %d", defcount))
	}
	for y := 0; y < int(ysize); y++ {
		for x := 0; x < int(xsize); x++ {

			popcount[y][x] = countTuple{count: defcount, silent: defsilent}
		}
	}

	return popcount
}

/*
Example file
500 R 0; 1 100 200 400; 0 5 5000
250 F 0; 2 100 400;
250 M 0;;

func loadPopulationFromFile(file string, targetpopsize int64) *fly.Population {
	flies := make([]fly.Fly, 0)
	readFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		tmp := strings.Split(line, ";")
		if len(tmp) != 3 {
			panic(fmt.Sprintf("Invalid base population entry %s", line))
		}
		tempsplit := strings.Split(tmp[0], " ")
		if len(tempsplit) != 3 {
			panic(fmt.Sprintf("Invalid base population entry %s", line))
		}
		femhap := []int64{}
		malehap := []int64{}
		if tmp[1] != "" {
			femsplit := strings.Split(strings.TrimSpace(tmp[1]), " ")
			femsslice := sslice2islice(femsplit)
			femhap = util.UniqueSort(femsslice)
		}
		if tmp[2] != "" {
			malesplit := strings.Split(strings.TrimSpace(tmp[2]), " ")
			malesslice := sslice2islice(malesplit)
			malehap = util.UniqueSort(malesslice)
		}

		count, errcount := strconv.ParseInt(tempsplit[0], 10, 64)
		matpi, errmatpi := strconv.ParseInt(tempsplit[2], 10, 64)
		if errcount != nil || errmatpi != nil {
			panic(fmt.Sprintf("Invalid base population entry %s", line))
		}
		for i := int64(0); i < count; i++ {
			sex := getSex(tempsplit[1])
			f := fly.NewFly(femhap, malehap, sex, matpi)
			flies = append(flies, *f)
		}

	}
	readFile.Close()

	if len(flies) != int(targetpopsize) {
		panic("Invalid base population; population size does not match user specificiations")
	}
	return fly.InitializePopulation(flies)
}

func sslice2islice(sslice []string) []int64 {
	toret := make([]int64, 0)
	for _, s := range sslice {
		si, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Invalid base population character %s", s))
		}
		toret = append(toret, si)
	}
	return toret

}
*/
