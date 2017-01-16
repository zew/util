package util

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestHumanizeFloat6(t *testing.T) {

	type tcT struct {
		idx  int
		in   float64
		want string
	}

	tcs := []tcT{
		{0, float64(0.0) * float64(10.21), "0"},
		{1, float64(10.0) * float64(10.21), "102.1"},
		{2, float64(-10.0) * float64(10.21), "-102.1"},
		{3, float64(100.0) * float64(10.21), "1021"},
		{4, float64(-100.0) * float64(10.21), "-1021"},
		{5, float64(11111111.00001) * float64(10.21), "113444443.310102"},
		{6, float64(11111111.0000001) * float64(10.21), "113444443.310001"},
		{7, float64(.0119999999), "0.012"},
		{8, float64(-.0119999999), "-0.012"},
		{9, float64(.9999999), "1"},
		{10, float64(-.9999999), "-1"},
		{11, float64(5.9999999), "6"},
		{12, float64(-5.9999999), "-6"},
		{13, float64(1234.599999), "1234.6"},
		{14, float64(-1234.599999), "-1234.6"},
		{15, float64(14285.71428571428571428571428), "14285.714286"}, // rounding up
		{16, float64(0.1428576428571428571428571428), "0.142858"},
		{17, float64(0.0142857142857142857142857142), "0.0142857"},
		{18, float64(0.0014285714285714285714285714), "0.00142857"},
		{19, float64(0.0000014285714285714285714285), "0.00000142857"},
		{20, float64(0.0000014285796285714285714285), "0.00000142858"}, // rounding up
		{21, float64(0.9999914285714285714285714285), "1"},
		{22, float64(0.0099999142857162857142857142), "0.01"},
		{23, float64(0.0099142887162857142857142333), "0.00991429"}, // rounding up
		{24, float64(0.0000099999142857162857142857), "0.00001"},
		{25, float64(0.0000000000999991428571628571), "0.0000000001"},
		{26, float64(0.9999900000142857142857142857), "1"},
	}

	fh = FloatHumanizer(6)

	for _, tc := range tcs {
		got := fh.Humanize(tc.in)
		inAsStr := strconv.FormatFloat(tc.in, 'f', -1, 64)
		// t.Logf("idx %2v inp %26v - wnt %20v got %20v\n",tc.idx, inAsStr, tc.want, got)
		if got != tc.want {
			t.Errorf("idx %2v inp %26v - wnt %20v got %20v\n", tc.idx, inAsStr, tc.want, got)
		}
	}

}

func TestHumanizeFloat4(t *testing.T) {

	type tcT struct {
		idx  int
		in   float64
		want string
	}

	tcs := []tcT{
		{0, float64(0.0) * float64(10.21), "0"},
		{1, float64(10.0) * float64(10.21), "102.1"},
		{2, float64(-10.0) * float64(10.21), "-102.1"},
		{3, float64(100.0) * float64(10.21), "1021"},
		{4, float64(-100.0) * float64(10.21), "-1021"},
		{5, float64(11111111.00001) * float64(10.21), "113444443.3101"},
		{6, float64(11111111.0000001) * float64(10.21), "113444443.31"},
		{7, float64(.0119999999), "0.012"},
		{8, float64(-.0119999999), "-0.012"},
		{9, float64(.9999999), "1"},
		{10, float64(-.9999999), "-1"},
		{11, float64(5.9999999), "6"},
		{12, float64(-5.9999999), "-6"},
		{13, float64(1234.599999), "1234.6"},
		{14, float64(-1234.599999), "-1234.6"},
		{15, float64(14285.71428571428571428571428), "14285.7143"},   // rounding up
		{16, float64(0.1428576428571428571428571428), "0.1429"},      // rounding up
		{17, float64(0.0142857142857142857142857142), "0.01429"},     // rounding up
		{18, float64(0.0014285714285714285714285714), "0.001429"},    // rounding up
		{19, float64(0.0000014285714285714285714285), "0.000001429"}, // rounding up
		{20, float64(0.0000014285796285714285714285), "0.000001429"}, // rounding up
		{21, float64(0.9999914285714285714285714285), "1"},
		{22, float64(0.0099999142857162857142857142), "0.01"},
		{23, float64(0.0099142887162857142857142333), "0.009914"}, // rounding up
		{24, float64(0.0000099999142857162857142857), "0.00001"},
		{25, float64(0.0000000000999991428571628571), "0.0000000001"},
		{26, float64(0.9999900000142857142857142857), "1"},
	}

	fh = FloatHumanizer(4)

	for _, tc := range tcs {
		got := fh.Humanize(tc.in)
		inAsStr := strconv.FormatFloat(tc.in, 'f', -1, 64)
		// t.Logf("idx %2v inp %26v - wnt %20v got %20v\n",tc.idx, inAsStr, tc.want, got)
		if got != tc.want {
			t.Errorf("idx %2v inp %26v - wnt %20v got %20v\n", tc.idx, inAsStr, tc.want, got)
		}
	}

}

func TestHumanizeFloat2(t *testing.T) {

	type tcT struct {
		idx  int
		in   float64
		want string
	}

	tcs := []tcT{
		{0, float64(0.0) * float64(10.21), "0"},
		{1, float64(10.0) * float64(10.21), "102.1"},
		{2, float64(-10.0) * float64(10.21), "-102.1"},
		{3, float64(100.0) * float64(10.21), "1021"},
		{4, float64(-100.0) * float64(10.21), "-1021"},
		{5, float64(11111111.00001) * float64(10.21), "113444443.31"},
		{6, float64(11111111.0000001) * float64(10.21), "113444443.31"},
		{7, float64(.0119999999), "0.012"},
		{8, float64(-.0119999999), "-0.012"},
		{9, float64(.9999999), "1"},
		{10, float64(-.9999999), "-1"},
		{11, float64(5.9999999), "6"},
		{12, float64(-5.9999999), "-6"},
		{13, float64(1234.599999), "1234.6"},
		{14, float64(-1234.599999), "-1234.6"},
		{15, float64(14285.71428571428571428571428), "14285.71"},   // rounding up
		{16, float64(0.1428576428571428571428571428), "0.14"},      // rounding up
		{17, float64(0.0142857142857142857142857142), "0.014"},     // rounding up
		{18, float64(0.0014285714285714285714285714), "0.0014"},    // rounding up
		{19, float64(0.0000014285714285714285714285), "0.0000014"}, // rounding up
		{20, float64(0.0000014285796285714285714285), "0.0000014"}, // rounding up
		{21, float64(0.9999914285714285714285714285), "1"},
		{22, float64(0.0099999142857162857142857142), "0.01"},
		{23, float64(0.0099142887162857142857142333), "0.01"}, // rounding up
		{24, float64(0.0000099999142857162857142857), "0.00001"},
		{25, float64(0.0000000000999991428571628571), "0.0000000001"},
		{26, float64(0.9999900000142857142857142857), "1"},
	}

	fh = FloatHumanizer(2)

	for _, tc := range tcs {
		got := fh.Humanize(tc.in)
		inAsStr := strconv.FormatFloat(tc.in, 'f', -1, 64)
		// t.Logf("idx %2v inp %26v - wnt %20v got %20v\n",tc.idx, inAsStr, tc.want, got)
		if got != tc.want {
			t.Errorf("idx %2v inp %26v - wnt %20v got %20v\n", tc.idx, inAsStr, tc.want, got)
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
