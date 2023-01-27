package env

// command line, run all tests "go test ./..." yes three points
import (
	"fmt"
	"invade/util"
	"math"
	"testing"
)

func TestGenomicInterval(t *testing.T) {

	gi := GenomicInterval{Start: 1, End: 100}
	if gi.Start != 1 {
		t.Error("Start wrong")
	}
	if gi.End != 100 {
		t.Error("End wrong")
	}
}
func TestGenomicLandscape(t *testing.T) {
	cl := newGenomicLandscape([]int64{100, 200, 300, 400})
	if len(cl.intervals) != 4 || len(cl.chrmSizes) != 4 || len(cl.offsets) != 4 {
		t.Error("incorrect number of chromosomes")
	}
	if cl.chrmSizes[0] != 100 || cl.chrmSizes[1] != 200 || cl.chrmSizes[2] != 300 || cl.chrmSizes[3] != 400 {
		t.Error("incorrect size of chromosome ")
	}
	if cl.offsets[0] != 0 || cl.offsets[1] != 100 || cl.offsets[2] != 300 || cl.offsets[3] != 600 {
		t.Error("incorrect offset ")
	}
	if cl.intervals[0].Start != 0 || cl.intervals[0].End != 99 {
		t.Error("incorrect first interval")
	}
	if cl.intervals[1].Start != 100 || cl.intervals[1].End != 299 {
		t.Error("incorrect second interval")
	}
	if cl.intervals[2].Start != 300 || cl.intervals[2].End != 599 {
		t.Error("incorrect third interval")
	}
	if cl.intervals[3].Start != 600 || cl.intervals[3].End != 999 {
		t.Error("incorrect fourth interval")
	}
}

func TestByte(t *testing.T) { // TODO test non cluster
	bi := byte(0)
	bii2 := float64(bi)/100.0 - 1.0
	bii := (float64(bi) - 100.0) / 100.0
	fmt.Print(bii, bii2)
}

func TestNewClusterNonCluster(t *testing.T) { // TODO test non cluster
	gl := newGenomicLandscape([]int64{100, 200, 300, 400})
	cl, ncl := newClusterNonClusters([]int64{10, 20, 30, 40}, gl)
	if cl.Count() != 4 {
		t.Error("incorrect number of cluster")
	}
	if cl.Intervals[0].Start != 0 || cl.Intervals[0].End != 9 {
		t.Error("incorrect first cluster")
	}
	if cl.Intervals[1].Start != 100 || cl.Intervals[1].End != 119 {
		t.Error("incorrect second cluster")
	}
	if cl.Intervals[2].Start != 300 || cl.Intervals[2].End != 329 {
		t.Error("incorrect third cluster")
	}
	if cl.Intervals[3].Start != 600 || cl.Intervals[3].End != 639 {
		t.Error("incorrect third cluster")
	}

	if ncl.Count() != 4 {
		t.Error("incorrect number of cluster")
	}
	if ncl.Intervals[0].Start != 10 || ncl.Intervals[0].End != 99 {
		t.Error("incorrect first cluster")
	}
	if ncl.Intervals[1].Start != 120 || ncl.Intervals[1].End != 299 {
		t.Error("incorrect first cluster")
	}
	if ncl.Intervals[2].Start != 330 || ncl.Intervals[2].End != 599 {
		t.Error("incorrect first cluster")
	}
	if ncl.Intervals[3].Start != 640 || ncl.Intervals[3].End != 999 {
		t.Error("incorrect first cluster")
	}
}
func TestIsClusterInsertion(t *testing.T) {
	// PRIME example on how tests should be implemented in Go, according to Kerninghan
	gl := newGenomicLandscape([]int64{100, 200, 300, 400})
	cl, _ := newClusterNonClusters([]int64{10, 20, 30, 40}, gl)
	env = Environment{genome: gl,
		clusters: *cl}

	var tests = []struct {
		position int64
		want     bool
	}{
		{10, false},
		{99, false},
		{120, false},
		{299, false},
		{330, false},
		{599, false},
		{640, false},
		{0, true},
		{9, true},
		{100, true},
		{119, true},
		{300, true},
		{329, true},
		{600, true},
		{639, true}}

	// very concise testing syntax with clear message
	for _, test := range tests {
		if IsClusterInsertion(test.position) != test.want {
			t.Errorf("IsClusterInsertion(%d)!=%t", test.position, test.want)
		}
	}

}

