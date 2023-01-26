package env

import (
	"fmt"
	"math/rand"
	"sort"
)

type TEInsertion byte

// report the insertion bias into piRNA clusters in %;
// possible values range from -100 over 0 to +100
// -100 no insertion will ever go into clusters
// 0 no insertion bias, insertion probability into piRNA cluster is just the size of piRNA cluster
// 100 every insertion will go into cluster
func (te TEInsertion) BiasPercent() int64 {
	bias := int64(te) - 100
	if bias < -100 || bias > 100 {
		panic(fmt.Sprintf("Invalid insertion bias, must be between -100 and 100, got %d", bias))
	}
	return bias

}

// report the insertion bias into piRNA cluster as fraction;
// possible values range from -1.0 over 0.0 to 1.0
func (te TEInsertion) BiasFraction() float64 {
	bias := float64(te)/100.0 - 1.0
	if bias < -1.0 || bias > 1.0 {
		panic(fmt.Sprintf("Invalid insertion bias, must be between -1.0 and 1.0, got %f", bias))
	}
	return bias
}

/*
introduces a stepwise mutation of the insertion bias; either randomly increase or decresase it by 1
*/
func (te TEInsertion) Mutate() TEInsertion {

	if rand.Float64() < 0.5 {
		return te.increase()

	} else {
		return te.decrease()
	}

}

func (te TEInsertion) increase() TEInsertion {
	current := byte(te)
	// solely increase if the maximum (200) has not yet been reached
	if current < 200 {
		current++
	}
	return TEInsertion(current)
}

func (te TEInsertion) decrease() TEInsertion {
	current := byte(te)
	// solely decrease if the minimum (0) has not yet been reached
	if current > 0 {
		current--
	}
	return TEInsertion(current)
}

/*
get the sorted keys for a map
*/
func GetSortedKeys(hap map[int64]TEInsertion) []int64 {
	// sort the recombinatin events by position
	keys := make([]int64, len(hap)) // if leng == capacity
	i := 0
	for k := range hap {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}

func MergeUniqueSort(sites1 map[int64]TEInsertion, sites2 map[int64]TEInsertion) []int64 {
	// add them and make them unique
	var insertionsites = make(map[int64]bool)
	for k, _ := range sites1 {
		insertionsites[k] = true
	}
	for k, _ := range sites2 {
		insertionsites[k] = true
	}
	// sort the recombinatin events by position
	keys := make([]int64, len(insertionsites)) // if leng == capacity
	i := 0
	for k := range insertionsites {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
