package gql

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/fafeitsch/Horologium/pkg/util"
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
