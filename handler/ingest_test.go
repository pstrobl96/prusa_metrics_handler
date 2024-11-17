package handler

import (
	"net"
	"testing"
	"time"
)

func TestStartSyslogServer(t *testing.T) {
	listenUDP := "localhost:514"

	channel, server := startSyslogServer(listenUDP)

	if channel == nil {
		t.Errorf("Expected non-nil LogPartsChannel, got nil")
	}

	if server == nil {
		t.Errorf("Expected non-nil syslog.Server, got nil")
	}

	// Clean up
	server.Kill()
}

func TestMetricsListener(t *testing.T) {
	listenUDP := "localhost:514"
	influxURL := "http://localhost:8086"
	influxToken := "my-token"
	influxBucket := "my-bucket"
	influxOrg := "my-org"

	go MetricsListener(listenUDP, influxURL, influxToken, influxBucket, influxOrg)

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("udp", listenUDP)
	if err != nil {
		t.Fatalf("Failed to connect to syslog server: %v", err)
	}
	defer conn.Close()

	message := "<34>1 2023-10-10T14:48:00Z mymachine.example.com su - ID47 - BOM'su root' failed for lonvick on /dev/pts/8"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Fatalf("Failed to send syslog message: %v", err)
	}

	// Since the process function is not defined here and it is not supposed to return anything, I cannot verify its behavior here.
}
