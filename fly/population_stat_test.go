package fly

import (
	"invade/env"
	"math"
	"testing"
)

func testhelper_generatewithclusterfrequency(clustercount []int64) Population {
	flies := make([]Fly, 0, len(clustercount))
	for _, cc := range clustercount {
		f := Fly{FlyStat: &FlyStatistic{CountCluster: cc}}
		flies = append(flies, f)
	}
	return Population{Flies: flies}
}
func testhelper_generatetotalcount(clustercount []int64) Population {
	flies := make([]Fly, 0, len(clustercount))
	for _, cc := range clustercount {
		f := Fly{FlyStat: &FlyStatistic{CountTotal: cc}}
		flies = append(flies, f)
	}
	return Population{Flies: flies}
}

func testhelper_setdefaultenv() {
	env.SetupEnvironment([]int64{100, 100}, // two chromosomes of size 100
		[]int64{0, 0}, // two clusters of size 100
		[]int64{0, 0}, // two reference regions of size 100
		[]bool{false}, //  trigger -> 0
		[]bool{false}, // para - > 1
		[]float64{0, 0}, 0.1)

}

func testhelper_hapmerger(haps [][]int64) *Population {
	flies := make([]Fly, 0)
	for i := 0; i < len(haps); i += 2 {
		femgam := haps[i]
		malegam := haps[i+1]
		f := NewFly(femgam, malegam, MALE, 0)
		flies = append(flies, *f)
	}
	return &Population{Flies: flies}
}

func TestGetWithClusterInsertionFrequency(test *testing.T) {

	var tests = []struct {
		clustercounts []int64
		want          float64
	}{{clustercounts: []int64{1, 1, 1, 1, 1, 0, 0, 0, 0, 0}, want: 0.5},
		{clustercounts: []int64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}, want: 0.1},
		{clustercounts: []int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 0}, want: 0.9},
		{clustercounts: []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, want: 0.0},
		{clustercounts: []int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, want: 1.0},
		{clustercounts: []int64{2, 3, 4, 5, 0, 0, 0, 0, 0, 10}, want: 0.5},
	}

	for _, t := range tests {
		pop := testhelper_generatewithclusterfrequency(t.clustercounts)
		got := pop.GetWithClusterInsertionFrequency()
		if math.Abs(got-t.want) > 0.001 {
			test.Errorf("Incorrect population frequency of cluster insertions; got %f, want %f", got, t.want)
		}

	}

}

func TestGetAverageClusterInsertions(test *testing.T) {

	var tests = []struct {
		clustercounts []int64
		want          float64
	}{{clustercounts: []int64{1, 1, 1, 1, 1, 0, 0, 0, 0, 0}, want: 0.5},
		{clustercounts: []int64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}, want: 0.1},
		{clustercounts: []int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 0}, want: 0.9},
		{clustercounts: []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, want: 0.0},
		{clustercounts: []int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, want: 1.0},
		{clustercounts: []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 10}, want: 1.0},
		{clustercounts: []int64{0, 0, 0, 0, 0, 10, 10, 10, 10, 10}, want: 5.0},
	}

	for _, t := range tests {
		pop := testhelper_generatewithclusterfrequency(t.clustercounts)
		got := pop.GetAverageClusterInsertions()
		if math.Abs(got-t.want) > 0.001 {
			test.Errorf("Incorrect average number of cluster insertions; got %f, want %f", got, t.want)
		}

	}

}

