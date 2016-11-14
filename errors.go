package util

import (
	"os"
	"strings"
	"sync"

	"github.com/zew/logx"
)

var occurred = map[string]int{}
var l sync.Mutex

// This should be moved to log package
func CheckErr(err error, tolerate ...string) {
	defer logx.SL().Incr().Decr()
	if err != nil {
		errStr := strings.ToLower(err.Error())
		for _, tol := range tolerate {
			tol = strings.ToLower(tol)
			if strings.Contains(errStr, tol) {
				l.Lock()
				occurred[err.Error()]++
				l.Unlock()
				if occurred[err.Error()] < 2 {
					logx.Printf("tolerated error: %v", err)
				}
				return
			}
		}
		logx.Printf("%v", err)
		str := logx.SPrintStackTrace(1, 4, 3)
		logx.Printf("\n\t%s\n", str)
		os.Exit(1)
	}
}

func SqlAlreadyExists(e error) bool {
	if e != nil {
		canTolerate := []string{"Duplicate entry", "UNIQUE constraint failed"}
		errStr := strings.ToLower(e.Error())
		for _, tol := range canTolerate {
			tol = strings.ToLower(tol)
			if strings.Contains(errStr, tol) {
				return true
			}
		}
	}
	return false
}
