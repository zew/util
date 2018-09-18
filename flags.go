package util

import (
	"flag"
	"os"
	"strings"

	"github.com/zew/logx"
)

// Example: FlagVal("configfile", "config.json", "cfg", "JSON file containing config data")
type FlagT struct {
	Long       string // key, such as config_file, overrides Short key
	Short      string // short key, such as cfg
	DefaultVal string
	Desc       string // description, printed on executable -h

	Val   string // computed
	valSh string // possible value from the short key;
}

type FlagsT []FlagT

func NewFlags() (f *FlagsT) {
	f = &FlagsT{}
	return
}

func (f *FlagsT) Add(a FlagT) {
	if f == nil {
		new := FlagsT{}
		*f = new
	}
	*f = append(*f, a)
}

//
// Command line flag overrides environment variable
func (f *FlagsT) Gen() {

	// Oh god!
	// But there is no other way get pointers
	// into	our slice of FlagT structs.
	for idx := range *f {
		//    -cfg        config.json
		//    -lgn        logins.json
		flag.StringVar(&(*f)[idx].Val, (*f)[idx].Long, "", (*f)[idx].Desc)
		//    -config_file config.json
		//    -logins_file logins.json
		flag.StringVar(&(*f)[idx].valSh, (*f)[idx].Short, "", (*f)[idx].Desc+", shorthand")
	}

	flag.Parse() // Parse and register all flags

	//

	if len(flag.Args()) > 0 {
		logx.Printf("UNRECOGNIZED command line arguments: %v", flag.Args())
	}

	// Loop again for ENV variable fallback
	// and default value fallback.
	for idx, ff := range *f {

		val := (*f)[idx].Val
		valSh := (*f)[idx].valSh

		if val == "" && valSh != "" {
			val = valSh
			logx.Printf("Taking short key val, since long key is empty %v / %v => %q", ff.Long, ff.Short, val)
		}

		if val == "" {
			uLong := strings.ToUpper(ff.Long)
			val = string(os.Getenv(uLong))
			if val != "" {
				logx.Printf("Taking %v from ENV %v: %v", ff.Long, uLong, val)
			}
		}

		if val == "" {
			logx.Printf("Taking %v from DEFAULT: %v", ff.Long, ff.DefaultVal)
			val = ff.DefaultVal
		}

		logx.Printf("Effective %v / %v => %q", ff.Long, ff.Short, val)

		(*f)[idx].Val = val
	}

}
