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

type IConsoleFormater interface {
	FormatInfo() string
	FormatPopulation(p *fly.Population, reincos int64, generation int64, statstr string, sampleid string) string
}

var formater IConsoleFormater

func (fo FormaterSummary) FormatInfo() string {
	// General info about the columns
	buf := new(bytes.Buffer)
	buf.WriteString("# ")
	buf.WriteString("rep\t")     // replicate
	buf.WriteString("gen\t")     // generation
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

	buf.WriteString("|\t")
	buf.WriteString("sampleids")
	return buf.String()
}

func (fo FormaterSummary) FormatPopulation(p *fly.Population, reincos int64, generation int64, statstr string, sampleid string) string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("%d\t", reincos))    // replicate
	buf.WriteString(fmt.Sprintf("%d\t", generation)) // generation
	buf.WriteString(fmt.Sprintf("%s\t", statstr))    // status
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
	buf.WriteString(fmt.Sprintf("%.2f\t", p.GetSilencedFrequency()))          // fw piRNAs (either cluster or para)

	if len(outman.sampleparsed) > 0 {
		buf.WriteString("|\t")
		for _, sid := range outman.sampleparsed {
			buf.WriteString(fmt.Sprintf("%s\t", sid))
		}
	}
	return buf.String()
}

func (fo FormaterGridOverview) FormatInfo() string {
	return ""
}
func (fo FormaterGridOverview) FormatPopulation(p *fly.Population, reincos int64, generation int64, statstr string, sampleid string) string {
	buf := new(bytes.Buffer)
	buf.WriteString(">")
	buf.WriteString(fmt.Sprintf("%d\t", reincos))    // replicate
	buf.WriteString(fmt.Sprintf("%d\t", generation)) // generation
	buf.WriteString(fmt.Sprintf("%s\t", statstr))    // status
	if len(outman.sampleparsed) > 0 {
		buf.WriteString("|\t")
		for _, sid := range outman.sampleparsed {
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
	buf.WriteString(">\n")
	return buf.String()
}

func (fo FormaterGridCount) FormatInfo() string {
	return ""
}
func (fo FormaterGridCount) FormatPopulation(p *fly.Population, reincos int64, generation int64, statstr string, sampleid string) string {
	buf := new(bytes.Buffer)
	buf.WriteString(">")
	buf.WriteString(fmt.Sprintf("%d\t", reincos))    // replicate
	buf.WriteString(fmt.Sprintf("%d\t", generation)) // generation
	buf.WriteString(fmt.Sprintf("%s\t", statstr))    // status
	if len(outman.sampleparsed) > 0 {
		buf.WriteString("|\t")
		for _, sid := range outman.sampleparsed {
			buf.WriteString(fmt.Sprintf("%s\t", sid))
		}
	}
	buf.WriteString("\n")
	sg := p.GetSilencedGrid()
	for yc, t := range p.GetCountGrid() {
		for xc, x := range t {
			buf.WriteString(fmt.Sprint(x))
			if sg[yc][xc] {
				buf.WriteString("'")
			}
			buf.WriteString(" ")

		}
		buf.WriteString("\n")
	}
	buf.WriteString(">\n")
	return buf.String()
}
