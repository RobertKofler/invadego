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
	EpiSilencing   float64 // fraction of offspring that ends up with epigenetic silencing
	TriggerSilence int64   // TE copy number per diploid where silencing will be triggered
	MateRadius     int64   // radius in grid where mates will be found
	SelfingRate    float64 // selfing rate between 0 and 1.0
	FileSpatial    string

	ArgString       string
	Silent          bool
	Genome          string
	RecRate         string
	U               float64 // transposition rate
	UC              float64 // transposition rate in the presence of host defence
	X               float64 // deleterious effect of a TE insertion
	Steps           int64   // report output each Steps generations
	Generations     int64
	SampleID        string
	ReplicateOffset int64
	Replicates      int64
	Seed            int64
	Threads         int64
	MinFitness      float64
	MaxInsertions   int64
	FileMHP         string
	FileTally       string
	FileDebug       string
	FileSFS         string
}

func ParseCommandLine() *CommandLineParameters {
	test := os.Args[1:]
	argstring := strings.Join(test, " ")
	// Mandatory parameters
	gridx := flag.Int64("grid-x", -1, "mandatory; the spatial grid size on X")
	gridy := flag.Int64("grid-y", -1, "mandatory; the spatial grid size on X")
	episilence := flag.Float64("epi-inherit", -1.0, "mandatory;  fraction of offspring with epigenetic silencing")

	genome := flag.String("genome", "", "mandatory; the genomic landscape; e.g. 'MB:2,3,1,5' specifiies four chromosomes with sizes of 2,3,1,5 Mb")
	generations := flag.Int64("gen", -1, "mandatory; run the simulations for '--gen' generations")
	basepop := flag.String("basepop", "", "mandatory; the individual(s) with the segregating insertions in the starting population; CoordX,CoordY,N")

	// Optional parameters
	triggerSilence := flag.Int64("trigger-defense", -1, "at which threshold value should the host defense been triggerd")
	mateRadius := flag.Int64("mate-radius", 2, "at which spatial distance in grid should mates be found")
	transrate := flag.Float64("u", 0.0, "the transposition rate")
	selfrate := flag.Float64("selfing-rate", 0.0, "the selfing rate")
	sampleid := flag.String("sampleid", "", "the ID of the sample; will be a help in R to group samples like with facete_grid()")
	rr := flag.String("rr", "", "the recombination rate per chromosome in cm/Mb; e.g. '3,4,4,5' ")
	x := flag.Float64("x", 0.0, "the deleterious effect of a single TE insertions")
	//t := flag.Float64("t", 1.0, "the synergistic effect of TE insertions")
	transrateResidual := flag.Float64("uc", 0.0, "the transposition rate in the presence of piRNAs")
	steps := flag.Int64("steps", 20, "report the output at each '--steps' generations")
	replicates := flag.Int64("rep", 1, "the number of replicates")
	reploffset := flag.Int64("replicate-offset", 1, "starting index of the replicates; may be used for pseudo-parallelization)")
	fileSpatial := flag.String("file-spatial", "", "optional output file: info about individual specimens including coordinates")
	fileMHP := flag.String("file-mhp", "", "optional output file: position and population frequency of each insertion")
	fileDebug := flag.String("file-debug", "", "optional output file for debugging various aspects")
	fileSFS := flag.String("file-sfs", "", "optional output file: site frequency spectra of TE insertions")
	fileTally := flag.String("file-tally", "", "optional output file: count of insertions per individual")
	maxins := flag.Int64("max-insertions", 10000, "the maximum number of insertions")
	minw := flag.Float64("min-w", 0.1, "the minimum frequency of an average individual in the population")
	seed := flag.Int64("seed", -1, "seed for the random number generator")
	threads := flag.Int64("threads", 1, "number of threads")
	silent := flag.Bool("silent", false, "suppress output")
	flag.Parse()

	// basic checks if parameters are suitable
	if *gridx < 2 {
		panic("Provide a suitable grid size --grid-x; must be larger than 1")
	}
	if *gridy < 2 {
		panic("Provide a suitable grid size --grid-y; must be larger than 1")
	}
	if *mateRadius < 1 {
		panic("Proivde a suitable mate radius --mate-radius; must be larger than 0")
	}
	if *selfrate < 0.0 || *selfrate > 1.0 {
		panic("Provide a suitable selfing rate --self-rate; must be between 0.0 and 1.0")
	}
	if *episilence < 0.0 || *episilence > 1.0 {
		panic("Provide a suitable epigenetic inheritance rate --epi-inherit; must be between 0.0 and 1.0")
	}
	if *transrate < 0.0 {
		panic("Provide a suitable transposition rate --u; must be larger or equal to 0.0")
	}
	if *transrateResidual < 0.0 {
		panic("Provide a suitable residual transposition rate --uc; must be larger or equal to 0.0")
	}
	if *x < 0.0 {
		panic("Provide a suitable deleterious effect of TEs --x; must be larger or equal to 0.0")
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
		EpiSilencing:   *episilence,
		MateRadius:     *mateRadius,

		RecRate:         *rr,
		BasePop:         *basepop,
		U:               *transrate,
		UC:              *transrateResidual,
		X:               *x,
		T:               *t,
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
		FileTally:       *fileTally,
		FileSFS:         *fileSFS,
		Generations:     *generations,
		SampleID:        *sampleid} //TODO implement as output
}
