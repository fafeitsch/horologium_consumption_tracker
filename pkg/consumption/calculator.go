package consumption

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
	"math"
	"sort"
	"time"
)

type Parameters struct {
	Start    time.Time
	End      time.Time
	readings []domain.MeterReading
	plans    []domain.PricingPlan
}

func lastReadingBefore(date time.Time, readings []domain.MeterReading) domain.MeterReading {
	sort.Slice(readings, func(i, j int) bool {
		return readings[i].Date.Before(readings[j].Date)
	})
	index := 0
	for index < len(readings) && (readings[index].Date.Equal(date) || readings[index].Date.Before(date)) {
		index = index + 1
	}
	return readings[index-1]
}

func firstReadingAfter(date time.Time, readings []domain.MeterReading) domain.MeterReading {
	sort.Slice(readings, func(i, j int) bool {
		return readings[j].Date.Before(readings[i].Date)
	})
	index := 0
	for index < len(readings) && (readings[index].Date.Equal(date) || readings[index].Date.After(date)) {
		index = index + 1
	}
	return readings[index-1]
}

func Consumption(params Parameters) float64 {
	firstReading := lastReadingBefore(params.Start, params.readings)
	lastReading := firstReadingAfter(params.End, params.readings)
	differenceDays := math.Round(lastReading.Date.Sub(firstReading.Date).Hours() / 24)
	slope := (lastReading.Count - firstReading.Count) / differenceDays
	startDays := math.Round(params.Start.Sub(firstReading.Date).Hours() / 24)
	endDays := math.Round(params.End.Sub(firstReading.Date).Hours() / 24)
	valueStart := slope*startDays + firstReading.Count
	valueEnd := slope*endDays + firstReading.Count
	return valueEnd - valueStart
}

func Costs(params Parameters) float64 {
	return 0
}
