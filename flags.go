package util

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
FlagT is a single setting with keys and description.
Val holds the value.
There is a short and a long key for giving a value.

	myexecutable  -cfg=config.json  -config_file=config.json

The equal signs are optional.
The hyphen is required.
The long key takes precedence over the short key.

	SET CONFIG_FILE=config.json

The environment variable is only used, if the short and the long key are empty.

Usage:

	fl := util.NewFlags()
	fl.Add(
		util.FlagT{
			Long:       "config_file",
			Short:      "cfg",
			DefaultVal: "config.json",
			Desc:       "JSON file containing config data",
		},
	)
	fl.Add(
		...
	)
	fl.Gen()  // filling Val(s)

	cfg.CfgPath = fl.ByKey("cfg").Val  // Reading some value
*/

// FlagT reads some command line value
type FlagT struct {
	Long       string // key, such as config_file, overrides Short key
	Short      string // short key, such as cfg
	DefaultVal string
	Desc       string // description, printed on executable -h

	Val   string // computed
	valSh string // possible value from the short key; only if long key was empty
}

// FlagsT is a slice of settings
type FlagsT []FlagT

// NewFlags returns a slice of settings
func NewFlags() (f *FlagsT) {
	f = &FlagsT{}
	return
}

// Add adds a setting to the slice of settings
func (f *FlagsT) Add(a FlagT) {
	if f == nil {
		new := FlagsT{}
		*f = new
	}
	*f = append(*f, a)
}

// Gen fills the Val members of the slice of flags.
// Command line flag overrides environment variable.
func (f *FlagsT) Gen() {

	// Oh god! But there is no other way get pointers
	// into	our slice of FlagT structs.
	for idx := range *f {
		//    -cfg        config.json
		flag.StringVar(&(*f)[idx].Val, (*f)[idx].Long, "", (*f)[idx].Desc)
		//    -config_file config.json
		flag.StringVar(&(*f)[idx].valSh, (*f)[idx].Short, "", (*f)[idx].Desc+", shorthand")
	}

	flag.Parse() // Parse and register all flags

	//

	if len(flag.Args()) > 0 {
		log.Printf("UNRECOGNIZED command line arguments: %v", flag.Args())
	}

	// Loop again
	//  - for long overriding short
	//  - for ENV variable fallback
	//  - for default value fallback.
	for idx, ff := range *f {

		val := (*f)[idx].Val
		valSh := (*f)[idx].valSh

		if val == "" && valSh != "" {
			val = valSh
			log.Printf("Taking short key val, since long key is empty %v / %v => %q", ff.Long, ff.Short, val)
		}

		if val == "" {
			uLong := strings.ToUpper(ff.Long)
			val = string(os.Getenv(uLong))
			if val != "" {
				log.Printf("Taking %v from ENV %v: %v", ff.Long, uLong, val)
			}
		}

		if val == "" {
			log.Printf("Taking %v from DEFAULT: %v", ff.Long, ff.DefaultVal)
			val = ff.DefaultVal
		}

		log.Printf("Effective %v / %v => %q", ff.Long, ff.Short, val)

		(*f)[idx].Val = val
	}

}

// ByKey returns one setting by key.
// No error for better chaining. But panic.
func (f *FlagsT) ByKey(longOrShort string) (a FlagT) {
	fl := FlagT{}
	for _, fl = range *f {
		if fl.Short == longOrShort || fl.Long == longOrShort {
			return fl
		}
	}
	panic(fmt.Sprintf("%v is not a defined flag", longOrShort))
}
