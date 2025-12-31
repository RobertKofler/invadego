package cmdparser

import (
	"flag"
	"os"
	"strings"
)

type CommandLineParameters struct {
	GridX          int64
	GridY          int64
	BasePop        string  // coordX,coordY,tecopynumber
	EpiSilencing   string  //  epigenetic silencing mode
	TriggerSilence int64   // TE copy number per diploid where silencing will be triggered
	MateRadius     string  // radius in grid where mates will be found
	SelfingRate    float64 // selfing rate between 0 and 1.0

	ConsoleFormat   string
	ArgString       string
	Silent          bool
	ClosedY         bool
	ClosedX         bool
	CondInvasion    bool
	Genome          string
	RecRate         string
	U               float64 // transposition rate
	UC              float64 // transposition rate in the presence of host defence
	S               float64 // deleterious effect of a TE insertion
	H               float64 // heterozygous effect of a TE insertion, new
	Steps           int64   // report output each Steps generations
	MateBarriers    string  // Mating Barriers; population structure
	BarrierStrength float64 // Barrier effect size; 1.0 no effect; 0.0 complete isolation
	PSiLoss         float64 // probabiltiy for loosing epigenetic silencing
	Generations     int64
	SampleID        string
	ReplicateOffset int64
	Replicates      int64
	Seed            int64
	Threads         int64
	MinFitness      float64
	MaxInsertions   int64
	FileMHP         string
	FileDebug       string
	FileSpatial     string
}

