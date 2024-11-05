package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/wakatime/wakatime-cli/pkg/version"

	"github.com/sirupsen/logrus"
	jww "github.com/spf13/jwalterweatherman"
)

// Logger is the log entry.
type Logger struct {
	entry             *logrus.Entry
	metrics           bool
	sendDiagsOnErrors bool
	verbose           bool
}

// New creates a new Logger.
func New(verbose, sendDiagsOnErrors, metrics bool) *Logger {
	logger := &Logger{
		entry:             new(),
		metrics:           metrics,
		sendDiagsOnErrors: sendDiagsOnErrors,
		verbose:           verbose,
	}

	logger.SetVerbose(verbose)

	return logger
}

func new() *logrus.Entry {
	entry := logrus.NewEntry(&logrus.Logger{
		Out: os.Stdout,
		Formatter: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime: "now",
				logrus.FieldKeyFile: "caller",
				logrus.FieldKeyMsg:  "message",
			},
			DisableHTMLEscape: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				// Simplifies function description by removing dangling func name from it.
				lastSlash := strings.LastIndexByte(f.Function, '/')
				if lastSlash < 0 {
					lastSlash = 0
				}
				parts := strings.Split(f.Function[lastSlash+1:], ".")

				// Simplifies file path by removing base path from it.
				lastPath := strings.LastIndex(f.File, "wakatime-cli/")
				if lastPath < 0 {
					lastPath = 0
				}
				file := f.File[lastPath+13:]

				return fmt.Sprintf("%s.%s", parts[0], parts[1]),
					fmt.Sprintf("%s:%d", file, f.Line)
			},
		},
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: true,
	})
	entry.Data["version"] = version.Version
	entry.Data["os/arch"] = fmt.Sprintf("%s/%s", version.OS, version.Arch)

	return entry
}

// IsMetricsEnabled returns true if it should collect metrics.
func (l *Logger) IsMetricsEnabled() bool {
	return l.metrics
}

// IsVerboseEnabled returns true if debug is enabled.
func (l *Logger) IsVerboseEnabled() bool {
	return l.verbose
}

// Output returns the current log output.
func (l *Logger) Output() io.Writer {
	return l.entry.Logger.Out
}

// SendDiagsOnErrors returns true if diagnostics should be sent on errors.
func (l *Logger) SendDiagsOnErrors() bool {
	return l.sendDiagsOnErrors
}

// SetOutput defines sets the log output to io.Writer.
func (l *Logger) SetOutput(w io.Writer) {
	l.entry.Logger.Out = w
}

// SetVerbose sets log level to debug if enabled.
func (l *Logger) SetVerbose(verbose bool) {
	if verbose {
		l.entry.Logger.SetLevel(logrus.DebugLevel)
	} else {
		l.entry.Logger.SetLevel(logrus.InfoLevel)
	}
}

// Flush flushes the log output and closes the file.
func (l *Logger) Flush() {
	if file, ok := l.entry.Logger.Out.(*os.File); ok {
		if err := file.Sync(); err != nil {
			l.entry.Debugf("failed to flush log file: %s", err)
		}

		if err := file.Close(); err != nil {
			l.entry.Debugf("failed to close log file: %s", err)
		}
	}
}

// SetJww sets jww log when debug enabled.
func SetJww(verbose bool, w io.Writer) {
	if verbose {
		jww.SetLogThreshold(jww.LevelDebug)
		jww.SetStdoutThreshold(jww.LevelDebug)

		jww.SetLogOutput(w)
		jww.SetStdoutOutput(w)
	}
}

// Debugf logs a message at level Debug.
func (l *Logger) Debugf(format string, args ...any) {
	l.entry.Debugf(format, args...)
}

// Infof logs a message at level Info.
func (l *Logger) Infof(format string, args ...any) {
	l.entry.Infof(format, args...)
}

// Warnf logs a message at level Warn.
func (l *Logger) Warnf(format string, args ...any) {
	l.entry.Warnf(format, args...)
}

// Errorf logs a message at level Error.
func (l *Logger) Errorf(format string, args ...any) {
	l.entry.Errorf(format, args...)
}

// Fatalf logs a message at level Fatal then the process will exit with status set to 1.
func (l *Logger) Fatalf(format string, args ...any) {
	l.entry.Fatalf(format, args...)
}

// Debugln logs a message at level Debug.
func (l *Logger) Debugln(args ...any) {
	l.entry.Debugln(args...)
}

// Infoln logs a message at level Info.
func (l *Logger) Infoln(args ...any) {
	l.entry.Infoln(args...)
}

// Warnln logs a message at level Warn.
func (l *Logger) Warnln(args ...any) {
	l.entry.Warnln(args...)
}

// Errorln logs a message at level Error.
func (l *Logger) Errorln(args ...any) {
	l.entry.Errorln(args...)
}

// Fatalln logs a message at level Fatal then the process will exit with status set to 1.
func (l *Logger) Fatalln(args ...any) {
	l.entry.Fatalln(args...)
}

// WithField adds a single field to the Logger.
func (l *Logger) WithField(key string, value any) {
	l.entry.Data[key] = value
}
