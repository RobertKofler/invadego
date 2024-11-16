package outman

import (
	"fmt"
	"invade/fly"
	"invade/io/writer"
	"strings"
)

var outman OutputManager

func SetupOutputManager(consoleFormat string, steps int64, replicateOffset int64,
	fileSpatial string, fileMHP string, fileDebug string, sampleid string) {

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
	if fileSpatial != "" {
		writer.SetupSpatialWriter(fileSpatial)
	}
	cf := strings.ToLower(consoleFormat)
	if cf == "summary" {
		formater = FormaterSummary{}

	} else if cf == "gridoverview" {
		formater = FormaterGridOverview{}

	} else if cf == "gridcount" {
		formater = FormaterGridCount{}
	} else {
		panic(fmt.Sprintf("Unknown console mode %s", cf))
	}

	outman = OutputManager{
		steps:           steps,
		replicateOffset: replicateOffset,
		fileMHP:         fileMHP,
		fileDebug:       fileDebug,
		fileSpatial:     fileSpatial,
		sampleid:        sampleid,
		sampleparsed:    sampleparsed,
	}

}

type OutputManager struct {
	steps           int64
	replicateOffset int64
	fileMHP         string
	fileSpatial     string
	fileDebug       string
	sampleid        string
	sampleparsed    []string
}

func WriteInfo(userargs string, seed int64, version string) {
	fmt.Println(fmt.Sprintf("# args: %s", userargs))
	fmt.Println(fmt.Sprintf("# version %s, seed: %d", version, seed))
	fmt.Println(formater.FormatInfo())
}

// Let the output manager know the job is done
// eg close open file handles
func End() {
	writer.CloseMHPWriter()
	writer.CloseDebugWriter()
	writer.CloseSpatialWriter()

}

func RecordPopulation(p *fly.Population, replicate int64, generation int64, popstat fly.PopStatus) {
	// Write populations if it is failure (including base population!)
	// or else if the generation has the required step (modulo == 0, hence including base population)
	if popstat == fly.FAIL0 || popstat == fly.FAILW || popstat == fly.FAILMAX {
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
	if outman.fileSpatial != "" {
		writer.WriteSpatialEntry(p, replicate+outman.replicateOffset, generation)
	}

	// INVADE
	stats := getStatusString(popstat)
	reos := replicate + outman.replicateOffset
	towrite := formater.FormatPopulation(p, reos, generation, stats, outman.sampleid)
	fmt.Println(towrite)
}

func getStatusString(popstat fly.PopStatus) string {
	if popstat == fly.BASEPOP {
		return "base"
	} else if popstat == fly.FAIL0 {
		return "fail-0"
	} else if popstat == fly.FAILW {
		return "fail-w"
	} else if popstat == fly.FAILMAX {
		return "fail-max"
	} else if popstat == fly.OK {
		return "ok"
	} else {
		panic("unknown population status ")
	}
}