/*
Super important priority test; Checking whether Flystat are correctly computed!
FlyStat is key for most reported statistics
*/
func TestGetFlyStat(test *testing.T) {
	env.SetupEnvironment([]int64{1000, 1000}, // two chromosomes of size 1000
		[]int64{100, 100}, // two clusters of size 100
		[]int64{100, 100}, // two reference regions of size 100
		[]bool{true, false, false, false, false, false, false, false, false, false}, //  trigger -> 0
		[]bool{false, true, false, false, false, false, false, false, false, false}, // para - > 1
		[]float64{0, 0}, 0.1)

	var tests = []struct {
		male          []int64
		female        []int64
		wantCluster   int64
		wantReference int64
		wantTotal     int64
		wantTrigger   int64
		wantPara      int64
	}{{male: []int64{1, 2}, female: []int64{1001, 1002}, wantCluster: 4, wantTotal: 4},
		{male: []int64{999, 990}, female: []int64{1999, 1990}, wantReference: 4, wantTotal: 4},
		{male: []int64{999, 990, 1999, 1990}, female: []int64{}, wantReference: 4, wantTotal: 4},
		{male: []int64{101, 111}, female: []int64{1111, 1121}, wantPara: 4, wantTotal: 4},
		{male: []int64{100, 110}, female: []int64{1110, 1120}, wantTrigger: 4, wantTotal: 4},
		{male: []int64{1, 2, 1999}, female: []int64{999, 990, 1001, 1002}, wantCluster: 4, wantReference: 3, wantTotal: 7},
		{male: []int64{1, 2, 1999}, female: []int64{999, 990, 1001, 1002, 101, 1121}, wantCluster: 4, wantReference: 3, wantPara: 2, wantTotal: 9},
		{male: []int64{1, 2, 1999, 110}, female: []int64{999, 990, 1001, 1002, 101, 1121}, wantCluster: 4, wantReference: 3, wantPara: 2, wantTrigger: 1, wantTotal: 10},
	}

	for _, t := range tests {
		fsm := getFlyStat(t.female, t.male)
		fsf := getFlyStat(t.male, t.female)

		if t.wantCluster != fsm.CountCluster || fsm.CountCluster != fsf.CountCluster {
			test.Errorf("Incorrect number of cluster insertions; want %d, got %d, %d", t.wantCluster, fsm.CountCluster, fsf.CountCluster)
		}
		if t.wantReference != fsm.CountReference || fsm.CountReference != fsf.CountReference {
			test.Errorf("Incorrect number of cluster insertions; want %d, got %d, %d", t.wantReference, fsm.CountReference, fsf.CountReference)
		}
		if t.wantPara != fsm.CountPara || fsm.CountPara != fsf.CountPara {
			test.Errorf("Incorrect number of cluster insertions; want %d, got %d, %d", t.wantPara, fsm.CountPara, fsf.CountPara)
		}
		if t.wantTrigger != fsm.CountTrigger || fsm.CountTrigger != fsf.CountTrigger {
			test.Errorf("Incorrect number of cluster insertions; want %d, got %d, %d", t.wantPara, fsm.CountPara, fsf.CountPara)
		}
		if t.wantTotal != fsm.CountTotal || fsm.CountTotal != fsf.CountTotal {
			test.Errorf("Incorrect number of total insertions; want %d, got %d, %d", t.wantTotal, fsm.CountTotal, fsf.CountTotal)

		}

	}
}

func TestGetWithTEFrequency(test *testing.T) {
	var tests = []struct {
		totc []int64
		want float64
	}{{totc: []int64{1, 1, 1, 1, 1, 0, 0, 0, 0, 0}, want: 0.5},
		{totc: []int64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}, want: 0.1},
		{totc: []int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 0}, want: 0.9},
		{totc: []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, want: 0.0},
		{totc: []int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, want: 1.0},
		{totc: []int64{10, 10, 10, 10, 0, 0, 0, 0, 0, 10}, want: 0.5},
	}

	for _, t := range tests {
		pop := testhelper_generatetotalcount(t.totc)
		got := pop.GetWithTEFrequency()
		if math.Abs(got-t.want) > 0.001 {
			test.Errorf("Incorrect population frequency of cluster insertions; got %f, want %f", got, t.want)
		}

	}
}