func TestPositionTranslationCluster(t *testing.T) {
	// PRIME example on how tests should be implemented in Go, according to Kerninghan
	gl := newGenomicLandscape([]int64{10, 10, 10})
	cl, ncl := newClusterNonClusters([]int64{5, 5, 5}, gl)
	env = Environment{genome: gl,
		clusters:    *cl,
		nonClusters: *ncl}
	var tests = []struct {
		index int64
		want  int64
	}{
		{0, 0},
		{4, 4},
		{5, 10},
		{9, 14},
		{10, 20},
		{14, 24}}

	for _, test := range tests {
		if cl.Positions[test.index] != test.want {
			t.Errorf("Position translation of clusters screwed up (%d)!=%d", cl.Positions[test.index], test.want)
		}
	}
}

func TestPositionTranslationNonCluster(t *testing.T) {
	// PRIME example on how tests should be implemented in Go, according to Kerninghan
	gl := newGenomicLandscape([]int64{10, 10, 10})
	cl, ncl := newClusterNonClusters([]int64{5, 5, 5}, gl)
	env = Environment{genome: gl,
		clusters:    *cl,
		nonClusters: *ncl}
	var tests = []struct {
		index int64
		want  int64
	}{
		{0, 5},
		{4, 9},
		{5, 15},
		{9, 19},
		{10, 25},
		{14, 29}}

	for _, test := range tests {
		if ncl.Positions[test.index] != test.want {
			t.Errorf("Position translation of non-clusters screwed up (%d)!=%d", cl.Positions[test.index], test.want)
		}
	}
}

func TestTEBiasIncrease(test *testing.T) {
	var tests = []struct {
		start     byte
		wantend   byte
		wantstart byte
	}{
		{4, 5, 4},
		{0, 1, 0},
		{199, 200, 199},
		{200, 200, 200}, // no increase at maximum
	}

	for _, t := range tests {
		ste := TEInsertion(t.start)
		ete := ste.increase()
		if byte(ste) != byte(t.wantstart) {
			test.Errorf("Insertionbias of parent changed (%d)!=%d", ste, t.wantstart)
		}
		if byte(ete) != byte(t.wantend) {
			test.Errorf("Insertionbias of child incorrect (%d)!=%d", ete, t.wantend)
		}
	}
}

// Test if direction of mutation is stochastic, half up and half down
func TestStochasticMutationDirection(test *testing.T) {
	tei := NewTEInsertion(0)
	biasmap := make(map[TEInsertion]int64)

	for i := 0; i < 10000; i++ {
		tem := tei.Mutate()
		biasmap[tem]++
	}
	print(biasmap)
	if biasmap[99] < 4800 || biasmap[99] > 5200 {
		test.Errorf("Non stochastic number of down mutations: %d does not approx %d", 5000, biasmap[99])
	}
	if biasmap[101] < 4800 || biasmap[101] > 5200 {
		test.Errorf("Non stochastic number of up mutations: %d does not approx %d", 5000, biasmap[101])
	}
}

func TestStochasticIntroduceMutations(test *testing.T) {
	gl := newGenomicLandscape([]int64{100, 100})
	env = Environment{genome: gl, mubias: 1.0}
	util.SetSeed(199)
	tei := NewTEInsertion(0)

	pos := make([]int64, 0)
	gamete := make(map[int64]TEInsertion)
	for i := int64(0); i < 50; i++ {
		gamete[i] = tei
		pos = append(pos, i)
	}

	// check if every position is mutated
	gamete = helperIntroduceMutations(gamete, pos)

	biasmap := make(map[TEInsertion]int64)

	for _, te := range gamete {
		biasmap[te]++
	}

	if biasmap[100] != 0 {
		test.Errorf("Incorrect number of non-mutated sites; should be zero, found: %d", biasmap[100])
	}

}

