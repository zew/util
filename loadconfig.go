package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// Returns the the source file path.
// Good to read file inside a library,
// completely independent of working dir
// or application dir.
func PathToSourceFile(levelsUp ...int) string {
	lvlUp := 1
	if len(levelsUp) > 0 {
		lvlUp = 1 + levelsUp[0]
	}
	_, srcFile, _, ok := runtime.Caller(lvlUp)
	if !ok {
		log.Panic("runtime caller not found")
	}
	p := path.Dir(srcFile)
	return p
}

func LoadConfig(names ...string) io.Reader {
	log.Fatal("use LoadConfigFile() instead")
	var xx io.Reader
	return xx
}

// Loads a file relative to app dir.
func LoadConfigFile(fName string, optSubdir ...string) (io.ReadCloser, error) {

	workDir, err := os.Getwd()
	CheckErr(err)
	srcDir := PathToSourceFile(1) // assuming caller main() in app root
	log.Printf("fileName: %v", fName)
	log.Printf("work dir: %v", workDir)
	log.Printf("src  dir: %v", srcDir)

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
			log.Printf("Not found in  %v", v)
			continue
		} else if err != nil {
			log.Fatalf("Error opening config file: %v", err)
		}
		log.Printf("Found: %v", v)
		found = true
		break
	}

	if !found {
		return file, fmt.Errorf("Could not load %v (%v)", fName, subDir)
	}
	return file, nil

}
