package cmdparser

import (
	"fmt"
	"invade/util"
	"strconv"
	"strings"
)

/*
Parse the genome definition string into a slice of chromosome sizes
*/
func ParseRegions(s string) []int64 {

	// check empty string
	if s == "" {
		return nil
	}

	var multiplier int64 = 1
	toproc := s
	if strings.Contains(s, ":") {

		tmp := strings.Split(s, ":")
		if len(tmp) > 2 {
			panic(fmt.Sprintf("incorrect command line argument; only a single ':' is allowed %s", s))
		}
		multiplier = getMultiplier(tmp[0])
		toproc = tmp[1]
	}
	var strs []string = []string{toproc}
	if strings.Contains(toproc, ",") {
		strs = strings.Split(toproc, ",")
	}

	toret := []int64{}
	for _, ss := range strs {
		nv, _ := strconv.ParseInt(ss, 10, 64) // 10 = base of int, 64 = int64
		nv *= multiplier
		toret = append(toret, nv)
	}
	return toret
}

/*
translate bp, kb or mb into 1, 1000, 1000000 respecitvely
*/
func getMultiplier(ms string) int64 {
	ms = strings.ToLower(ms)
	if ms == "bp" {
		return 1
	} else if ms == "kb" {
		return 1000
	} else if ms == "mb" {
		return 1000000
	} else {
		util.InvadeLogger.Fatal("Do not recognize unit " + ms)
	}
	return 0
}

/*
parse command line argument of recurrent sites eg 23:1,2,10,22 ;
these are sites occuring at regular intervals in the genome;
parser will return nil if an empty string was provided
*/
func ParseRecurrentRegions(s string) []bool {
	if s == "" {
		return nil
	}
	if !strings.Contains(s, ":") {
		panic("incorrect parameter; must contain a ':'")
	}
	tmp := strings.Split(s, ":")
	modulo, _ := strconv.ParseInt(tmp[0], 10, 64)
	var sites []int64
	if strings.Contains(tmp[1], ",") {
		tmpsite := strings.Split(tmp[1], ",")
		for _, i := range tmpsite {
			sit, _ := strconv.ParseInt(i, 10, 64)
			if sit >= modulo {
				panic(fmt.Sprintf("Invalid recurrent site %d; must not be larger or equal to modulo %d", sit, modulo))
			}
			sites = append(sites, sit)
		}

	} else {
		sit, _ := strconv.ParseInt(tmp[1], 10, 64)
		sites = append(sites, sit)
	}
	toret := make([]bool, modulo)
	for _, s := range sites {
		toret[s] = true
	}
	return toret

}

/*
	Parses recombination rate argument
	returns the recombination rate in cm/Mb for each window
*/
func ParseRecombination(s string) []float64 {
	if s == "" {
		util.InvadeLogger.Printf("no definition of the recombination rate provided; will assume 0 for each chromosome")
		return nil
	}
	var tmp []string
	if strings.Contains(s, ",") {
		tmp = strings.Split(s, ",")

	} else {
		tmp = append(tmp, s)
	}
	toret := make([]float64, len(tmp))
	for idx, i := range tmp {
		rit, _ := strconv.ParseFloat(i, 64)
		if rit > 49 {
			panic(fmt.Sprintf("Invalid recombination rate in cM/Mbp %f; must be smaller than 49cM/Mb", rit)) // the map function breaks down for high values
		}
		toret[idx] = rit
	}
	util.InvadeLogger.Printf("Parsed recombination rate; will proceed with a recombination rate of %v cM/Mb for the different chormosomes", toret)
	return toret

}