func TestTEBiasDecrease(test *testing.T) {
	var tests = []struct {
		start     byte
		wantend   byte
		wantstart byte
	}{
		{4, 3, 4},
		{0, 0, 0},
		{1, 0, 1},
		{199, 198, 199},
		{200, 199, 200}, // no increase at maximum
	}

	for _, t := range tests {
		ste := TEInsertion(t.start)
		ete := ste.decrease()
		if byte(ste) != byte(t.wantstart) {
			test.Errorf("Insertionbias of parent changed (%d)!=%d", ste, t.wantstart)
		}
		if byte(ete) != byte(t.wantend) {
			test.Errorf("Insertionbias of child incorrect (%d)!=%d", ete, t.wantend)
		}
	}
}

func TestInsertionBias(test *testing.T) {
	var tests = []struct {
		bias  float64
		clusi float64
		want  float64
	}{
		{0.0, 0.03, 0.03},
		{-1.0, 0.03, 0.0},
		{1.0, 0.03, 1.0},
		{0.0, 0.5, 0.5},
		{-1.0, 0.5, 0.0},
		{1.0, 0.5, 1.0},
		{0.1, 0.03, 0.03642384},
		{0.5, 0.03, 0.08490566},
		{-0.5, 0.03, 0.01020408},
		{0.9, 0.03, 0.3701299},
	}

	for _, t := range tests {
		got := getProbForBiasAndClusi(t.bias, t.clusi)
		want := t.want
		if math.Abs(got-want) > 0.000001 {
			test.Errorf("getProbForBiasAndClusi(%f,%f)!=%f", t.bias, t.clusi, t.want)
		}

	}
}

func TestNilClusterReference(t *testing.T) {
	gl := newGenomicLandscape([]int64{100, 100, 100, 100})
	cl, ncl := newClusterNonClusters(nil, gl)
	env = Environment{genome: gl,
		clusters: *cl}
	for i := int64(0); i < gl.totalGenome; i++ {
		if IsClusterInsertion(i) {
			t.Errorf("incorrect cluster insertion in nil cluster, position %d", i)
		}
	}
	if ncl.Count() != 4 {
		t.Error("incorrect number of non cluster regions")
	}
	if ncl.Size() != 400 {
		t.Error("incorrect size of non cluster regions")
	}
	if ncl.Intervals[0].Start != 0 || ncl.Intervals[0].End != 99 {
		t.Error("incorrect first cluster")
	}
	if ncl.Intervals[1].Start != 100 || ncl.Intervals[1].End != 199 {
		t.Error("incorrect first cluster")
	}
	if ncl.Intervals[2].Start != 200 || ncl.Intervals[2].End != 299 {
		t.Error("incorrect first cluster")
	}
	if ncl.Intervals[3].Start != 300 || ncl.Intervals[3].End != 399 {
		t.Error("incorrect first cluster")
	}

}

func TestStochasticGetNovelClusterSites(t *testing.T) {
	util.SetSeed(5)
	genome := newGenomicLandscape([]int64{10, 10, 10})

	cl, ncl := newClusterNonClusters([]int64{5, 5, 5}, genome)
	env = Environment{
		genome:      genome,
		clusters:    *cl,
		nonClusters: *ncl,
	}

	for i := 0; i < 1000; i++ {
		pos := GetRandomClusterSite()
		if !IsClusterInsertion(pos) {
			t.Errorf("genomic site is OUTSIDE of piRNA clusters %d", pos)
		}
	}
}

func TestStochasticGetNovelNonClusterSites(t *testing.T) {
	util.SetSeed(5)
	genome := newGenomicLandscape([]int64{10, 10, 10})

	cl, ncl := newClusterNonClusters([]int64{5, 5, 5}, genome)
	env = Environment{
		genome:      genome,
		clusters:    *cl,
		nonClusters: *ncl,
	}

	for i := 0; i < 1000; i++ {
		pos := GetRandomNonClusterSite()
		if IsClusterInsertion(pos) {
			t.Errorf("genomic site is INSIDE of piRNA clusters %d", pos)
		}
	}
}

