package daas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

var (
	logNone  *log.Logger
	logDebug *log.Logger
	logInfo  *log.Logger
	logWarn  *log.Logger
	logError *log.Logger
)

func init() {
	IsPrd := os.Getenv(ENV) == "production"

	logNone = log.New(os.Stdout, "", 0)
	logDebug = log.New(os.Stdout, "[Debug] ", 0)
	logInfo = log.New(os.Stdout, " [Info] ", 0)
	logError = log.New(os.Stdout, "[Error] ", 0)
	logWarn = log.New(os.Stdout, " [Warn] ", 0)

	logNone.SetOutput(ioutil.Discard)
	logDebug.SetOutput(ioutil.Discard)

	if !IsPrd {
		logNone.SetOutput(os.Stdout)
		logDebug.SetOutput(os.Stdout)
		logDebug.SetFlags(log.Ltime)
		logInfo.SetFlags(log.Ltime)
		logWarn.SetFlags(log.Ltime)
		logError.SetFlags(log.Ltime)
	}
}

func DisableOutput() {
	logNone.SetOutput(ioutil.Discard)
	logDebug.SetOutput(ioutil.Discard)
	logInfo.SetOutput(ioutil.Discard)
	logWarn.SetOutput(ioutil.Discard)
	logError.SetOutput(ioutil.Discard)
}

func DebugNewline() {
	logNone.Print("\n")
}

func Debug(v ...interface{}) {
	logDebug.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	logDebug.Printf(format, v...)
}
func Debugv(v ...interface{}) {
	for i := range v {
		logDebug.Println("Inspect :", reflect.TypeOf(v[i]).String())
		result, _ := json.MarshalIndent(v[i], "", "  ")
		logNone.Print(string(result))
	}
	logNone.Println("")
}

func Info(v ...interface{}) {
	logInfo.Println(v...)
}

func Infof(format string, v ...interface{}) {
	logInfo.Printf(format, v...)
}

func Warn(v ...interface{}) {
	logWarn.Println(v...)
}

func Warnf(format string, v ...interface{}) {
	logWarn.Printf(format, v...)
}

func Error(v ...interface{}) {
	logError.Println(v...)
}
func Errorf(format string, v ...interface{}) {
	logError.Printf(format, v...)
}

func HttpErrorf(code int, err error) *HTTP {
	return &HTTP{Code: code, Error: err}
}
func HttpErrorPrint(code int, format string, v ...interface{}) *HTTP {
	return &HTTP{Code: code, Error: fmt.Errorf(format, v...)}
}

func HttpErrorPrintf(code int, v ...interface{}) *HTTP {
	return &HTTP{Code: code, Error: fmt.Errorf("%s", v...)}
}

func Fatal(v ...interface{}) {
	logError.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	logError.Fatalf(format, v...)
}
