package fly

import (
	"invade/env"
	"invade/util"
	"math"
	"testing"
)

/*
Check if candidates contains totest
*/
func contains(totest int64, candidates []*Fly) bool {
	for _, c := range candidates {
		if totest == c.FlyNumber {
			return true
		}
	}
	return false

}

/* helper get single fly
 */
func hGetstandardFly() *Fly {
	return NewFly([]int64{}, []int64{}, false)
}

func hGetStandardPopulation(size int64) [][]*Fly {
	env.SetupEnvironment(size, size, false, false, []int64{}, []float64{}, []int64{}, 0.0, 40, "dros", 0.0, 0.1, 1000)
	tr := make([][]*Fly, size)
	for i := int64(0); i < size; i++ {
		tr[i] = make([]*Fly, 0)
		for k := int64(0); k < size; k++ {
			tr[i] = append(tr[i], hGetstandardFly())
		}
	}
	return tr
}

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

func TestFitness(test *testing.T) {
	var tests = []struct {
		hap1    []int64
		hap2    []int64
		s       float64
		h       float64
		exp_fit float64
	}{
		{hap1: []int64{}, hap2: []int64{}, s: 1.0, h: 0.0, exp_fit: 1.0}, // no insertions
		{hap1: []int64{}, hap2: []int64{}, s: 1.0, h: 0.5, exp_fit: 1.0}, // no insertions
		{hap1: []int64{}, hap2: []int64{}, s: 1.0, h: 1.0, exp_fit: 1.0}, // no insertions

		{hap1: []int64{1}, hap2: []int64{1}, s: 0.0, h: 0.0, exp_fit: 1.0}, // a single homoz insertion; no effect
		{hap1: []int64{1}, hap2: []int64{1}, s: 0.0, h: 0.5, exp_fit: 1.0}, // a single homoz insertion; no effect
		{hap1: []int64{1}, hap2: []int64{1}, s: 0.0, h: 1.0, exp_fit: 1.0}, // a single homoz insertion; no effect

		{hap1: []int64{1}, hap2: []int64{1}, s: 0.1, h: 0.5, exp_fit: 0.9}, // a single homoz insertion; minor effect
		{hap1: []int64{1}, hap2: []int64{1}, s: 0.1, h: 0.0, exp_fit: 0.9}, // a single homoz insertion; minor effect
		{hap1: []int64{1}, hap2: []int64{1}, s: 0.1, h: 1.0, exp_fit: 0.9}, // a single homoz insertion; minor effect

		{hap1: []int64{1}, hap2: []int64{}, s: 0.1, h: 0.5, exp_fit: 0.95}, // a single het insertion; addititive
		{hap1: []int64{1}, hap2: []int64{}, s: 0.1, h: 0.0, exp_fit: 1.0},  // a single het insertion; recessive
		{hap1: []int64{1}, hap2: []int64{}, s: 0.1, h: 1.0, exp_fit: 0.9},  // a single het insertion; dominant

		{hap1: []int64{1}, hap2: []int64{}, s: 0.2, h: 0.5, exp_fit: 0.9}, // a single het insertion; addititive
		{hap1: []int64{1}, hap2: []int64{}, s: 0.2, h: 0.0, exp_fit: 1.0}, // a single het insertion; recessive
		{hap1: []int64{1}, hap2: []int64{}, s: 0.2, h: 1.0, exp_fit: 0.8}, // a single het insertion; dominant

		{hap1: []int64{1}, hap2: []int64{2}, s: 0.1, h: 0.5, exp_fit: 0.9025}, // two  het insertion; addititive
		{hap1: []int64{1}, hap2: []int64{2}, s: 0.1, h: 0.0, exp_fit: 1.0},    // two  het insertion; recessive
		{hap1: []int64{1}, hap2: []int64{2}, s: 0.1, h: 1.0, exp_fit: 0.81},   // two  het insertion; dominant

		{hap1: []int64{1}, hap2: []int64{1, 2}, s: 0.1, h: 0.5, exp_fit: 0.855}, // het + homo
		{hap1: []int64{1}, hap2: []int64{1, 2}, s: 0.1, h: 0.0, exp_fit: 0.9},   // het + homo
		{hap1: []int64{1}, hap2: []int64{1, 2}, s: 0.1, h: 1.0, exp_fit: 0.81},  // het + homo

		{hap1: []int64{1, 4}, hap2: []int64{2, 3}, s: 0.1, h: 0.5, exp_fit: 0.8145062}, // four
		{hap1: []int64{1, 4}, hap2: []int64{2, 3}, s: 0.1, h: 0.0, exp_fit: 1.0},       // four
		{hap1: []int64{1, 4}, hap2: []int64{2, 3}, s: 0.1, h: 1.0, exp_fit: 0.6561},    // four

	}

	for _, t := range tests {
		f := Fly{
			Hap1: t.hap1,
			Hap2: t.hap2}
		SetupFitness(t.s, t.h)
		got_fit := GetFitness(&f)

		if math.Abs(got_fit-t.exp_fit) > 0.0001 {
			test.Errorf("Incorrect fitness; got %f, want %f", got_fit, t.exp_fit)
		}
	}
}

