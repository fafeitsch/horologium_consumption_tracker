package horologium

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// MeterReading represents the counter on a certain meter at a certain date.
type MeterReading struct {
	Count float64   // The count showing on the meter.
	Date  time.Time // The date the meter showed the Count. The time part of the date should always be 0:00.
}

// A slice of meter readings.
type MeterReadings []MeterReading

func (m MeterReadings) interpolateValueAtDate(date time.Time) float64 {
	firstReading := m.lastReadingBefore(date)
	lastReading := m.firstReadingAfter(date)
	if firstReading == lastReading {
		lastReading = m.firstReadingAfter(firstReading.Date.Add(24 * time.Hour))
	}
	if date == lastReading.Date || date.After(lastReading.Date) {
		return lastReading.Count
	}
	differenceDays := math.Round(lastReading.Date.Sub(firstReading.Date).Hours() / 24)
	slope := (lastReading.Count - firstReading.Count) / differenceDays
	xValue := math.Round(date.Sub(firstReading.Date).Hours() / 24)
	return slope*xValue + firstReading.Count
}

func (m MeterReadings) lastReadingBefore(date time.Time) MeterReading {
	index := 0
	for index < len(m) && (m[index].Date.Equal(date) || m[index].Date.Before(date)) {
		index = index + 1
	}
	return m[index-1]
}

func (m MeterReadings) firstReadingAfter(date time.Time) MeterReading {
	index := len(m) - 1
	for index >= 0 && (m[index].Date.Equal(date) || m[index].Date.After(date)) {
		index = index - 1
	}
	if index == len(m)-1 {
		return m[len(m)-1]
	}
	return m[index+1]
}

// Consumption computes how much units were needed between start and end. Both
// the start and the end should be dates, with their time set to 0:00 (start of the day).
// Important: the meter readings should be sorted (see Sort function)
//
// Between two consecutive meter readings the consumption is assumed to be linear.
func (m MeterReadings) Consumption(start time.Time, end time.Time) float64 {
	valueStart := m.interpolateValueAtDate(start)
	valueEnd := m.interpolateValueAtDate(end)
	return valueEnd - valueStart
}

// Sort sorts the meter readings in ascending order by the date. Most functions
// rely on the meter readings to be sorted.
func (m MeterReadings) Sort() {
	sort.Slice(m, func(i, j int) bool {
		return m[i].Date.Before(m[j].Date)
	})
}

// Convenience method to create a date based on a date with time 0:00.
//
// Do not enter trailing zeros for months and days because then Go will treat the numbers octal.
func CreateDate(year int, month int, day int) time.Time {
	result, _ := time.Parse(DateFormat, fmt.Sprintf("%04d-%02d-%02d", year, month, day))
	return result
}

// The date format used for all dates in this package.
const DateFormat = "2006-01-02"
