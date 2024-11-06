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
		chrm, chrpos := env.TranslateCoordinates(pos)
		printline := fmt.Sprintf("%d\t%d\t%d\t%d\t%f", replicate, generation, chrm, chrpos, freq)
		mhpwriter.WriteString(printline + "\n")
	}

}

func CloseMHPWriter() {
	if mhpwriter != nil {
		mhpwriter.Close()
	}
}
