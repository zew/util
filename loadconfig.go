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
	logx.Println("workDir: ", workDir)

	_, srcFile, _, ok := runtime.Caller(1)
	if !ok {
		logx.Fatalf("runtime caller not found")
	}

	paths := []string{
		path.Join(path.Dir(srcFile), "config.json"),
		path.Join(".", "config.json"),
		path.Join(".", "config", "config.json"),
	}

	found := false
	var file *os.File
	for _, v := range paths {
		file, err = os.Open(v)
		if err != nil {
			logx.Printf("could not open v: %v", v, err)
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
