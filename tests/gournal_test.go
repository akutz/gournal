package gournal

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	. "github.com/emccode/gournal"
	"github.com/emccode/gournal/iowriter"
)

func TestMain(m *testing.M) {
	DefaultLevel = DebugLevel
	os.Exit(m.Run())
}

func TestDefaultAppender(t *testing.T) {
	w := &bytes.Buffer{}
	DefaultAppender = iowriter.NewWithOptions(w)

	assert.NotPanics(t, func() {
		ctx := context.Background()
		Error(ctx, "Hello %s", "Bob")

		ctx = context.WithValue(ctx, LevelKey, ErrorLevel)
		Info(ctx, "Hello %s", "Alice")
		Error(ctx, "Hello %s %s", "Mary", "Kay")

		Debug(nil, "Goodbye %d", 3)
	})

	expected := "[ERROR] Hello Bob\n" +
		"[ERROR] Hello Mary Kay\n" +
		"[DEBUG] Goodbye 3\n"

	assert.Equal(t, expected, w.String())
}

func TestNilDefaultAppender(t *testing.T) {
	assert.Panics(t, func() {
		DefaultAppender = nil
		Error(context.Background(), "Hello %s", "Bob")
	})
}

func TestInfo(t *testing.T) {
	buf, ctx := newTestContext()
	Info(ctx, "Hello %s", "Bob")
	assert.Equal(t, "[INFO] Hello Bob\n", buf.String())
}

func TestNoOpLog(t *testing.T) {
	buf, ctx := newTestContext()
	ctx = context.WithValue(ctx, LevelKey, ErrorLevel)
	Info(ctx, "Hello %s", "Alice")
	assert.Zero(t, buf.Len())
}

func TestChildContextLevelLog(t *testing.T) {
	buf, ctx := newTestContext()
	ctx = context.WithValue(ctx, LevelKey, ErrorLevel)
	Info(ctx, "Hello %s", "Bob")
	assert.Zero(t, buf.Len())

	ctx = context.WithValue(ctx, LevelKey, DebugLevel)
	Debug(ctx, "Hello %s", "Alice")
	assert.Equal(t, "[DEBUG] Hello Alice\n", buf.String())
}

func TestContextFieldsMap(t *testing.T) {
	buf, ctx := newTestContext()
	ctx = context.WithValue(ctx, LevelKey, InfoLevel)

	ctxFields := map[string]interface{}{
		"point": struct {
			x int
			y int
		}{1, -1},
	}
	ctx = context.WithValue(ctx, FieldsKey, ctxFields)

	Info(ctx, "Discovered planet")
	assert.Equal(
		t, "[INFO] Discovered planet map[point:{1 -1}]\n", buf.String())
}

func TestContextFieldsFunc(t *testing.T) {
	buf, ctx := newTestContext()
	ctx = context.WithValue(ctx, LevelKey, InfoLevel)

	ctxFieldsFunc := func() map[string]interface{} {
		return map[string]interface{}{
			"point": struct {
				x int
				y int
			}{1, -1},
		}
	}
	ctx = context.WithValue(ctx, FieldsKey, ctxFieldsFunc)

	Info(ctx, "Discovered planet")
	assert.Equal(
		t, "[INFO] Discovered planet map[point:{1 -1}]\n", buf.String())
}

func TestContextFieldsFuncEx(t *testing.T) {
	buf, ctx := newTestContext()
	ctx = context.WithValue(ctx, LevelKey, InfoLevel)

	ctxFieldsFunc := func(
		ctx context.Context,
		lvl Level,
		fields map[string]interface{},
		args ...interface{}) map[string]interface{} {

		if v, ok := fields["size"].(int); ok {
			delete(fields, "size")
			return map[string]interface{}{
				"point": struct {
					x int
					y int
					z int
				}{1, -1, v},
			}
		}

		return map[string]interface{}{
			"point": struct {
				x int
				y int
			}{1, -1},
		}
	}
	ctx = context.WithValue(ctx, FieldsKey, ctxFieldsFunc)
	ctxLogger := New(ctx)

	Info(ctx, "Discovered planet")
	assert.Equal(
		t, "[INFO] Discovered planet map[point:{1 -1}]\n", buf.String())

	buf.Reset()
	assert.Equal(t, buf.Len(), 0)

	ctxLogger.Info("Discovered planet")
	assert.Equal(
		t, "[INFO] Discovered planet map[point:{1 -1}]\n", buf.String())

	buf.Reset()
	assert.Equal(t, buf.Len(), 0)

	WithField("size", 3).Info(ctx, "Discovered planet")
	assert.Equal(
		t, "[INFO] Discovered planet map[point:{1 -1 3}]\n", buf.String())
}

func newTestContext() (*bytes.Buffer, context.Context) {
	w := &bytes.Buffer{}
	a := iowriter.NewWithOptions(io.MultiWriter(w, os.Stdout))
	return w, context.WithValue(context.Background(), AppenderKey, a)
}
