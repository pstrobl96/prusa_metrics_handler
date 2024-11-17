package syslog

import (
	"fmt"
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
	tmValue, ok := data["tm"].(int64)
	if !ok {
		return "", 0, fmt.Errorf("tm value is not an int64")
	}
	if !ok {
		return "", 0, fmt.Errorf("hostname is not a string")
	}
	timedelta := tmValue
	timestamp := time.Now().UnixNano()
	printerInterface, _ := printerStates.LoadOrStore(mac, &PrinterStatus{})
	state := printerInterface.(*PrinterStatus)
	defer state.mu.Unlock()

	// Handle the case when the printer is restarted.
	// This is not ideal as in rare cases the value can be the same,
	// and then it will not work correctly. At the moment, it's better than nothing.

	if state.LastDelta < timedelta || state.FirstTimestamp == 0 {
		state.FirstTimestamp = received.Add(-time.Duration(timedelta)).UnixNano()
		timestamp = received.Add(-time.Duration(timedelta)).UnixNano()
		state.LastDelta = timedelta
	} else {
		timestamp = state.FirstTimestamp
		state.LastDelta = timedelta
	}
	return mac, timestamp, nil
}
