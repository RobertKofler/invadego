package env

import "fmt"

type TEInsertion byte

// report the insertion bias into piRNA clusters in %;
// possible values range from -100 over 0 to +100
// -100 no insertion will ever go into clusters
// 0 no insertion bias, insertion probability into piRNA cluster is just the size of piRNA cluster
// 100 every insertion will go into cluster
func (te TEInsertion) BiasPercent() int64 {
	bias := int64(te) - 100
	if bias < -100 || bias > 100 {
		panic(fmt.Sprintf("Invalid insertion bias, must be between -100 and 100, got %d", bias))
	}
	return bias

}

// report the insertion bias into piRNA cluster as fraction;
// possible values range from -1.0 over 0.0 to 1.0
func (te TEInsertion) BiasFraction() float64 {
	bias := float64(te)/100.0 - 1.0
	if bias < -1.0 || bias > 1.0 {
		panic(fmt.Sprintf("Invalid insertion bias, must be between -1.0 and 1.0, got %f", bias))
	}
	return bias
}
