package outman

import (
	"bytes"
	"fmt"
	"invade/fly"
	"invade/io/writer"
	"strings"
)

var outman OutputManager

func SetupOutputManager(steps int64, replicateOffset int64,
	fileMHP string, fileTally string, fileSFS string, fileDebug string, sampleid string) {
	sampleparsed := []string{}
	if strings.Contains(sampleid, ",") {
		sampleparsed = strings.Split(sampleid, ",")
	} else if len(sampleid) > 0 {
		sampleparsed = []string{sampleid}
	}

	if fileMHP != "" {
		writer.SetupMHPWriter(fileMHP)
	}
	if fileDebug != "" {
		writer.SetupDebugWriter(fileDebug)
	}

	outman = OutputManager{
		steps:           steps,
		replicateOffset: replicateOffset,
		fileMHP:         fileMHP,
		fileTally:       fileTally,
		fileDebug:       fileDebug,
		fileSFS:         fileSFS,
		sampleid:        sampleid,
		sampleparsed:    sampleparsed,
	}

}

type OutputManager struct {
	steps           int64
	replicateOffset int64
	fileMHP         string
	fileSFS         string
	fileTally       string
	fileDebug       string
	sampleid        string
	sampleparsed    []string
}

func WriteInfo(userargs string, usedseed int64, version string) {
	fmt.Println(fmt.Sprintf("# args: %s", userargs))
	fmt.Println(fmt.Sprintf("# version %s, seed: %d", version, usedseed))
	// General info about the columns
	buf := new(bytes.Buffer)
	buf.WriteString("# ")
	buf.WriteString("rep\t")       // replicate
	buf.WriteString("gen\t")       // generation
	buf.WriteString("popstat\t")   // population status
	buf.WriteString("|\t")         // |
	buf.WriteString("fwte\t")      // fraction of individuals with at least one TE insertion
	buf.WriteString("avw\t")       //  fitness
	buf.WriteString("minw\t")      //  minimum fitness during the invasion
	buf.WriteString("avtes\t")     //  TE insertions per diploid
	buf.WriteString("avpopfreq\t") //  population frquency of a TE insertion
	buf.WriteString("fixed\t")     // number of fixed TE insertions
	buf.WriteString("|\t")         // |
	buf.WriteString("phase\t")     // phase of the invasion; rapi, shot, inac
	buf.WriteString("fwcli\t")     // fraction of individuals with a cluster insertion
	buf.WriteString("avcli\t")     //  number of cluster insertions per individual
	buf.WriteString("fixcli\t")    // number of fixed cluster insertions
	buf.WriteString("|\t")
	buf.WriteString("sampleids")
	fmt.Println(buf.String())

}

// Let the output manager know the job is done
// eg close open file handles
func Done() {
	writer.CloseMHPWriter()
	writer.CloseDebugWriter()

}

func RecordPopulation(p *fly.Population, replicate int64, generation int64, popstat fly.PopStatus) {
	// Write populations if it is failure (including base population!)
	// or else if the generation has the required step (modulo == 0, hence including base population)
	if popstat == fly.FAIL0 || popstat == fly.FAILW {
		writePopulation(p, replicate, generation, popstat)
	} else if popstat == fly.OK && generation%outman.steps == 0 {
		writePopulation(p, replicate, generation, popstat)
	}
	// Ignore if neither an unusual status or the requested recording generation
}

func writePopulation(p *fly.Population, replicate int64, generation int64, popstat fly.PopStatus) {
	if outman.fileMHP != "" {
		writer.WriteMHPEntry(p, replicate+outman.replicateOffset, generation)
	}
	if outman.fileDebug != "" {
		writer.WriteDebugEntry(p, replicate+outman.replicateOffset, generation)
	}
	if outman.fileSFS != "" {

	}
	if outman.fileTally != "" {

	}
	// INVADE
	// #replicate	generation	| fwt	w	tes	popfreq	fixed	| fwcli	cluins	cluins_popfreq	cluins_fixed	phase	| novel	sites	clusites	tes_stdev	cluins_stdev	fw0	w_min	popsize

	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("%d\t", replicate+outman.replicateOffset))       // replicate
	buf.WriteString(fmt.Sprintf("%d\t", generation))                             // generation
	buf.WriteString(fmt.Sprintf("%s\t", getStatusString(popstat)))               // status
	buf.WriteString("|\t")                                                       // |
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetWithTEFrequency()))               // fwte
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetAverageFitness()))                // w
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetMinimumFitness()))                // minw
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetAverageInsertions()))             // avtes
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetAveragePopulationFrequency()))    //  popfreq all
	buf.WriteString(fmt.Sprintf("%d\t", len(p.GetFixedInsertions())))            // fixed insertions
	buf.WriteString("|\t")                                                       // |
	buf.WriteString(fmt.Sprintf("%s\t", getPhaseString(p.GetPhase())))           // Phase
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetWithClusterInsertionFrequency())) // fw cluster insertions
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetAverageClusterInsertions()))      //  number of cluster insertions
	buf.WriteString(fmt.Sprintf("%d\t", p.GetFixedClusterInsertionCount()))      // get fixed cluster insertions                                                    // |

	if len(outman.sampleparsed) > 0 {
		buf.WriteString("|\t")
		for _, sid := range outman.sampleparsed {
			buf.WriteString(fmt.Sprintf("%s\t", sid))
		}
	}

	// fwcli
	// fixcli
	// avcli
	// |
	// fwpara
	// fw
	fmt.Println(buf.String())
}

func getStatusString(popstat fly.PopStatus) string {
	if popstat == fly.BASEPOP {
		return "base"
	} else if popstat == fly.FAIL0 {
		return "fail-0"
	} else if popstat == fly.FAILW {
		return "fail-w"
	} else if popstat == fly.OK {
		return "ok"
	} else {
		panic("unknown population status ")
	}
}
func getPhaseString(status fly.Phase) string {
	if status == fly.RAPIDINVASION {
		return "rapi"
	} else if status == fly.SHOTGUN {
		return "shot"
	} else if status == fly.INACTIVE {
		return "inac"
	} else {
		panic("unknown population status ")
	}
}
