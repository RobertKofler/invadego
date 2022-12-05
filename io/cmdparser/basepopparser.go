package cmdparser

import (
	"bufio"
	"fmt"
	"invade/env"
	"invade/fly"
	"invade/util"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func ParseBasePop(basepop string, popsize int64) *fly.Population {
	if inscount, err := strconv.ParseInt(basepop, 10, 64); err == nil {
		return loadPopulation(inscount, popsize)
	} else {
		return loadPopulationFromFile(basepop, popsize)
	}
}

/*
 Load a fly population of a given popsize;
 randomly inserts 'inscount' TE insertions;
 multiple insertions at the same site are ignored
*/
func loadPopulation(inscount int64, popsize int64) *fly.Population {
	fhaps := make([][]int64, 2*popsize)
	for i := int64(0); i < 2*popsize; i++ {
		fhaps[i] = []int64{}
	}
	for i := int64(0); i < inscount; i++ {
		ri := rand.Int63n(2 * popsize)
		genpos := env.GetRandomSite()
		fhaps[ri] = append(fhaps[ri], genpos)
	}
	flies := make([]fly.Fly, popsize)
	for i := int64(0); i < popsize; i++ {
		hap1 := util.UniqueSort(fhaps[2*i])
		hap2 := util.UniqueSort(fhaps[2*i+1])
		sex := fly.GetRandomSex()
		nf := fly.NewFly(hap1, hap2, sex, 0)
		flies[i] = *nf
	}

	return fly.InitializePopulation(flies)
}

/*
Example file
500 R 0; 1 100 200 400; 0 5 5000
250 F 0; 2 100 400;
250 M 0;;
*/
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

func getSex(s string) fly.Sex {
	s = strings.ToUpper(s)
	if s == "M" {
		return fly.MALE
	} else if s == "F" {
		return fly.FEMALE
	} else if s == "R" {
		return fly.GetRandomSex()
	} else {
		panic("unknown sex specified")
	}
}
