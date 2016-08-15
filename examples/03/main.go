package main

import (
	"golang.org/x/net/context"

	log "github.com/emccode/gournal"
	glogrus "github.com/emccode/gournal/logrus"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, log.LevelKey(), log.InfoLevel)
	ctx = context.WithValue(ctx, log.AppenderKey(), glogrus.New())

	ctx = context.WithValue(
		ctx,
		log.FieldsKey(),
		map[string]interface{}{
			"name":  "Venus",
			"color": 0x00ff00,
		})

	// The following log entry will print the message and the name and color
	// of the planet.
	log.Info(ctx, "Discovered planet")

	ctx = context.WithValue(
		ctx,
		log.FieldsKey(),
		func() map[string]interface{} {
			return map[string]interface{}{
				"galaxy":   "Milky Way",
				"distance": 42,
			}
		})

	// The following log entry will print the message and the galactic location
	// and distance of the planet.
	log.Info(ctx, "Discovered planet")

	// Create a Context with the FieldsKey that points to a function which
	// returns a Context's derived fields based upon what data was provided
	// to a the log function.
	ctx = context.WithValue(
		ctx,
		log.FieldsKey(),
		func(ctx context.Context,
			lvl log.Level,
			fields map[string]interface{},
			args ...interface{}) map[string]interface{} {

			if v, ok := fields["z-value"].(int); ok {
				delete(fields, "z-value")
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
		})

	// The following log entry will print the message and two-dimensional
	// location information about the planet.
	log.Info(ctx, "Discovered planet")

	// This log entry, however, will print the message and the same location
	// information, however, because the function used to derive the Context's
	// fields inspects the field's "z-value" key, it will add that data to the
	// location information, making it three-dimensional.
	log.WithField("z-value", 3).Info(ctx, "Discovered planet")
}
