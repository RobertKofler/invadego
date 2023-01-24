package env

// command line, run all tests "go test ./..." yes three points
import (
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
func TestNewCluster(t *testing.T) { // TODO test non cluster
	gl := newGenomicLandscape([]int64{100, 200, 300, 400})
	cl, _ := newClusterNonClusters([]int64{10, 20, 30, 40}, gl)
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

func TestNilClusterReference(t *testing.T) {
	gl := newGenomicLandscape([]int64{100, 100, 100, 100})
	cl, _ := newClusterNonClusters(nil, gl)
	env = Environment{genome: gl,
		clusters: *cl}
	for i := int64(0); i < gl.totalGenome; i++ {
		if IsClusterInsertion(i) {
			t.Errorf("incorrect cluster insertion in nil cluster, position %d", i)
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
		tot    int64
		cluins int64 // number of cluster insertions
		jump   Jumper
		want   float64
	}{
		{tot: 10, cluins: 0, jump: Jumper{u: 0.1, uc: 0.0}, want: 1.0},
		{tot: 10, cluins: 1, jump: Jumper{u: 0.1, uc: 0.0}, want: 0.0},
		{tot: 100, cluins: 0, jump: Jumper{u: 0.1, uc: 0.0}, want: 10.0},
		{tot: 100, cluins: 1, jump: Jumper{u: 0.1, uc: 0.0}, want: 0.0},
		{tot: 100, cluins: 1, jump: Jumper{u: 0.1, uc: 0.01}, want: 1.0},
	}

	for _, test := range tests {
		ju := test.jump
		got := ju.getNovelInsertionCount(test.tot, test.cluins)
		dif := math.Abs(test.want - got)
		if dif > 0.00001 {
			t.Errorf("getNovelInsertionCount(); got %f wanted %f", got, test.want)
		}

	}
}

func TestStochasticGetNovelInsertionSites(test *testing.T) {
	util.SetSeed(5)
	SetJumper(0.1, 0.0)
	genome := newGenomicLandscape([]int64{10, 10})
	env = Environment{
		genome: genome,
	}

	var sitecounter = make(map[int64]int64)
	totcounter := 0
	for i := 0; i < 1000; i++ {

		newsites := GetNewTranspositionSites(100, 0)
		totcounter += len(newsites)
		for _, n := range newsites {
			sitecounter[n]++
		}
	}
	if totcounter < 4900 || totcounter > 5100 {
		// why 5000? 1000 * 0.1*100 / 2 (division by two because insertions in a haploid gamete
		test.Errorf("Invalid number of novel TE insertion events, should be around 5000; observed %d ", totcounter)
	}
	for site, count := range sitecounter {
		if count < 200 || count > 3000 {
			test.Errorf("Invalid number of TE insertion events for site %d, should be around 250; observed %d ", site, count)
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
		[]float64{4, 4, 4, 4}, 0.1)
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