/*
func TestSeparateInsertions(t *testing.T) {
	gl := newGenomicLandscape([]int64{100, 100, 100, 100, 100})
	cl, _ := newCluster([]int64{10, 10, 10, 10, 10}, gl)
	env = Environment{genome: gl,
		clusters: cl}
	sites := make([]int64, 500)
	for i := int64(0); i < 500; i++ {
		sites[i] = i
	}
	ret := SeparateInsertions(sites)
	if len(ret.Cluster) != 50 {
		t.Error("incorrect number of cluster sites")
	}
	if len(ret.RefRegion) != 75 {
		t.Error("incorrect number of ref regions sites")
	}
	if len(ret.NOE) != 375 {
		t.Errorf("incorrect number of NOE sites: 375!=%d", len(ret.NOE))
	}
}
*/

func TestGetNovelInsertionCount(t *testing.T) {
	var tests = []struct {
		tot  int64
		totc int64 // number of cluster insertions
		scc  int64
		jump Jumper
		want float64
	}{
		{tot: 10, totc: 0, jump: Jumper{u: 0.1, uc: 0.0}, want: 1.0},
		{tot: 10, totc: 1, jump: Jumper{u: 0.1, uc: 0.0}, want: 0.0},
		{tot: 100, totc: 0, jump: Jumper{u: 0.1, uc: 0.0}, want: 10.0},
		{tot: 100, totc: 1, jump: Jumper{u: 0.1, uc: 0.0}, want: 0.0},
		{tot: 100, totc: 1, jump: Jumper{u: 0.1, uc: 0.01}, want: 1.0},
		// selective cluster insertions test
		{tot: 10, totc: 1, jump: Jumper{u: 0.1, uc: 0.0, sci: true}, want: 1.0},
		{tot: 10, totc: 1, scc: 1, jump: Jumper{u: 0.1, uc: 0.0, sci: true}, want: 0.0},
		{tot: 10, totc: 0, scc: 1, jump: Jumper{u: 0.1, uc: 0.0, sci: true}, want: 0.0},
		{tot: 10, totc: 1, scc: 1, jump: Jumper{u: 0.1, uc: 0.0, sci: false}, want: 0.0},
		{tot: 10, totc: 0, scc: 1, jump: Jumper{u: 0.1, uc: 0.0, sci: false}, want: 1.0},
	}

	for _, test := range tests {
		ju := test.jump
		got := ju.getNovelInsertionCount(test.tot, test.scc, test.totc)
		dif := math.Abs(test.want - got)
		if dif > 0.00001 {
			t.Errorf("getNovelInsertionCount(); got %f wanted %f", got, test.want)
		}

	}
}

func TestCentiMorgan2Lambda(test *testing.T) {
	var tests = []struct {
		winsize int64
		cm      float64
		want    float64
	}{{winsize: 1000000, cm: 4.0, want: 0.0416908},
		{winsize: 100000, cm: 4.0, want: 0.00416908},
		{winsize: 10000000, cm: 4.0, want: 0.416908},
		{winsize: 100000000, cm: 4.0, want: 4.16908},
		{winsize: 1000000, cm: 10.0, want: 0.1115718},
		{winsize: 1000000, cm: 40.0, want: 0.804719},
	}

	for _, t := range tests {
		got := cmpmb2lambda(t.cm, t.winsize)
		if math.Abs(got-t.want) > 0.001 {
			test.Errorf("cmpmb2lambda(%f,%d)!=%f", t.cm, t.winsize, t.want)
		}
	}
}

