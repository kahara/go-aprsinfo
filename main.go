package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"iter"
	"maps"
	"sort"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	source()

	iter.Pull(maps.Keys(bins))

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
		log.Info().Time("timestamp", key).Int("callsigns", len(bin.Callsigns)).Uint64("packets", bin.Packets).Uint64("duplicates", bin.Duplicates).Send()
	}
}
