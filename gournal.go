/*
Package gournal (pronounced "Journal") is a Context-aware logging framework.

Gournal introduces the Google Context type (https://blog.golang.org/context) as
a first-class parameter to all common log functions such as Info, Debug, etc.

Instead of being Yet Another Go Log library, Gournal actually takes its
inspiration from the Simple Logging Facade for Java (SLF4J). Gournal is not
attempting to replace anyone's favorite logger, rather existing logging
frameworks such as Logrus, Zap, etc. can easily participate as a Gournal
Appender.

For more information on Gournal's features or how to use it, please refer
to the project's README file or https://github.com/emccode/gournal.
*/
package gournal

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

var debug, _ = strconv.ParseBool(os.Getenv("GOURNAL_DEBUG"))

var (
	// ErrorKey defines the key when adding errors using WithError.
	ErrorKey = "error"

	// DefaultLevel is used when a Level is not present in a Context.
	DefaultLevel = ErrorLevel

	// DefaultAppender is used when an Appender is not present in a Context.
	DefaultAppender Appender

	// DefaultContext is used when a log method is invoked with a nil Context.
	DefaultContext = context.Background()
)

// Key is the key type used for storing gournal-related objects in a
// Context.
type Key uint8

const (
	// AppenderKey is the key for storing the implementation of the Appender
	// interface in a Context.
	AppenderKey Key = iota

	// LevelKey is the key for storing the log Level constant in a Context.
	LevelKey

	// FieldsKey is the key used to store/retrieve the Context-specific field
	// data to append with each log entry. Three different types of data are
	// inspected for this context key:
	//
	//     * map[string]interface{}
	//
	//     * func() map[string]interface{}
	//
	//     * func(ctx context.Cotext,
	//            fields map[string]interface{},
	//            args ...interface{}) map[string]interface{}
	FieldsKey
)

// Level is a log level.
type Level uint8

// These are the different logging levels.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic
	// with the message passed to Debug, Info, ...
	PanicLevel Level = iota

	// FatalLevel level. Logs and then calls os.Exit(1). It will exit even
	// if the logging level is set to Panic.
	FatalLevel

	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel

	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel

	// InfoLevel level. General operational entries about what's going on
	// inside the application.
	InfoLevel

	// DebugLevel level. Usually only enabled when debugging. Very verbose
	// logging.
	DebugLevel
)

var (
	lvlStrsToVals = map[string]Level{
		"DEBUG":   DebugLevel,
		"INFO":    InfoLevel,
		"WARN":    WarnLevel,
		"WARNING": WarnLevel,
		"ERROR":   ErrorLevel,
		"FATAL":   FatalLevel,
		"PANIC":   PanicLevel,
	}

	lvlValsToStrs = map[Level]string{
		DebugLevel: "DEBUG",
		InfoLevel:  "INFO",
		WarnLevel:  "WARN",
		ErrorLevel: "ERROR",
		FatalLevel: "FATAL",
		PanicLevel: "PANIC",
	}
)

// String returns string representation of a Level.
func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	case PanicLevel:
		return "PANIC"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel parses a string and returns its constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}
	return 0, fmt.Errorf("invalid level: %v", lvl)
}

// Logger provides backwards-compatibility for code that does not yet use
// context-aware logging.
type Logger interface {

	// Debug emits a log entry at the DEBUG level.
	Debug(args ...interface{})
	// Debugf is an alias for Debug.
	Debugf(format string, args ...interface{})
	// Debugln is an alias for Debug.
	Debugln(args ...interface{})

	// Info emits a log entry at the INFO level.
	Info(args ...interface{})
	// Infof is an alias for Info.
	Infof(format string, args ...interface{})
	// Infoln is an alias for Info.
	Infoln(args ...interface{})

	// Print emits a log entry at the INFO level.
	Print(args ...interface{})
	// Printf is an alias for Print.
	Printf(format string, args ...interface{})
	// Println is an alias for Print.
	Println(args ...interface{})

	// Warn emits a log entry at the WARN level.
	Warn(args ...interface{})
	// Warnf is an alias for Warn.
	Warnf(format string, args ...interface{})
	// Warnln is an alias for Warn.
	Warnln(args ...interface{})

	// Error emits a log entry at the ERROR level.
	Error(args ...interface{})
	// Errorf is an alias for Error.
	Errorf(format string, args ...interface{})
	// Errorln is an alias for Error.
	Errorln(args ...interface{})

	// Fatal emits a log entry at the FATAL level.
	Fatal(args ...interface{})
	// Fatalf is an alias for Fatal.
	Fatalf(format string, args ...interface{})
	// Fatalln is an alias for Fatal.
	Fatalln(args ...interface{})

	// Panic emits a log entry at the PANIC level.
	Panic(args ...interface{})
	// Panicf is an alias for Panic.
	Panicf(format string, args ...interface{})
	// Panicln is an alias for Panic.
	Panicln(args ...interface{})
}