/*
func TestFitnessOmxnMultiplicative(t *testing.T) {
	var tests = []struct {
		x    float64
		ct   int64 // count total
		want float64
	}{
		{x: 0.1, ct: 2, want: 0.81},
		{x: 0.1, ct: 10, want: 0.3486784},
		{x: 0.1, ct: 100, want: 0.0000265},
	}

	for _, test := range tests {
		iff = FitnessFunctionMultiplicative{x: test.x}
		want := test.want
		got := iff.ComputeFitness(test.ct)
		if math.Abs(want-got) > 0.0001 {
			t.Errorf("ff.ComputeFitness(%d) != %f; got = %f", test.ct, test.want, got)
		}
	}
}
*/

/*
test
fitness function w=1-xn^t
cluster insertions (+reference insertions) may be considered

func TestFitnessOmxntLinear(t *testing.T) {
	var tests = []struct {
		x  float64
		ct int64 // count total

		t    float64
		want float64
	}{
		{x: 0.1, ct: 2, t: 1.0, want: 0.8},
		{x: 0.1, ct: 10, t: 1.0, want: 0.0},
		{x: 0.1, ct: 100, t: 1.0, want: 0.0}, // ceck whether min w is 0.0
		{x: 0.1, ct: 2, t: 1.5, want: 0.7171573},
	}

	for _, test := range tests {
		iff = FitnessFunctionLinear{x: test.x, t: test.t}
		want := test.want
		got := iff.ComputeFitness(test.ct)
		if math.Abs(want-got) > 0.0001 {
			t.Errorf("ff.ComputeFitness(%d) != %f; got = %f", test.ct, test.want, got)
		}
	}
}
*/

