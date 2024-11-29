package handler

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/rs/zerolog/log"
)

var (
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
)

func sentOtlp(message []string) {
	log.Trace().Msg("Sending to OTLP")
}

func sentToInflux(message []string, writeAPI api.WriteAPIBlocking) (err error) {
	log.Trace().Msg("Sending to InfluxDB")

	for _, line := range message {
		err = writeAPI.WriteRecord(context.Background(), line)
		if err != nil {
			log.Error().Err(err).Msg("Error while sending to InfluxDB")
			return err
		}
	}

	return nil
}
