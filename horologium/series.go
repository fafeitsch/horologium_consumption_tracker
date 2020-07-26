package horologium

import (
	"fmt"
	"time"
)

type PricingPlan struct {
	Name      string
	BasePrice float64
	UnitPrice float64
	ValidFrom *time.Time
	ValidTo   *time.Time
}

type PricingPlans []PricingPlan

type Series struct {
	Name          string
	PricingPlans  PricingPlans
	MeterReadings MeterReadings
}

// TODO:
// Validation-Method for Parameters:
// * beginning of the first plan must be before Parameters.Start
// * Plans must be continous
// * Plans must not overlap
// * Plans have to start at the first of the month

func (s *Series) CostsAndConsumption(start time.Time, end time.Time) (float64, float64) {
	costs := 0.0
	index := 0
	for index < len(s.PricingPlans) && s.PricingPlans[index].ValidTo != nil && (s.PricingPlans[index].ValidTo.Before(start) || s.PricingPlans[index].ValidTo.Equal(start)) {
		index = index + 1
	}
	totalConsumption := 0.0
	for index < len(s.PricingPlans) && s.PricingPlans[index].ValidFrom.Before(end) {
		plan := s.PricingPlans[index]
		tmpEnd := end
		if plan.ValidTo != nil && plan.ValidTo.Before(end) {
			tmpEnd = (*plan.ValidTo).Add(24 * time.Hour)
		}
		consumption := s.MeterReadings.Consumption(start, tmpEnd)
		totalConsumption = totalConsumption + consumption
		costs = costs + consumption*plan.UnitPrice + plan.BasePrice*float64(monthsBetween(start, tmpEnd))
		if index < len(s.PricingPlans)-1 {
			start = *s.PricingPlans[index+1].ValidFrom
		}
		index = index + 1
	}
	return costs, totalConsumption
}

// Returns the number of different months between the given times,
// no matter whether they are contained completely or not. However, the end is exclusive. For example,
// 2019-30-4 and 2019-05-01 contains one months; and 2019-30-4 and 2019-05-02 contains two months.
func monthsBetween(start time.Time, end time.Time) int {
	yearMonths := (end.Year() - start.Year()) * 12
	months := yearMonths + int(end.Month()-start.Month()+1)
	if end.Day() == 1 {
		return months - 1
	}
	return months
}

type Statistics struct {
	ValidFrom   time.Time
	ValidTo     time.Time
	Costs       float64
	Consumption float64
}

func (s *Series) MonthlyCosts(start time.Time, end time.Time) []Statistics {
	nextTime := func(date time.Time) time.Time {
		addedStart := date.AddDate(0, 1, 0)
		month := addedStart.Month()
		year := addedStart.Year()
		monthEnd, _ := time.Parse(DateFormat, fmt.Sprintf("%d-%02d-%02d", year, month, 1))
		return monthEnd
	}
	return s.granularCosts(start, end, nextTime)
}

func (s *Series) granularCosts(start time.Time, end time.Time, nextTime func(date time.Time) time.Time) []Statistics {
	result := make([]Statistics, 0, 0)
	monthStart := start
	for monthStart.Before(end) {
		monthEnd := nextTime(monthStart)
		if end.Before(monthEnd) {
			monthEnd = end
		}
		costs, cons := s.CostsAndConsumption(monthStart, monthEnd)
		stats := Statistics{
			ValidFrom:   monthStart,
			ValidTo:     monthEnd,
			Costs:       costs,
			Consumption: cons,
		}
		result = append(result, stats)
		monthStart = monthEnd
	}
	return result
}
