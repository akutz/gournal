package gournal_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"

	"github.com/emccode/gournal"
	glogrus "github.com/emccode/gournal/logrus"
	glog "github.com/emccode/gournal/stdlib"
	gzap "github.com/emccode/gournal/zap"
)

func BenchmarkNativeStdLibWithoutFields(b *testing.B) {
	l := log.New(os.Stderr, "", log.LstdFlags)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Println("Run Barry, run.")
		}
	})
}

func BenchmarkNativeLogrusWithoutFields(b *testing.B) {
	l := &logrus.Logger{
		Out:       os.Stderr,
		Level:     logrus.DebugLevel,
		Formatter: logrus.StandardLogger().Formatter,
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Println("Run Barry, run.")
		}
	})
}

func BenchmarkNativeZapWithoutFields(b *testing.B) {
	l := zap.New(zap.NewJSONEncoder(), zap.Output(os.Stderr))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Info("Run Barry, run.")
		}
	})
}

func BenchmarkGournalStdLibWithoutFields(b *testing.B) {
	benchmarkWithoutFields(
		b, glog.NewWithOptions(os.Stderr, "", log.LstdFlags))
}

func BenchmarkGournalLogrusWithoutFields(b *testing.B) {
	benchmarkWithoutFields(b, glogrus.NewWithOptions(
		os.Stderr, logrus.DebugLevel, logrus.StandardLogger().Formatter))
}

func BenchmarkGournalZapWithoutFields(b *testing.B) {
	benchmarkWithoutFields(b, gzap.NewWithOptions(
		zap.NewJSONEncoder(), zap.Output(os.Stderr)))
}

func BenchmarkGournalStdLibWithFields(b *testing.B) {
	benchmarkWithFields(
		b, glog.NewWithOptions(os.Stderr, "", log.LstdFlags))
}

func BenchmarkGournalLogrusWithFields(b *testing.B) {
	benchmarkWithFields(b, glogrus.NewWithOptions(
		os.Stderr, logrus.DebugLevel, logrus.StandardLogger().Formatter))
}

func BenchmarkGournalZapWithFields(b *testing.B) {
	benchmarkWithFields(b, gzap.NewWithOptions(
		zap.NewJSONEncoder(), zap.Output(os.Stderr)))
}

func newContext(a gournal.Appender) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, gournal.LevelKey, gournal.InfoLevel)
	ctx = context.WithValue(ctx, gournal.AppenderKey, a)
	return ctx
}

func benchmarkWithoutFields(b *testing.B, a gournal.Appender) {
	ctx := newContext(a)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gournal.Info(ctx, "Run Barry, run.")
		}
	})
}

func benchmarkWithFields(b *testing.B, a gournal.Appender) {
	ctx := newContext(a)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gournal.WithFields(map[string]interface{}{
				"name": "Bob",
				"size": 10,
				"when": time.Now(),
			}).Info(ctx, "Run Barry, run.")
		}
	})
}
