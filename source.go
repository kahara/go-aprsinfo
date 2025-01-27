package main

import (
	"bufio"
	"github.com/kahara/go-canner"
	"github.com/pd0mz/go-aprs"
	"os"
)

func source() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		text := scanner.Text()
		record, err := canner.NewRecord(text)
		if err != nil {
			//log.Err(err).Str("source", text).Msg("Error parsing record")
			continue
		}

		//log.Debug().Str("record", string(record.Payload)).Msg("Processing payload")

		packet, err := aprs.ParsePacket(string(record.Payload))
		if err != nil {
			//log.Err(err).Str("payload", string(record.Payload)).Msg("Error parsing packet")
			continue
		}
		//log.Debug().Any("packet", packet).Msg("Parsed packet")
		timebinPacket(record.Timestamp.UTC(), packet)
	}
}
