package cmdparser

// command line, run all tests "go test ./..." yes three points
import (
	"invade/env"
	"testing"
)

func TestParseRecombination(t *testing.T) {

	var tests = []struct {
		toparse string
		want    []float64
	}{
		{toparse: "", want: nil},
		{toparse: "2", want: []float64{2}},
		{toparse: "2,3,1,4", want: []float64{2, 3, 1, 4}}}
	for _, test := range tests {
		got := ParseRecombination(test.toparse)
		want := test.want
		if len(got) != len(want) {
			t.Errorf("Incorrect length %v %v", got, want)
		}
		for i, k := range got {
			if k != want[i] {
				t.Errorf("Incorrect entry %f %f", k, want[i])
			}
		}

	}
}

func TestParseRecombinationMultiple(t *testing.T) {

	res := ParseRecombination("2,3,2,4,5")

	if len(res) != 5 {
		t.Error("wrong length")
	}
	if res[4] != 5 {
		t.Error("wrong recombination")
	}
}

func TestParseCoordinateRanges(t *testing.T) {
	var tests = []struct {
		tpx       string
		tpy       string
		wantcount int64
		winc      []env.XY
	}{
		{tpx: "1", tpy: "1", wantcount: 1, winc: []env.XY{{X: 1, Y: 1}}},
		{tpx: "1", tpy: "1-1", wantcount: 1, winc: []env.XY{{X: 1, Y: 1}}},
		{tpx: "1-1", tpy: "1", wantcount: 1, winc: []env.XY{{X: 1, Y: 1}}},
		{tpx: "10", tpy: "23", wantcount: 1, winc: []env.XY{{X: 10, Y: 23}}},
		{tpx: "1-10", tpy: "1", wantcount: 10, winc: []env.XY{{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, {X: 4, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1}, {X: 8, Y: 1}, {X: 9, Y: 1}, {X: 10, Y: 1}}},
		{tpx: "10-1", tpy: "1", wantcount: 10, winc: []env.XY{{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}, {X: 4, Y: 1}, {X: 5, Y: 1}, {X: 6, Y: 1}, {X: 7, Y: 1}, {X: 8, Y: 1}, {X: 9, Y: 1}, {X: 10, Y: 1}}},
		{tpx: "1", tpy: "1-10", wantcount: 10, winc: []env.XY{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}, {X: 1, Y: 4}, {X: 1, Y: 5}, {X: 1, Y: 6}, {X: 1, Y: 7}, {X: 1, Y: 8}, {X: 1, Y: 9}, {X: 1, Y: 10}}},
		{tpx: "1", tpy: "10-1", wantcount: 10, winc: []env.XY{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}, {X: 1, Y: 4}, {X: 1, Y: 5}, {X: 1, Y: 6}, {X: 1, Y: 7}, {X: 1, Y: 8}, {X: 1, Y: 9}, {X: 1, Y: 10}}},
		{tpx: "5-50", tpy: "1", wantcount: 46, winc: []env.XY{{X: 5, Y: 1}, {X: 50, Y: 1}}},
		{tpx: "1", tpy: "5-50", wantcount: 46, winc: []env.XY{{X: 1, Y: 5}, {X: 1, Y: 50}}},
		{tpx: "1-10", tpy: "1-10", wantcount: 100, winc: []env.XY{{X: 1, Y: 1}, {X: 1, Y: 10}, {X: 10, Y: 1}, {X: 10, Y: 10}}},
		{tpx: "6-10", tpy: "6-10", wantcount: 25, winc: []env.XY{{X: 6, Y: 6}, {X: 6, Y: 10}, {X: 10, Y: 6}, {X: 10, Y: 10}}},
	}
	for _, test := range tests {
		got := getCoordinates(test.tpy, test.tpx)

		if len(got) != int(test.wantcount) {
			t.Errorf("Incorrect length %v; want %v", len(got), test.wantcount)
		}
		for _, xy := range test.winc {
			inside := false
			for _, item := range got {
				if item == xy {
					inside = true
				}
			}
			if !inside {
				t.Errorf("Incorrect parsing does not contain XY %v; got %v", xy, got)
			}

		}

	}

}

func TestParseGenome(t *testing.T) {
	var tests = []struct {
		toparse string
		want    []int64
	}{
		{toparse: "", want: nil},
		{toparse: "2", want: []int64{2}},
		{toparse: "kb:2", want: []int64{2000}},
		{toparse: "2,3,1,4", want: []int64{2, 3, 1, 4}},
		{toparse: "kb:2,3,1,4", want: []int64{2000, 3000, 1000, 4000}},
		{toparse: "mb:2,3", want: []int64{2000000, 3000000}}}
	for _, test := range tests {
		got := ParseRegions(test.toparse)
		want := test.want
		if len(got) != len(want) {
			t.Errorf("Incorrect length %v %v", got, want)
		}
		for i, k := range got {
			if k != want[i] {
				t.Errorf("Incorrect entry %d %d", k, want[i])
			}
		}

	}

}

/*
func TestLoadGenome(t *testing.T) {
	util.SetSeed(7)
	env.SetupEnvironment([]int64{5000, 5000}, []int64{0, 0}, []int64{0, 0}, []bool{}, []bool{}, []float64{1, 1}, 0.1, 1000.0)
	fly.SetupFitness(0, 0, true, false)
	var tests = []struct {
		popsize   int64
		inscount  int64
		wantsites int64
		wantpop   int64
	}{
		{popsize: 10, inscount: 100, wantsites: 100, wantpop: 10},
		{popsize: 1000, inscount: 0, wantsites: 0, wantpop: 1000},
		{popsize: 1000, inscount: 100, wantsites: 100, wantpop: 1000},
	}

	for _, test := range tests {
		got := loadPopulation(test.inscount, test.popsize)
		sites := got.GetInsertionSites()
		gotsites := len(sites)

		if len(got.Flies) != int(test.wantpop) {
			t.Errorf("Incorrect population size %d vs %d", len(got.Flies), test.wantpop)
		}
		if gotsites != int(test.wantsites) {
			t.Errorf("Incorrect number of sites %d vs %d", gotsites, test.wantpop)
		}

		for _, s := range sites {
			if s < 0 || s > 9999 {
				t.Errorf("Incorrect position of site; must be between 0-9999; got %d", s)
			}
		}
		for _, f := range got.Flies {
			if f.FlyNumber < 1 {
				t.Errorf("Incorrect Fly number got %d", f.FlyNumber)
			}
			if f.Matpirna != 0 {
				t.Errorf("Incorrect maternal piRNAs; got %d", f.Matpirna)
			}
			if f.Sex != fly.FEMALE && f.Sex != fly.MALE {
				t.Errorf("Incorrect sex; got %d", f.Sex)
			}
		}

	}

}
*/

func TestRecurrentRegion(t *testing.T) {
	var tests = []struct {
		toparse string
		want    []bool
	}{
		{toparse: "", want: nil},
		{toparse: "3:1", want: []bool{false, true, false}},
		{toparse: "3:1,2", want: []bool{false, true, true}},
		{toparse: "5:1,2", want: []bool{false, true, true, false, false}}}
	for _, test := range tests {
		got := ParseRecurrentRegions(test.toparse)
		want := test.want
		if len(got) != len(want) {
			t.Errorf("Incorrect length %v %v", got, want)
		}
		for i, k := range got {
			if k != want[i] {
				t.Errorf("Incorrect entry %t %t", k, want[i])
			}
		}

	}

}
