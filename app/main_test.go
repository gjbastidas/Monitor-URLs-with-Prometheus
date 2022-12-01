package main

import (
	"testing"
)

func TestRecordMetrics(t *testing.T) {
	tests := []struct {
		input         string
		expectedError bool
	}{
		{input: "https://httpstat.us/200", expectedError: false},
		{input: "http://BadURL", expectedError: true},
	}

	for _, tt := range tests {
		if err := recordMetrics(tt.input); err != nil && tt.expectedError == false {
			t.Errorf("got (%v), expectedError(%v)", err.Error(), tt.expectedError)
		}
	}
}