func ParseCommandLine() *CommandLineParameters {
	test := os.Args[1:]
	argstring := strings.Join(test, " ")
	// Mandatory parameters
	gridx := flag.Int64("grid-x", -1, "mandatory; the spatial grid size on X")
	gridy := flag.Int64("grid-y", -1, "mandatory; the spatial grid size on X")
	episilence := flag.String("epi-inherit", "", "mandatory;  either dros, ara, none")
	genome := flag.String("genome", "", "mandatory; the genomic landscape; e.g. 'MB:2,3,1,5' specifiies four chromosomes with sizes of 2,3,1,5 Mb")
	generations := flag.Int64("gen", -1, "mandatory; run the simulations for '--gen' generations")
	basepop := flag.String("basepop", "", "mandatory; coordinates of the individual(s) with the TE insertions in the first generation ; 1-based coordinates!; entries need to be in form Y,X,N:Y,X,N; where N' is feasible and means N silenced insertions; can also be in form  Ystart-Yend,Xstart-Xend,N")

	// Optional parameters
	triggerSilence := flag.Int64("trigger-defense", -1, "at which threshold value should the host defense been triggerd")
	psiloss := flag.Float64("p-sil-loss", 0.0, "probabiltiy that epigenetic silencing can be lost")
	mateRadius := flag.String("mate-radius", "2", "at which spatial distance in grid should mates be found; either integer eg 2 or 'pan' for panmictici")
	transrate := flag.Float64("u", 0.0, "the transposition rate")
	selfrate := flag.Float64("selfing-rate", 0.0, "the selfing rate")
	sampleid := flag.String("sampleid", "", "the ID of the sample; will be a help in R to group samples like with facete_grid()")
	rr := flag.String("rr", "", "the recombination rate per chromosome in cm/Mb; e.g. '3,4,4,5' ")
	s := flag.Float64("s", 0.0, "the deleterious effect of a single homozygous TE insertion")
	h := flag.Float64("h", 0.0, "the heterozygous effect of a TE insertion")
	closedx := flag.Bool("closed-x", false, "closed circle at x-coordinate")
	closedy := flag.Bool("closed-y", false, "closed circle at y-coordinate")
	//t := flag.Float64("t", 1.0, "the synergistic effect of TE insertions")
	transrateResidual := flag.Float64("uc", 0.0, "the transposition rate in the presence of piRNAs")
	steps := flag.Int64("steps", 20, "report the output at each '--steps' generations")
	replicates := flag.Int64("rep", 1, "the number of replicates")
	mateBarriers := flag.String("mate-barrier", "", "a comma separated list about x-coordiantes where mating barriers will be introduced eg. '50,100' will introduce mate barriers at x-coordinates 50 and 100")
	barStrength := flag.Float64("barrier-strength", 0.0, "the strength of mate barrier between 0.0 and 1.0; 1.0= total exclusion effect and 0.0= no effect")

	reploffset := flag.Int64("replicate-offset", 1, "starting index of the replicates; may be used for pseudo-parallelization)")
	fileSpatial := flag.String("file-spatial", "", "optional output file: info about individual specimens including coordinates")
	fileMHP := flag.String("file-mhp", "", "optional output file: position and population frequency of each insertion")
	fileDebug := flag.String("file-debug", "", "optional output file for debugging various aspects")
	consoleFormat := flag.String("console-format", "summary", "formating of console output, default summary: summary|gridoverview|gridcount|transectx")

	maxins := flag.Int64("max-insertions", 10000, "the maximum number of insertions")
	minw := flag.Float64("min-w", 0.1, "the minimum frequency of an average individual in the population")
	seed := flag.Int64("seed", -1, "seed for the random number generator")
	threads := flag.Int64("threads", 1, "number of threads")
	silent := flag.Bool("silent", false, "suppress output")
	conditionalInvasion := flag.Bool("condinv", false, "only show replicates where at least one TE copy was present until the final generation")
	flag.Parse()

	// basic checks if parameters are suitable
	if *gridx < 2 {
		panic("Provide a suitable grid size --grid-x; must be larger than 1")
	}
	if *gridy < 2 {
		panic("Provide a suitable grid size --grid-y; must be larger than 1")
	}
	if *selfrate < 0.0 || *selfrate > 1.0 {
		panic("Provide a suitable selfing rate --self-rate; must be between 0.0 and 1.0")
	}
	if *psiloss < 0.0 || *psiloss > 1.0 {
		panic("Provide a suitable probability --p-sil-loss; must be between 0.0 and 1.0")
	}
	if *barStrength < 0.0 || *barStrength > 1.0 {
		panic("Provide a suitable barrier effect --barrier-effect; must be between 0.0 and 1.0")
	}
	lepisilence := strings.ToLower(*episilence)
	if lepisilence != "none" && lepisilence != "ara" && lepisilence != "dros" {
		panic("Provide a suitable epigenetic inheritance mode --epi-inherit; must one of none, ara, dros")
	}
	lconsoleFormat := strings.ToLower(*consoleFormat)
	if lconsoleFormat != "summary" && lconsoleFormat != "gridoverview" && lconsoleFormat != "gridcount" && lconsoleFormat != "transectx" {
		panic("Provide a suitable console format: summary, gridoverview, gridcount, transect")
	}
	if *transrate < 0.0 {
		panic("Provide a suitable transposition rate --u; must be larger or equal to 0.0")
	}
	if *transrateResidual < 0.0 {
		panic("Provide a suitable residual transposition rate --uc; must be larger or equal to 0.0")
	}
	if *s < 0.0 || *s > 1.0 {
		panic("Provide a suitable deleterious effect of TEs --s; must be between 0.0 and 1.0")
	}
	if *h < 0.0 || *h > 1.0 {
		panic("Provide a suitable heterozygous effect of TEs --h; must be between 0.0 and 1.0")
	}

	//if *t < 1.0 {
	//	panic("Provide a suitable epistatic effect of TEs --t; must be larger or equal to 1.0")
	//}
	if *genome == "" {
		panic("Provide a suitable genome --genome")
	}
	if *basepop == "" {
		panic("Provide a suitable base population --basepop")
	}
	if *generations < 1 {
		panic("Provide a suitable number of generations --gen")
	}
	if *steps < 1 {
		panic("Provide suitable steps --steps; must be larger or equal to 1")
	}
	return &CommandLineParameters{
		ArgString:      argstring,
		Silent:         *silent,
		GridX:          *gridx,
		GridY:          *gridy,
		Genome:         *genome,
		TriggerSilence: *triggerSilence,
		SelfingRate:    *selfrate,
		EpiSilencing:   lepisilence,
		MateRadius:     *mateRadius,
		ConsoleFormat:  lconsoleFormat,
		CondInvasion:   *conditionalInvasion,

		RecRate:         *rr,
		PSiLoss:         *psiloss,
		MateBarriers:    *mateBarriers,
		BarrierStrength: *barStrength,
		BasePop:         *basepop,
		U:               *transrate,
		UC:              *transrateResidual,
		S:               *s,
		H:               *h,
		ClosedX:         *closedx,
		ClosedY:         *closedy,
		Steps:           *steps,
		ReplicateOffset: *reploffset,
		Seed:            *seed,
		Threads:         *threads,
		MinFitness:      *minw,
		Replicates:      *replicates,
		MaxInsertions:   *maxins,
		FileSpatial:     *fileSpatial,
		FileMHP:         *fileMHP,
		FileDebug:       *fileDebug,
		Generations:     *generations,
		SampleID:        *sampleid} //TODO implement as output
}
