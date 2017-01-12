package util

import (
	"fmt"
	"strings"
	"testing"
)

func TestHumanizeFloat(t *testing.T) {

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
		{float64(0.1428571428571428571428571428), "0.1428571"},
		{float64(0.0142857142857142857142857142), "0.01428571"},
		{float64(0.0014285714285714285714285714), "0.001428571"},
		{float64(14285.71428571428571428571428), "14285.714286"},
	}

	for _, tc := range tcs {
		got := HumanizeFloat(tc.in)
		t.Logf("%20v - want %20v got %20v\n", tc.in, tc.want, got)
		if got != tc.want {
			t.Errorf("%20v - want %20v got %20v\n", tc.in, tc.want, got)
		}
	}

}

func TestTrimRedundant(t *testing.T) {

	type tcT struct {
		in   []string
		want []string
	}

	tcs := []tcT{
		{[]string{"", ""}, []string{}},
		{[]string{"Bello", ""}, []string{"Bello"}},
		{[]string{"Bello", "Heino"}, []string{"Bello", "Heino"}},
		{[]string{"Bello", "ello"}, []string{"Bello"}},
		{[]string{"Bello", "Ello"}, []string{"Bello", "Ello"}},
		{[]string{"ino", "Heino"}, []string{"Heino"}},

		{[]string{"Bello", "Heino", "Bello not Heino"}, []string{"Bello not Heino"}},
		{[]string{"Bello", "Heino", "Bello not Harald"}, []string{"Heino", "Bello not Harald"}},

		{[]string{"Bello", "Heino", "Bello not Heino", "Cardigan"}, []string{"Bello not Heino", "Cardigan"}},
	}

	for i, tc := range tcs {
		got := TrimRedundant(tc.in)
		// t.Logf("%2v: inp %-20v - want %-20v got %-20v\n", i, strings.Join(tc.in, ","), strings.Join(tc.want, ","), strings.Join(got, ","))
		if strings.Join(got, ",") != strings.Join(tc.want, ",") {
			t.Errorf("%2v: inp %-20v - want %-20v got %-20v\n", i, strings.Join(tc.in, ","), strings.Join(tc.want, ","), strings.Join(got, ","))
		}
	}

}

func ExampleHumanizeFloat() {
	fmt.Println(HumanizeFloat(float64(-.0119999999)))
	// Output: -0.012
}
