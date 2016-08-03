// Package iowriter provides a Gournal Appender that writes to any io.Writer
// object.
package iowriter

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/context"

	"github.com/emccode/gournal"
)

// New returns a Gournal Appender that writes to any io.Writer object.
// This function selects os.Stdout as the io.Writer instance to use.
func New() gournal.Appender {
	return &appender{os.Stdout}
}

// NewWithOptions returns a Gournal Appender that writes to the provided
// io.Writer object.
func NewWithOptions(w io.Writer) gournal.Appender {
	return &appender{w}
}

type appender struct {
	w io.Writer
}

func (a *appender) Append(
	ctx context.Context,
	lvl gournal.Level,
	fields map[string]interface{},
	msg string) {

	var (
		w        = a.w
		panicBuf *bytes.Buffer
	)

	if lvl == gournal.PanicLevel {
		panicBuf = &bytes.Buffer{}
		w = io.MultiWriter(a.w, panicBuf)
	}

	if len(fields) == 0 {
		fmt.Fprintf(w, "[%s] %s\n", lvl, msg)
	} else {
		fmt.Fprintf(w, "[%s] %s %v\n", lvl, msg, fields)
	}

	if lvl == gournal.FatalLevel {
		os.Exit(1)
	}

	if lvl == gournal.PanicLevel {
		panic(panicBuf.String())
	}
}
