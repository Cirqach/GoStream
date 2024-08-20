package logger

import (
	"log"
	"runtime"
	"strings"
)

// LogError function    logging error message
func LogError(funcName, message string) {
	funcName = cleanFuncName(funcName)
	log.Println("ERROR in func " + funcName + " : " + message)
}

// LogMessage function    logging information message
func LogMessage(funcName, message string) {
	funcName = cleanFuncName(funcName)
	log.Println("MESSAGE from " + funcName + " : " + message)
}

// GetFuncName function    return name of function where it called
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
	funcName = cleanFuncName(funcName)
	log.Fatal("FATAL: in func " + funcName + " : " + message)
}

func cleanFuncName(funcName string) string {
	if strings.Contains(funcName, "func") {
		names := strings.Split(funcName, ".")
		return names[len(names)-2]
	}
	return funcName[strings.LastIndex(funcName, ".")+1:]
}
