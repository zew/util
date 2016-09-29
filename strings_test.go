package util

import "testing"

func TestB(t *testing.T) {

	type tcT struct {
		in   float64
		want string
	}

	tcs := []tcT{
		{float64(10.0) * float64(10.21), "102.1"},
		{float64(-10.0) * float64(10.21), "-102.1"},
		{float64(100.0) * float64(10.21), "1021"},
		{float64(-100.0) * float64(10.21), "-1021"},
		{float64(.0119999999), "0.012"},
		{float64(-.0119999999), "-0.012"},
		{float64(.9999999), "1"},
		{float64(-.9999999), "-1"},
		{float64(5.9999999), "6"},
		{float64(-5.9999999), "-6"},
		{float64(1234.599999), "1234.6"},
		{float64(-1234.599999), "-1234.6"},
	}

	for _, tc := range tcs {
		got := HumanizeFloat(tc.in)
		// t.Logf("%20v - want %20v got %20v\n", tc.in, tc.want, got)
		if got != tc.want {
			t.Errorf("%20v - want %20v got %20v\n", tc.in, tc.want, got)
		}
	}

}
