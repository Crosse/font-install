//Package gosimplelogger implements a simple logging package that
//includes the notion of logging levels.
package gosimplelogger

import (
	"fmt"
	"log"
	"os"
)

const (
	// Log Levels.
	LogPanic = iota
	LogFatal
	LogError
	LogInfo
	LogVerbose
	LogDebug
)

var (
	// LogLevel is the highest level that will be logged.  The
	// default, LogInfo, means that Info, Error, Fatal, and Panic
	// log messages will be emitted, but Verbose and Debug messages
	// will not.
	LogLevel      = LogInfo
	verboseLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
)

func init() {
	verboseLogger = log.New(os.Stdout, "", 0)
	infoLogger = log.New(os.Stdout, "", 0)
	errorLogger = log.New(os.Stderr, "", 0)
}

func logf(logger *log.Logger, logLevel int, format string, v ...interface{}) {
	if LogLevel >= logLevel {
		logger.Println(fmt.Sprintf(format, v...))
	}
}

func logln(logger *log.Logger, logLevel int, v ...interface{}) {
	logf(logger, logLevel, fmt.Sprint(v...))
}

// Debug prints to the verbose logger, but only if LogLevel >= LogDebug.
// Arguments are handled in the manner of fmt.Println.
func Debug(v ...interface{}) {
	logln(verboseLogger, LogDebug, fmt.Sprint(v...))
}

// Debugf prints to the verbose logger, but only if LogLevel >=
// LogDebug.  Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, v ...interface{}) {
	logf(verboseLogger, LogDebug, format, v...)
}

// Verbose prints to the verbose logger, but only if LogLevel >=
// LogVerbose.  Arguments are handled in the manner of fmt.Println.
func Verbose(v ...interface{}) {
	logln(verboseLogger, LogVerbose, fmt.Sprint(v...))
}

// Verbosef prints to the verbose logger, but only if LogLevel >=
// LogVerbose.  Arguments are handled in the manner of fmt.Printf.
func Verbosef(format string, v ...interface{}) {
	logf(verboseLogger, LogVerbose, format, v...)
}

// Info prints to the info logger, but only if LogLevel >= LogInfo.
// Arguments are handled in the manner of fmt.Println.
func Info(v ...interface{}) {
	logln(infoLogger, LogInfo, v...)
}

// Infof prints to the info logger, but only if LogLevel >= LogInfo.
// Arguments are handled in the manner of fmt.Printf.
func Infof(format string, v ...interface{}) {
	logf(infoLogger, LogInfo, format, v...)
}

// Println prints to the info logger, but only if LogLevel >= LogInfo.
// Arguments are handled in the manner of fmt.Println.  Println simply
// calls Info() and is included for compatibility with Go's log.Println.
func Println(v ...interface{}) {
	Info(v...)
}

// Printf prints to the info logger, but only if LogLevel >= LogInfo.
// Arguments are handled in the manner of fmt.Printf.  Printf simply
// calls Infof() and is included for compatibility with Go's log.Printf.
func Printf(format string, v ...interface{}) {
	Infof(format, v...)
}

// Error prints to the error logger (which logs to standard error by
// default), but only if LogLevel >= LogError.  Arguments are handled in
// the manner of fmt.Println.
func Error(v ...interface{}) {
	logln(errorLogger, LogError, v...)
}

// Errorf prints to the error logger (which logs to standard error by
// default), but only if LogLevel >= LogError.  Arguments are handled in
// the manner of fmt.Printf.
func Errorf(format string, v ...interface{}) {
	logf(errorLogger, LogError, format, v...)
}

// Fatal prints to the error logger (which logs to standard error by
// default), but only if LogLevel >= LogFatal.  Arguments are handled in
// the manner of fmt.Println. Fatal is equivalent to Error() followed by
// a call to os.Exit(1).
func Fatal(v ...interface{}) {
	logln(errorLogger, LogFatal, v...)
	os.Exit(1)
}

// Fatalf prints to the error logger (which logs to standard error by
// default), but only if LogLevel >= LogFatal.  Arguments are handled in
// the manner of fmt.Printf. Fatal is equivalent to Error() followed by
// a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	logf(errorLogger, LogFatal, format, v...)
	os.Exit(1)
}

// Panic prints to the error logger (which logs to standard error by
// default), but only if LogLevel >= LogPanic.  Arguments are handled in
// the manner of fmt.Println. Panic is equivalent to calling Error()
// followed by a call to panic().
func Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	logln(errorLogger, LogPanic, s)
	panic(s)
}

// Panicf prints to the error logger (which logs to standard error by
// default), but only if LogLevel >= LogPanic.  Arguments are handled in
// the manner of fmt.Printf. Panicf is equivalent to calling Errorf()
// followed by a call to panic().
func Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	logln(errorLogger, LogPanic, s)
	panic(s)
}
