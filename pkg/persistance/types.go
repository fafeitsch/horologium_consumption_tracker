package orm

import (
	"github.com/fafeitsch/Horologium/pkg/domain"
	"time"
)

const DateFormat = "2006-01-02"

type seriesEntity struct {
	Id   uint `gorm:"primary_key"`
	Name string
}

func (s *seriesEntity) toDomainSeries() domain.Series {
	return domain.Series{
		Id:   s.Id,
		Name: s.Name,
	}
}

func toSeriesEntity(series domain.Series) seriesEntity {
	return seriesEntity{
		Id:   series.Id,
		Name: series.Name,
	}
}

type meterReadingEntity struct {
	Id       uint `gorm:"primary_key"`
	Count    float64
	Date     time.Time
	Series   seriesEntity
	SeriesID uint `sql:"type:integer REFERENCES series(id) ON DELETE RESTRICT ON UPDATE CASCADE"`
}

func (m *meterReadingEntity) toDomainMeterReadingEntity() domain.MeterReading {
	series := m.Series.toDomainSeries()
	return domain.MeterReading{
		Id:     m.Id,
		Count:  m.Count,
		Date:   m.Date,
		Series: &series,
	}
}

func toMeterReadingEntity(reading domain.MeterReading) meterReadingEntity {
	series := toSeriesEntity(*reading.Series)
	return meterReadingEntity{
		Id:       reading.Id,
		Count:    reading.Count,
		Date:     reading.Date,
		Series:   series,
		SeriesID: series.Id,
	}
}

type pricingPlanEntity struct {
	Id        uint `gorm:"primary_key"`
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
		Id:        p.Id,
		Name:      p.Name,
		BasePrice: p.BasePrice,
		UnitPrice: p.UnitPrice,
		ValidFrom: p.ValidFrom,
		ValidTo:   p.ValidTo,
		Series:    &series,
	}
}

func toPricingPlanEntity(plan domain.PricingPlan) pricingPlanEntity {
	series := toSeriesEntity(*plan.Series)
	return pricingPlanEntity{
		Id:        plan.Id,
		Name:      plan.Name,
		BasePrice: plan.BasePrice,
		UnitPrice: plan.UnitPrice,
		ValidFrom: plan.ValidFrom,
		ValidTo:   plan.ValidTo,
		Series:    series,
		SeriesID:  plan.Series.Id,
	}
}
