package writer

import (
	"fmt"
	"invade/env"
	"invade/fly"
	"os"
)

var mhpwriter *os.File

func SetupMHPWriter(file string) {
	tmp, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	mhpwriter = tmp
}

func WriteMHPEntry(p *fly.Population, replicate int64, generation int64) {
	insfreq := p.GetMHPPopulationFrequency()
	for pos, freq := range insfreq {
		score := env.ScoreInsertion(pos)
		chrm, chrpos := env.TranslateCoordinates(pos)
		printline := fmt.Sprintf("%d\t%d\t%d\t%d\t%s\t%f", replicate, generation, chrm, chrpos, score, freq)
		mhpwriter.WriteString(printline + "\n")
	}

}

func CloseMHPWriter() {
	if mhpwriter != nil {
		mhpwriter.Close()
	}
}
