package util

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
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

	paths := []string{
		// path.Join(workDir, fName),                // same as next below
		path.Join(".", fName),  // main.go
		path.Join("..", fName), // one fallback from one directory higher - usually because integrationtests runs from one directory deeper
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
	Long       string
	Short      string
	DefaultVal string
	Desc       string

	Val string // computed
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

	for idx, ff := range *f {
		val := ""
		// i.e. --configfile file1.json
		// or   --configfile=file1.json
		flag.StringVar(&val, ff.Long, "", ff.Desc)
		// i.e. -cfg file2.json
		// or   -cfg=file2.json
		flag.StringVar(&val, ff.Short, "", ff.Desc+", shorthand")
		(*f)[idx].Val = val
	}

	flag.Parse() // Parse and register all flags

	// Loop again for ENV variable fallback
	// and default value fallback.
	for idx, ff := range *f {

		val := (*f)[idx].Val

		if val == "" {
			uLong := strings.ToLower(ff.Long)
			val = string(os.Getenv(uLong))
			if val != "" {
				logx.Printf("Taking %v from ENV %v: %v", ff.Long, uLong, val)
			}
		}

		if val == "" {
			val = ff.DefaultVal
		}

		logx.Printf("flag %v = %v", ff.Long, val)

		(*f)[idx].Val = val
	}

}
