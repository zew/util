package util

import (
	"testing"
	"time"
)

const parseFmt string = "2006-01-02 15:04:05"

// https://github.com/golang/go/commit/1d9f67daf0e8ba950da75f68f1f3f2650b13cd67
func TestTimeZoneParsing(t *testing.T) {

	t.Logf("\nTesting no time zone...\n")
	do(
		t,
		[]string{"2012-07-31  22:34:45", "2012-12-31  22:34:45"},
		[]string{"2012-07-31 22:34:45 +0000 UTC", "2012-12-31 22:34:45 +0000 UTC"},
		parseFmt,
	)

	t.Logf("... explicit time zone - CET/CEST - zone AND hour are adapted according to month... ")
	ts := []string{
		"2012-07-31  22:34:45 (CET)",
		"2012-12-31  22:34:45 (CET)",
		"2012-07-31  22:34:45 (CEST)",
		"2012-12-31  22:34:45 (CEST)",
	}
	expect := []string{
		"2012-07-31 23:34:45 +0200 CEST",
		"2012-12-31 22:34:45 +0100 CET",
		"2012-07-31 22:34:45 +0200 CEST",
		"2012-12-31 21:34:45 +0100 CET",
	}
	do(t, ts, expect, parseFmt+" (MST)")

	t.Logf("... explicit location - CET is dynchanged to CEST depending on month; hour unchanged.")
	explicitLocation(t)

}

func do(tst *testing.T, ts, expect []string, f string) {
	for i, v := range ts {
		t, err := time.Parse(f, v)
		if err != nil {
			tst.Errorf("%v", err)
		}
		got := t.String()
		wnt := expect[i]
		if got != wnt {
			tst.Errorf("idx %2v - \nwnt %20v \ngot %20v\n", i, wnt, got)
		}
	}
}

func explicitLocation(tst *testing.T) {

	ts := []*time.Location{}

	ts = append(ts, time.Now().Location())

	loc, err := time.LoadLocation(time.Now().Location().String())
	if err != nil {
		tst.Errorf("%v", err)
	}
	ts = append(ts, loc)

	loc, err = time.LoadLocation("Local")
	if err != nil {
		tst.Errorf("%v", err)
	}
	ts = append(ts, loc)

	loc, err = time.LoadLocation("CET")
	if err != nil {
		tst.Errorf("%v", err)
	}
	ts = append(ts, loc)

	// loc, err = time.LoadLocation("UTC")
	// if err != nil {
	// 	tst.Errorf("%v", err)
	// }
	// ts = append(ts, loc)

	//
	loc, err = time.LoadLocation("Europe/Berlin")
	if err != nil {
		tst.Errorf("%v", err)
	}
	ts = append(ts, loc)

	//
	//
	for i, v := range ts {

		{
			t1, err := time.ParseInLocation(parseFmt, "2012-07-31 22:34:45", v)
			if err != nil {
				tst.Errorf("%v", err)
			}
			got := t1.String()
			wnt := "2012-07-31 22:34:45 +0200 CEST"
			if got != wnt {
				tst.Errorf("idx %2va - \nwnt %20v \ngot %20v\n", i, wnt, got)
			}

		}
		{
			t2, err := time.ParseInLocation(parseFmt, "2012-12-31 22:34:45", v)
			if err != nil {
				tst.Errorf("%v", err)
			}
			got := t2.String()
			wnt := "2012-12-31 22:34:45 +0100 CET"
			if got != wnt {
				tst.Errorf("idx %2va - \nwnt %20v \ngot %20v\n", i, wnt, got)
			}

		}
	}
	// fmt.Println(time.Now().String())

}
