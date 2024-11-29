package handler

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"

	"github.com/rs/zerolog/log"
	"gopkg.in/mcuadros/go-syslog.v2"
)

func startSyslogServer(listenUDP string) (syslog.LogPartsChannel, *syslog.Server) {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)
	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	server.ListenUDP(listenUDP)
	server.Boot()
	return channel, server
}

// MetricsListener is a function to handle syslog metrics and sent them to processor
func MetricsListener(listenUDP string, influxURL string, influxToken string, influxBucket string, influxOrg string) {
	client = influxdb2.NewClient(influxURL, influxToken)
	writeAPI = client.WriteAPIBlocking(influxOrg, influxBucket)
	channel, server := startSyslogServer(listenUDP)

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			received := time.Now()
			log.Trace().Msg(fmt.Sprintf("%v", logParts))

			message, err := process(logParts, received)

			if err != nil {
				log.Error().Msg(fmt.Sprintf("Error processing message: %v", err))
				continue
			}

			result, err := sentToInflux(message, writeAPI)

			if err != nil {
				log.Error().Msg(fmt.Sprintf("Error sending to InfluxDB: %v", err))
				continue
			}

			if result {
				log.Info().Msg("Metrics sent to InfluxDB")
			}

			/*
				if influx {
					sentToInflux(message, writeAPI)
				} else {
					sentOtlp(message)
				}
			*/
		}
	}(channel)

	server.Wait()

}
