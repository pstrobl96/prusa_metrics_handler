package handler

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

// PrinterStatus is a structure to hold printer status
type PrinterStatus struct {
	FirstTimestamp int64
	LastDelta      int64
	mu             sync.Mutex
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

func process(data format.LogParts, received time.Time) {
	mac, timestamp, err := processTimestamp(data, received)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error processing timestamp: %v", err))
		return
	}
	log.Debug().Msg(fmt.Sprintf("Processing timestamp %d", timestamp))
	log.Debug().Msg(fmt.Sprintf("Processing data for printer %s", mac))
}

// processTimestamp returns the MAC address and timestamp from the ingested data
// it is basically used for the synchronization of time between handler and the printer
func processTimestamp(data format.LogParts, received time.Time) (string, int64, error) {
	mac, ok := data["hostname"].(string)
	if !ok {
		return "", 0, fmt.Errorf("mac is not an string")
	}
	re := regexp.MustCompile(`tm=(\d+)`)
	matches := re.FindAllStringSubmatch(data["message"].(string), -1)

	for _, match := range matches {
		if len(match) < 1 {
			return "", 0, fmt.Errorf("none time delta value found")
		}
	}
	tmValue, err := strconv.ParseInt(matches[0][1], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("time delta cannot be converted to int64")
	}

	timedelta := tmValue
	timestamp := time.Now().UnixNano()
	printerInterface, _ := printerStates.LoadOrStore(mac, &PrinterStatus{})
	state := printerInterface.(*PrinterStatus)
	state.mu.Lock()
	defer state.mu.Unlock()

	if state.FirstTimestamp == 0 || state.LastDelta > timedelta {
		log.Debug().Msg("First timestamp recorded for printer " + mac)
		state.FirstTimestamp = received.Add(-time.Duration(timedelta)).UnixNano()
		timestamp = received.Add(-time.Duration(timedelta)).UnixNano()
		state.LastDelta = timedelta
	} else {
		log.Debug().Msg("Not the timestamp recorded for printer " + mac)
		state.LastDelta = timedelta
	}
	return mac, timestamp, nil
}