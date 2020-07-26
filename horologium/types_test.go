package horologium

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInterpolateValueAtDate(t *testing.T) {
	readings := MeterReadings{
		{Date: FormatDate(2019, 9, 12), Count: 45},
		{Date: FormatDate(2019, 9, 23), Count: 80},
		{Date: FormatDate(2019, 9, 27), Count: 134},
		{Date: FormatDate(2019, 9, 30), Count: 178},
	}
	tests := []struct {
		date string
		want float64
	}{
		{date: "2019-09-27", want: 134},
		{date: "2019-09-30", want: 178},
		{date: "2020-01-02", want: 178},
		{date: "2019-09-26", want: 120.5},
	}
	for _, tt := range tests {
		t.Run(tt.date, func(t *testing.T) {
			date, _ := time.Parse(DateFormat, tt.date)
			got := readings.interpolateValueAtDate(date)
			assert.Equal(t, tt.want, got, "interpolated value is wrong")
		})
	}
}

func TestLastReadingBefore(t *testing.T) {
	readings := MeterReadings{
		{Date: FormatDate(2019, 9, 12), Id: 5},
		{Date: FormatDate(2019, 9, 23), Id: 6},
		{Date: FormatDate(2019, 9, 27), Id: 7},
		{Date: FormatDate(2019, 9, 30), Id: 8},
	}
	t.Run("simple", func(t *testing.T) {
		got := readings.lastReadingBefore(FormatDate(2019, 9, 26))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("same day", func(t *testing.T) {
		got := readings.lastReadingBefore(FormatDate(2019, 9, 23))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("all smaller", func(t *testing.T) {
		got := readings.lastReadingBefore(FormatDate(2019, 10, 1))
		assert.Equal(t, readings[3], got, "last reading calculated incorrectly")
	})
}

func TestFirstReadingAfter(t *testing.T) {
	readings := MeterReadings{
		{Date: FormatDate(2019, 9, 12), Id: 5},
		{Date: FormatDate(2019, 9, 23), Id: 6},
		{Date: FormatDate(2019, 9, 27), Id: 7},
		{Date: FormatDate(2019, 9, 30), Id: 8},
	}
	cpy := make(MeterReadings, len(readings))
	copy(cpy, readings)
	t.Run("simple", func(t *testing.T) {
		got := cpy.firstReadingAfter(FormatDate(2019, 9, 17))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("same day", func(t *testing.T) {
		got := cpy.firstReadingAfter(FormatDate(2019, 9, 23))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("all bigger", func(t *testing.T) {
		got := cpy.firstReadingAfter(FormatDate(2019, 9, 1))
		assert.Equal(t, readings[0], got, "last reading calculated incorrectly")
	})
}
