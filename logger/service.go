package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	logger             *logrus.Logger
	getLogMutex        sync.Mutex
	defaultLogFileName = "daily.log"
	defaultLevel       = "debug"
)

type Config struct {
	Level string

	Formatter string // only support json and text

	DisableConsole bool
	Write          bool
	Path           string
	FileName       string

	MaxAge       time.Duration
	RotationTime time.Duration

	Debug bool // if set true, separate
}

func defaultConfig() *Config {
	return &Config{
		Level:          defaultLevel,
		Formatter:      "text",
		DisableConsole: false,
		Write:          false,
		Path:           os.TempDir(),
		FileName:       defaultLogFileName,
		MaxAge:         time.Duration(24) * time.Hour,
		RotationTime:   time.Duration(7*24) * time.Hour,
		Debug:          false,
	}
}

func NewService(config *Config) *logrus.Logger {
	getLogMutex.Lock()
	defer getLogMutex.Unlock()

	if config == nil {
		config = defaultConfig()
	}

	if logger != nil {
		return logger
	}

	log := logrus.New()

	// get logLevel
	level := config.Level
	if level == "" {
		level = defaultLevel
	}
	logLevel := GetLogLevel(level)

	logDir := config.Path
	if logDir == "" {
		logDir = os.TempDir()
	}

	logFileName := config.FileName
	if logFileName == "" {
		logFileName = defaultLogFileName
	}

	printLog := !config.DisableConsole

	maxAge := config.MaxAge

	rotationTime := config.RotationTime

	log.SetLevel(logLevel)

	if config.Write {
		storeLogDir := logDir

		err := os.MkdirAll(storeLogDir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("creating log file failed: %s", err.Error()))
		}

		path := filepath.Join(storeLogDir, logFileName)
		writer, err := rotatelogs.New(
			path+".%Y-%m-%d",
			rotatelogs.WithClock(rotatelogs.Local),
			rotatelogs.WithMaxAge(time.Duration(maxAge)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Hour),
		)
		if err != nil {
			panic(fmt.Sprintf("rotatelogs log failed: %s", err.Error()))
		}

		var formatter logrus.Formatter

		formatter = &logrus.TextFormatter{}
		if config.Formatter == "json" {
			formatter = &logrus.JSONFormatter{}
		}
		if config.Debug {
			log.AddHook(lfshook.NewHook(
				lfshook.WriterMap{
					logrus.DebugLevel: writer,
					logrus.InfoLevel:  writer,
					logrus.WarnLevel:  writer,
					logrus.ErrorLevel: writer,
					logrus.FatalLevel: writer,
				},
				formatter,
			))

			defaultLogFilePrex := logFileName + "."
			pathMap := lfshook.PathMap{
				logrus.DebugLevel: fmt.Sprintf("%s/%sdebug", storeLogDir, defaultLogFilePrex),
				logrus.InfoLevel:  fmt.Sprintf("%s/%sinfo", storeLogDir, defaultLogFilePrex),
				logrus.WarnLevel:  fmt.Sprintf("%s/%swarn", storeLogDir, defaultLogFilePrex),
				logrus.ErrorLevel: fmt.Sprintf("%s/%serror", storeLogDir, defaultLogFilePrex),
				logrus.FatalLevel: fmt.Sprintf("%s/%sfatal", storeLogDir, defaultLogFilePrex),
			}
			log.AddHook(lfshook.NewHook(
				pathMap,
				formatter,
			))
		} else {
			log.Out = writer
		}

	} else {
		if printLog {
			log.Out = os.Stdout
		}

	}
	return log
}

func FormatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {
			//format string
		} else {
			//do not contain format char
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}

// Debug wrapper Debug logger
func (l *Logger) Debug(f interface{}, args ...interface{}) {
	l.log.Debug(FormatLog(f, args...))
}

// Info wrapper Info logger
func (l *Logger) Info(f interface{}, args ...interface{}) {
	l.log.Info(FormatLog(f, args...))
}

// Warn wrapper Warn logger
func (l *Logger) Warn(f interface{}, args ...interface{}) {
	l.log.Warn(FormatLog(f, args...))
}

// Printf wrapper Printf logger
func (l *Logger) Printf(f interface{}, args ...interface{}) {
	l.log.Print(FormatLog(f, args...))
}

// Panic wrapper Panic logger
func (l *Logger) Panic(f interface{}, args ...interface{}) {
	l.log.Panic(FormatLog(f, args...))
}

// Fatal wrapper Fatal logger
func (l *Logger) Fatal(f interface{}, args ...interface{}) {
	l.log.Fatal(FormatLog(f, args...))
}

// Error wrapper Error logger
func (l *Logger) Error(f interface{}, args ...interface{}) {
	l.log.Error(FormatLog(f, args...))
}

// Debugln wrapper Debugln logger
func (l *Logger) Debugln(v ...interface{}) {
	l.log.Debug(fmt.Sprintln(v...))
}

// Infoln wrapper Infoln logger
func (l *Logger) Infoln(args ...interface{}) {
	l.log.Info(fmt.Sprintln(args...))
}

// Warnln wrapper Warnln logger
func (l *Logger) Warnln(args ...interface{}) {
	l.log.Warn(fmt.Sprintln(args...))
}

// Printfln wrapper Printfln logger
func (l *Logger) Printfln(args ...interface{}) {
	l.log.Print(fmt.Sprintln(args...))
}

// Panicln wrapper Panicln logger
func (l *Logger) Panicln(args ...interface{}) {
	l.log.Panic(fmt.Sprintln(args...))
}

// Fatalln wrapper Fatalln logger
func (l *Logger) Fatalln(args ...interface{}) {
	l.log.Fatal(fmt.Sprintln(args...))
}

// Errorln wrapper Errorln logger
func (l *Logger) Errorln(args ...interface{}) {
	l.log.Error(fmt.Sprintln(args...))
}
