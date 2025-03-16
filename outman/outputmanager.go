package outman

import (
	"fmt"
	"invade/fly"
	"invade/io/writer"
	"strings"
)

func SetupOutputManager(consoleFormat string, steps int64, replicateOffset int64,
	fileSpatial string, fileMHP string, fileDebug string, repsNeeded int64, sampleid string) *OutputManager {

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
	} else if cf == "transectx" {
		formater = FormaterTransect{}
	} else {
		panic(fmt.Sprintf("Unknown console mode %s", cf))
	}

	return &OutputManager{
		steps:            steps,
		replicateOffset:  replicateOffset,
		fileMHP:          fileMHP,
		fileDebug:        fileDebug,
		fileSpatial:      fileSpatial,
		sampleid:         sampleid,
		currentReplicate: 1,
		replicatesNeeded: repsNeeded,
		sampleparsed:     sampleparsed,
	}

}

type OutputManager struct {
	steps            int64
	replicateOffset  int64
	fileMHP          string
	fileSpatial      string
	fileDebug        string
	sampleid         string
	replicatesNeeded int64
	currentReplicate int64
	sampleparsed     []string
}

func (om *OutputManager) WriteInfo(userargs string, seed int64, version string) {
	fmt.Println(fmt.Sprintf("# args: %s", userargs))
	fmt.Println(fmt.Sprintf("# version %s, seed: %d", version, seed))
	fmt.Println(formater.FormatInfo())
}
func (om *OutputManager) FinaliseReplicate(popstat fly.PopStatus) {
	om.currentReplicate = om.currentReplicate + 1
}

func (om *OutputManager) NeedMoreReplicates() bool {
	if om.currentReplicate <= om.replicatesNeeded {
		return true
	} else {
		return false
	}
}

// Let the output manager know the job is done
// eg close open file handles
func (om *OutputManager) End() {
	writer.CloseMHPWriter()
	writer.CloseDebugWriter()
	writer.CloseSpatialWriter()

}

func (om *OutputManager) RecordPopulation(p *fly.Population, generation int64, popstat fly.PopStatus) {
	// Write populations if it is failure (including base population!)
	// or else if the generation has the required step (modulo == 0, hence including base population)
	if popstat == fly.FAIL0 || popstat == fly.FAILW || popstat == fly.FAILMAX {
		om.writePopulation(p, generation, popstat)

	} else if popstat == fly.OK && generation%om.steps == 0 {
		om.writePopulation(p, generation, popstat)

	}
	// Ignore if neither an unusual status or the requested recording generation
}

func (om *OutputManager) writePopulation(p *fly.Population, generation int64, popstat fly.PopStatus) {

	replicatetoprint := om.currentReplicate + om.replicateOffset - 1

	if om.fileMHP != "" {

		writer.WriteMHPEntry(p, replicatetoprint, generation)
	}
	if om.fileDebug != "" {
		writer.WriteDebugEntry(p, replicatetoprint, generation)
	}
	if om.fileSpatial != "" {
		writer.WriteSpatialEntry(p, replicatetoprint, generation)
	}
	// INVADE
	stats := getStatusString(popstat)
	towrite := formater.FormatPopulation(p, replicatetoprint, generation, stats, om.sampleparsed)
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
