package fly

import (
	"invade/util"
	"math"
	"testing"
)

func TestRecombine(t *testing.T) {
	var tests = []struct {
		hap1 []int64
		hap2 []int64
		rec  []int64
		want []int64
	}{
		{hap1: []int64{1, 2, 3}, hap2: []int64{1, 2, 3}, rec: []int64{1, 2, 3}, want: []int64{1, 2, 3}},
		{hap1: []int64{1, 2, 3}, hap2: []int64{4, 5, 6}, rec: []int64{}, want: []int64{1, 2, 3}},
		{hap1: []int64{1, 2, 3}, hap2: []int64{4, 5, 6}, rec: []int64{1}, want: []int64{4, 5, 6}},
		{hap1: []int64{4, 5, 6}, hap2: []int64{1, 2, 3}, rec: []int64{1}, want: []int64{1, 2, 3}},
		{hap1: []int64{1, 2, 3}, hap2: []int64{4, 5, 6}, rec: []int64{7}, want: []int64{1, 2, 3}},
		{hap1: []int64{4, 5, 6}, hap2: []int64{1, 2, 3}, rec: []int64{7}, want: []int64{4, 5, 6}},
		{hap1: []int64{1, 2, 3}, hap2: []int64{4, 5, 6}, rec: []int64{4}, want: []int64{1, 2, 3, 4, 5, 6}},
		{hap1: []int64{4, 5, 6}, hap2: []int64{1, 2, 3}, rec: []int64{4}, want: []int64{}},
		{hap1: []int64{1, 2, 3}, hap2: []int64{4, 5, 6}, rec: []int64{6}, want: []int64{1, 2, 3, 6}},
		{hap1: []int64{1, 3, 5}, hap2: []int64{2, 4, 6}, rec: []int64{}, want: []int64{1, 3, 5}},
		{hap1: []int64{1, 3, 5}, hap2: []int64{2, 4, 6}, rec: []int64{2, 3}, want: []int64{1, 2, 3, 5}},
		{hap1: []int64{1, 3, 5}, hap2: []int64{2, 4, 6}, rec: []int64{2, 3, 4, 5, 6}, want: []int64{1, 2, 3, 4, 5, 6}},
	}

	for _, test := range tests {
		f := Fly{
			Hap1: test.hap1,
			Hap2: test.hap2}
		got := f.recombine(test.rec)
		want := test.want
		if len(got) != len(want) {
			t.Errorf("recombine(); unequal length %v vs %v", got, want)
		} else {
			for i, val := range got {
				if want[i] != val {
					t.Errorf("recombine(); Different entries at position %d;  %v vs %v", i, got, want)
				}
			}
		}

	}
}

/*
fitness function w=1-xn^t
cluster insertions (+reference insertions) may be considered
*/
func TestFitnessOmxnt(t *testing.T) {
	var tests = []struct {
		x    float64
		ct   int64 // count total
		cc   int64 // count cluster
		cr   int64 // count reference
		nx   bool  // noxclusterinsertion
		t    float64
		want float64
	}{
		{x: 0.1, ct: 2, cc: 0, cr: 0, nx: false, t: 1.0, want: 0.8},
		{x: 0.1, ct: 10, cc: 0, cr: 0, nx: false, t: 1.0, want: 0.0},
		{x: 0.1, ct: 100, cc: 0, cr: 0, nx: false, t: 1.0, want: 0.0}, // ceck whether min w is 0.0
		{x: 0.1, ct: 2, cc: 0, cr: 0, nx: false, t: 1.5, want: 0.7171573},
		{x: 0.1, ct: 10, cc: 8, cr: 0, nx: false, t: 1.0, want: 0.0},
		{x: 0.1, ct: 10, cc: 8, cr: 0, nx: true, t: 1.0, want: 0.8},
		{x: 0.1, ct: 10, cc: 0, cr: 8, nx: true, t: 1.0, want: 0.8},
		{x: 0.1, ct: 10, cc: 4, cr: 4, nx: true, t: 1.0, want: 0.8},
		{x: 0.1, ct: 10, cc: 4, cr: 4, nx: true, t: 1.5, want: 0.7171573},
	}

	for _, test := range tests {
		ff = FitnessFunction{x: test.x, t: test.t, noxincluins: test.nx}
		want := test.want
		got := ff.ComputeFitness(test.ct, test.cc, test.cr)
		if math.Abs(want-got) > 0.0001 {
			t.Errorf("ff.ComputeFitness(%d,%d,%d) != %f; got = %f", test.ct, test.cc, test.cr, test.want, got)
		}
	}
}

