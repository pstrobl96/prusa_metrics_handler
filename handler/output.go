package handler

import (
	"context"

	"github.com/influxdata/influxdb-client-go/api"
	"github.com/rs/zerolog/log"
)

func sentToInflux(message []string, writeAPI api.WriteAPIBlocking) (result bool, err error) {
	log.Trace().Msg("Sending to InfluxDB")

	for _, line := range message {
		err = writeAPI.WriteRecord(context.Background(), line)
		if err != nil {
			log.Trace().Err(err).Msg("Error while sending to InfluxDB")
			return false, err
		}
	}

	return false, nil
}
