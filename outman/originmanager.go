package outman

import (
	"fmt"
	"invade/fly"
	"sort"
	"strings"
)

type OriginManager struct {
	keytable map[int64]int64
	counter  int64
}

var originman *OriginManager

func newOriginManger() *OriginManager {
	toret := OriginManager{keytable: make(map[int64]int64), counter: 1}
	return &toret
}

func (om *OriginManager) GetShortOriginID(longid int64) int64 {
	if val, ok := om.keytable[longid]; ok {
		return val
	} else {
		toret := om.counter
		om.keytable[longid] = toret
		om.counter++
		return toret
	}
}

func formatOriginFreq(ostat []fly.OriginFreq, minfreq float64) string {
	sort.Slice(ostat, func(i, j int) bool { return ostat[i].Freq > ostat[j].Freq })

	tojoin := []string{}
	for _, o := range ostat {
		if o.Freq >= minfreq {
			key := originman.GetShortOriginID(o.Origin)
			tojoin = append(tojoin, fmt.Sprintf("%d:%.2f", key, o.Freq))
		}
	}

	return strings.Join(tojoin, ",")

}
