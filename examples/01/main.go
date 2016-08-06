package main

import (
	"golang.org/x/net/context"

	log "github.com/emccode/gournal"
	glogrus "github.com/emccode/gournal/logrus"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, log.LevelKey, log.InfoLevel)
	ctx = context.WithValue(ctx, log.AppenderKey, glogrus.New())

	log.Info(ctx, "Hello %s", "Bob")

	log.WithFields(map[string]interface{}{
		"size":     1,
		"location": "Austin",
	}).Warn(ctx, "Hello %s", "Mary")
}
