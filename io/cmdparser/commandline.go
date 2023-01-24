package cmdparser

import (
	"flag"
	"os"
	"strings"
)

type CommandLineParameters struct {
	ArgString       string
	Silent          bool
	Popsize         int64
	Genome          string
	Cluster         string
	RecRate         string
	U               float64 // transposition rate
	UC              float64 // transposition rate in the presence of piRNAs
	X               float64 // deleterious effect of a TE insertion
	T               float64 // exponential deleterious effect of a TE insertion
	Steps           int64   // report output each Steps generations
	Generations     int64
	BasePop         string
	Noxcluins       bool
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
	popsize := flag.Int64("N", -1, "mandatory; the population size")
	genome := flag.String("genome", "", "mandatory; the genomic landscape; e.g. 'MB:2,3,1,5' specifiies four chromosomes with sizes of 2,3,1,5 Mb")
	generations := flag.Int64("gen", -1, "mandatory; run the simulations for '--gen' generations")
	basepop := flag.String("basepop", "", "mandatory; the segregating insertions in the starting population; either number (e.g. 100) or file")

	// Optional parameters
	transrate := flag.Float64("u", 0.0, "the transposition rate")
	cluster := flag.String("cluster", "", "piRNA clusters; e.g. 'kb:1,1,1,1' specifies a cluster of 1kb at the beginning of each chromosome")
	sampleid := flag.String("sampleid", "", "the ID of the sample; will be a help in R to group samples like with facete_grid()")
	rr := flag.String("rr", "", "the recombination rate per chromosome in cm/Mb; e.g. '3,4,4,5' ")
	x := flag.Float64("x", 0.0, "the deleterious effect of a single TE insertions")
	t := flag.Float64("t", 1.0, "the synergistic effect of TE insertions")
	noxcluins := flag.Bool("no-x-cluins", false, "cluster insertions incur no negative effects")
	//ignoreFailed := flag.Bool("ignored-failed", false, "ignore invasions where the TE did not get established")
	transrateResidual := flag.Float64("uc", 0.0, "the transposition rate in the presence of piRNAs")
	steps := flag.Int64("steps", 20, "report the output at each '--steps' generations")
	replicates := flag.Int64("rep", 1, "the number of replicates")
	reploffset := flag.Int64("replicate-offset", 1, "starting index of the replicates; may be used for pseudo-parallelization)")
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
	if *popsize < 2 {
		panic("Provide a suitable population size --N; must be larger than 1")
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
	if *t < 1.0 {
		panic("Provide a suitable epistatic effect of TEs --t; must be larger or equal to 1.0")
	}
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
		ArgString:       argstring,
		Silent:          *silent,
		Popsize:         *popsize,
		Genome:          *genome,
		Cluster:         *cluster,
		RecRate:         *rr,
		BasePop:         *basepop,
		U:               *transrate,
		UC:              *transrateResidual,
		X:               *x,
		T:               *t,
		Steps:           *steps,
		Noxcluins:       *noxcluins,
		ReplicateOffset: *reploffset,
		Seed:            *seed,
		Threads:         *threads,
		MinFitness:      *minw,
		Replicates:      *replicates,
		MaxInsertions:   *maxins,
		FileMHP:         *fileMHP,
		FileDebug:       *fileDebug,
		FileTally:       *fileTally,
		FileSFS:         *fileSFS,
		Generations:     *generations,
		SampleID:        *sampleid} //TODO implement as output
}
