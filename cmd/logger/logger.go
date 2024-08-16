package logger

import (
	"log"
	"runtime"
)

func LogError(funcName, message string) {
	log.Println("ERROR: in func " + " : " + message)
}

func LogMessage(funcName, message string) {
	log.Println("MESSAGE " + ": " + message)
}

func GetFuncName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown"
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return "unknown"
	}
	return f.Name()
}
