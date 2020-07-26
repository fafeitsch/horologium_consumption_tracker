package horologium

import (
	"fmt"
	"math"
	"sort"
	"time"
)

type MeterReading struct {
	Count float64
	Date  time.Time
}

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
	sort.Slice(m, func(i, j int) bool {
		return m[i].Date.Before(m[j].Date)
	})
	index := 0
	for index < len(m) && (m[index].Date.Equal(date) || m[index].Date.Before(date)) {
		index = index + 1
	}
	return m[index-1]
}

func (m MeterReadings) firstReadingAfter(date time.Time) MeterReading {
	sort.Slice(m, func(i, j int) bool {
		return m[j].Date.Before(m[i].Date)
	})
	index := 0
	for index < len(m) && (m[index].Date.Equal(date) || m[index].Date.After(date)) {
		index = index + 1
	}
	if index == 0 {
		return m[0]
	}
	return m[index-1]
}

func (m MeterReadings) Consumption(start time.Time, end time.Time) float64 {
	valueStart := m.interpolateValueAtDate(start)
	valueEnd := m.interpolateValueAtDate(end)
	return valueEnd - valueStart
}

func FormatDate(year int, month int, day int) time.Time {
	result, _ := time.Parse(DateFormat, fmt.Sprintf("%04d-%02d-%02d", year, month, day))
	return result
}

func FormatDatePtr(year int, month int, day int) *time.Time {
	result := FormatDate(year, month, day)
	return &result
}

const DateFormat = "2006-01-02"
