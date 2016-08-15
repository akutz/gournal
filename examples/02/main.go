package main

import (
	"golang.org/x/net/context"

	log "github.com/emccode/gournal"
	glogrus "github.com/emccode/gournal/logrus"
	gzap "github.com/emccode/gournal/zap"
)

func main() {
	// Make a Zap-based Appender the default appender for when one is not
	// present in a Context, or when a nill Context is provided to a logging
	// function.
	log.DefaultAppender = gzap.New()

	// The following call fails to provide a valid Context argument. In this
	// case the DefaultAppender is used.
	log.WithFields(map[string]interface{}{
		"size":     2,
		"location": "Boston",
	}).Error(nil, "Hello %s", "Bob")

	ctx := context.Background()
	ctx = context.WithValue(ctx, log.LevelKey(), log.InfoLevel)

	// Even though this next call provides a valid Context, there is no
	// Appender present in the Context so the DefaultAppender will be used.
	log.Info(ctx, "Hello %s", "Mary")

	ctx = context.WithValue(ctx, log.AppenderKey(), glogrus.New())

	// This last log function uses a Context that has been created with a
	// Logrus Appender. Even though the DefaultAppender is assigned and is a
	// Zap-based logger, this call will utilize the Context Appender instance,
	// a Logrus Appender.
	log.WithFields(map[string]interface{}{
		"size":     1,
		"location": "Austin",
	}).Warn(ctx, "Hello %s", "Alice")
}
