package gql

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
	orm "github.com/fafeitsch/Horologium/pkg/persistance"
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
	start := plan.ValidFrom.Format(orm.DateFormat)
	var end *string
	if plan.ValidTo != nil {
		tmp := plan.ValidTo.Format(orm.DateFormat)
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