// New returns a Logger for the provided context.
func New(ctx context.Context) Logger {
	return &logger{ctx}
}

type logger struct {
	ctx context.Context
}

func (l *logger) Debug(args ...interface{}) {
	Debug(l.ctx, args...)
}
func (l *logger) Debugf(format string, args ...interface{}) {
	l.Debug(append([]interface{}{format}, args...)...)
}
func (l *logger) Debugln(args ...interface{}) {
	l.Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	Info(l.ctx, args...)
}
func (l *logger) Infof(format string, args ...interface{}) {
	l.Info(append([]interface{}{format}, args...)...)
}
func (l *logger) Infoln(args ...interface{}) {
	l.Info(args...)
}

func (l *logger) Print(args ...interface{}) {
	Print(l.ctx, args...)
}
func (l *logger) Printf(format string, args ...interface{}) {
	l.Print(append([]interface{}{format}, args...)...)
}
func (l *logger) Println(args ...interface{}) {
	l.Print(args...)
}

func (l *logger) Warn(args ...interface{}) {
	Warn(l.ctx, args...)
}
func (l *logger) Warnf(format string, args ...interface{}) {
	l.Warn(append([]interface{}{format}, args...)...)
}
func (l *logger) Warnln(args ...interface{}) {
	l.Warn(args...)
}

func (l *logger) Error(args ...interface{}) {
	Error(l.ctx, args...)
}
func (l *logger) Errorf(format string, args ...interface{}) {
	l.Error(append([]interface{}{format}, args...)...)
}
func (l *logger) Errorln(args ...interface{}) {
	l.Error(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	Fatal(l.ctx, args...)
}
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Error(append([]interface{}{format}, args...)...)
}
func (l *logger) Fatalln(args ...interface{}) {
	l.Fatal(args...)
}

func (l *logger) Panic(args ...interface{}) {
	Panic(l.ctx, args...)
}
func (l *logger) Panicf(format string, args ...interface{}) {
	l.Panic(append([]interface{}{format}, args...)...)
}
func (l *logger) Panicln(args ...interface{}) {
	l.Panic(args...)
}

// Entry is the interface for types that contain information to be emmitted
// to a log appender.
type Entry interface {

	// WithField adds a single field to the Entry. The provided key will
	// override an existing, equivalent key in the Entry.
	WithField(key string, value interface{}) Entry

	// WithFields adds a map to the Entry. Keys in the provided map will
	// override existing, equivalent keys in the Entry.
	WithFields(fields map[string]interface{}) Entry

	// WithError adds the provided error to the Entry using the ErrorKey value
	// as the key.
	WithError(err error) Entry

	// Debug emits a log entry at the DEBUG level.
	Debug(ctx context.Context, args ...interface{})

	// Info emits a log entry at the INFO level.
	Info(ctx context.Context, args ...interface{})

	// Print emits a log entry at the INFO level.
	Print(ctx context.Context, args ...interface{})

	// Warn emits a log entry at the WARN level.
	Warn(ctx context.Context, args ...interface{})

	// Error emits a log entry at the ERROR level.
	Error(ctx context.Context, args ...interface{})

	// Fatal emits a log entry at the FATAL level.
	Fatal(ctx context.Context, args ...interface{})

	// Panic emits a log entry at the PANIC level.
	Panic(ctx context.Context, args ...interface{})
}

// Appender is the interface that must be implemented by the logging frameworks
// which are members of the Gournal facade.
type Appender interface {

	// Append is implemented by logging frameworks to accept the log entry
	// at the provided level, its message, and its associated field data.
	Append(
		ctx context.Context,
		lvl Level,
		fields map[string]interface{},
		msg string)
}

// WithField adds a single field to the Entry. The provided key will override
// an existing, equivalent key in the Entry.
func WithField(key string, value interface{}) Entry {
	return &entry{map[string]interface{}{key: value}}
}

// WithFields adds a map to the Entry. Keys in the provided map will override
// existing, equivalent keys in the Entry.
func WithFields(fields map[string]interface{}) Entry {
	return &entry{fields}
}

// WithError adds the provided error to the Entry using the ErrorKey value
// as the key.
func WithError(err error) Entry {
	return &entry{map[string]interface{}{ErrorKey: err.Error()}}
}

