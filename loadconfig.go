package util

import (
	"io"
	"os"
	"path"

	"github.com/zew/logx"
)

func LoadConfig(names ...string) io.Reader {

	workDir, err := os.Getwd()
	CheckErr(err)
	srcDir := logx.PathToSourceFile(1)
	logx.Println("work dir: ", workDir)
	logx.Println("src  dir: ", srcDir)

	fName := "config.json"
	if len(names) > 0 {
		fName = names[0]
	}

	paths := []string{
		path.Join(".", fName),
		path.Join(workDir, fName), // same as .
		path.Join(".", "config", fName),
		path.Join(srcDir, fName),                    // src file location
		path.Join(workDir, "appaccess-only", fName), // app engine:
	}

	found := false
	var file *os.File
	for _, v := range paths {
		file, err = os.Open(v)
		if err != nil {
			logx.Printf("- %v", err)
			continue
		}
		logx.Printf("found: %v", v)
		found = true
		break
	}

	if !found {
		logx.Fatalf("could not load config")
	}

	return file

}
