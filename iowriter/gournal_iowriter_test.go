package iowriter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/emccode/gournal"
)

func TestIOWriterAppenderNoFields(t *testing.T) {
	gournal.Info(ctx(), "Hello %s", "Bob")
}

func TestIOWriterAppenderWithField(t *testing.T) {
	gournal.WithField("size", 2).Info(ctx(), "Hello %s", "Alice")
}

func TestIOWriterAppenderWithFields(t *testing.T) {
	gournal.WithFields(map[string]interface{}{
		"size":     1,
		"location": "Austin",
	}).Warn(ctx(), "Hello %s", "Mary")
}

func TestIOWriterAppenderPanic(t *testing.T) {

	defer func() {
		r := recover()
		assert.NotNil(t, r, "no panic")
		assert.IsType(t, r, "")
		assert.Equal(t, "[PANIC] Hello Bob\n", r)
	}()

	gournal.Panic(ctx(), "Hello %s", "Bob")
}

func ctx() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, gournal.LevelKey, gournal.InfoLevel)
	ctx = context.WithValue(ctx, gournal.AppenderKey, New())
	return ctx
}
