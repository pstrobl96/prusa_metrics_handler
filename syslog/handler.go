package syslog

import (
	"fmt"

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

// MetricsListener is a function to handle syslog metrics and send them to InfluxDB
func MetricsListener(listenUDP string, influxAddress string) {
	channel, server := startSyslogServer(listenUDP)

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			log.Trace().Msg(fmt.Sprintf("%v", logParts))

			sentToInflux(influxAddress, processTimestamps(logParts["message"].(string)))
		}
	}(channel)

	server.Wait()

}

func processTimestamps(message string) (result string) {
	// Dummy function
	log.Trace().Msg("Processing timestamps for " + message)
	metricsHandlerTotal.Inc()
	return message
}

func sentToInflux(influxAddress string, message string) (result bool, err error) {
	// Dummy function
	log.Trace().Msg("Sending to " + influxAddress + ": " + message)
	return false, nil
}