func TestGenerateCumFitness(test *testing.T) {
	var tests = []struct {
		flies []*Fly
		want  []float64
	}{
		{flies: []*Fly{&Fly{Fitness: 1.0}}, want: []float64{1.0}},
		{flies: []*Fly{&Fly{Fitness: 10.0}}, want: []float64{1.0}},
		{flies: []*Fly{&Fly{Fitness: 10.0}, &Fly{Fitness: 10.0}}, want: []float64{0.5, 1.0}},
		{flies: []*Fly{&Fly{Fitness: 5.0}, &Fly{Fitness: 15.0}}, want: []float64{0.75, 1.0}},
		{flies: []*Fly{&Fly{Fitness: 15.0}, &Fly{Fitness: 5.0}}, want: []float64{0.75, 1.0}},
		{flies: []*Fly{&Fly{Fitness: 1.0}, &Fly{Fitness: 2.0}, &Fly{Fitness: 3.0}, &Fly{Fitness: 4.0}}, want: []float64{0.4, 0.7, 0.9, 1.0}},
		{flies: []*Fly{&Fly{Fitness: 4.0}, &Fly{Fitness: 2.0}, &Fly{Fitness: 1.0}, &Fly{Fitness: 3.0}}, want: []float64{0.4, 0.7, 0.9, 1.0}},
		{flies: []*Fly{&Fly{Fitness: 0.4}, &Fly{Fitness: 0.2}, &Fly{Fitness: 0.1}, &Fly{Fitness: 0.3}}, want: []float64{0.4, 0.7, 0.9, 1.0}},
		{flies: []*Fly{&Fly{Fitness: 0.04}, &Fly{Fitness: 0.02}, &Fly{Fitness: 0.01}, &Fly{Fitness: 0.03}}, want: []float64{0.4, 0.7, 0.9, 1.0}},
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

	fems := make([]*Fly, 0, 100)
	iff = FitnessFunctionMultiplicative{}
	for i := 0; i < 100; i++ {

		fems = append(fems, NewFly([]int64{}, []int64{}, false))
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

func TestGetNeighborhoodOpen(test *testing.T) {
	env.SetupEnvironment(100, 100, false, false, []int64{}, []float64{}, []int64{}, 0.0, 40, "dros", 0.0, 0.1, 1000)
	var tests = []struct {
		x      int64
		y      int64
		radius int64
		w_minx int64
		w_maxx int64
		w_miny int64
		w_maxy int64
		w_size int64
	}{
		// []int64{1, 2, 3, 5, 6, 7, 9, 10, 11}
		{x: 0, y: 0, radius: 2, w_minx: 0, w_miny: 0, w_maxx: 2, w_maxy: 2, w_size: 9},
		{x: 10, y: 10, radius: 2, w_minx: 8, w_maxx: 12, w_miny: 8, w_maxy: 12, w_size: 25},
		{x: 1, y: 1, radius: 1, w_minx: 0, w_miny: 0, w_maxx: 2, w_maxy: 2, w_size: 9},
		{x: 10, y: 10, radius: 1, w_minx: 9, w_miny: 9, w_maxx: 11, w_maxy: 11, w_size: 9},
		{x: 10, y: 10, radius: 5, w_minx: 5, w_miny: 5, w_maxx: 15, w_maxy: 15, w_size: 121},
		{x: 2, y: 0, radius: 1, w_minx: 1, w_miny: 0, w_maxx: 3, w_maxy: 1, w_size: 6},
		{x: 0, y: 2, radius: 1, w_minx: 0, w_miny: 1, w_maxx: 1, w_maxy: 3, w_size: 6},
		{x: 0, y: 99, radius: 2, w_minx: 0, w_miny: 97, w_maxx: 2, w_maxy: 99, w_size: 9},
		{x: 99, y: 0, radius: 2, w_minx: 97, w_miny: 0, w_maxx: 99, w_maxy: 2, w_size: 9},
		{x: 97, y: 0, radius: 1, w_minx: 96, w_miny: 0, w_maxx: 98, w_maxy: 1, w_size: 6},
		{x: 0, y: 97, radius: 1, w_minx: 0, w_miny: 96, w_maxx: 1, w_maxy: 98, w_size: 6},
		{x: 99, y: 99, radius: 2, w_minx: 97, w_miny: 97, w_maxx: 99, w_maxy: 99, w_size: 9},
		// NOTHING
	}
	for _, t := range tests {
		got := env.GetNeighborhoodCoordinates(t.y, t.x, t.radius)
		gminx, gmaxx, gminy, gmaxy := t.x, int64(0), t.y, int64(0)
		for _, c := range got {
			if c.X < gminx {
				gminx = c.X
			}
			if c.X > gmaxx {
				gmaxx = c.X
			}
			if c.Y < gminy {
				gminy = c.Y
			}
			if c.Y > gmaxy {
				gmaxy = c.Y
			}
		}

		if gminx != t.w_minx {
			test.Errorf("Incorrect x min; got %d, want %d", gminx, t.w_minx)
		}
		if gmaxx != t.w_maxx {
			test.Errorf("Incorrect x max; got %d, want %d", gmaxx, t.w_maxx)
		}
		if gminy != t.w_miny {
			test.Errorf("Incorrect y min; got %d, want %d", gminy, t.w_miny)
		}
		if gmaxy != t.w_maxy {
			test.Errorf("Incorrect y max; got %d, want %d", gmaxy, t.w_maxy)
		}
		if len(got) != int(t.w_size) {
			test.Errorf("Incorrect size; got %d, want %d", len(got), t.w_size)
		}
	}
}

func TestGetNeighborhoodClosedX(test *testing.T) {
	env.SetupEnvironment(100, 100, true, false, []int64{}, []float64{}, []int64{}, 0.0, 40, "dros", 0.0, 0.1, 1000)
	var tests = []struct {
		x      int64
		y      int64
		radius int64
		w_minx int64
		w_maxx int64
		w_miny int64
		w_maxy int64
		w_size int64
	}{
		// []int64{1, 2, 3, 5, 6, 7, 9, 10, 11}
		{x: 0, y: 0, radius: 2, w_minx: 0, w_miny: 0, w_maxx: 99, w_maxy: 2, w_size: 15},
		{x: 10, y: 10, radius: 2, w_minx: 8, w_maxx: 12, w_miny: 8, w_maxy: 12, w_size: 25},
		{x: 1, y: 1, radius: 1, w_minx: 0, w_miny: 0, w_maxx: 2, w_maxy: 2, w_size: 9},
		{x: 10, y: 10, radius: 1, w_minx: 9, w_miny: 9, w_maxx: 11, w_maxy: 11, w_size: 9},
		{x: 10, y: 10, radius: 5, w_minx: 5, w_miny: 5, w_maxx: 15, w_maxy: 15, w_size: 121},
		{x: 2, y: 0, radius: 1, w_minx: 1, w_miny: 0, w_maxx: 3, w_maxy: 1, w_size: 6},
		{x: 0, y: 2, radius: 1, w_minx: 0, w_miny: 1, w_maxx: 99, w_maxy: 3, w_size: 9},
		{x: 0, y: 99, radius: 2, w_minx: 0, w_miny: 97, w_maxx: 99, w_maxy: 99, w_size: 15},
		{x: 99, y: 0, radius: 2, w_minx: 0, w_miny: 0, w_maxx: 99, w_maxy: 2, w_size: 15},
		{x: 97, y: 0, radius: 1, w_minx: 96, w_miny: 0, w_maxx: 98, w_maxy: 1, w_size: 6},
		{x: 0, y: 97, radius: 1, w_minx: 0, w_miny: 96, w_maxx: 99, w_maxy: 98, w_size: 9},
		{x: 99, y: 99, radius: 2, w_minx: 0, w_miny: 97, w_maxx: 99, w_maxy: 99, w_size: 15},
		// NOTHING
	}
	for _, t := range tests {
		got := env.GetNeighborhoodCoordinates(t.y, t.x, t.radius)
		gminx, gmaxx, gminy, gmaxy := t.x, int64(0), t.y, int64(0)
		for _, c := range got {
			if c.X < gminx {
				gminx = c.X
			}
			if c.X > gmaxx {
				gmaxx = c.X
			}
			if c.Y < gminy {
				gminy = c.Y
			}
			if c.Y > gmaxy {
				gmaxy = c.Y
			}
		}

		if gminx != t.w_minx {
			test.Errorf("Incorrect x min; got %d, want %d", gminx, t.w_minx)
		}
		if gmaxx != t.w_maxx {
			test.Errorf("Incorrect x max; got %d, want %d", gmaxx, t.w_maxx)
		}
		if gminy != t.w_miny {
			test.Errorf("Incorrect y min; got %d, want %d", gminy, t.w_miny)
		}
		if gmaxy != t.w_maxy {
			test.Errorf("Incorrect y max; got %d, want %d", gmaxy, t.w_maxy)
		}
		if len(got) != int(t.w_size) {
			test.Errorf("Incorrect size; got %d, want %d", len(got), t.w_size)
		}
	}
}

func TestGetNeighborhoodClosedY(test *testing.T) {
	env.SetupEnvironment(100, 100, false, true, []int64{}, []float64{}, []int64{}, 0.0, 40, "dros", 0.0, 0.1, 1000)
	var tests = []struct {
		x      int64
		y      int64
		radius int64
		w_minx int64
		w_maxx int64
		w_miny int64
		w_maxy int64
		w_size int64
	}{
		// []int64{1, 2, 3, 5, 6, 7, 9, 10, 11}
		{x: 0, y: 0, radius: 2, w_minx: 0, w_miny: 0, w_maxx: 2, w_maxy: 99, w_size: 15},
		{x: 10, y: 10, radius: 2, w_minx: 8, w_maxx: 12, w_miny: 8, w_maxy: 12, w_size: 25},
		{x: 1, y: 1, radius: 1, w_minx: 0, w_miny: 0, w_maxx: 2, w_maxy: 2, w_size: 9},
		{x: 10, y: 10, radius: 1, w_minx: 9, w_miny: 9, w_maxx: 11, w_maxy: 11, w_size: 9},
		{x: 10, y: 10, radius: 5, w_minx: 5, w_miny: 5, w_maxx: 15, w_maxy: 15, w_size: 121},
		{x: 2, y: 0, radius: 1, w_minx: 1, w_miny: 0, w_maxx: 3, w_maxy: 99, w_size: 9},
		{x: 0, y: 2, radius: 1, w_minx: 0, w_miny: 1, w_maxx: 1, w_maxy: 3, w_size: 6},
		{x: 0, y: 99, radius: 2, w_minx: 0, w_miny: 0, w_maxx: 2, w_maxy: 99, w_size: 15},
		{x: 99, y: 0, radius: 2, w_minx: 97, w_miny: 0, w_maxx: 99, w_maxy: 99, w_size: 15},
		{x: 97, y: 0, radius: 1, w_minx: 96, w_miny: 0, w_maxx: 98, w_maxy: 99, w_size: 9},
		{x: 0, y: 97, radius: 1, w_minx: 0, w_miny: 96, w_maxx: 1, w_maxy: 98, w_size: 6},
		{x: 99, y: 99, radius: 2, w_minx: 97, w_miny: 0, w_maxx: 99, w_maxy: 99, w_size: 15},
		// NOTHING
	}
	for _, t := range tests {
		got := env.GetNeighborhoodCoordinates(t.y, t.x, t.radius)
		gminx, gmaxx, gminy, gmaxy := t.x, int64(0), t.y, int64(0)
		for _, c := range got {
			if c.X < gminx {
				gminx = c.X
			}
			if c.X > gmaxx {
				gmaxx = c.X
			}
			if c.Y < gminy {
				gminy = c.Y
			}
			if c.Y > gmaxy {
				gmaxy = c.Y
			}
		}

		if gminx != t.w_minx {
			test.Errorf("Incorrect x min; got %d, want %d", gminx, t.w_minx)
		}
		if gmaxx != t.w_maxx {
			test.Errorf("Incorrect x max; got %d, want %d", gmaxx, t.w_maxx)
		}
		if gminy != t.w_miny {
			test.Errorf("Incorrect y min; got %d, want %d", gminy, t.w_miny)
		}
		if gmaxy != t.w_maxy {
			test.Errorf("Incorrect y max; got %d, want %d", gmaxy, t.w_maxy)
		}
		if len(got) != int(t.w_size) {
			test.Errorf("Incorrect size; got %d, want %d", len(got), t.w_size)
		}
	}
}

func TestGetNeighborhoodClosedXY(test *testing.T) {
	env.SetupEnvironment(100, 100, true, true, []int64{}, []float64{}, []int64{}, 0.0, 40, "dros", 0.0, 0.1, 1000)
	var tests = []struct {
		x      int64
		y      int64
		radius int64
		w_minx int64
		w_maxx int64
		w_miny int64
		w_maxy int64
		w_size int64
	}{
		// []int64{1, 2, 3, 5, 6, 7, 9, 10, 11}
		{x: 0, y: 0, radius: 2, w_minx: 0, w_miny: 0, w_maxx: 99, w_maxy: 99, w_size: 25},
		{x: 10, y: 10, radius: 2, w_minx: 8, w_maxx: 12, w_miny: 8, w_maxy: 12, w_size: 25},
		{x: 1, y: 1, radius: 1, w_minx: 0, w_miny: 0, w_maxx: 2, w_maxy: 2, w_size: 9},
		{x: 10, y: 10, radius: 1, w_minx: 9, w_miny: 9, w_maxx: 11, w_maxy: 11, w_size: 9},
		{x: 10, y: 10, radius: 5, w_minx: 5, w_miny: 5, w_maxx: 15, w_maxy: 15, w_size: 121},
		{x: 2, y: 0, radius: 1, w_minx: 1, w_miny: 0, w_maxx: 3, w_maxy: 99, w_size: 9},
		{x: 0, y: 2, radius: 1, w_minx: 0, w_miny: 1, w_maxx: 99, w_maxy: 3, w_size: 9},
		{x: 0, y: 99, radius: 2, w_minx: 0, w_miny: 0, w_maxx: 99, w_maxy: 99, w_size: 25},
		{x: 99, y: 0, radius: 2, w_minx: 0, w_miny: 0, w_maxx: 99, w_maxy: 99, w_size: 25},
		{x: 97, y: 0, radius: 1, w_minx: 96, w_miny: 0, w_maxx: 98, w_maxy: 99, w_size: 9},
		{x: 0, y: 97, radius: 1, w_minx: 0, w_miny: 96, w_maxx: 99, w_maxy: 98, w_size: 9},
		{x: 99, y: 99, radius: 2, w_minx: 0, w_miny: 0, w_maxx: 99, w_maxy: 99, w_size: 25},
		// NOTHING
	}
	for _, t := range tests {
		got := env.GetNeighborhoodCoordinates(t.y, t.x, t.radius)
		gminx, gmaxx, gminy, gmaxy := t.x, int64(0), t.y, int64(0)
		for _, c := range got {
			if c.X < gminx {
				gminx = c.X
			}
			if c.X > gmaxx {
				gmaxx = c.X
			}
			if c.Y < gminy {
				gminy = c.Y
			}
			if c.Y > gmaxy {
				gmaxy = c.Y
			}
		}

		if gminx != t.w_minx {
			test.Errorf("Incorrect x min; got %d, want %d", gminx, t.w_minx)
		}
		if gmaxx != t.w_maxx {
			test.Errorf("Incorrect x max; got %d, want %d", gmaxx, t.w_maxx)
		}
		if gminy != t.w_miny {
			test.Errorf("Incorrect y min; got %d, want %d", gminy, t.w_miny)
		}
		if gmaxy != t.w_maxy {
			test.Errorf("Incorrect y max; got %d, want %d", gmaxy, t.w_maxy)
		}
		if len(got) != int(t.w_size) {
			test.Errorf("Incorrect size; got %d, want %d", len(got), t.w_size)
		}
	}
}

func TestGetSilencingStatus(test *testing.T) {
	env.SetupEnvironment(10, 10, false, false, []int64{}, []float64{}, []int64{}, 0.0, 40, "dros", 0.0, 0.1, 1000)
	var tests = []struct {
		fs       FlyStatistic
		silenced bool
		want     bool
	}{
		{fs: FlyStatistic{}, silenced: false, want: false},               // sanity
		{fs: FlyStatistic{CountTotal: 39}, silenced: false, want: false}, // NOTHING
		{fs: FlyStatistic{CountTotal: 41}, silenced: false, want: true},  //GAIN
		{fs: FlyStatistic{CountTotal: 1}, silenced: true, want: true},    // RETAIN
		{fs: FlyStatistic{CountTotal: 0}, silenced: true, want: false},   // LOSS
	}

	for _, t := range tests {
		got := getSilencingStatus(t.fs.CountTotal, t.silenced, 1)

		if got != t.want {
			test.Errorf("Incorrect getSilencingStatus(); got %v, want %v", got, t.want)

		}

	}
}

func TestStochasticGetRandomSex(test *testing.T) {
	util.SetSeed(6)
	cmale := 0
	cfem := 0
	for i := 0; i < 10000; i++ {
		sex := GetRandomSex()
		if sex == MALE {
			cmale++
		} else if sex == FEMALE {
			cfem++
		} else {
			panic("unknown sex")
		}
	}
	// should be around 5000 for both sexes
	if cmale < 4900 || cmale > 5900 {
		test.Errorf("Problematic number of males %d", cmale)
	}
	if cfem < 4900 || cfem > 5900 {
		test.Errorf("Problematic number of males %d", cfem)
	}

}

func TestGetNeighbors(test *testing.T) {

	FLYCOUNTER = 1
	iff = FitnessFunctionMultiplicative{}
	//  1  2  3  4
	//  5  6  7  8
	//  9 10 11 12
	// 13 14 15 16
	p := hGetStandardPopulation(4)
	var tests = []struct {
		x    int64
		y    int64
		r    int64
		want []int64
	}{
		{x: 0, y: 0, r: 1, want: []int64{1, 2, 5, 6}},
		{x: 3, y: 0, r: 1, want: []int64{3, 4, 7, 8}},
		{x: 0, y: 3, r: 1, want: []int64{9, 10, 13, 14}},
		{x: 3, y: 3, r: 1, want: []int64{11, 12, 15, 16}},
		{x: 1, y: 1, r: 1, want: []int64{1, 2, 3, 5, 6, 7, 9, 10, 11}},
	}

	for _, t := range tests {
		got := getNeighbors(p, t.y, t.x, t.r)
		for _, f := range t.want {
			if !contains(f, got) {
				test.Errorf("getNeighbors(%d,%d,%d) does not contain fly %d but should %v", t.y, t.x, t.r, f, t.want)
			}
			if len(got) != len(t.want) {
				test.Errorf("getNeighbors(%d,%d,%d)", t.y, t.x, t.r)
			}
		}

	}
}
func TestStochasticMaterNeighborhood(test *testing.T) {

	//  1  2  3  4  5
	//  6  7  8  9 10
	// 11 12 13 14 15
	counter := make([]int64, 26)
	expect := map[int64]bool{1: true, 2: true, 3: true, 6: true, 7: true, 8: true, 11: true, 12: true, 13: true}
	FLYCOUNTER = 1
	iff = FitnessFunctionMultiplicative{}
	m := MaterNeighborhood{selfingRate: 0.0, neighborDistance: 1}
	p := hGetStandardPopulation(5)
	for i := 0; i < 4500; i++ {
		mp := m.getMatePairCoord(p, 1, 1)
		counter[mp.female.FlyNumber]++
		counter[mp.male.FlyNumber]++
	}

	for fid, counts := range counter {
		_, pres := expect[int64(fid)]
		if pres {

			if counts > 1100 || counts < 900 {
				test.Errorf("Error in stochasit random neighborhood, expected 1000; got %d for fly %d", counts, fid)
			}
		} else {
			if counts > 0 {
				test.Errorf("Error in stochasit random neighborhood, expected 0; got %d for fly %d", counts, fid)
			}
		}
	}

}

func TestStochasticMaterNeighborhoodSelfing(test *testing.T) {

	FLYCOUNTER = 1
	iff = FitnessFunctionMultiplicative{}
	m := MaterNeighborhood{selfingRate: 0.95, neighborDistance: 1}
	p := hGetStandardPopulation(10)
	self, nonself := 0, 0
	for i := 0; i < 100; i++ {
		mpss := m.GetMatePairs(p)
		for _, mps := range mpss {
			for _, mp := range mps {
				if mp.female.FlyNumber == mp.male.FlyNumber {
					self++
				} else {
					nonself++
				}

			}
		}

	}

	if self > 9700 || self < 9300 {
		test.Errorf("Error in selfing of stochastic random neighborhood, expected around 9500; got %d", self)
	}

}

func TestStochasticMaterPanmictic(test *testing.T) {

	//  1  2  3  4  5
	//  6  7  8  9 10
	// 11 12 13 14 15
	counter := make([]int64, 25)
	FLYCOUNTER = 1
	iff = FitnessFunctionMultiplicative{}
	m := MaterPanmictic{selfingRate: 0.0}
	p := hGetStandardPopulation(5)
	for i := 0; i < 500; i++ {
		mpss := m.GetMatePairs(p)
		for _, mps := range mpss {
			for _, mp := range mps {
				counter[mp.female.FlyNumber-1]++
				counter[mp.male.FlyNumber-1]++
			}
		}

	}

	for fid, counts := range counter {

		if counts > 1100 || counts < 900 {
			test.Errorf("Error in stochasit random neighborhood, expected 1000; got %d for fly %d", counts, fid+1)
		}
	}

}

func TestStochasticMaterPanmicticSelfing(test *testing.T) {

	FLYCOUNTER = 1
	iff = FitnessFunctionMultiplicative{}
	m := MaterPanmictic{selfingRate: 0.95}
	p := hGetStandardPopulation(10)
	self, nonself := 0, 0
	for i := 0; i < 100; i++ {
		mpss := m.GetMatePairs(p)
		for _, mps := range mpss {
			for _, mp := range mps {
				if mp.female.FlyNumber == mp.male.FlyNumber {
					self++
				} else {
					nonself++
				}

			}
		}

	}

	if self > 9700 || self < 9300 {
		test.Errorf("Error in selfing of stochastic random neighborhood, expected around 9500; got %d", self)
	}

}

func TestStochasticLossOfSilencing(test *testing.T) {
	env.SetupEnvironment(100, 100, false, false, []int64{}, []float64{}, []int64{10, 20}, 0.2, 40, "dros", 0.1, 0.0, 1000)
	SetupFitness(0.0, 0.0)
	FLYCOUNTER = 1
	lost := 0
	for i := 0; i < 10000; i++ {
		nf := env.OffspringIsSilenced(true, true)
		if nf == false {
			lost++
		}

	}

	if lost < 800 || lost > 1200 {
		test.Errorf("Error in loss of epigenetic silencing; got %d", lost)
	}

}

//getNeighbors(flies [][]*Fly, ycord int64, xcord int64, radius int64) []*Fly
