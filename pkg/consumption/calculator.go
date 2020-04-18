package consumption

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
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

func Consumption(params Parameters) float64 {
	valueStart := params.Readings.InterpolateValueAtDate(params.Start)
	valueEnd := params.Readings.InterpolateValueAtDate(params.End)
	return valueEnd - valueStart
}

func Costs(params Parameters) (costs float64, totalConsumption float64) {
	plans := params.Plans
	costs = 0.0
	index := 0
	for index < len(plans) && plans[index].ValidTo != nil && plans[index].ValidTo.Before(params.Start) {
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
		consumption := Consumption(Parameters{Start: start, End: end, Readings: params.Readings})
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
	yearMonths := (start.Year() - end.Year()) * 12
	months := yearMonths + int(end.Month()-start.Month()+1)
	if end.Day() == 1 {
		return months - 1
	}
	return months
}
