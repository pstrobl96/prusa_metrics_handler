package handler

import (
	"testing"
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
