package outman

import (
	"bytes"
	"fmt"
	"invade/fly"
)

type FormaterGridOverview struct {
}

type FormaterGridCount struct {
}

type FormaterSummary struct {
}

type FormaterTransect struct {
}

type IConsoleFormater interface {
	FormatInfo() string
	FormatPopulation(p *fly.Population, reincos int64, generation int64, statstr string, sampleparsed []string) string
}

var formater IConsoleFormater

func (fo FormaterSummary) FormatInfo() string {
	// General info about the columns
	buf := new(bytes.Buffer)
	buf.WriteString("# ")
	buf.WriteString("rep\t")     // replicate
	buf.WriteString("gen\t")     // generation
	buf.WriteString("status\t")  // status
	buf.WriteString("phase\t")   // phase
	buf.WriteString("|\t")       // |
	buf.WriteString("cte\t")     // count with TE
	buf.WriteString("ctens\t")   // count with te and not silenced
	buf.WriteString("ctes\t")    // count with te and silenced
	buf.WriteString("avt_wte\t") // average count for those having the TE
	buf.WriteString("|\t")
	buf.WriteString("fwte\t")      // fraction of individuals with at leats one TE insertion
	buf.WriteString("avw\t")       //  fitness
	buf.WriteString("avtes\t")     //  TE insertions per diploid
	buf.WriteString("avpopfreq\t") //  population frquency of a TE insertion
	buf.WriteString("fixed\t")     // number of fixed TE insertions      // |
	buf.WriteString("fsilenced\t") // fraction of silenced
	buf.WriteString("sites\t")     // fraction of silenced
	buf.WriteString("|\t")
	buf.WriteString("denovo\t") // denovo count
	buf.WriteString("actWX\t")  // median x-position of active wave
	buf.WriteString("silWX\t")  // median x-position of silencing wave

	buf.WriteString("|\t")
	buf.WriteString("sampleids")
	return buf.String()
}

func (fo FormaterSummary) FormatPopulation(p *fly.Population, reincos int64, generation int64, statstr string, sampleparsed []string) string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("%d\t", reincos))                   // replicate
	buf.WriteString(fmt.Sprintf("%d\t", generation))                // generation
	buf.WriteString(fmt.Sprintf("%s\t", statstr))                   // status
	buf.WriteString(fmt.Sprintf("%s\t", getPhaseString(p.Phase()))) // phase
	buf.WriteString("|\t")
	buf.WriteString(fmt.Sprintf("%d\t", p.GetWithTECount()))                  // count with TE
	buf.WriteString(fmt.Sprintf("%d\t", p.GetWithTEAndNotSilencedCount()))    // count with TE and not silenced
	buf.WriteString(fmt.Sprintf("%d\t", p.GetWithTEAndSilencedCount()))       //  count with TE and silenced
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetAverageWithTE()))              //  average count for individuals having the TE
	buf.WriteString("|\t")                                                    // |
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetWithTEFrequency()))            // fwte
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetAverageFitness()))             // w
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetAverageInsertions()))          // avtes
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetAveragePopulationFrequency())) //  popfreq all
	buf.WriteString(fmt.Sprintf("%d\t", len(p.GetFixedInsertions())))         // fixed insertions                                                  // |
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetSilencedFrequency()))          // fw piRNAs
	buf.WriteString(fmt.Sprintf("%d\t", len(p.GetInsertionSites())))          // sites
	buf.WriteString("|\t")
	buf.WriteString(fmt.Sprintf("%d\t", p.GetDenovo()))                    // de novo silencing count
	buf.WriteString(fmt.Sprintf("%d\t", p.GetMedianXPosActiveWave()+1))    // x-position of active wave; 1-based output and input
	buf.WriteString(fmt.Sprintf("%d\t", p.GetMedianXPosSilencingWave()+1)) // x-position of silencing wave; 1-based output and input

	if len(sampleparsed) > 0 {
		buf.WriteString("|\t")
		for _, sid := range sampleparsed {
			buf.WriteString(fmt.Sprintf("%s\t", sid))
		}
	}
	tr := buf.String()
	return tr
}

