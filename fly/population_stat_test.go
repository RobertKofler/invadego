package fly

import (
	"invade/env"
	"math"
	"testing"
)

func testhelper_generatetotalcount(clustercount []int64) *Population {
	flies := make([][]*Fly, 0)
	for i, cc := range clustercount {
		flies = append(flies, make([]*Fly, 0))

		f := &Fly{FlyStat: &FlyStatistic{CountTotal: cc}}
		flies[i] = append(flies[i], f)
	}
	return NewPopulation(flies, INVASION)
}

func testhelper_setdefaultenv() {
	env.SetupEnvironment(100,
		100, false, false,
		[]int64{100, 100}, // two chromosomes of size 100
		[]float64{0, 0},
		[]int64{}, 0.0,
		40, "dros", 0.0,
		0.1, 1000.0)
	SetupFitness(0.0, 0.0)

}

func testhelper_hapmerger(haps [][]int64) *Population {
	flies := make([][]*Fly, 0)
	flies = append(flies, make([]*Fly, 0))
	for i := 0; i < len(haps); i += 2 {
		femgam := haps[i]
		malegam := haps[i+1]
		f := NewFly(femgam, malegam, false)
		flies[0] = append(flies[0], f)
	}
	return NewPopulation(flies, INVASION)

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

func TestXTransect(test *testing.T) {
	testhelper_setdefaultenv()
	ysize, xsize := env.GetYSize(), env.GetXSize()
	pop := make([][]*Fly, ysize)
	for i, _ := range pop {
		pop[i] = make([]*Fly, xsize)
	}
	for y := 0; y < int(ysize); y++ {
		for x := 0; x < int(xsize); x++ {
			counter := float64(x) / 2.0
			pop[y][x] = &Fly{FlyStat: &FlyStatistic{CountTotal: int64(counter)}}
		}
	}
	newpop := NewPopulation(pop, INACTIVE)
	transect := newpop.GetXTransectAvCount()

	var tests = []struct {
		x    int64
		want float64
	}{
		{x: 1, want: 0.0},
		{x: 0, want: 0.0},
		{x: 50, want: 25.0},
		{x: 98, want: 49.0},
		{x: 99, want: 49.0},
	}

	for _, t := range tests {

		got := transect[t.x]
		if math.Abs(got-t.want) > 0.001 {
			test.Errorf("Incorrect x-transect at position %d; got %f, want %f", t.x, got, t.want)
		}

	}
}
