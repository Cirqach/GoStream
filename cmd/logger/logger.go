package logger

import (
	"log"
	"runtime"
)

func LogError(funcName, message string) {
	log.Println("ERROR: in func " + funcName + " : " + message)
}

func LogMessage(funcName, message string) {
	log.Println("MESSAGE: " + funcName + " : " + message)
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

func Fatal(funcName, message string) {
	log.Fatal("FATAL: in func " + funcName + " : " + message)
}
