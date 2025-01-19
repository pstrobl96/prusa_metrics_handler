package handler

import (
	"fmt"
	"pstrobl96/prusa_metrics_handler/prometheus"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

// PrinterStatus is a structure to hold printer status
type PrinterStatus struct {
	FirstTimestamp int64
	LastDelta      int64
	mutex          sync.Mutex // mutex because Mutex is sync.Mutex
}

var printerStates sync.Map

// PrinterMetric is a structure to hold printer metrics
type PrinterMetric struct {
	Mac         string
	Timestamp   int64
	MetricName  string
	MetricValue float64
	Labels      map[string]string
}

func process(data format.LogParts, received time.Time, prefix string) {
	mac, ip, timestamp, err := processTimestamp(data, received)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error processing timestamp: %v", err))
		return
	}
	log.Debug().Msg(fmt.Sprintf("Processing data for printer %s with timestamp %d", mac, timestamp))
	metrics, err := processMessage(data["message"].(string), timestamp, mac, prefix, ip)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error processing message: %v", err))
		return
	}

	sentToInflux(metrics, writeAPI)
}

// processTimestamp returns the MAC address and timestamp from the ingested data
// it is basically used for the synchronization of time between handler and the printer
func processTimestamp(data format.LogParts, received time.Time) (string, string, int64, error) {
	mac, ok := data["hostname"].(string)
	if !ok {
		return "", "", 0, fmt.Errorf("mac is not an string")
	}

	message, ok := data["message"].(string)
	if !ok {
		return "", "", 0, fmt.Errorf("message is not an string")
	}

	ip, ok := data["client"].(string)
	if !ok {
		return "", "", 0, fmt.Errorf("ip is not an string")
	}

	re := regexp.MustCompile(`tm=(\d+)`)
	matches := re.FindAllStringSubmatch(message, -1)

	for _, match := range matches {
		if len(match) < 1 {
			return "", "", 0, fmt.Errorf("none time delta value found")
		}
	}
	tmValue, err := strconv.ParseInt(matches[0][1], 10, 64)
	if err != nil {
		return "", "", 0, fmt.Errorf("time delta cannot be converted to int64")
	}

	timedelta := tmValue
	timestamp := time.Now().UnixNano()
	printerInterface, _ := printerStates.LoadOrStore(mac, &PrinterStatus{})
	state := printerInterface.(*PrinterStatus)
	state.mutex.Lock()
	defer state.mutex.Unlock()

	if state.FirstTimestamp == 0 || state.LastDelta > timedelta {
		log.Debug().Msg("First timestamp recorded for printer " + mac)
		state.FirstTimestamp = received.Add(-time.Duration(timedelta)).UnixNano()
		timestamp = received.Add(-time.Duration(timedelta)).UnixNano()
		state.LastDelta = timedelta
	} else {
		log.Debug().Msg("Not the timestamp recorded for printer " + mac)
		state.LastDelta = timedelta
		return mac, ip, state.FirstTimestamp + timedelta, nil
	}
	return mac, ip, timestamp, nil
}

func processMessage(message string, timestamp int64, mac string, prefix string, ip string) ([]string, error) {
	messageSplit := strings.Split(message, "\n")

	if len(messageSplit) == 0 {
		return nil, fmt.Errorf("message is empty")
	}

	firstMessage, err := parseFirstMessage(messageSplit[0])

	if err != nil {
		return nil, fmt.Errorf("error parsing first message: %v", err)
	}

	messageSplit = append(messageSplit[1:], firstMessage)

	for i, line := range messageSplit {
		splitted := strings.Split(line, " ")
		delta, err := strconv.ParseInt(splitted[len(splitted)-1], 10, 64)
		if err != nil {
			log.Error().Msg("Expected error while parsing time delta for metric: " + splitted[0] + " error:" + err.Error())
			continue
		}
		splitted[len(splitted)-1] = strconv.FormatInt(timestamp+delta, 10)
		log.Debug().Msg("Processing timestamps for " + message)
		splitted, err = updateMetric(splitted, prefix, mac, ip)
		if err != nil {
			log.Error().Msg("Expected error while adding mac label for metric: " + splitted[0] + " error:" + err.Error())
			continue
		}
		messageSplit[i] = strings.Join(splitted, " ")
	}
	prometheus.MetricsHandlerTotal.Inc()
	return messageSplit, nil
}

func parseFirstMessage(message string) (string, error) {
	splitted := strings.Split(message, " ")
	if len(splitted) == 0 {
		return "", fmt.Errorf("splitted message is empty")
	}
	firstMsg := splitted[1:]
	return strings.Join(firstMsg, " "), nil
}

func updateMetric(splitted []string, prefix string, mac string, ip string) ([]string, error) {
	if len(splitted) == 0 {
		return nil, fmt.Errorf("splitted message is empty")
	}
	splitted[0] = fmt.Sprintf("%s%s,mac=%s,ip=%s", prefix, splitted[0], mac, ip)

	return splitted, nil
}
