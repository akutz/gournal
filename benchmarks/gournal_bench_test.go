package benchmarks

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/uber-go/zap"
	//gaetest "google.golang.org/appengine/aetest"
	//gae "google.golang.org/appengine/log"

	"github.com/akutz/gournal"
	//ggae "github.com/akutz/gournal/gae"
	glogrus "github.com/akutz/gournal/logrus"
	glog "github.com/akutz/gournal/stdlib"
	gzap "github.com/akutz/gournal/zap"
)

var gaeCtx context.Context

func TestMain(m *testing.M) {
	gournal.DefaultLevel = gournal.DebugLevel

	/*var (
		err  error
		done func()
	)

	if gaeCtx, done, err = gaetest.NewContext(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}*/
	gaeCtx = context.Background()

	ec := m.Run()

	//done()
	os.Exit(ec)
}

func BenchmarkNativeStdLibWithoutFields(b *testing.B) {
	l := log.New(os.Stderr, "", log.LstdFlags)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Printf("Run Barry, run.")
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
			l.Printf("Run Barry, run.")
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

/*func BenchmarkNativeGAEWithoutFields(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gae.Infof(gaeCtx, "Run Barry, run.")
		}
	})
}*/

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

/*func BenchmarkGournalGAEWithoutFields(b *testing.B) {
	benchmarkWithoutFields(b, ggae.New())
}*/

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

/*func BenchmarkGournalGAEWithFields(b *testing.B) {
	benchmarkWithFields(b, ggae.New())
}*/

func newContext(a gournal.Appender) context.Context {
	/*var ctx context.Context
	if a == ggae.New() {
		ctx = gaeCtx
	} else {
		ctx = context.Background()
	}*/
	ctx := context.Background()
	ctx = context.WithValue(ctx, gournal.LevelKey(), gournal.InfoLevel)
	ctx = context.WithValue(ctx, gournal.AppenderKey(), a)
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