// Debug emits a log entry at the DEBUG level.
func Debug(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, DebugLevel, nil, args...)
}

// Info emits a log entry at the INFO level.
func Info(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, InfoLevel, nil, args...)
}

// Print emits a log entry at the INFO level.
func Print(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, InfoLevel, nil, args...)
}

// Warn emits a log entry at the WARN level.
func Warn(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, WarnLevel, nil, args...)
}

// Error emits a log entry at the ERROR level.
func Error(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, ErrorLevel, nil, args...)
}

// Fatal emits a log entry at the FATAL level.
func Fatal(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, FatalLevel, nil, args...)
}

// Panic emits a log entry at the PANIC level.
func Panic(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, PanicLevel, nil, args...)
}

func swapFields(appendFields, ctxFields *map[string]interface{}) {
	if len(*ctxFields) == 0 {
		return
	}

	if len(*appendFields) == 0 {
		*appendFields = *ctxFields
		return
	}

	for k, v := range *ctxFields {
		(*appendFields)[k] = v
	}
}

var defaultCtx = context.Background()

func sendToAppender(
	ctx context.Context,
	lvl Level,
	fields map[string]interface{},
	args ...interface{}) {

	if ctx == nil {
		ctx = defaultCtx
	}

	// do not append if the provided log level is less than that of the
	// provided context's log level
	if getLevel(ctx) < lvl {
		return
	}

	// do not proceed without an appender
	a := getAppender(ctx)

	// grab any of the context fields to append alongside each new log entry
	switch tv := ctx.Value(FieldsKey).(type) {
	case map[string]interface{}:
		swapFields(&fields, &tv)
	case func() map[string]interface{}:
		ctxFields := tv()
		swapFields(&fields, &ctxFields)
	case func(
		ctx context.Context,
		lvl Level,
		fields map[string]interface{},
		args ...interface{}) map[string]interface{}:

		ctxFields := tv(ctx, lvl, fields, args...)
		swapFields(&fields, &ctxFields)
	}

	if len(args) == 0 {
		traceAppend(a, ctx, lvl, fields, "")
		return
	}

	var (
		arg0IsString bool
		msg          string
	)

	switch tv := args[0].(type) {
	case string:
		msg = tv
		arg0IsString = true
	case error:
		msg = tv.Error()
	case fmt.Stringer:
		msg = tv.String()
	default:
		msg = fmt.Sprintf("%v", tv)
	}

	if len(args) < 2 {
		traceAppend(a, ctx, lvl, fields, msg)
		return
	}

	if arg0IsString {
		msg = fmt.Sprintf(msg, args[1:]...)
	} else {
		w := &bytes.Buffer{}
		fmt.Fprint(w, args[1:]...)
		msg = w.String()
	}

	traceAppend(a, ctx, lvl, fields, msg)
}

func traceAppend(
	a Appender,
	ctx context.Context,
	lvl Level,
	fields map[string]interface{},
	msg string) {

	if debug {
		fmt.Fprintf(os.Stderr,
			"GOURNAL: append: a=%T, lvl=%v, fields=%v, msg=%v\n",
			a, lvl, fields, msg)
	}

	a.Append(ctx, lvl, fields, msg)
}

func getLevel(ctx context.Context) Level {
	if v, ok := ctx.Value(LevelKey).(Level); ok {
		return v
	}
	return DefaultLevel
}

func getAppender(ctx context.Context) Appender {
	if ctx == defaultCtx {
		return DefaultAppender
	}

	if v, ok := ctx.Value(AppenderKey).(Appender); ok {
		return v
	}

	return DefaultAppender
}

type entry struct {
	fields map[string]interface{}
}

func (e *entry) WithField(key string, value interface{}) Entry {
	e.fields[key] = value
	return e
}
func (e *entry) WithFields(fields map[string]interface{}) Entry {
	for k, v := range fields {
		e.fields[k] = v
	}
	return e
}
func (e *entry) WithError(err error) Entry {
	e.fields[ErrorKey] = err.Error()
	return e
}

func (e *entry) Debug(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, DebugLevel, e.fields, args...)
}

func (e *entry) Info(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, InfoLevel, e.fields, args...)
}

func (e *entry) Print(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, InfoLevel, e.fields, args...)
}

func (e *entry) Warn(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, WarnLevel, e.fields, args...)
}

func (e *entry) Error(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, ErrorLevel, e.fields, args...)
}

func (e *entry) Fatal(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, FatalLevel, e.fields, args...)
}

func (e *entry) Panic(ctx context.Context, args ...interface{}) {
	sendToAppender(ctx, PanicLevel, e.fields, args...)
}
