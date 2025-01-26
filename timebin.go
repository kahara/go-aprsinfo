package main

import (
	"github.com/cespare/xxhash"
	"github.com/pd0mz/go-aprs"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

const (
	// FIXME these would come from outside
	DestinationFilter = ""
	TimebinDuration   = time.Duration(time.Hour)
	DuplicateLookback = time.Duration(time.Second * 30)
)

type Timebin struct {
	Callsigns  map[string]bool
	Prefixes   map[string]uint64
	Packets    uint64
	Duplicates uint64
}

var (
	hashes = make(map[uint64]time.Time) // FIXME This should get pruned periodically
	bins   = make(map[time.Time]*Timebin)
)

func timebinPacket(timestamp time.Time, packet aprs.Packet) {
	if DestinationFilter != "" && packet.Dst.String() != DestinationFilter {
		return
	}

	binTime := timestamp.Truncate(TimebinDuration)
	if _, ok := bins[binTime]; !ok {
		bins[binTime] = &Timebin{
			Callsigns:  make(map[string]bool),
			Prefixes:   make(map[string]uint64),
			Packets:    0,
			Duplicates: 0,
		}
	}
	bins[binTime].Callsigns[packet.Src.String()] = true
	bins[binTime].Prefixes[strings.ToUpper(packet.Src.String())[:2]]++
	bins[binTime].Packets++

	var sb strings.Builder

	sb.WriteString(packet.Src.String())
	sb.WriteString(string(packet.Payload))
	//log.Debug().Str("packetHash", sb.String()).Send()
	packetHash := xxhash.Sum64String(sb.String())

	if _, ok := hashes[packetHash]; ok && timestamp.Sub(hashes[packetHash]) < DuplicateLookback {
		log.Debug().Any("packethash", packetHash).Time("timestamp", timestamp).Any("packet", packet).Msg("Duplicate")
		bins[binTime].Duplicates++
	}
	hashes[packetHash] = timestamp
}
