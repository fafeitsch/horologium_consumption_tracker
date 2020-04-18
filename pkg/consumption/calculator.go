package consumption

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
	"time"
)

type Parameters struct {
	Start    time.Time
	End      time.Time
	readings domain.MeterReadings
	plans    []domain.PricingPlan
}

// TODO:
// Validation-Method for Parameters:
// * beginning of the first plan must be before Parameters.Start
// * plans must be continous
// * plans must not overlap

func Consumption(params Parameters) float64 {
	valueStart := params.readings.InterpolateValueAtDate(params.Start)
	valueEnd := params.readings.InterpolateValueAtDate(params.End)
	return valueEnd - valueStart
}

func Costs(params Parameters) float64 {
	plans := params.plans
	result := 0.0
	index := 0
	for index < len(plans) && plans[index].ValidTo != nil && plans[index].ValidTo.Before(params.Start) {
		index = index + 1
	}
	check := 0.0
	start := params.Start
	for index < len(plans) && plans[index].ValidFrom.Before(params.End) {
		plan := plans[index]
		end := params.End
		if plan.ValidTo != nil && plan.ValidTo.Before(end) {
			end = (*plan.ValidTo).Add(24 * time.Hour)
		}
		consumption := Consumption(Parameters{Start: start, End: end, readings: params.readings})
		check = check + consumption
		result = result + consumption*plan.UnitPrice + plan.BasePrice*float64(monthsBetween(start, end))
		if index < len(plans)-1 {
			start = *plans[index+1].ValidFrom
		}
		index = index + 1
	}
	return result
}

//Returns the number of different months between the given times,
//no matter whether they are contained completely or not. For example,
//2019-30-4 and 2019-05-01 contain two months.
func monthsBetween(start time.Time, end time.Time) int {
	yearMonths := (start.Year() - end.Year()) * 12
	return yearMonths + int(end.Month()-start.Month()+1)
}