func TestGetAverageInsertions(test *testing.T) {
	var tests = []struct {
		totc []int64
		want float64
	}{{totc: []int64{1, 1, 1, 1, 1, 0, 0, 0, 0, 0}, want: 0.5},
		{totc: []int64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}, want: 0.1},
		{totc: []int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 0}, want: 0.9},
		{totc: []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, want: 0.0},
		{totc: []int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, want: 1.0},
		{totc: []int64{10, 10, 10, 10, 0, 0, 0, 0, 0, 10}, want: 5.0},
	}

	for _, t := range tests {
		pop := testhelper_generatetotalcount(t.totc)
		got := pop.GetAverageInsertions()
		if math.Abs(got-t.want) > 0.001 {
			test.Errorf("Incorrect population frequency of cluster insertions; got %f, want %f", got, t.want)
		}

	}
}
func TestGetAveragePopulationFrequency(test *testing.T) {
	testhelper_setdefaultenv()
	var tests = []struct {
		haps [][]int64
		want float64
	}{{haps: [][]int64{[]int64{1, 10}, []int64{1, 10}}, want: 1.0},
		{haps: [][]int64{[]int64{1, 10}, []int64{}}, want: 0.5},
		{haps: [][]int64{[]int64{1, 10}, []int64{1, 10}, []int64{1, 10}, []int64{1, 10}}, want: 1.0},
		{haps: [][]int64{[]int64{1, 10}, []int64{3, 105}}, want: 0.5},
		{haps: [][]int64{[]int64{1, 10}, []int64{1, 10}, []int64{}, []int64{1, 10}}, want: 0.75},
		{haps: [][]int64{[]int64{1, 10}, []int64{1, 10}, []int64{}, []int64{}}, want: 0.5},
		{haps: [][]int64{[]int64{1, 10}, []int64{}, []int64{}, []int64{}}, want: 0.25},
	}

	for _, t := range tests {
		pop := testhelper_hapmerger(t.haps)
		got := pop.GetAveragePopulationFrequency()
		if math.Abs(got-t.want) > 0.001 {
			test.Errorf("Incorrect population frequency of cluster insertions; got %f, want %f", got, t.want)
		}

	}
}

func TestMHP(test *testing.T) {
	testhelper_setdefaultenv()
	var tests = []struct {
		haps [][]int64
		pos  int64
		want float64
	}{
		{haps: [][]int64{[]int64{10}, []int64{10}}, pos: 10, want: 1.0},
		{haps: [][]int64{[]int64{10}, []int64{10}, []int64{10}, []int64{10}}, pos: 10, want: 1.0},
		{haps: [][]int64{[]int64{10}, []int64{}, []int64{}, []int64{10}}, pos: 10, want: 0.5},
		{haps: [][]int64{[]int64{10}, []int64{}, []int64{}, []int64{}}, pos: 10, want: 0.25},
		{haps: [][]int64{[]int64{10}, []int64{10}, []int64{10}, []int64{10}}, pos: 0, want: 0.0},
	}

	for _, t := range tests {
		pop := testhelper_hapmerger(t.haps)
		got := pop.GetMHPPopulationFrequency()
		gotfreq := got[t.pos]
		if math.Abs(gotfreq-t.want) > 0.001 {
			test.Errorf("Incorrect population frequency of cluster insertions at position %d; got %f, want %f", t.pos, gotfreq, t.want)
		}

	}
}

func TestGetHaplotypes(test *testing.T) {
	testhelper_setdefaultenv()
	var tests = []struct {
		haps [][]int64
		want int64
	}{
		{haps: [][]int64{[]int64{10}, []int64{10}}, want: 2},
		{haps: [][]int64{[]int64{10}, []int64{10}, []int64{10}, []int64{10}}, want: 4},
		{haps: [][]int64{[]int64{}, []int64{}, []int64{}, []int64{}}, want: 4},
	}

	for _, t := range tests {
		pop := testhelper_hapmerger(t.haps)
		haps := pop.GetHaplotypes()
		got := int64(len(haps))
		if got != t.want {
			test.Errorf("Error, incorrect number of haploytpes; want %d got %d", t.want, got)
		}

	}
}
