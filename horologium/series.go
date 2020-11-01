package horologium

import (
	"fmt"
	"time"
)

// Pricing plan defines the costs of one unit in a certain time interval.
// Additionally, a base price per month can be given, as well as a name.
type PricingPlan struct {
	Name      string     // a name for the pricing plan
	BasePrice float64    // the monthly base price
	UnitPrice float64    // the price for one unit
	ValidFrom *time.Time // the start time from which the pricing plan is valid
	ValidTo   *time.Time // the end time from which the pricing plan is not valid any more
}

// PricingPlans is a slice of pricing plans
type PricingPlans []PricingPlan

// Series combines pricing plans and meter readings. It offers methods to calculate the
// costs and consumption in a certain time interval.
type Series struct {
	Name              string        // the name of the series
	ConsumptionFormat string        // the format used for the consumption, e.g. %.2f kWh
	CurrencyFormat    string        // the format used for the currency, e.g. %.2f Euro,
	PricingPlans      PricingPlans  // the collection of pricing plans
	MeterReadings     MeterReadings // the collection of meter readings.
}

// CostsAndConsumption computes the costs and consumption of a certain series.
// Between to meter readings, the consumption is calculated as if it would change linearly.
//
// All dates are treated with time 0:00. Thus, the start day is always inclusive and the end day is exclusive.
//
// Please note that the monthly base price of pricing plans is only applied if the first day of the month is included
// in the time between start and end (see example).
//
// This method assumes the following about the series. It may panic or deliver wrong results if the bullet points
// are not fulfilled (in future, the will be a method to check these prerequisites)
// * the time spans defined in the pricing plans must not overlap and be continuous
// * the first pricing plans's validFrom must either be before start or be nil
// * the last pricing plan's validTo must eiether be after end or be nil
// * plans do have to start at the first of the month
// * the series must have both pricing plans and meter readings initialized
func (s *Series) CostsAndConsumption(start time.Time, end time.Time) (float64, float64) {
	costs := 0.0
	index := 0
	plans := make(PricingPlans, len(s.PricingPlans))
	copy(plans, s.PricingPlans)
	if len(plans) > 0 && plans[0].ValidFrom == nil {
		plans[0].ValidFrom = &start
	}
	for index < len(plans) && plans[index].ValidTo != nil && (plans[index].ValidTo.Before(start) || plans[index].ValidTo.Equal(start)) {
		index = index + 1
	}
	totalConsumption := 0.0
	for index < len(plans) && plans[index].ValidFrom.Before(end) {
		plan := plans[index]
		tmpEnd := end
		if plan.ValidTo != nil && plan.ValidTo.Before(end) {
			tmpEnd = *plan.ValidTo
		}
		consumption := s.MeterReadings.Consumption(start, tmpEnd)
		totalConsumption = totalConsumption + consumption
		costs = costs + consumption*plan.UnitPrice + plan.BasePrice*float64(monthsBetween(start, tmpEnd))
		if index < len(plans)-1 {
			start = *plans[index+1].ValidFrom
		}
		index = index + 1
	}
	return costs, totalConsumption
}

// Returns the number of different months between the given times,
// no matter whether they are contained completely or not. However, the end is exclusive. For example,
// 2019-30-4 and 2019-05-01 contains one month; and 2019-30-4 and 2019-05-02 contains two months.
func monthsBetween(start time.Time, end time.Time) int {
	yearMonths := (end.Year() - start.Year()) * 12
	months := yearMonths + int(end.Month()-start.Month()+1)
	if end.Day() == 1 {
		return months - 1
	}
	return months
}

// Statistics contain information about costs in consumption in a certain time interval.
type Statistics struct {
	ValidFrom         time.Time
	ValidTo           time.Time
	Costs             float64
	Consumption       float64
	ConsumptionFormat string
	CurrencyFormat    string
}

// FormatConsumption formats the consumption of the statistics
// according to the Statistics's ConsumptionFormat field.
// Uses a reasonable default format if the ConsumptionFormat is empty.
func (s *Statistics) FormatConsumption() string {
	if len(s.ConsumptionFormat) == 0 {
		return fmt.Sprintf("%.2f", s.Consumption)
	}
	return fmt.Sprintf(s.ConsumptionFormat, s.Consumption)
}

// CurrencyFormat formats the costs of the statistics
// according to the Statistic's CostsFormat field.
// Uses a reasonable default format if the CurrencyFormat is empty.
func (s *Statistics) FormatCosts() string {
	if len(s.CurrencyFormat) == 0 {
		return fmt.Sprintf("%.2f", s.Costs)
	}
	return fmt.Sprintf(s.CurrencyFormat, s.Costs)
}

// Monthly statistics computes costs and consumption for every month in the specified time span.
// The result is returned as Monthly Statistics, which can be rendered as table.
// The monthly statistics are sorted ascendingly with earliest months first.
func (s *Series) MonthlyStatistics(start time.Time, end time.Time) MonthlyStatistics {
	nextTime := func(date time.Time) time.Time {
		addedStart := date.AddDate(0, 1, 0)
		month := addedStart.Month()
		year := addedStart.Year()
		monthEnd := CreateDate(year, int(month), 1)
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
			ValidFrom:         monthStart,
			ValidTo:           monthEnd,
			Costs:             costs,
			Consumption:       cons,
			ConsumptionFormat: s.ConsumptionFormat,
			CurrencyFormat:    s.CurrencyFormat,
		}
		result = append(result, stats)
		monthStart = monthEnd
	}
	return result
}
