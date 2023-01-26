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

func TestParseBasePopulationParser(test *testing.T) {

	var tests = []struct {
		toparse   string
		wantindex int64
		wantcount int64
		wantbias  byte
	}{
		{toparse: "100(0)", wantindex: 0, wantcount: 100, wantbias: 100}, // Note bias +100 (byte can not be negative, therefore i use a 100 shift)
		{toparse: "100(-100)", wantindex: 0, wantcount: 100, wantbias: 0},
		{toparse: "1(100)", wantindex: 0, wantcount: 1, wantbias: 200},
		{toparse: "1(1)", wantindex: 0, wantcount: 1, wantbias: 101},
		{toparse: "10(100),100(0)", wantindex: 0, wantcount: 10, wantbias: 200},
		{toparse: "10(100),100(0)", wantindex: 1, wantcount: 100, wantbias: 100},
		{toparse: "10(-100),100000(-50)", wantindex: 1, wantcount: 100000, wantbias: 50},
		{toparse: "10(-100),100000(-50)", wantindex: 0, wantcount: 10, wantbias: 0},
		{toparse: "1(1),1(1),100(0)", wantindex: 2, wantcount: 100, wantbias: 100},
	}
	for _, t := range tests {
		got := parseBasepopString(t.toparse)
		g := got[t.wantindex]

		if t.wantcount != g.count {
			test.Errorf("Incorrect TE count %d != %d", g.count, t.wantcount)
		}
		if t.wantbias != g.bias {
			test.Errorf("Incorrect bias count %d != %d", g.bias, t.wantbias)
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
	env.SetupEnvironment([]int64{10000, 10000}, []int64{0, 0}, []float64{1, 1}, 0.1, 0.0)
	var tests = []struct {
		popsize       int64
		basestring    string
		wantsites     int64
		wantpop       int64
		wantbias      byte
		wantbiascount int64
	}{
		{popsize: 10, basestring: "100(0)", wantsites: 100, wantpop: 10, wantbias: 100, wantbiascount: 100},
		{popsize: 1000, basestring: "0(0)", wantsites: 0, wantpop: 1000, wantbias: 0, wantbiascount: 0},
		{popsize: 1000, basestring: "50(0),50(100)", wantsites: 100, wantpop: 1000, wantbias: 200, wantbiascount: 50},
		{popsize: 1000, basestring: "50(0),50(100)", wantsites: 100, wantpop: 1000, wantbias: 100, wantbiascount: 50},
	}

	for _, test := range tests {
		got := loadPopulation(test.basestring, test.popsize)
		sites := got.GetInsertionSites()
		gotsites := len(sites)

		if len(got.Flies) != int(test.wantpop) {
			t.Errorf("Incorrect population size %d vs %d", len(got.Flies), test.wantpop)
		}
		if gotsites != int(test.wantsites) {
			t.Errorf("Incorrect number of sites %d vs %d", gotsites, test.wantsites)
		}
		gotbm := got.GetBiasMapTotal()
		gb := gotbm[env.TEInsertion(test.wantbias)]
		if gb != test.wantbiascount {
			t.Errorf("incorrect number of insertions with the requested bias; %d != %d", gb, test.wantbiascount)
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
