package outman

import "testing"

// command line, run all tests "go test ./..." yes three points

func TestOriginManager(test *testing.T) {
	orm := newOriginManger()
	var tests = []struct {
		long int64
		want int64
	}{
		{long: 1000, want: 1},
		{long: 500, want: 2},
		{long: 10000, want: 3},
		{long: 1000, want: 1},
		{long: 500, want: 2},
		{long: 10000, want: 3},
		{long: 1, want: 4},
	}

	for _, t := range tests {
		got := orm.GetShortOriginID(t.long)
		if got != t.want {
			test.Errorf("orm.GetShortOriginID(%d)!=%d; got %d", t.long, t.want, got)
		}

	}
}
