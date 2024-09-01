package logger

import (
	"log"
	"runtime"
	"strings"
)

var (
	colors = map[string]string{
		"red":     "\033[31m",
		"green":   "\033[32m",
		"cyan":    "\033[36m",
		"reset":   "\033[0m",
		"yellow":  "\033[33m",
		"gray":    "\033[37m",
		"magenta": "\033[35m",
	}
)

// LogError function    logging error message
func LogError(funcName, message string) {
	funcName = cleanFuncName(funcName)
	log.Println(colors["red"] + "ERROR" + colors["reset"] + " in func " + colors["magenta"] + funcName + colors["reset"] + " : " + colors["gray"] + message + colors["reset"])
}

// LogMessage function    logging information message
func LogMessage(funcName, message string) {
	funcName = cleanFuncName(funcName)
	log.Println(colors["cyan"] + "MESSAGE" + colors["reset"] + " from " + colors["magenta"] + funcName + colors["reset"] + " : " + colors["grey"] + message + colors["reset"])
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
	log.Fatal(colors["red"] + "FATAL: in func " + funcName + " : " + message + colors["reset"])
}

func cleanFuncName(funcName string) string {
	if strings.Contains(funcName, "func") {
		names := strings.Split(funcName, ".")
		return names[len(names)-2]
	}
	return funcName[strings.LastIndex(funcName, ".")+1:]
}
