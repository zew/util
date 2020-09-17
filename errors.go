package util

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var occurred = map[string]int{}
var l sync.Mutex

// This should be moved to log package
func CheckErr(err error, tolerate ...string) {
	if err != nil {
		errStr := strings.ToLower(err.Error())
		for _, tol := range tolerate {
			tol = strings.ToLower(tol)
			if strings.Contains(errStr, tol) {
				l.Lock()
				occurred[err.Error()]++
				l.Unlock()
				if occurred[err.Error()] < 2 {
					log.Printf("tolerated error: %v", err)
				}
				return
			}
		}
		log.Printf("%v", err)
		str := StackTraceStr(1, 4, 3)
		log.Printf("\n\t%s\n", str)
		os.Exit(1)
	}
}

// This should be moved to log package
func BubbleUp(err error, tolerate ...string) {
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
}

func SqlAlreadyExists(e error) bool {
	if e != nil {
		canTolerate := []string{
			"Duplicate entry",
			"UNIQUE constraint failed",
			"duplicate key value violates",
		}
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
