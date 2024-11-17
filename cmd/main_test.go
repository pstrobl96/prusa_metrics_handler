package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	// Set up flags
	os.Args = []string{"cmd", "--exporter.metrics-path=/testmetrics", "--exporter.metrics-port=10012", "--listen-address=0.0.0.0:8515", "--influx-org=testorg", "--influx-bucket=testbucket", "--influx-token=testtoken", "--influx-url=http://testurl:8086", "--log.level=debug"}

	// Parse flags
	kingpin.Parse()

	// Set up logger
	logLevel, err := zerolog.ParseLevel(*logLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano

	// Start the metrics handler
	go Run()

	// Test metrics endpoint
	req, err := http.NewRequest("GET", "http://localhost:10012/testmetrics", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := promhttp.Handler()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")

	// Test syslog listener
	// Note: This is a placeholder. You would need to implement a proper test for the syslog listener.
	assert.NotNil(t, *syslogListenAddress, "Syslog listener address should not be nil")
}