func TestGenerateCumFitness(test *testing.T) {
	var tests = []struct {
		flies []Fly
		want  []float64
	}{
		{flies: []Fly{Fly{Fitness: 1.0}}, want: []float64{1.0}},
		{flies: []Fly{Fly{Fitness: 10.0}}, want: []float64{1.0}},
		{flies: []Fly{Fly{Fitness: 10.0}, Fly{Fitness: 10.0}}, want: []float64{0.5, 1.0}},
		{flies: []Fly{Fly{Fitness: 5.0}, Fly{Fitness: 15.0}}, want: []float64{0.75, 1.0}},
		{flies: []Fly{Fly{Fitness: 15.0}, Fly{Fitness: 5.0}}, want: []float64{0.75, 1.0}},
		{flies: []Fly{Fly{Fitness: 1.0}, Fly{Fitness: 2.0}, Fly{Fitness: 3.0}, Fly{Fitness: 4.0}}, want: []float64{0.4, 0.7, 0.9, 1.0}},
		{flies: []Fly{Fly{Fitness: 4.0}, Fly{Fitness: 2.0}, Fly{Fitness: 1.0}, Fly{Fitness: 3.0}}, want: []float64{0.4, 0.7, 0.9, 1.0}},
		{flies: []Fly{Fly{Fitness: 0.4}, Fly{Fitness: 0.2}, Fly{Fitness: 0.1}, Fly{Fitness: 0.3}}, want: []float64{0.4, 0.7, 0.9, 1.0}},
		{flies: []Fly{Fly{Fitness: 0.04}, Fly{Fitness: 0.02}, Fly{Fitness: 0.01}, Fly{Fitness: 0.03}}, want: []float64{0.4, 0.7, 0.9, 1.0}},
	}
	for _, t := range tests {
		cf := generateCumFitness(t.flies)
		for i, w := range t.want {
			if math.Abs(w-cf[i].cumFit) > 0.001 {
				test.Errorf("Incorrect cumulative fitness; got %f wanted %f", cf[i].cumFit, w)
			}
		}

	}
}

