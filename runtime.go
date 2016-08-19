package util

import "runtime"

func ThisFunc() *runtime.Func {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc)
}
func ThisFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
