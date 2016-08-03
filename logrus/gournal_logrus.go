// Package logrus provides a Logrus logger that implements the Gournal Appender
// interface.
package logrus

import (
	"io"

	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/emccode/gournal"
)

type appender struct {
	logger *logrus.Logger
}

// New returns a logrus logger that implements the Gournal Appender interface.
func New() gournal.Appender {
	return &appender{logrus.New()}
}

// NewWithOptions returns a logrus logger that implements the Gournal Appender
// interface.
func NewWithOptions(
	out io.Writer,
	lvl logrus.Level,
	formatter logrus.Formatter) gournal.Appender {

	return &appender{&logrus.Logger{Out: out, Level: lvl, Formatter: formatter}}
}

func (a *appender) Append(
	ctx context.Context,
	lvl gournal.Level,
	fields map[string]interface{},
	msg string) {

	switch lvl {
	case gournal.DebugLevel:
		a.logger.WithFields(fields).Debug(msg)
	case gournal.InfoLevel:
		a.logger.WithFields(fields).Info(msg)
	case gournal.WarnLevel:
		a.logger.WithFields(fields).Warn(msg)
	case gournal.ErrorLevel:
		a.logger.WithFields(fields).Error(msg)
	case gournal.FatalLevel:
		a.logger.WithFields(fields).Fatal(msg)
	case gournal.PanicLevel:
		a.logger.WithFields(fields).Panic(msg)
	}
}