func (fo FormaterGridOverview) FormatInfo() string {
	return ""
}
func (fo FormaterGridOverview) FormatPopulation(p *fly.Population, reincos int64, generation int64, statstr string, sampleparsed []string) string {
	buf := new(bytes.Buffer)
	buf.WriteString(">")
	buf.WriteString(fmt.Sprintf("%d\t", reincos))                   // replicate
	buf.WriteString(fmt.Sprintf("%d\t", generation))                // generation
	buf.WriteString(fmt.Sprintf("%s\t", statstr))                   // status
	buf.WriteString(fmt.Sprintf("%s\t", getPhaseString(p.Phase()))) // phase
	if len(sampleparsed) > 0 {
		buf.WriteString("|\t")
		for _, sid := range sampleparsed {
			buf.WriteString(fmt.Sprintf("%s\t", sid))
		}
	}
	buf.WriteString("\n")
	for _, t := range p.GetCountGrid() {
		for _, x := range t {
			if x > 0 {
				buf.WriteString("x")
			} else {
				buf.WriteString("o")
			}
		}
		buf.WriteString("\n")
	}
	buf.WriteString("<\n")
	return buf.String()
}

func (fo FormaterGridCount) FormatInfo() string {
	return ""
}
func (fo FormaterGridCount) FormatPopulation(p *fly.Population, reincos int64, generation int64, statstr string, sampleparsed []string) string {
	buf := new(bytes.Buffer)
	buf.WriteString(">")
	buf.WriteString(fmt.Sprintf("%d\t", reincos))                          // replicate
	buf.WriteString(fmt.Sprintf("%d\t", generation))                       // generation
	buf.WriteString(fmt.Sprintf("%s\t", statstr))                          // status
	buf.WriteString(fmt.Sprintf("%d\t", p.GetWithTEAndNotSilencedCount())) // with TE and NOT silenced
	buf.WriteString(fmt.Sprintf("%d\t", p.GetWithTEAndSilencedCount()))    // with TE and silenced
	buf.WriteString(fmt.Sprintf("%s\t", getPhaseString(p.Phase())))        // phase
	if len(sampleparsed) > 0 {
		buf.WriteString("|\t")
		for _, sid := range sampleparsed {
			buf.WriteString(fmt.Sprintf("%s\t", sid))
		}
	}
	buf.WriteString("\n")
	sg := p.GetSilencedGrid()
	dg := p.GetDenovoGrid()
	for yc, t := range p.GetCountGrid() {
		for xc, x := range t {
			buf.WriteString(fmt.Sprint(x))
			if dg[yc][xc] {
				buf.WriteString("*")
			} else if sg[yc][xc] {
				buf.WriteString("'")
			}
			buf.WriteString(" ")

		}
		buf.WriteString("\n")
	}
	buf.WriteString("<\n")
	return buf.String()
}

func getPhaseString(status fly.Phase) string {
	if status == fly.INVASION {
		return "inv"
	} else if status == fly.DECLINE {
		return "dec"
	} else if status == fly.INACTIVE {
		return "sil"
	} else {
		panic(fmt.Sprintf("unknown population status %d ", status))
	}
}

func (fo FormaterTransect) FormatInfo() string {
	return ""
}
func (fo FormaterTransect) FormatPopulation(p *fly.Population, reincos int64, generation int64, statstr string, sampleparsed []string) string {

	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("%d\t", reincos))                   // replicate
	buf.WriteString(fmt.Sprintf("%d\t", generation))                // generation
	buf.WriteString(fmt.Sprintf("%s\t", statstr))                   // status
	buf.WriteString(fmt.Sprintf("%s\t", getPhaseString(p.Phase()))) // phase
	transect := p.GetXTransectAvCount()
	for _, ts := range transect {
		buf.WriteString(fmt.Sprintf("%.2f\t", ts))
		buf.WriteString(" ")
	}
	return buf.String()
}