func TestGetFlyForRandomNumber(test *testing.T) {
	var tests = []struct {
		flies []cumFitFly
		index float64
		want  int64
	}{
		{flies: []cumFitFly{cumFitFly{cumFit: 1.0, fly: &Fly{FlyNumber: 1}}}, index: 0.5, want: 1},
		{flies: []cumFitFly{cumFitFly{cumFit: 1.0, fly: &Fly{FlyNumber: 1}}}, index: 0.0, want: 1},
		{flies: []cumFitFly{cumFitFly{cumFit: 0.1, fly: &Fly{FlyNumber: 1}}, cumFitFly{cumFit: 1.0, fly: &Fly{FlyNumber: 2}}}, index: 0.01, want: 1},
		{flies: []cumFitFly{cumFitFly{cumFit: 0.1, fly: &Fly{FlyNumber: 1}}, cumFitFly{cumFit: 1.0, fly: &Fly{FlyNumber: 2}}}, index: 0.099, want: 1},
		{flies: []cumFitFly{cumFitFly{cumFit: 0.1, fly: &Fly{FlyNumber: 1}}, cumFitFly{cumFit: 1.0, fly: &Fly{FlyNumber: 2}}}, index: 0.1, want: 2},
		{flies: []cumFitFly{cumFitFly{cumFit: 0.1, fly: &Fly{FlyNumber: 1}}, cumFitFly{cumFit: 1.0, fly: &Fly{FlyNumber: 2}}}, index: 0.99, want: 2},
		{flies: []cumFitFly{cumFitFly{cumFit: 0.25, fly: &Fly{FlyNumber: 1}}, cumFitFly{cumFit: 0.5, fly: &Fly{FlyNumber: 2}}, cumFitFly{cumFit: 0.75, fly: &Fly{FlyNumber: 3}}, cumFitFly{cumFit: 1.0, fly: &Fly{FlyNumber: 4}}}, index: 0.24, want: 1},
		{flies: []cumFitFly{cumFitFly{cumFit: 0.25, fly: &Fly{FlyNumber: 1}}, cumFitFly{cumFit: 0.5, fly: &Fly{FlyNumber: 2}}, cumFitFly{cumFit: 0.75, fly: &Fly{FlyNumber: 3}}, cumFitFly{cumFit: 1.0, fly: &Fly{FlyNumber: 4}}}, index: 0.49, want: 2},
		{flies: []cumFitFly{cumFitFly{cumFit: 0.25, fly: &Fly{FlyNumber: 1}}, cumFitFly{cumFit: 0.5, fly: &Fly{FlyNumber: 2}}, cumFitFly{cumFit: 0.75, fly: &Fly{FlyNumber: 3}}, cumFitFly{cumFit: 1.0, fly: &Fly{FlyNumber: 4}}}, index: 0.74, want: 3},
		{flies: []cumFitFly{cumFitFly{cumFit: 0.25, fly: &Fly{FlyNumber: 1}}, cumFitFly{cumFit: 0.5, fly: &Fly{FlyNumber: 2}}, cumFitFly{cumFit: 0.75, fly: &Fly{FlyNumber: 3}}, cumFitFly{cumFit: 1.0, fly: &Fly{FlyNumber: 4}}}, index: 0.99, want: 4},
	}
	for _, t := range tests {
		cf := getFlyForRandomNumber(t.flies, t.index)

		if cf.fly.FlyNumber != t.want {
			test.Errorf("Incorrect fly returned")

		}

	}
}

/*
Tests for pointer bug in generateCumFitness;
generateCumFitness was returning a slice of the same fly repeated all over;
shity range pointer problem
cumFit
*/
func TestGetFlyForRandomNumberLargePop(test *testing.T) {
	fems := make([]Fly, 0, 100)
	for i := 0; i < 100; i++ {

		fems = append(fems, *NewFly([]int64{}, []int64{}))
	}
	var tests = []struct {
		index float64
		want  int64
	}{
		{index: 0.0, want: 1},
		{index: 0.01, want: 2},
		{index: 0.015, want: 2},
		{index: 0.025, want: 3},
		{index: 0.495, want: 50},
		{index: 0.505, want: 51},
		{index: 0.995, want: 100},
	}
	cumfems := generateCumFitness(fems)

	for _, t := range tests {
		cf := getFlyForRandomNumber(cumfems, t.index)

		if cf.fly.FlyNumber != t.want {
			test.Errorf("Incorrect fly returned; wanted %d, got %d", t.want, cf.fly.FlyNumber)
		}

	}

}

/*
With equal fitness, all flies (males and females) should participate in a similar number of matings;
ie around 200 matings in the following scenario
*/
func TestStochasticGetMatePairs(test *testing.T) {
	util.SetSeed(8)
	flies := make([]Fly, 0, 100)
	for i := 0; i < 100; i++ {

		flies = append(flies, *NewFly([]int64{}, []int64{}))
	}

	matep := getMatePairs(flies, 10000)
	var flycounter = make(map[int64]int64)
	for _, mp := range matep {
		flycounter[mp.mate1.FlyNumber]++
		flycounter[mp.mate2.FlyNumber]++
	}

	for _, val := range flycounter {
		if val < 160 || val > 240 {
			test.Errorf("Problematic number of matings %d", val)
		}
	}
}
