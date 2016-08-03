package stdlib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/emccode/gournal"
)

func TestStdLibAppenderNoFields(t *testing.T) {
	gournal.Info(ctx(), "Hello %s", "Bob")
}

func TestStdLibAppenderWithField(t *testing.T) {
	gournal.WithField("size", 2).Info(ctx(), "Hello %s", "Alice")
}

func TestStdLibAppenderWithFields(t *testing.T) {
	gournal.WithFields(map[string]interface{}{
		"size":     1,
		"location": "Austin",
	}).Warn(ctx(), "Hello %s", "Mary")
}

func TestStdLibAppenderPanic(t *testing.T) {

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
