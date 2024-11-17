package handler

import (
	"testing"
	"time"

	"gopkg.in/mcuadros/go-syslog.v2/format"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name     string
		data     format.LogParts
		received time.Time
		wantErr  bool
	}{
		{
			name: "Valid data",
			data: format.LogParts{
				"hostname": "00:11:22:33:44:55",
				"message":  "tm=123456 metric1 10 metric2 20",
			},
			received: time.Now(),
			wantErr:  false,
		},
		{
			name: "Invalid hostname",
			data: format.LogParts{
				"hostname": 12345,
				"message":  "tm=123456 metric1 10 metric2 20",
			},
			received: time.Now(),
			wantErr:  true,
		},
		{
			name: "Invalid message",
			data: format.LogParts{
				"hostname": "00:11:22:33:44:55",
				"message":  12345,
			},
			received: time.Now(),
			wantErr:  true,
		},
		{
			name: "No time delta in message",
			data: format.LogParts{
				"hostname": "00:11:22:33:44:55",
				"message":  "metric1 10 metric2 20",
			},
			received: time.Now(),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						t.Errorf("process() panicked: %v", r)
					}
				}
			}()

			process(tt.data, tt.received)
		})
	}
}
