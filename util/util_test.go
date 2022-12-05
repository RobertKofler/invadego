package util

import (
	"testing"
)

func TestMean(t *testing.T) {
	var tests = []struct {
		data []int64
		want float64
	}{
		{data: []int64{2}, want: 2},
		{data: []int64{2, 4}, want: 3},
		{data: []int64{2, 1, 3}, want: 2},
		{data: []int64{12, 100}, want: 56}}

	for _, test := range tests {

		if Mean(test.data) != test.want {
			t.Errorf("Mean(%v)!=%f", test.data, test.want)
		}
	}
}

func TestVariance(t *testing.T) {
	// https://www.calculatorsoup.com/calculators/statistics/variance-calculator.php
	var tests = []struct {
		data []int64
		want float64
	}{
		{data: []int64{2, 4}, want: 1},
		{data: []int64{2, 8}, want: 9},
		{data: []int64{2, 4, 6, 8}, want: 5}}

	for _, test := range tests {

		if Variance(test.data) != test.want {
			t.Errorf("Mean(%v)!=%f", test.data, test.want)
		}
	}
}

// command line, run all tests "go test ./..." yes three points
func TestPoissonLow(t *testing.T) {
	SetSeed(2)
	tol := 0.1
	lambda := 4.0
	rands := make([]int64, 10000)
	for i := int64(0); i < 10000; i++ {
		rands[i] = Poisson(lambda)
	}
	mean := Mean(rands)
	vari := Variance(rands)
	if mean < lambda-tol || mean > lambda+tol {
		t.Error("incorrect Poisson distribution")
	}
	if vari < lambda-tol || vari > lambda+tol {
		t.Error("incorrect Poisson distribution")
	}

}

func TestPoissonHigh(t *testing.T) {
	SetSeed(5)
	tol := 31.0
	lambda := 1500.0
	rands := make([]int64, 10000)
	for i := int64(0); i < 10000; i++ {
		rands[i] = Poisson(lambda)
	}
	mean := Mean(rands)
	vari := Variance(rands)
	if mean < lambda-tol || mean > lambda+tol {
		t.Error("incorrect Poisson distribution")
	}
	if vari < lambda-tol || vari > lambda+tol {
		t.Error("incorrect Poisson distribution")
	}
}
