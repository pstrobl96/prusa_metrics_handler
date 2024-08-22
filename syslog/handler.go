package syslog

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
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

// MetricsListener is a function to handle syslog metrics and send them to InfluxDB v3
func MetricsListener(listenUDP string, influxHost string, influxToken string, influxDb string) {
	channel, server := startSyslogServer(listenUDP)

	client, err := influxdb3.New(influxdb3.ClientConfig{
		Host:     influxHost,
		Token:    influxToken,
		Database: influxDb,
	})

	if err != nil {
		log.Panic().Msg("Error creating InfluxDB client: " + err.Error())
	}

	defer func(client *influxdb3.Client) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			log.Trace().Msg(fmt.Sprintf("%v", logParts))

			sentToInflux(processTimestamps(logParts["message"].(string)), client)
		}
	}(channel)

	server.Wait()

}

func processTimestamps(message string) (result string) {
	splitted := strings.Split(message, " ")
	delta, err := strconv.ParseInt(splitted[len(splitted)], 10, 64)
	if err != nil {
		log.Error().Msg("Error parsing timestamp: " + err.Error())
	}

	splitted[len(splitted)] = strconv.FormatInt(time.Now().Unix()+delta, 10)
	log.Trace().Msg("Processing timestamps for " + message)
	metricsHandlerTotal.Inc()
	return message
}

func sentToInflux(message string, client *influxdb3.Client) (result bool, err error) {
	err = client.Write(context.Background(), []byte(message))
	if err != nil {
		log.Panic().Msg("Error sending to InfluxDB: " + err.Error())
	}

	log.Trace().Msg("Sending to InfluxDB: " + message)
	return false, nil
}
