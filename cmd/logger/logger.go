package logger

import (
	"log"
	"reflect"
	"runtime"
)

func LogError(funcName, message string) {
	log.Println("ERROR: in func " + funcName + " : " + message)
}

func LogMessage(funcName, message string) {
	log.Println("MESSAGE " + funcName + ": " + message)
}

func GetFuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
