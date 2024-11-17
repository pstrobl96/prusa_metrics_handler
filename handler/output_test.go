package handler

import (
	"context"
	"testing"

	"github.com/influxdata/influxdb-client-go/api"
	"github.com/stretchr/testify/assert"
)

type mockWriteAPI struct {
	api.WriteAPIBlocking
	writeRecordFunc func(ctx context.Context, line string) error
}

func (m *mockWriteAPI) WriteRecord(ctx context.Context, lines ...string) error {
	for _, line := range lines {
		if err := m.writeRecordFunc(ctx, line); err != nil {
			return err
		}
	}
	return nil
}

func TestSentToInflux(t *testing.T) {
	tests := []struct {
		name        string
		message     []string
		writeRecord func(ctx context.Context, line string) error
		expected    bool
		expectError bool
	}{
		{
			name:    "successful write",
			message: []string{"line1", "line2"},
			writeRecord: func(ctx context.Context, line string) error {
				return nil
			},
			expected:    false,
			expectError: false,
		},
		{
			name:    "write error",
			message: []string{"line1", "line2"},
			writeRecord: func(ctx context.Context, line string) error {
				return assert.AnError
			},
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &mockWriteAPI{
				writeRecordFunc: tt.writeRecord,
			}

			result, err := sentToInflux(tt.message, mockAPI)
			assert.Equal(t, tt.expected, result)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
