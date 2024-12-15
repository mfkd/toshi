package logger

import (
	"log"
	"os"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	logLevel = LevelInfo // Shared log level
	logger   = log.New(os.Stdout, "", log.LstdFlags)
)

// Configure sets up the global log level and output.
func Configure(level int, output *os.File) {
	logLevel = level
	if output != nil {
		logger.SetOutput(output)
	}
}

// SetLogLevel updates the global log level.
func SetLogLevel(level int) {
	logLevel = level
}

// SetOutput updates the global output destination.
func SetOutput(output *os.File) {
	if output != nil {
		logger.SetOutput(output)
	}
}

// Debug logs a plain or formatted DEBUG-level message.
func Debug(v ...interface{}) {
	if logLevel <= LevelDebug {
		logger.SetPrefix("[DEBUG] ")
		logger.Println(v...)
	}
}

// Debugf logs a formatted DEBUG-level message.
func Debugf(format string, v ...interface{}) {
	if logLevel <= LevelDebug {
		logger.SetPrefix("[DEBUG] ")
		logger.Printf(format, v...)
	}
}

// Info logs a plain or formatted INFO-level message.
func Info(v ...interface{}) {
	if logLevel <= LevelInfo {
		logger.SetPrefix("[INFO] ")
		logger.Println(v...)
	}
}

// Infof logs a formatted INFO-level message.
func Infof(format string, v ...interface{}) {
	if logLevel <= LevelInfo {
		logger.SetPrefix("[INFO] ")
		logger.Printf(format, v...)
	}
}

// Warn logs a plain or formatted WARN-level message.
func Warn(v ...interface{}) {
	if logLevel <= LevelWarn {
		logger.SetPrefix("[WARN] ")
		logger.Println(v...)
	}
}

// Warnf logs a formatted WARN-level message.
func Warnf(format string, v ...interface{}) {
	if logLevel <= LevelWarn {
		logger.SetPrefix("[WARN] ")
		logger.Printf(format, v...)
	}
}

// Error logs a plain or formatted ERROR-level message.
func Error(v ...interface{}) {
	if logLevel <= LevelError {
		logger.SetPrefix("[ERROR] ")
		logger.Println(v...)
	}
}

// Errorf logs a formatted ERROR-level message.
func Errorf(format string, v ...interface{}) {
	if logLevel <= LevelError {
		logger.SetPrefix("[ERROR] ")
		logger.Printf(format, v...)
	}
}

// Fatal logs an ERROR-level message and exits the program.
func Fatal(v ...interface{}) {
	logger.SetPrefix("[FATAL] ")
	logger.Println(v...)
	os.Exit(1)
}

// Fatalf logs a formatted ERROR-level message and exits the program.
func Fatalf(format string, v ...interface{}) {
	logger.SetPrefix("[FATAL] ")
	logger.Printf(format, v...)
	os.Exit(1)
}
