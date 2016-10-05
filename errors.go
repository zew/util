package util

import (
	"os"
	"strings"
	"sync"

	"github.com/zew/logx"
)

var occurred = map[string]int{}
var l sync.Mutex

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
		str := strings.Join(logx.StackTrace(2, 3, 2), "\n\t")
		logx.Printf("\n\t%s\n", str)
		os.Exit(1)
	}
}
