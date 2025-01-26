package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sort"
	"time"
)

// FIXME this would come from outside, like the parameter consts in timebin.go
const Command = "prefixes"

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	source()

	// Massage ingested data
	switch Command {
	case "prefixes":
		prefixes()
	}

	for _, key := range func() []time.Time {
		// There has to be a more straight-forward way to do this
		var keys []time.Time
		for key, _ := range bins {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i].Before(keys[j]) })
		return keys
	}() {
		bin := bins[key]

		switch Command {
		case "counts":
			log.Info().Time("timestamp", key).Int("callsigns", len(bin.Callsigns)).Uint64("packets", bin.Packets).Uint64("duplicates", bin.Duplicates).Send()
		case "prefixes":
			log.Info().Time("timestamp", key).Any("prefixes", bin.Prefixes).Send()
		}
	}
}
