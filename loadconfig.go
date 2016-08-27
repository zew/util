package util

import (
	"io"
	"os"
	"path"
	"runtime"

	"github.com/zew/logx"
)

func LoadConfig() io.Reader {

	workDir, err := os.Getwd()
	CheckErr(err)
	logx.Println("workDir is: ", workDir)

	_, srcFile, _, ok := runtime.Caller(1)
	if !ok {
		logx.Fatalf("runtime caller not found")
	}

	fName := "config.json"
	paths := []string{
		path.Join(".", fName),
		path.Join(".", "config", fName),
		path.Join(workDir, fName),
		path.Join(path.Dir(srcFile), fName), // src file location as last option
	}

	found := false
	var file *os.File
	for _, v := range paths {
		file, err = os.Open(v)
		if err != nil {
			logx.Printf("could not open: %v %v", v, err)
			continue
		}
		found = true
		break
	}

	if !found {
		logx.Fatalf("could not load config")
	}

	return file

}
