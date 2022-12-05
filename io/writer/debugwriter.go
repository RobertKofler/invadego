package writer

import (
	"fmt"
	"invade/fly"
	"os"
)

var debugwriter *os.File

func SetupDebugWriter(file string) {
	tmp, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	debugwriter = tmp
}

func WriteDebugEntry(p *fly.Population, replicate int64, generation int64) {
	d := p.GetD(0, 999999) // debugging LD decay
	printline := fmt.Sprintf("%d\t%d\t%f", replicate, generation, d)
	debugwriter.WriteString(printline + "\n")

}

func CloseDebugWriter() {
	if debugwriter != nil {
		debugwriter.Close()
	}
}
