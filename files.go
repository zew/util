package util

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// Checks whether file exists, checks whether it is a directory.
// Dont forget to close it
// defer dir.Close()
func Directory(p string) (dir *os.File, err error) {

	dir, err = os.Open(p)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Error opening %v", p))
		return
	}

	s, err := dir.Stat()
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("Error in dir.Stat() for %v", p))
		return
	}

	if !s.IsDir() {
		err = errors.Wrap(err, fmt.Sprintf("%v is not a directory", p))
		return
	}
	return
}

// Checks whether file exists
func FileExists(p string) (bool, error) {
	_, err := os.Stat(p)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
