package cmdparser

import (
	"bufio"
	"fmt"
	"invade/env"
	"invade/fly"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type biasparse struct {
	count int64
	bias  int64
}

func ParseBasePop(basepop string, popsize int64) *fly.Population {

	if strings.HasPrefix(basepop, "file:") {

		return loadPopulationFromFile(basepop[5:], popsize)

	} else {
		return loadPopulation(basepop, popsize)
	}
}

/*
Load a fly population of a given popsize;
randomly inserts 'inscount' TE insertions;
multiple insertions at the same site are ignored
*/
func loadPopulation(basepop string, popsize int64) *fly.Population {
	fhaps := make([]map[int64]env.TEInsertion, 2*popsize)
	for i := int64(0); i < 2*popsize; i++ {
		fhaps[i] = make(map[int64]env.TEInsertion)
	}
	parsed := parseBasepopString(basepop)
	for _, ic := range parsed { // insertion class
		teinsertion := env.NewTEInsertion(ic.bias)
		biasf := teinsertion.BiasFraction()
		inssites := env.GetSitesForBias(ic.count, biasf)
		for _, is := range inssites {
			ri := rand.Int63n(2 * popsize)
			fhaps[ri][is] = teinsertion

		}
	}
	flies := make([]fly.Fly, popsize)
	for i := int64(0); i < popsize; i++ {
		hap1 := fhaps[2*i]
		hap2 := fhaps[2*i+1]
		nf := fly.NewFly(hap1, hap2)
		flies[i] = *nf
	}

	return fly.InitializePopulation(flies)
}

/*
parse base population string; eg 100(-100),200(0);
code nice example on how to use anonymous struct, but uff quite cumbersome; note bias = value+100 (byte can not be negative)
*/
func parseBasepopString(basepop string) []biasparse {
	reg := regexp.MustCompile(`(?P<Count>\d+)\((?P<Bias>-?\d+)\)`)
	var toret = make([]biasparse, 0)

	toparse := []string{basepop}
	if strings.Contains(basepop, ",") {
		toparse = strings.Split(basepop, ",")
	}

	for _, sus := range toparse {
		match := reg.FindStringSubmatch(sus)
		count, _ := strconv.ParseInt(match[1], 10, 64)
		bias, _ := strconv.ParseInt(match[2], 10, 64)
		if bias < -100 || bias > 100 {
			panic(fmt.Sprintf("Invalid insertion bias, must be between -100 and 100, got %d", bias))
		}
		if count < 0 {
			panic(fmt.Sprintf("Invalid insertion count, must not be smaller than zero; got %d", count))
		}
		toret = append(toret, biasparse{count: count, bias: bias})

	}
	return toret
}

func parseTEEntry(entry string) (int64, env.TEInsertion) {

	reg := regexp.MustCompile(`(?P<Count>\d+)\((?P<Bias>-?\d+)\)`)
	match := reg.FindStringSubmatch(entry)
	position, _ := strconv.ParseInt(match[1], 10, 64)
	bias, _ := strconv.ParseInt(match[2], 10, 64)
	te := env.NewTEInsertion(bias)
	return position, te
}

func cloneHapmap(toclone map[int64]env.TEInsertion) map[int64]env.TEInsertion {
	toret := make(map[int64]env.TEInsertion)
	for p, te := range toclone {
		toret[p] = te
	}
	return toret
}

/*
Example file
500; 1 100 200 400; 0 5 5000
250; 2 100 400;
250;;
*/
func loadPopulationFromFile(file string, targetpopsize int64) *fly.Population {
	flies := make([]fly.Fly, 0)
	readFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() { // for each line
		line := fileScanner.Text()
		tmp := strings.Split(line, ";")
		if len(tmp) != 3 {
			panic(fmt.Sprintf("Invalid base population entry %s", line))
		}

		count, errcount := strconv.ParseInt(tmp[0], 10, 64)
		if errcount != nil {
			panic(fmt.Sprintf("Invalid number of flies; must be integer; got %s", tmp[0]))
		}

		femhap := make(map[int64]env.TEInsertion)
		malehap := make(map[int64]env.TEInsertion)
		if tmp[1] != "" {
			femsplit := strings.Split(strings.TrimSpace(tmp[1]), " ")
			for _, f := range femsplit {
				pos, te := parseTEEntry(f)
				femhap[pos] = te
			}
		}
		if tmp[2] != "" {
			malesplit := strings.Split(strings.TrimSpace(tmp[2]), " ")
			for _, m := range malesplit {
				pos, te := parseTEEntry(m)
				malehap[pos] = te
			}
		}

		for i := int64(0); i < count; i++ {
			f := fly.NewFly(cloneHapmap(femhap), cloneHapmap(malehap))
			flies = append(flies, *f)
		}

	}
	readFile.Close()

	if len(flies) != int(targetpopsize) {
		panic("Invalid base population; population size does not match user specificiations")
	}
	return fly.InitializePopulation(flies)
}
