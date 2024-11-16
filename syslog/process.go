package syslog

import (
	"sync"
	"time"

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
	processTimestamp(data, received)
}

// processTimestamp returns the MAC address and timestamp from the ingested data
// it is basically used for the synchronization of time between handler and the printer
func processTimestamp(data format.LogParts, received time.Time) (string, int64, error) {
	mac := data["hostname"].(string)
	timedelta := data["tm"].(int64)
	timestamp := time.Now().UnixNano()
	printerInterface, _ := printerStates.LoadOrStore(mac, &PrinterStatus{})
	state := printerInterface.(*PrinterStatus)

	state.mu.Lock()
	defer state.mu.Unlock()

	// handling the case when the printer is restarted - not ideal can in rare case can be value the same and then it will not work correctly. Atm it's better than anything
	if state.LastDelta < timedelta || state.FirstTimestamp == 0 {
		state.FirstTimestamp = received.Add(-time.Duration(timedelta)).UnixNano()
		timestamp = received.Add(-time.Duration(timedelta)).UnixNano()
	} else {
		timestamp = state.FirstTimestamp
	}

	return mac, timestamp, nil
}
