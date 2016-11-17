package util

import

// the package url

(
	"strings"
	"testing"
)

func TestHostCore(t *testing.T) {

	type tcT struct {
		inp   string
		want1 string
		want2 []string
	}

	tcs := []tcT{
		{"sd1.sd2.sd3.zew.de:332",
			"zew.de",
			[]string{"sd1", "sd2", "sd3"},
		},
		{"www.zew.de:332",
			"zew.de",
			[]string{},
		},
		{"sd1.www.zew.de:332",
			"zew.de",
			[]string{"sd1"},
		},
	}

	for i, tc := range tcs {
		got1, gots2 := HostCore(tc.inp)
		got2 := strings.Join(gots2, "")
		want2 := strings.Join(tc.want2, "")
		if got1 != tc.want1 || got2 != want2 {
			t.Errorf("%2v: \ninp     %-20v \nwant1-2 %-20v %-20v \ngot1-2  %-20v %-20v\n", i, tc.inp, tc.want1, want2, got1, got2)
		}
	}

}

func TestNormalizeSubdomainsToPath(t *testing.T) {

	type tcT struct {
		inp  string
		want string
	}

	tcs := []tcT{
		{"http://iche:secret@sd1.sd2.sd3.zew.de:332/dir1/dir2/file.ext?p1=v1#aaa",
			"zew.de/sd1/sd2/sd3/dir1/dir2/file.ext?p1=v1#aaa"},
		{"www.zew.de:332/dir1/dir2/file.ext?p1=v1#aaa",
			"zew.de/dir1/dir2/file.ext?p1=v1#aaa"},
		{"www.subd1.zew.de:33332/dir1/dir2/file.ext?p1=v1#aaa",
			"zew.de/www/subd1/dir1/dir2/file.ext?p1=v1#aaa"},
	}

	for i, tc := range tcs {
		u, err := UrlParseImproved(tc.inp)
		CheckErr(err)
		got := NormalizeSubdomainsToPath(u)
		t.Logf("%2v: \ninp  %-20v - \nwant %-20v \ngot  %-20v\n", i, tc.inp, tc.want, got)
		if got != tc.want {
			t.Errorf("%2v: \ninp  %-20v - \nwant %-20v \ngot  %-20v\n", i, tc.inp, tc.want, got)
		}
	}

}
