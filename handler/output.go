package handler

import (
	"context"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	// Channel is used to send messages to be processed
	Channel chan []string
)

// InitInfluxDB initializes the InfluxDB client and writeAPI for send metrics to InfluxDB
func InitInfluxDB(influxURL string, influxToken string, influxBucket string, influxOrg string) {
	client = influxdb2.NewClient(influxURL, influxToken)
	writeAPI = client.WriteAPIBlocking(influxOrg, influxBucket)
}

// InitOTLP initializes the OpenTelemetry client
func InitOTLP(endpoint string) {
	log.Info().Msg("OTLP enabled")
}

// SentToInflux sends the messages array in a loop to InfluxDB
func SentToInflux(message []string, writeAPI api.WriteAPIBlocking) (err error) {
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

// WriteToStdout writes the messages to stdout
func WriteToStdout(message []string) {
	stdoutLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	for _, line := range message {
		stdoutLogger.Info().Msg(line)
	}
}
