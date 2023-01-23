package cmdparser

// command line, run all tests "go test ./..." yes three points
import (
	"invade/env"
	"invade/util"
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

func TestLoadGenome(t *testing.T) {
	util.SetSeed(8)
	env.SetupEnvironment([]int64{10000, 10000}, []int64{0, 0}, []int64{0, 0}, []float64{1, 1}, 0.1)
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
			t.Errorf("Incorrect number of sites %d vs %d", gotsites, test.wantsites)
		}

		for _, s := range sites {
			if s < 0 || s > 19999 {
				t.Errorf("Incorrect position of site; must be between 0-9999; got %d", s)
			}
		}
		for _, f := range got.Flies {
			if f.FlyNumber < 1 {
				t.Errorf("Incorrect Fly number got %d", f.FlyNumber)
			}
		}

	}

}

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
