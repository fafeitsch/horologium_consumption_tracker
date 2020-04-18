package domain

import (
	"github.com/fafeitsch/Horologium/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLastReadingBefore(t *testing.T) {
	readings := MeterReadings{
		{Date: util.FormatDate(2019, 9, 12), Id: 5},
		{Date: util.FormatDate(2019, 9, 23), Id: 6},
		{Date: util.FormatDate(2019, 9, 27), Id: 7},
		{Date: util.FormatDate(2019, 9, 30), Id: 8},
	}
	t.Run("simple", func(t *testing.T) {
		got := readings.lastReadingBefore(util.FormatDate(2019, 9, 26))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("same day", func(t *testing.T) {
		got := readings.lastReadingBefore(util.FormatDate(2019, 9, 23))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("all smaller", func(t *testing.T) {
		got := readings.lastReadingBefore(util.FormatDate(2019, 10, 1))
		assert.Equal(t, readings[3], got, "last reading calculated incorrectly")
	})
}

func TestFirstReadingAfter(t *testing.T) {
	readings := MeterReadings{
		{Date: util.FormatDate(2019, 9, 12), Id: 5},
		{Date: util.FormatDate(2019, 9, 23), Id: 6},
		{Date: util.FormatDate(2019, 9, 27), Id: 7},
		{Date: util.FormatDate(2019, 9, 30), Id: 8},
	}
	cpy := make(MeterReadings, len(readings))
	copy(cpy, readings)
	t.Run("simple", func(t *testing.T) {
		got := cpy.firstReadingAfter(util.FormatDate(2019, 9, 17))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("same day", func(t *testing.T) {
		got := cpy.firstReadingAfter(util.FormatDate(2019, 9, 23))
		assert.Equal(t, readings[1], got, "last reading calculated incorrectly")
	})
	t.Run("all bigger", func(t *testing.T) {
		got := cpy.firstReadingAfter(util.FormatDate(2019, 9, 1))
		assert.Equal(t, readings[0], got, "last reading calculated incorrectly")
	})
}
