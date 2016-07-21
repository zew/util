package util

import (
	"os"
	"strings"

	"github.com/zew/logx"
)

func CheckErr(err error, tolerate ...string) {
	defer logx.SL().Incr().Decr()
	if err != nil {
		errStr := strings.ToLower(err.Error())
		for _, tol := range tolerate {
			tol = strings.ToLower(tol)
			if strings.Contains(errStr, tol) {
				logx.Printf("tolerated error: %v", err)
				return
			}
		}
		logx.Printf("%v", err)
		str := strings.Join(logx.StackTrace(2, 3, 2), "\n")
		logx.Printf("\nStacktrace: \n%s", str)
		os.Exit(1)
	}
}
