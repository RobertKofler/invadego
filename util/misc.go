package util

import (
	"log"
	"os"
	"sort"
)

// declare the invade logger
var InvadeLogger *log.Logger = log.New(os.Stdout, "Invade: ", log.Ltime)

// no more error logger => panic
// var InvadeLoggerError *log.Logger = log.New(os.Stderr, "Invade: ", log.Ltime)

/*
Usefull utility function;
merge sites from multiple slices, make them unique and sort them
*/
func MergeUniqueSort(sites1 []int64, sites2 []int64) []int64 {
	// add them and make them unique
	var insertionsites = make(map[int64]bool)
	for _, is := range sites1 {
		insertionsites[is] = true
	}
	for _, is := range sites2 {
		insertionsites[is] = true
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

func UniqueSort(sites []int64) []int64 {
	// add them and make them unique
	var insertionsites = make(map[int64]bool)
	for _, is := range sites {
		insertionsites[is] = true
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
