package util

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zew/logx"
)

func LoadConfig(names ...string) io.Reader {
	logx.Fatal("use LoadConfigFile() instead")
	var xx io.Reader
	return xx
}

// Loads a file relative to app dir.
func LoadConfigFile(fName string, optSubdir ...string) (io.ReadCloser, error) {

	workDir, err := os.Getwd()
	CheckErr(err)
	srcDir := logx.PathToSourceFile(1)
	logx.Printf("work dir: %v", workDir)
	logx.Printf("src  dir: %v", srcDir)

	ext := filepath.Ext(fName)
	fNameExample := strings.TrimSuffix(fName, ext) + "-example" + ext

	paths := []string{
		// path.Join(workDir, fName),                // same as next below
		path.Join(".", fName), // main.go
		path.Join(".", fNameExample),

		// Fallback: Search for file one directory higher -
		// Was useful because systemtests runs from one directory deeper.
		// Discontinued because systemtests now changes directory
		// path.Join("..", fName),
		// path.Join("..", fNameExample),
	}
	subDir := ""
	if len(optSubdir) > 0 {
		subDir = optSubdir[0]
	}
	if subDir != "" {
		paths = append(paths, path.Join(".", subDir, fName))
		paths = append(paths, path.Join("..", subDir, fName))
	}
	paths = append(paths, path.Join(srcDir, fName))                    // caller src file location
	paths = append(paths, path.Join(workDir, "appaccess-only", fName)) // app engine

	//
	found := false
	var file *os.File
	for _, v := range paths {
		file, err = os.Open(v)
		if err != nil && os.IsNotExist(err) {
			logx.Printf("Not found in  %v", v)
			continue
		} else if err != nil {
			logx.Fatalf("Error opening config file: %v", err)
		}
		logx.Printf("Found: %v", v)
		found = true
		break
	}

	if !found {
		return file, fmt.Errorf("Could not load %v (%v)", fName, subDir)
	}
	return file, nil

}

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
