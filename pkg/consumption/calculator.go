package consumption

import (
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/fafeitsch/Horologium/pkg/util"
	"time"
)

type Parameters struct {
	Start    time.Time
	End      time.Time
	Readings domain.MeterReadings
	Plans    []domain.PricingPlan
}

// TODO:
// Validation-Method for Parameters:
// * beginning of the first plan must be before Parameters.Start
// * Plans must be continous
// * Plans must not overlap
// * Plans have to start at the first of the month

func Costs(params Parameters) (costs float64, totalConsumption float64) {
	plans := params.Plans
	costs = 0.0
	index := 0
	for index < len(plans) && plans[index].ValidTo != nil && (plans[index].ValidTo.Before(params.Start) || plans[index].ValidTo.Equal(params.Start)) {
		index = index + 1
	}
	totalConsumption = 0.0
	start := params.Start
	for index < len(plans) && plans[index].ValidFrom.Before(params.End) {
		plan := plans[index]
		end := params.End
		if plan.ValidTo != nil && plan.ValidTo.Before(end) {
			end = (*plan.ValidTo).Add(24 * time.Hour)
		}
		consumption := params.Readings.Consumption(start, end)
		totalConsumption = totalConsumption + consumption
		costs = costs + consumption*plan.UnitPrice + plan.BasePrice*float64(monthsBetween(start, end))
		if index < len(plans)-1 {
			start = *plans[index+1].ValidFrom
		}
		index = index + 1
	}
	return costs, totalConsumption
}

//Returns the number of different months between the given times,
//no matter whether they are contained completely or not. However, the end is exclusive. For example,
//2019-30-4 and 2019-05-01 contains one months; and 2019-30-4 and 2019-05-02 contains two months.
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

func MonthlyCosts(params Parameters) []Statistics {
	nextTime := func(date time.Time) time.Time {
		addedStart := date.AddDate(0, 1, 0)
		month := addedStart.Month()
		year := addedStart.Year()
		monthEnd, _ := time.Parse(util.DateFormat, fmt.Sprintf("%d-%02d-%02d", year, month, 1))
		return monthEnd
	}
	return granularCosts(params, nextTime)
}

func granularCosts(params Parameters, nextTime func(date time.Time) time.Time) []Statistics {
	result := make([]Statistics, 0, 0)
	monthStart := params.Start
	for monthStart.Before(params.End) {
		monthEnd := nextTime(monthStart)
		if params.End.Before(monthEnd) {
			monthEnd = params.End
		}
		params := Parameters{
			Start:    monthStart,
			End:      monthEnd,
			Readings: params.Readings,
			Plans:    params.Plans,
		}
		costs, cons := Costs(params)
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
