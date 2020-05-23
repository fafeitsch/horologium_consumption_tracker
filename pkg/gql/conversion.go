package gql

import (
	"fmt"
	"github.com/fafeitsch/Horologium/pkg/consumption"
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/fafeitsch/Horologium/pkg/util"
	"time"
)

func toQLSeries(series *domain.Series) *Series {
	if series == nil {
		return nil
	}
	return &Series{
		ID:   int(series.Id),
		Name: series.Name,
	}
}

func parseDates(dates ...*string) ([]*time.Time, error) {
	result := make([]*time.Time, 0, len(dates))
	for _, date := range dates {
		if date == nil {
			result = append(result, nil)
			continue
		}
		converted, err := time.Parse(util.DateFormat, *date)
		if err != nil {
			return nil, fmt.Errorf("could not parse \"%s\" as format YYYY-MM-DD", *date)
		}
		result = append(result, &converted)
	}
	return result, nil
}

func toQLPricingPlan(plan *domain.PricingPlan) *PricingPlan {
	if plan == nil {
		return nil
	}
	start := plan.ValidFrom.Format(util.DateFormat)
	var end *string
	if plan.ValidTo != nil {
		tmp := plan.ValidTo.Format(util.DateFormat)
		end = &tmp
	}
	return &PricingPlan{
		ID:        int(plan.Id),
		Name:      plan.Name,
		BasePrice: plan.BasePrice,
		UnitPrice: plan.UnitPrice,
		ValidFrom: start,
		ValidTo:   end,
		SeriesID:  int(plan.Series.Id),
	}
}

func toQlMeterReading(reading *domain.MeterReading) *MeterReading {
	if reading == nil {
		return nil
	}
	date := reading.Date.Format(util.DateFormat)
	return &MeterReading{
		ID:       int(reading.Id),
		Count:    float64(reading.Count),
		Date:     date,
		SeriesID: int(reading.Series.Id),
	}
}

func toQlStatistics(stastistics *consumption.Statistics) *Statistics {
	if stastistics == nil {
		return nil
	}
	start := stastistics.ValidFrom.Format(util.DateFormat)
	end := stastistics.ValidTo.Format(util.DateFormat)
	return &Statistics{
		ValidFrom:   start,
		ValidTo:     end,
		Costs:       stastistics.Costs,
		Consumption: stastistics.Consumption,
	}
}
