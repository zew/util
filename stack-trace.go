package util

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// We dont want 20 leading directories of a source file.
// But the filename alone is not enough.
// "main.go" does not help.
func leadDirsBeforeSourceFile(path string, dirsBeforeSourceFile int) string {
	rump := path // init
	dirs := make([]string, 0, dirsBeforeSourceFile)
	for i := 0; i < dirsBeforeSourceFile; i++ {
		rump = filepath.Dir(rump)
		dir := filepath.Base(rump)
		dirs = append([]string{dir}, dirs...)
	}
	lastDirs := filepath.Join(dirs...)
	lastDirs = filepath.Join(lastDirs, filepath.Base(path))
	return lastDirs
}

// StackTrace returns the call stack as slice of strings
// First  arg => level init
// Second arg => levels up
// Third  arg => dirs of before source file
func StackTrace(args ...int) []string {

	var (
		lvlInit              = 2 // One for this func, one since direct caller is already logged in sourceLocationPrefix()
		lvlsUp               = lvlInit + 2
		dirsBeforeSourceFile = 2 // How many dirs are shown before the source file.
	)

	if len(args) > 0 {
		lvlInit += args[0]
	}
	if len(args) > 1 {
		lvlsUp = args[1]
	}
	if len(args) > 2 {
		dirsBeforeSourceFile = args[2]
	}

	lines := make([]string, lvlsUp)
	for i := 0; i < lvlsUp; i++ {

		_, file, line, _ := runtime.Caller(i + lvlInit)
		if line == 0 && file == "." {
			break
		}
		file = leadDirsBeforeSourceFile(file, dirsBeforeSourceFile)

		lines[i] = fmt.Sprintf("%-42s:%d", file, line)
	}
	return lines
}

// StackTraceStr returns StackTrace as a string with new lines
func StackTraceStr(args ...int) string {
	lines := StackTrace(args...)
	return "\n\t" + strings.Join(lines, "\n\t")
}

// StackTraceHTML returns StackTrace as a string embeddable into HTML
func StackTraceHTML(args ...int) string {
	lines := StackTrace(args...)
	s := "\n\t" + strings.Join(lines, "\n\t")
	s = fmt.Sprintf("\n<pre style=\"margin: 0;font-family: 'Courier New', Courier, monospace; font-size: 0.8rem;\">%v</pre>\n", s)
	return s
}
