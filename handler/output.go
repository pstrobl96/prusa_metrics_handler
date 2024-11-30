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

func sentToInflux(message []string, writeAPI api.WriteAPIBlocking) (err error) {
	log.Trace().Msg("Sending to InfluxDB")

	for _, line := range message {
		err = writeAPI.WriteRecord(context.Background(), line)
		if err != nil {
			log.Error().Err(err)
			return err
		}
	}

	return nil
}
