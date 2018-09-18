package util

import (
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
