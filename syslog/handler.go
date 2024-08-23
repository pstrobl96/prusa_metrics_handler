package syslog

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/rs/zerolog/log"
	"gopkg.in/mcuadros/go-syslog.v2"
)

var (
	mutex             sync.RWMutex
	printersTimestamp = map[string]int64{}
	pritnersMsg       = map[string]int64{}
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
func MetricsListener(listenUDP string, influxURL string, influxToken string, influxBucket string, influxOrg string) {

	client := influxdb2.NewClient(influxURL, influxToken)
	writeAPI := client.WriteAPIBlocking(influxOrg, influxBucket)

	channel, server := startSyslogServer(listenUDP)

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			log.Trace().Msg(fmt.Sprintf("%v", logParts))

			sentToInflux(processTimestamps(logParts["message"].(string), logParts["hostname"].(string)), writeAPI)
		}
	}(channel)

	server.Wait()

}

func processTimestamps(message string, mac string) (result []string) {

	messageSplit := strings.Split(message, "\n")

	firstMsg := strings.Split(messageSplit[0], " ")
	messageSplit[0] = firstMsg[1]


	printerInfo := strings.Split(firstMsg[0], ",")
	msg := strings.Split(printerInfo[0], "=")[1]
	printerTm := strings.Split(printerInfo[1], "=")[1]


	var printerStart int64

	mutex.Lock()
	_, contains := printersTimestamp[mac]
	if !contains {
		if printerTm < 10000000 { // allow 10 ms drift
			printersTimestamp[mac] = time.Now().UnixNano() - printerTm
		} else {
			printersTimestamp[mac] = -1

		}
	} else {
		printerStart = printersTimestamp[mac]
	}
	mutex.Unlock()

	for i, line := range strings.Split(message, "\n") {
		splitted := strings.Split(line, " ")
		delta, err := strconv.ParseInt(splitted[len(splitted)-1], 10, 64)
		if err != nil {
			log.Info().Msg("Expected error while parsing timestamp: " + err.Error())
			continue
		}
		splitted[len(splitted)-1] = strconv.FormatInt(calculateTimestamp(printerStart, printerTm, delta), 10)
		log.Trace().Msg("Processing timestamps for " + message)
		messageSplit[i] = strings.Join(splitted, " ")
	}

	metricsHandlerTotal.Inc()
	return messageSplit
}

func sentToInflux(message []string, writeAPI api.WriteAPIBlocking) (result bool, err error) {
	log.Trace().Msg("Sending to InfluxDB")

	for _, line := range message {
		err = writeAPI.WriteRecord(context.Background(), line)
		if err != nil {
			log.Trace().Err(err).Msg("Error while sending to InfluxDB")
			for _, line := range fixLine(line) {
				log.Trace().Msg("Trying to send fixed line to InfluxDB")
				writeAPI.WriteRecord(context.Background(), line)

			}
			return false, err
		}
	}

	return false, nil
}

func fixLine(line string) (message []string) {
	if strings.Contains(line, "msg") {
		return strings.Split(line, " ")
	}
	return []string{}

}

func calculateTimestamp(printerStart int64, printerRunning int64, metricDelta int64) (result int64) {
	if printerStart == -1 {
		return time.Now().UnixNano() + metricDelta
	}
	return printerStart + printerRunning + metricDelta
}

func checkPrinterRestart(messageSplitted []string) (printerTm int64, result bool, err error) {
	msg := strings.Split(messageSplitted[0], ",")
	msgItem := strings.Split(msg[0], "=")[0]
	printerInterval := strings.Split(msg[1], "=")[1]

	if msgItem < 

	i, err := strconv.ParseInt(strings.Split(strings.Split(messageSplitted[0], ",")[1], "=")[1], 10, 64)
	if err != nil {
		return 0, false, err
	}
	return i, true, nil
}
