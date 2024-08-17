package apps

import (
	"testing"
	"time"
)

func TestNewApp(t *testing.T) {
	type inputs struct {
		id             string
		submissionTime time.Time
	}
	tests := []struct {
		caseName string
		input    inputs
	}{
		{"normal", inputs{"application-01", time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)}},
	}
	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			tmp := NewApp(test.input.id, test.input.submissionTime,"")
			if tmp.Id != test.input.id || !tmp.SubmissionTime.Equal(test.input.submissionTime) {
				t.Errorf("Expected %s %v, got %s %v",
					test.input.id,
					test.input.submissionTime,
					tmp.Id,
					tmp.SubmissionTime)
			}
		})
	}
}

