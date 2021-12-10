package zxcvbn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_displayTime(t *testing.T) {
	tests := []struct {
		seconds float64
		want    string
	}{
		{0, "less than a second"},
		{30, "30 seconds"},
		{89, "1 minute"},
		{90, "2 minutes"},
		{1905.8, "32 minutes"},
		{9047.062, "3 hours"},
		{1905800, "22 days"},
		{686088000, "21 years"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, displayTime(tt.seconds))
	}
}
