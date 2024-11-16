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
	spatialwriter.WriteString("#rep\tgen\tY\tX\tTEs\tsilenced\n")
}

func WriteSpatialEntry(p *fly.Population, replicate int64, generation int64) {
	fsums := p.GetFlySummaries()
	for _, fs := range fsums {
		printline := fmt.Sprintf("%d\t%d\t%d\t%d\t%d\t%v", replicate, generation, fs.Ycoord, fs.Xcoord, fs.CountTE, fs.Silenced)
		spatialwriter.WriteString(printline + "\n")
	}

}

func CloseSpatialWriter() {
	if spatialwriter != nil {
		spatialwriter.Close()
	}
}
