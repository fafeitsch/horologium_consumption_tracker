package orm

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
	"github.com/jinzhu/gorm"
	"time"
)

type seriesEntity struct {
	gorm.Model
	Name string
}

func (s *seriesEntity) toDomainSeries() domain.Series {
	return domain.Series{
		Id:   s.ID,
		Name: s.Name,
	}
}

func toSeriesEntity(series domain.Series) seriesEntity {
	return seriesEntity{
		Model: gorm.Model{ID: series.Id},
		Name:  series.Name,
	}
}

type meterReadingEntity struct {
	gorm.Model
	Count    float64
	Date     time.Time
	Series   seriesEntity
	SeriesID int
}

type pricingPlanEntity struct {
	gorm.Model
	Name      string
	BasePrice float64
	UnitPrice float64
	ValidFrom *time.Time
	ValidTo   *time.Time
	Series    seriesEntity
	SeriesID  uint `sql:"type:integer REFERENCES series(id) ON DELETE RESTRICT ON UPDATE CASCADE"`
}

func (p *pricingPlanEntity) toDomainPricingPlan() domain.PricingPlan {
	series := p.Series.toDomainSeries()
	return domain.PricingPlan{
		Id:        p.ID,
		Name:      p.Name,
		BasePrice: p.BasePrice,
		UnitPrice: p.UnitPrice,
		ValidFrom: p.ValidFrom,
		ValidTo:   p.ValidTo,
		Series:    series,
	}
}

func toPricingPlanEntity(plan domain.PricingPlan) pricingPlanEntity {
	series := toSeriesEntity(plan.Series)
	return pricingPlanEntity{
		Model:     gorm.Model{ID: plan.Id},
		Name:      plan.Name,
		BasePrice: plan.BasePrice,
		UnitPrice: plan.UnitPrice,
		ValidFrom: plan.ValidFrom,
		ValidTo:   plan.ValidTo,
		Series:    series,
		SeriesID:  plan.Series.Id,
	}
}
