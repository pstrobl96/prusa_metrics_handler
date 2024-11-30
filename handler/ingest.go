package handler

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/mcuadros/go-syslog.v2"
)

// StartSyslogServer starts a custom syslog server and returns the channel and server
func StartSyslogServer(listenUDP string) (syslog.LogPartsChannel, *syslog.Server) {
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
func MetricsListener(listenUDP string, prefix string) {
	channel, server := StartSyslogServer(listenUDP)

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			received := time.Now()
			log.Trace().Msg(fmt.Sprintf("%v", logParts))

			message, err := Process(logParts, received, prefix)

			if err != nil {
				log.Error().Msg(fmt.Sprintf("Error processing message: %v", err))
				continue
			}

			err = SentToInflux(message, writeAPI)

			if err != nil {
				log.Error().Msg(fmt.Sprintf("Error sending to InfluxDB: %v", err))
				continue
			}
		}
	}(channel)

	server.Wait()

}
