package fly

import (
	"invade/env"
	"math"
	"testing"
)

type BiasCounter struct {
	bias  int64
	count int64
}

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

func testhelper_generateBiasTot(biaslist []BiasCounter) Population {
	flies := make([]Fly, 0, len(biaslist))
	for _, bc := range biaslist {
		tm := make(map[env.TEInsertion]int64)
		te := env.NewTEInsertion(bc.bias)
		tm[te] = bc.count

		f := Fly{FlyStat: &FlyStatistic{TotMap: tm}}
		flies = append(flies, f)
	}
	return Population{Flies: flies}
}

func testhelper_setdefaultenv() {
	env.SetupEnvironment([]int64{100, 100}, // two chromosomes of size 100
		[]int64{0, 0}, // two clusters of size 100
		[]float64{0, 0}, 0.1, 0.0)

}

func testhelper_hapmerger(haps [][]int64) *Population {
	flies := make([]Fly, 0)
	for i := 0; i < len(haps); i += 2 {
		femgam := haps[i]
		malegam := haps[i+1]
		f := NewFly(hm(femgam), hm(malegam))
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

func TestGetFlyStat(test *testing.T) {
	env.SetupEnvironment([]int64{1000, 1000}, // two chromosomes of size 1000
		[]int64{100, 100}, // two clusters of size 100
		[]float64{0, 0}, 0.1, 0.0)

	var tests = []struct {
		male        []int64
		female      []int64
		wantCluster int64
		wantTotal   int64
	}{{male: []int64{1, 2}, female: []int64{1001, 1002}, wantCluster: 4, wantTotal: 4},
		{male: []int64{999, 990}, female: []int64{1999, 1990}, wantTotal: 4},
		{male: []int64{999, 990, 1999, 1990}, female: []int64{}, wantTotal: 4},
		{male: []int64{1, 2, 1999}, female: []int64{999, 990, 1001, 1002}, wantCluster: 4, wantTotal: 7},
		{male: []int64{1, 2, 1999}, female: []int64{999, 990, 1001, 1002, 101, 1121}, wantCluster: 4, wantTotal: 9},
		{male: []int64{1, 2, 1999, 110}, female: []int64{999, 990, 1001, 1002, 101, 1121}, wantCluster: 4, wantTotal: 10},
	}

	for _, t := range tests {
		fsm := getFlyStat(hm(t.female), hm(t.male))
		fsf := getFlyStat(hm(t.male), hm(t.female))

		if t.wantCluster != fsm.CountCluster || fsm.CountCluster != fsf.CountCluster {
			test.Errorf("Incorrect number of cluster insertions; want %d, got %d, %d", t.wantCluster, fsm.CountCluster, fsf.CountCluster)
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

func TestAverageInsertionBias(test *testing.T) {
	var tests = []struct {
		totb []BiasCounter
		want float64
	}{
		{totb: []BiasCounter{{50, 1}, {50, 1}, {50, 1}, {50, 1}, {50, 1}, {0, 1}, {0, 1}, {0, 1}, {0, 1}, {0, 1}}, want: 25},
		{totb: []BiasCounter{{50, 10}, {0, 10}}, want: 25},
		{totb: []BiasCounter{{50, 10}, {25, 10}, {0, 10}}, want: 25},
		{totb: []BiasCounter{{50, 10}, {50, 9}, {50, 8}, {50, 7}, {50, 6}, {0, 5}, {0, 4}, {0, 3}, {0, 2}, {0, 1}}, want: 36.36364}}

	for _, t := range tests {
		pop := testhelper_generateBiasTot(t.totb)
		got := pop.GetAverageBias()
		if math.Abs(got-t.want) > 0.001 {
			test.Errorf("Incorrect average bias; got %f, want %f", got, t.want)
		}
	}
}

func TestBiasMapTotal(test *testing.T) {
	var tests = []struct {
		totb      []BiasCounter
		wantbias  int64
		wantcount int64
		wantlen   int64
	}{
		{totb: []BiasCounter{{50, 1}, {50, 1}, {50, 1}, {50, 1}, {50, 1}, {0, 1}, {0, 1}, {0, 1}, {0, 1}, {0, 1}}, wantbias: 50, wantcount: 5, wantlen: 2},
		{totb: []BiasCounter{{50, 10}, {0, 10}}, wantbias: 50, wantcount: 10, wantlen: 2},
		{totb: []BiasCounter{{50, 10}, {0, 10}}, wantbias: 0, wantcount: 10, wantlen: 2},
		{totb: []BiasCounter{{50, 10}, {25, 2}, {0, 10}}, wantbias: 25, wantcount: 2, wantlen: 3},
	}

	for _, t := range tests {
		pop := testhelper_generateBiasTot(t.totb)
		got := pop.GetBiasMapTotal()
		if len(got) != int(t.wantlen) {
			test.Errorf("Incorrect length of biasmap; got %d, want %d", len(got), t.wantlen)
		}
		tei := env.NewTEInsertion(t.wantbias)
		if got[tei] != t.wantcount {
			test.Errorf("Incorrect count for bias; got %d, want %d", got[tei], t.wantcount)
		}
	}
}

func TestConvertBiasMap(test *testing.T) {
	popsize := int64(10)
	var tests = []struct {
		totb      []BiasCounter
		wantbias  int64
		wantcount float64
		pos       int64
		wantlen   int64
	}{
		{totb: []BiasCounter{{50, 100}, {0, 20}}, pos: 0, wantbias: 50, wantcount: 10, wantlen: 2},
		{totb: []BiasCounter{{50, 100}, {0, 20}}, pos: 1, wantbias: 0, wantcount: 2, wantlen: 2},
		{totb: []BiasCounter{{50, 100}, {-100, 2}}, pos: 1, wantbias: -100, wantcount: 0.2, wantlen: 2},
		{totb: []BiasCounter{{50, 100}, {-100, 2}, {-50, 40}, {90, 90}}, pos: 1, wantbias: 90, wantcount: 9, wantlen: 4},
	}

	for _, te := range tests {
		bm := make(map[env.TEInsertion]int64)
		for _, b := range te.totb {
			te := env.NewTEInsertion(b.bias)
			bm[te] = b.count
		}
		got := convertBiasMapToSortedBiasSlice(bm, popsize)

		if len(got) != int(te.wantlen) {
			test.Errorf("invalid length, got %d want %d", len(got), te.wantlen)
		}
		if got[te.pos].Bias != te.wantbias {
			test.Errorf("invalid bias, got %d want %d", got[te.pos].Bias, te.wantbias)
		}

		if math.Abs(got[te.pos].AvCount-te.wantcount) > 0.001 {
			test.Errorf("invalid average count, got %f want %f", got[te.pos].AvCount, te.wantcount)
		}
	}

}