func TestStochasticRecombinationWindow(test *testing.T) {
	util.SetSeed(5)
	rw := RecombinationWindow{genint: GenomicInterval{10, 19}, lambda: 4}
	eventcounter := 0
	var sitecounter = make(map[int64]int64)

	for i := 0; i < 1000; i++ {
		recevents := rw.getRecombinationNumber()
		eventcounter += int(recevents)
		for k := 0; k < int(recevents); k++ {
			site := rw.getRandomPosition()
			sitecounter[site]++
		}
	}

	if eventcounter < 3900 || eventcounter > 4100 {
		test.Errorf("Invalid number of recombination events, should be around 4000; observed %d ", eventcounter)
	}
	for key, value := range sitecounter {
		if key < 11 || key > 19 {
			test.Errorf("invalid site; only interval between 11-19 accepted; got  %d", key)

		}
		if value < 400 || value > 480 {
			test.Errorf("invalid number of recombination events for site %d; should be around 444; observed %d", key, value)
		}
	}
}

func TestStochasticRandomAssortment(test *testing.T) {
	util.SetSeed(5)
	gl := GenomicLandscape{offsets: []int64{10, 20, 30, 40, 50}}
	env = Environment{
		genome: &gl}

	sitecounter := make(map[int64]int64)
	for i := 0; i < 1000; i++ {
		recs := GetRecombinationEvents()
		for _, r := range recs {
			sitecounter[r]++
		}
	}
	if len(sitecounter) != 5 {
		test.Errorf("Incorrect number of sites; should be 5 is %d", len(sitecounter))
	}
	for _, key := range []int64{10, 20, 30, 40, 50} {
		got := sitecounter[key]
		if got < 450 || got > 550 {
			test.Errorf("Invalid number of random segregations for chromosome offset %d; should be around 500; got %d", key, got)
		}
	}
}

func TestStochasticRandomAssortmentAndRecombination(test *testing.T) {
	util.SetSeed(5)
	genome := newGenomicLandscape([]int64{10, 10})

	rwins := []*RecombinationWindow{&RecombinationWindow{genint: genome.intervals[0], lambda: 1},
		&RecombinationWindow{genint: genome.intervals[1], lambda: 1}}
	env = Environment{
		genome:               genome,
		recombinationWindows: rwins,
	}

	sitecounter := make(map[int64]int64)
	for i := 0; i < 1000; i++ {
		recs := GetRecombinationEvents()
		for _, r := range recs {
			sitecounter[r]++
		}
	}

	if len(sitecounter) != 20 {
		test.Errorf("Incorrect number of sites; should be 20 is %d", len(sitecounter))
	}
	for site, value := range sitecounter {
		if site < 0 || site > 19 {
			test.Errorf("Invalid site; must be between [0-19]; got %d", site)
		}
		if value < 70 {
			test.Errorf("Number of recombination events is too small; should be >70; got %d", value)
		}
		if site%10 == 0 && value < 480 {
			test.Errorf("Number of recombination events is too small (random assortement of chromosomes); should be >500; got %d", value)
		}
	}

}

func TestTranslateCoordinates(test *testing.T) {
	SetupEnvironment([]int64{100, 200, 300, 400}, // two chromosomes of size 1000
		nil, //
		[]float64{4, 4, 4, 4}, 0.1, 0.0)
	var tests = []struct {
		pos     int64
		wantchr int64
		wantpos int64
	}{{pos: 0, wantchr: 1, wantpos: 1},
		{pos: 1, wantchr: 1, wantpos: 2},
		{pos: 99, wantchr: 1, wantpos: 100},
		{pos: 100, wantchr: 2, wantpos: 1},
		{pos: 299, wantchr: 2, wantpos: 200},
		{pos: 300, wantchr: 3, wantpos: 1},
		{pos: 599, wantchr: 3, wantpos: 300},
		{pos: 600, wantchr: 4, wantpos: 1},
		{pos: 999, wantchr: 4, wantpos: 400},
	}

	for _, t := range tests {
		chrm, pos := TranslateCoordinates(t.pos)
		if chrm != t.wantchr {
			test.Errorf("Wrong chromosome; want %d got %d", t.wantchr, chrm)
		}
		if pos != t.wantpos {
			test.Errorf("Wrong position; want %d got %d", t.wantpos, pos)
		}

	}
}
