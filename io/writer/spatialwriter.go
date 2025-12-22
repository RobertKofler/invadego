package writer

import (
	"fmt"
	"invade/fly"
	"os"
)

var spatialwriter *os.File

func SetupSpatialWriter(file string) {
	tmp, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	spatialwriter = tmp
	spatialwriter.WriteString("#rep\tgen\tY\tX\tTEs\tsilenced\tdenovo\n")
}

func WriteSpatialEntry(p *fly.Population, replicate int64, generation int64) {
	fsums := p.GetFlySummaries()
	for _, fs := range fsums {
		var si string
		var de string

		if fs.Silenced {
			si = "s"
		} else {
			si = "n"
		}
		if fs.Denovo {
			de = "d"
		} else {
			de = "n"
		}

		printline := fmt.Sprintf("%d\t%d\t%d\t%d\t%d\t%v\t%v", replicate, generation, fs.Ycoord+1, fs.Xcoord+1, fs.CountTE, si, de)
		spatialwriter.WriteString(printline + "\n")
	}

}

func CloseSpatialWriter() {
	if spatialwriter != nil {
		spatialwriter.Close()
	}
}
