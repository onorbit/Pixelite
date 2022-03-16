package log

import (
	"fmt"
	"time"
)

func Info(format string, params ...interface{}) {
	logStr := fmt.Sprintf("INFO "+format, params...)
	printLog(logStr)
}

func Warn(format string, params ...interface{}) {
	logStr := fmt.Sprintf("WARN "+format, params...)
	printLog(logStr)
}

func Error(format string, params ...interface{}) {
	logStr := fmt.Sprintf("ERROR "+format, params...)
	printLog(logStr)
}

func printLog(log string) {
	timeStr := time.Now().Format("2006-01-02T15:04:05")
	fmt.Printf("%s %s\n", timeStr, log)
}
